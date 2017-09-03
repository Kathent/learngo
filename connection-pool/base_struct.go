package connection_pool

import (
	"sync"
	"unsafe"
)

type safeArray struct {
	arr []interface{}
	lock sync.RWMutex
}

type poolContainer interface {
	//remove(val interface{}) (pre interface{})
	add(val interface{}) bool
	take() interface{}
	len() int
}

func newArray(size int) *safeArray{
	return &safeArray{
		arr: make([]interface{}, 0, size),
		lock: sync.RWMutex{},
	}
}

func (s *safeArray) len() int{
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.arr)
}

func (s *safeArray)forEach(f func(interface{})){
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, val := range s.arr {
		f(val)
	}
}

//func (s *safeArray) remove(val interface{}) (pre interface{}){
//	s.lock.Lock()
//	defer s.lock.Unlock()
//
//	tmp := make([]interface{}, 0, len(s.arr))
//	index := 0
//	for _, v := range s.arr{
//		if val != v {
//			tmp = append(tmp, v)
//			index++
//		}else {
//			pre = v
//		}
//	}
//
//	s.arr = tmp
//
//	return
//}

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

type node struct {
	next *node
	data interface{}
}

type safeList struct {
	head *node
	length int
	lock sync.RWMutex
	capacity int
	compare func(node1, node2 *node) bool
}

func defaultCompare(node1, node2 *node) bool{
	return *(*uint64)(unsafe.Pointer(&node1.data)) < *(*uint64)(unsafe.Pointer(&node2.data))
}

func NewList(size int, f func(n1,n2 *node) bool) *safeList{
	tmp:= &safeList{}
	tmp.capacity = size
	tmp.lock = sync.RWMutex{}
	tmp.head = new(node)
	if f == nil {
		tmp.compare = defaultCompare
	}else {
		tmp.compare = f
	}
	return tmp
}

func (s *safeList) add(val interface{}) bool {
	s.lock.RLock()
	length := s.length
	capacity := s.capacity
	s.lock.RUnlock()
	if length >= capacity {
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	if s.length >= s.capacity {
		return false
	}

	newVal := &node{data: val}

	tmpNext := s.head
	for tmpNext != nil && tmpNext.next != nil && s.compare(tmpNext.next, newVal){
		tmpNext = tmpNext.next
	}

	tmpNext.next, newVal.next = newVal, tmpNext.next
	s.length++
	return true
}

func (s *safeList) take() interface{} {
	if s.len() <= 0 {
		return nil
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.length <= 0 {
		return false
	}

	res := s.head.next
	s.head = s.head.next
	s.length--
	return res
}

func (s *safeList) len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.length
}


