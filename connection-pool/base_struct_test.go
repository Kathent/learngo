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

func Test_Poll_Remove(t *testing.T){
	array := newArray(100)

	for i := 0; i < 100; i++ {
		array.add(i)
	}

	assert.Equal(t, false, array.add(1))

	for i := 0; i < 100; i++ {
		assert.Equal(t, i, array.remove(i))
	}
}
