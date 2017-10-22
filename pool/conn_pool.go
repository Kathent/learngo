package pool

import (
	"container/list"
	"sync"
	"context"
	"time"
	"sync/atomic"
)

const(
	SCALABLE_POOL_DECREASE_STATE_DIABLE = 0
	SCALABLE_POOL_DECREASE_STATE_ENABLE = 1
	SCALABLE_POOL_DEFAULT_CLEAN_FACTOR  = 0.8
	SCALABLE_POOL_MAX_IDLE_TIME         = 10
	SCALABLE_POOL_MAX_FORCE_RETURN_TIME = 6
)

type ScalablePool struct {
	freeCons *list.List  //空闲Conn容器
	usedCons *list.List	//正在使用的Conn容器
	sync.RWMutex

	initialSize int //初始化大小
	maxSize int //最大大小
	gen       Generator //obj工厂方法
	maxIdleTime int //最大闲置时间
	decState int32 //清理协程状态 0 未开启 1 开启
	factor float32 //清理系数 达到max*factor就停止清理
}

type Generator func() *Holder
type Options func(pool *ScalablePool)

func InitSize(initialSize int) Options{
	return func(pool *ScalablePool) {
		pool.initialSize = initialSize
	}
}

//WithMaxIdleTime 设置最大空闲时间
func MaxIdle(idleTime int) Options{
	return func(pool *ScalablePool) {
		pool.maxIdleTime = idleTime
	}
}

func CleanFactor(factor float32) Options{
	return func(pool *ScalablePool) {
		pool.factor = factor
	}
}

//NewScalablePool
//maxSize 最大容量(扩容后需要逐渐恢复)
//gen 元素生产工厂
//options 可选参数
func NewScalablePool(maxSize int, gen Generator, options ...Options) *ScalablePool{
	p := &ScalablePool{freeCons: list.New(),
		usedCons: list.New(),
		gen: gen,
		maxIdleTime: SCALABLE_POOL_MAX_IDLE_TIME,
		maxSize: maxSize,
		factor: SCALABLE_POOL_DEFAULT_CLEAN_FACTOR}

	for _, op := range options{
		op(p)
	}
	//预先丢几个进去
	for i := 0; i < p.initialSize; i++ {
		tmpVal := p.gen()
		if tmpVal != nil {
			p.freeCons.PushFront(tmpVal)
		}
	}
	return p
}

func (p *ScalablePool) Take() (t interface{}){
	defer func() {
		if t == nil {
			return
		}

		//无论是新生成的还是旧有的都要放入使用中列表中
		p.Lock()
		ele := p.usedCons.PushFront(t)
		p.Unlock()

		if holder, ok := t.(*Holder) ; ok {
			//更新访问时间
			holder.updateAccess()
			//做超期判断
			ctx, cancel := context.WithTimeout(context.Background(), SCALABLE_POOL_MAX_FORCE_RETURN_TIME* time.Second)
			holder.cl = cancel
			go func() {
				for {
					select {
					case <- ctx.Done():
						if dl, ok := ctx.Deadline(); ok {//超期了,还回去;
							if dl.Before(time.Now()) {
								//超期移除
								p.Lock()
								if p.usedCons.Remove(ele) != nil{
									//添加到空闲列表
									p.freeCons.PushFront(holder)
								}
								p.Unlock()
							}
						}
						return
					}
				}
			}()
		}
	}()

	//循环查找一个有效的元素;
	p.Lock()
	for {
		ele := p.freeCons.Front()
		if ele == nil || ele.Value == nil{
			break
		}

		t = ele.Value

		if h, ok := t.(*Holder); ok {
			if p.valid(h){
				break
			}else if h.abandoned != nil{
				//过期了或者失效了,要删除掉，size减1
				h.abandoned()
				p.freeCons.Remove(ele)
			}
		}
	}

	if t == nil {
		//当前list为空,新建连接
		t = p.gen()
		if t != nil {
			//是否超过了容量限制,超过了要启用定期减少元素协程
			if p.usedCons.Len() + p.freeCons.Len() >= p.maxSize {
				if atomic.CompareAndSwapInt32(&p.decState, SCALABLE_POOL_DECREASE_STATE_DIABLE,
					SCALABLE_POOL_DECREASE_STATE_ENABLE){
					go p.doCleanUp()
				}
			}

		}
	}

	p.Unlock()
	return t
}

func (p *ScalablePool) doCleanUp() {
	for {
		exitCh := make(chan byte, 1)
		func(){
			p.Lock()
			defer p.Unlock()

			ele := p.freeCons.Back()
			if ele == nil {
				return
			}

			if ele.Value == nil {
				return
			}

			if h, ok := ele.Value.(*Holder); ok {
				if !p.valid(h) {
					//失效了
					if h.abandoned != nil {
						h.abandoned()
					}

					p.freeCons.Remove(ele)
				}
			}

			cleanThresh := int(float32(p.maxSize) * p.factor)
			if cleanThresh < p.initialSize{
				cleanThresh = p.initialSize
			}
			if p.freeCons.Len() + p.usedCons.Len() < cleanThresh {
				exitCh <- 1
			}

		}()

		select {
			case <- exitCh:
				atomic.StoreInt32(&p.decState, SCALABLE_POOL_DECREASE_STATE_DIABLE)
				return
			default:
				//do nothing.

		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (p *ScalablePool) valid(h *Holder) bool{
	nowT := time.Now().Unix()
	if h.lastAccessed + int64(p.maxIdleTime) <= nowT{
		return p.freeCons.Len() + p.usedCons.Len() <= p.initialSize
	}

	if h.checkUseful != nil {
		return h.checkUseful()
	}else {
		return true
	}
}

func (p *ScalablePool) Return(val interface{}) bool{
	if val != nil {
		p.Lock()
		defer p.Unlock()
		//从使用中列表删除
		if remove(p.usedCons, val){
			//删除成功 说明之前在使用中列表中,需要取消定时器
			if h, ok := val.(*Holder) ; ok {
				h.cl()
			}

			//加入空闲列表
			return p.freeCons.PushFront(val) != nil
		}
	}

	return false
}

func remove(l *list.List, val interface{}) bool{
	for next := l.Front(); next != nil; next = next.Next() {
		if next.Value == val {
			l.Remove(next)
			return true
		}
	}
	return false
}
