package pool

import (
	"math"
	"time"
	"context"
)

const(
	CLIENT_POOL_MAX_IDLE_TIME 	= 120
	CLIENT_POOL_MAX_TAKE_TIME 	= time.Second * 10
)

type GenFactory func() *Holder

type ClientPool struct {
	freeCons Container  //空闲Conn容器
	usedCons Container	//正在使用的Conn容器

	initialSize int //初始化大小
	gen       GenFactory //obj工厂方法
	maxIdleTime int //最大闲置时间
}

type ConnOption func(pool *ClientPool)

//WithInitialSize 设置初始化大小
func InitialSize(initialSize int) ConnOption{
	return func(pool *ClientPool) {
		pool.initialSize = initialSize
	}
}

//WithMaxIdleTime 设置最大空闲时间
func MaxIdleTime(idleTime int) ConnOption{
	return func(pool *ClientPool) {
		pool.maxIdleTime = idleTime
	}
}

//NewConnPool 新建连接池
//gen 工厂方法
//compare 比较器
//option 可选参数
func NewConnPool(gen GenFactory, f func(data1, data2 interface{}) bool, option ...ConnOption) *ClientPool {
	p := &ClientPool{freeCons: NewSafeStackList(math.MaxInt64, f),
				     usedCons: NewSafeStackList(math.MaxInt64, nil),
						  gen: gen,
			      maxIdleTime: CLIENT_POOL_MAX_IDLE_TIME}

	for _, op := range option{
		op(p)
	}
	//预先丢几个进去
	for i := 0; i < p.initialSize; i++ {
		tmpVal := p.gen()
		if tmpVal != nil {
			p.freeCons.add(tmpVal)
		}
	}
	return p
}

func (p *ClientPool)isValidObj(h *Holder) bool{
	nowT := time.Now().Unix()
	if h.lastAccessed + int64(p.maxIdleTime) <= nowT{
		return p.freeCons.lLen() <= p.initialSize
	}

	if h.checkUseful != nil {
		return h.checkUseful()
	}else {
		return true
	}
}

//Take 取连接
func (p *ClientPool) Take() (t interface{}){
	defer func() {
		if t == nil {
			return
		}

		//无论是新生成的还是旧有的都要放入使用中列表中
		p.usedCons.add(t)
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
								if p.usedCons.remove(holder) {
									//添加到空闲列表
									p.freeCons.add(holder)
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
		t := p.freeCons.take()
		if t == nil {
			break
		}

		if h, ok := t.(*Holder); ok {
			if p.isValidObj(h){
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

//Return 扔回连接池
func (p *ClientPool) Return(val interface{}) bool{
	if val != nil {
		//从使用中列表删除
		if p.usedCons.remove(val){
			//删除成功 说明之前在使用中列表中,需要取消定时器
			if h, ok := val.(*Holder) ; ok {
				h.cl()
			}

			//加入空闲列表
			b := p.freeCons.add(val)
			return b
		}
	}

	return false
}


