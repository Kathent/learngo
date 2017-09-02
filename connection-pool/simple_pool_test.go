package connection_pool

import (
	"testing"
	"time"
	"math/rand"
	"github.com/stretchr/testify/assert"
)

func TestObjectPool_Take(t *testing.T) {
	c := newArray(10)
	p := NewSimplePool(c, func() interface{} {
		return &holder{obj:1, checkUseful: func() bool{
			time.Sleep(time.Second)
			return rand.Intn(2) == 1
		}, }
	})

	val := p.Take()
	assert.NotNil(t, val)
	assert.Equal(t, true, p.Return(val))
	assert.Equal(t, val, p.Take())
	assert.Equal(t, true, p.Return(val))
	time.Sleep(time.Second * 13)
	assert.NotEqual(t, val, p.Take())
}


func init(){
	rand.Seed(time.Now().Unix())
}
