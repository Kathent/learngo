package pool

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewSafeSortList(t *testing.T) {
	tmpList := NewSafeStackList(10, func(n1, n2 interface{}) bool {
		return n1.(int) < n2.(int)
	})

	val1, val2, val3, val4, val5, val6, val7, val8, val9, val10 := 1, 2, 3,4 ,5, 6,7 ,8,9,10
	//expected is 1,2 ,4,5,6,7
	tmpList.add(val1)
	tmpList.add(val10)
	tmpList.add(val2)
	tmpList.add(val9)
	tmpList.add(val3)
	tmpList.add(val4)
	tmpList.add(val6)
	tmpList.add(val7)
	tmpList.add(val5)
	tmpList.add(val8)
	assert.Equal(t, val1, tmpList.take())
	assert.Equal(t, val2, tmpList.take())
	assert.Equal(t, val3, tmpList.take())
	assert.Equal(t, val4, tmpList.take())
	assert.Equal(t, val5, tmpList.take())
	assert.Equal(t, val6, tmpList.take())
	assert.Equal(t, val7, tmpList.take())
	assert.Equal(t, val8, tmpList.take())
	assert.Equal(t, val9, tmpList.take())
	assert.Equal(t, val10, tmpList.take())
}

func TestNewSafeSortList_Add(t *testing.T) {
	tmpList := NewSafeStackList(10, nil)
	tmpList.add(1)
	assert.Equal(t, tmpList.lLen(), 1)
}

func TestNewSafeSortList_Take(t *testing.T) {
	tmpList := NewSafeStackList(10, nil)
	tmpList.add(1)
	assert.Equal(t, tmpList.lLen(), 1)
	assert.Equal(t, 1, tmpList.take())
	assert.Equal(t, tmpList.lLen(), 0)
	assert.Nil(t, tmpList.take())
}

func TestNewSafeSortList_Compare(t *testing.T) {
	tmpList := NewSafeStackList(10, func(n1, n2 interface{}) bool {
		return n1.(int) < n2.(int)
	})

	val1, val2, val3, val4, val5, val6, val7, val8, val9, val10 := 1, 2, 3,4 ,5, 6,7 ,8,9,10
	//expected is 1,2 ,4,5,6,7
	tmpList.add(val1)
	tmpList.add(val10)
	tmpList.add(val2)
	tmpList.add(val9)
	tmpList.add(val3)
	tmpList.add(val4)
	tmpList.add(val6)
	tmpList.add(val7)
	tmpList.add(val5)
	tmpList.add(val8)
	assert.Equal(t, val1, tmpList.take())
	assert.Equal(t, val2, tmpList.take())
	assert.Equal(t, val3, tmpList.take())
	assert.Equal(t, val4, tmpList.take())
	assert.Equal(t, val5, tmpList.take())
	assert.Equal(t, val6, tmpList.take())
	assert.Equal(t, val7, tmpList.take())
	assert.Equal(t, val8, tmpList.take())
	assert.Equal(t, val9, tmpList.take())
	assert.Equal(t, val10, tmpList.take())
}

func TestNewSafeSortList_Remove(t *testing.T) {
	tmpList := NewSafeStackList(10, func(n1, n2 interface{}) bool {
		return n1.(int) < n2.(int)
	})

	tmpList.add(1)
	tmpList.add(5)
	tmpList.add(4)
	tmpList.add(3)
	tmpList.add(2)
	tmpList.remove(5)
	assert.Equal(t, tmpList.lLen(), 4)
	assert.Equal(t, tmpList.take(), 1)
	assert.Equal(t, tmpList.take(), 2)
	assert.Equal(t, tmpList.take(), 3)
	assert.Equal(t, tmpList.take(), 4)
	assert.Equal(t, tmpList.lLen(), 0)

	tmpList.add(6)
	assert.Equal(t, false, tmpList.remove(5))
	assert.Equal(t, tmpList.take(), 6)
}
