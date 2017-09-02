package connection_pool

import "time"

type holder struct{
	obj interface{}
	lastAccessed int64
	checkUseful func() bool
}

func (h *holder) GetObj() interface{}{
	if h.obj != nil {
		h.lastAccessed = time.Now().Unix()
	}
	return h.obj
}
