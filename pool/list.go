package pool

type node struct {
	next *node
	data interface{}
}

type sortList struct {
	head *node
	length int
	capacity int
	compare func(data1, data2 interface{}) bool
}

//func defaultCompare(data1, data2 interface{}) bool{
//	return *(*uint64)(unsafe.Pointer(&data1)) < *(*uint64)(unsafe.Pointer(&data2))
//}

//NewList 新建safeList
//size 容量
//f 比较函数 值小的在list前面
func NewList(capacity int, f func(data1, data2 interface{}) bool) *sortList {
	tmp:= &sortList{}
	tmp.capacity = capacity
	tmp.head = new(node)
	tmp.compare = f
	return tmp
}

func (s *sortList) add(val interface{}) bool {
	length := s.length
	capacity := s.capacity
	if length >= capacity {
		return false
	}

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

func (s *sortList) lLen() int {
	return s.length
}

func (s *sortList) take() interface{} {
	if s.lLen() <= 0 {
		return nil
	}

	if s.length <= 0 {
		return false
	}

	res := s.head.next
	s.head = s.head.next
	s.length--
	return res.data
}

func (s *sortList) cap() int{
	return s.capacity
}