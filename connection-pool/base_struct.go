package connection_pool

import "sync"

type safeArray struct {
	arr []interface{}
	lock sync.RWMutex
}

type poolContainer interface {
	remove(val interface{}) (pre interface{})
	add(val interface{}) bool
	take() interface{}
}

func newArray(size int) safeArray{
	return safeArray{
		arr: make([]interface{}, 0, size),
		lock: sync.RWMutex{},
	}
}

func (s *safeArray) remove(val interface{}) (pre interface{}){
	s.lock.Lock()
	defer s.lock.Unlock()

	tmp := make([]interface{}, len(s.arr))
	index := 0
	var res interface{}
	for _, v := range s.arr{
		if val != v {
			tmp[index] = v
			index++
		}else {
			res = v
		}
	}

	s.arr = tmp
	return res
}

func (s *safeArray) add(val interface{}) bool{
	s.lock.RLock()
	length := len(s.arr)
	capacity := cap(s.arr)
	s.lock.RUnlock()

	if length >= capacity{
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.arr = append(s.arr, val)
	return true
}

func (s *safeArray) take() interface{} {
	s.lock.RLock()
	length := len(s.arr)
	s.lock.RUnlock()

	if length <= 0 {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.arr) <= 0{
		return nil
	}

	res := s.arr[0]
	s.arr = s.arr[1:]
	return res
}