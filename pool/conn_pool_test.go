package pool


import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewScalablePool(t *testing.T) {
	cp := NewScalablePool(5, func() *Holder {
		newVal := &struct {}{}
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, CleanFactor(1))

	tk := cp.Take()
	assert.NotNil(t, tk)
	assert.Equal(t, cp.usedCons.Len(), 1)
	//sleep 一段时间 应该强制回收了
	time.Sleep(time.Second * 11)
	assert.Equal(t, cp.usedCons.Len(), 0)
	assert.Equal(t, cp.freeCons.Len(), 1)
	//强制回收之后 Return失败
	assert.Equal(t, false, cp.Return(tk))
	assert.Equal(t, tk, cp.Take())
}

func TestNewScalablePool2(t *testing.T) {
	cp := NewScalablePool(5, func() *Holder {
		newVal := &struct {}{}
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, CleanFactor(1))

	//连取7个 超出上限 检查时间到了之后有没有清除掉
	arr := make([]interface{}, 0)
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	arr = append(arr, cp.Take())
	assert.Equal(t, cp.decState, int32(1))
	time.Sleep(time.Second * 12)
	assert.Equal(t, 4, cp.freeCons.Len())
}

func BenchmarkNewScalablePool(b *testing.B) {
	cp := NewScalablePool(5, func() *Holder {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tk := cp.Take()
			if tk != nil {
				cp.Return(tk)
			}
		}
	})

	b.Logf("used:%d, free:%d, state:%d", cp.usedCons.Len(), cp.freeCons.Len(), cp.decState)
}