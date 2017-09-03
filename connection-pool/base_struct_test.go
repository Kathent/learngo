package connection_pool

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_Poll_Get(t *testing.T){
	array := newArray(100)
	addVal := 1
	array.add(addVal)

	assert.Equal(t, addVal,array.take())
	assert.Equal(t, 0, array.len())
}

func Test_Poll_ZeroSize(t *testing.T){
	array := newArray(10)

	assert.Nil(t, array.take())
}

func TestNewList_Add(t *testing.T) {
	newList := NewList(0, nil)
	assert.Equal(t, false, newList.add(1))

	newList = NewList(1, nil)
	assert.Equal(t, true, newList.add(1))
}

func TestNewList_Take(t *testing.T) {
	newList := NewList(2, nil)

	assert.Equal(t, nil, newList.take())

	newVal := 1
	newList.add(newVal)
	assert.Equal(t, newVal, newList.take().(*node).data)
}

func TestNewList_Func(t *testing.T) {
	newList := NewList(2, func(n1, n2 *node) bool {
		return n1.data.(int) < n2.data.(int)
	})

	little := 1
	big := 3
	newList.add(big)
	newList.add(little)

	assert.Equal(t, little, newList.take().(*node).data)
	assert.Equal(t, big, newList.take().(*node).data)
}
