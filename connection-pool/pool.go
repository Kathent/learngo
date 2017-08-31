package connection_pool

type pool interface {
	Take() interface{}
	Return(interface{}) interface{}
}

type ObjectPool struct {
	container poolContainer
}

func (p *ObjectPool) Take() interface{}{
	return p.container.take()
}

func (p *ObjectPool) Return(val interface{}) interface{}{
	return p.container.add(val)
}




