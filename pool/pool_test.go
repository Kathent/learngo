package pool

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewSimplePool(t *testing.T) {
	tmpList := NewList(2, nil)
	idleTime := 3
	tmpOp := NewSimplePool(tmpList, func() interface{} {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, WithMaxIdleTime(idleTime))

	tmpObj := tmpOp.Take()
	tmpVal := tmpObj.(*Holder).GetObj()
	assert.Equal(t, 1, tmpVal)
	assert.Equal(t, 0, tmpOp.container.lLen())
	assert.Equal(t, true, tmpOp.Return(tmpObj))

	//超时了取出的应该是新的
	time.Sleep(time.Duration(idleTime + 1) * time.Second)
	tmp := tmpOp.Take()
	assert.NotEqual(t, tmpObj, tmp)
	//不还回去,连续取俩个,只能取出来一个
	assert.NotNil(t, tmpOp.Take())
	assert.Nil(t, tmpOp.Take())
}

func BenchmarkObjectPool_Take2(b *testing.B) {
	tmpList := NewList(2, nil)
	idleTime := 3
	tmpOp := NewSimplePool(tmpList, func() interface{} {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, WithMaxIdleTime(idleTime))

	for i := 0; i < b.N; i++ {
		tmpOp.Return(tmpOp.Take())
	}
}

func BenchmarkObjectPool_Take(b *testing.B) {
	tmpList := NewList(2, nil)
	idleTime := 3
	tmpOp := NewSimplePool(tmpList, func() interface{} {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, WithMaxIdleTime(idleTime))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			func(){
				tmpOp.Return(tmpOp.Take())
			}()
		}
	})
}
