package pool

import (
	"sync"
	"time"
	"context"
)

//IPool 连接池接口
type IPool interface {
	Take() interface{}
	Return(interface{}) bool
}

//Container 容器接口
type Container interface {
	add(val interface{}) bool
	take() interface{}
	lLen() int
	cap() int
}

const(
	OBJECT_POOL_MAX_IDLE_TIME 	= 120
)

//ObjectPool 一个简单的连接池实现
type ObjectPool struct {
	initialSize int //初始化大小
	container Container  //实际容器
	m         *sync.RWMutex  //锁
	gen       func() interface{} //obj工厂方法
	maxIdleTime int //最大闲置时间
}

type Option func(pool *ObjectPool)

//WithInitialSize 设置初始化大小
func WithInitialSize(initialSize int) Option{
	return func(pool *ObjectPool) {
		pool.initialSize = initialSize
	}
}

//WithMaxIdleTime 设置最大空闲时间
func WithMaxIdleTime(idleTime int) Option{
	return func(pool *ObjectPool) {
		pool.maxIdleTime = idleTime
	}
}

//NewSimplePool 新建简单连接池
//c 连接容器
//gen 连接工厂方法
func NewSimplePool(c Container, gen func() interface{}, option ...Option) *ObjectPool{
	p := &ObjectPool{container:c, m: new(sync.RWMutex), gen: gen, maxIdleTime: OBJECT_POOL_MAX_IDLE_TIME}

	for _, op := range option{
		op(p)
	}
	//预先丢几个进去
	for i := 0; i < p.initialSize; i++ {
		tmpVal := p.gen()
		if tmpVal != nil {
			p.container.add(tmpVal)
		}
	}
	return p
}

//Take 取连接
func (p *ObjectPool) Take() interface{}{
	p.m.Lock()
	defer p.m.Unlock()

	for {
		t := p.container.take()
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
	return p.gen()
}

func (p *ObjectPool)isValidObj(h *Holder) bool{
	nowT := time.Now().Unix()
	if h.lastAccessed + int64(p.maxIdleTime) <= nowT{
		return p.container.lLen() <= p.initialSize
	}

	if h.checkUseful != nil {
		return h.checkUseful()
	}else {
		return true
	}
}

//Return 扔回连接池
func (p *ObjectPool) Return(val interface{}) bool{
	p.m.Lock()
	defer p.m.Unlock()
	return p.container.add(val)
}


func CompareHolder(n1, n2 interface{}) bool {
	if n1 == nil {
		return true
	}else if n2 == nil {
		return false
	}
	return n1.(*Holder).Less(n2.(*Holder))
}


type Factory func() interface{}

//GetHolder 从连接池中取连接 取完一定要Return.
//sleepTime 重试时的休眠时间
func (op *ObjectPool)GetHolder(ctx context.Context, sleepTime time.Duration) (*Holder, error){
	if op == nil {
		return nil, nil
	}else {
		for {
			select {
			case <- ctx.Done():
				return nil, ctx.Err()
			default:
				tmp := op.Take()
				if tmp != nil {
					//拿到了Holder,直接return.
					h, ok := tmp.(*Holder)
					if ok {
						return h, nil
					}
				}
				time.Sleep(sleepTime)
			}
		}
	}
}

//ReturnHolder 返回Grpc client holder
func (op *ObjectPool)ReturnHolder(h *Holder) {
	if h != nil{
		op.Return(h)
	}
}


