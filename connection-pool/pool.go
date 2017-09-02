package connection_pool

type pool interface {
	Take() interface{}
	Return(interface{}) bool
}




