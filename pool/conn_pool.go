package pool

import (
	"container/list"
	"sync"
	"context"
	"time"
)

type ScalablePool struct {
	freeCons *list.List  //空闲Conn容器
	usedCons *list.List	//正在使用的Conn容器
	sync.RWMutex

	initialSize int //初始化大小
	maxSize int //最大大小
	gen       Generator //obj工厂方法
	maxIdleTime int //最大闲置时间
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

//NewScalablePool
//maxSize 最大容量(扩容后需要逐渐恢复)
//gen 元素生产工厂
//options 可选参数
func NewScalablePool(maxSize int, gen Generator, options ...Options) *ScalablePool{
	p := &ScalablePool{freeCons: list.New(),
		usedCons: list.New(),
		gen: gen,
		maxIdleTime: CLIENT_POOL_MAX_IDLE_TIME,
		maxSize: maxSize}

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
		ele := p.usedCons.PushFront(t)
		if holder, ok := t.(*Holder) ; ok {
			//做超期判断
			ctx, cancel := context.WithTimeout(context.Background(), CLIENT_POOL_MAX_TAKE_TIME)
			holder.cl = cancel
			go func() {
				for {
					select {
					case <- ctx.Done():
						if dl, ok := ctx.Deadline(); ok {//超期了,还回去;
							if dl.Before(time.Now()) {
								//超期移除
								if p.usedCons.Remove(ele) != nil{
									//添加到空闲列表
									p.freeCons.PushFront(holder)
								}
							}
						}
						return
					}
				}
			}()
		}
	}()

	for {
		ele := p.freeCons.Front()
		if ele == nil {
			break
		}

		t = ele.Value
		if t == nil {
			break
		}

		if h, ok := t.(*Holder); ok {
			if p.valid(h){
				return t
			}else if h.abandoned != nil{
				//过期了或者失效了,要删除掉，size减1
				h.abandoned()
			}
		}
	}
	//当前list为空,新建连接
	t = p.gen()
	return t
}

func (p *ScalablePool) valid(h *Holder) bool{
	nowT := time.Now().Unix()
	if h.lastAccessed + int64(p.maxIdleTime) <= nowT{
		return p.freeCons.Len() <= p.initialSize
	}

	if h.checkUseful != nil {
		return h.checkUseful()
	}else {
		return true
	}
}

func (p *ScalablePool) Return() (val interface{}){
	if val != nil {
		//从使用中列表删除
		if remove(p.usedCons, val){
			//删除成功 说明之前在使用中列表中,需要取消定时器
			if h, ok := val.(*Holder) ; ok {
				h.cl()
			}

			//加入空闲列表
			b := p.freeCons.PushFront(val)
			return b
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
