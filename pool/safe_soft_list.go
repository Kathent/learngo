package pool

import "sync"

type safeSortList struct {
	head *node
	length int
	capacity int
	compare func(data1, data2 interface{}) bool
	sync.RWMutex
}

//func defaultCompare(data1, data2 interface{}) bool{
//	return *(*uint64)(unsafe.Pointer(&data1)) < *(*uint64)(unsafe.Pointer(&data2))
//}

//NewList 新建safeList
//size 容量
//f 比较函数 值小的在list前面
func NewSafeSortList(capacity int, f func(data1, data2 interface{}) bool) *safeSortList {
	tmp:= &safeSortList{}
	tmp.capacity = capacity
	tmp.head = new(node)
	tmp.compare = f
	return tmp
}

func (s *safeSortList) add(val interface{}) bool {
	s.RLock()
	length := s.length
	capacity := s.capacity
	s.RUnlock()
	if length >= capacity {
		return false
	}

	s.Lock()
	defer s.Unlock()

	if s.length >= s.capacity {
		return false
	}

	newVal := &node{data: val}

	tmpNext := s.head
	for tmpNext != nil && tmpNext.next != nil && (s.compare == nil || s.compare(tmpNext.next.data, newVal.data)){
		tmpNext = tmpNext.next
	}

	tmpNext.next, newVal.next = newVal, tmpNext.next
	s.length++
	return true
}

func (s *safeSortList) lLen() int {
	s.RLock()
	defer s.RUnlock()
	return s.length
}

func (s *safeSortList) take() interface{} {
	if s.lLen() <= 0 {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	if s.length <= 0 {
		return nil
	}

	res := s.head.next
	s.head = s.head.next
	s.length--
	return res.data
}

func (s *safeSortList) cap() int{
	s.RLock()
	defer s.RUnlock()
	return s.capacity
}