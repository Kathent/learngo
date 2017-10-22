package pool

import (
	"time"
	"context"
)

//Holder 对象Holder struct
type Holder struct{
	obj interface{}
	//上次访问时间
	lastAccessed int64
	//有效性检查方法
	checkUseful func() bool
	//对象销毁方法
	abandoned func() bool
	//定时回收方法
	cl context.CancelFunc
}

//Less 对象比较方法 按访问时间比较
func (h *Holder)Less(val *Holder) bool{
	return h.lastAccessed < val.lastAccessed
}

//GetObj 获取实际对象 并更新访问时间
func (h *Holder) GetObj() interface{}{
	return h.obj
}

func (h *Holder) updateAccess() {
	h.lastAccessed = time.Now().Unix()
}

//NewHolder 对象holder struct
//val 实际池中对象
//fs[0] 对象有效性检查方法
//fs[1] 对象销毁方法
func NewHolder(val interface{}, fs ...func()bool) *Holder{
	res := new(Holder)
	res.obj = val
	res.lastAccessed = time.Now().Unix()
	if len(fs) >= 1 {
		res.checkUseful = fs[0]
	}
	if len(fs) >= 2 {
		res.abandoned = fs[1]
	}
	return res
}
