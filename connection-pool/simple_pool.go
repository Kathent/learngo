package connection_pool

import (
	"sync"
	"time"
	"log"
	"fmt"
	"reflect"
)

const(
	OBJECT_POOL_SELF_CHECK_TIME = time.Second * 3
	OBJECT_POOL_MIN_CHECK_TIME = time.Second * 1
	OBJECT_POOL_MAX_IDLE_TIME 	= time.Second * 10
)

type ObjectPool struct {
	container poolContainer
	m         sync.Locker
	gen       func() interface{}
}

func NewSimplePool(c poolContainer, gen func() interface{}) *ObjectPool{
	pool := &ObjectPool{container:c, m: new(sync.Mutex), gen: gen}

	go pool.selfCheck()
	return pool
}

func (p *ObjectPool) selfCheck() {
	size := p.container.len()
	nowT := time.Now().Unix()
	for i:= 0; i < size; i++{
		tmp := p.container.take()
		if tmp != nil {
			if h, ok := tmp.(*holder); ok {
				if h.lastAccessed + int64(OBJECT_POOL_MIN_CHECK_TIME) >= nowT{
					//accessed, return
					p.container.add(tmp)
					continue
				}else if h.lastAccessed + int64(OBJECT_POOL_MAX_IDLE_TIME) <= nowT{
					//exceed max idle time
					continue
				}else if h.checkUseful != nil && h.checkUseful() {
					p.container.add(tmp)
				}
			}else {
				log.Println(fmt.Sprintf("selfCheck element is not a holder...%v %v", tmp, reflect.TypeOf(tmp)))
			}
		}
	}

	time.AfterFunc(OBJECT_POOL_SELF_CHECK_TIME, p.selfCheck)
}

func (p *ObjectPool) Take() interface{}{
	t := p.container.take()
	if t != nil {
		return t
	}

	p.m.Lock()
	defer p.m.Unlock()

	tt := p.container.take()
	if tt != nil{
		return tt
	}else {
		newVal := p.gen()
		if p.Return(newVal){
			return newVal
		}

		return nil
	}
}

func (p *ObjectPool) Return(val interface{}) bool{
	return p.container.add(val)
}
