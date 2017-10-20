package pool

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewConnPool(t *testing.T) {
	cp := NewConnPool(func() *Holder {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, nil)

	tk := cp.Take()
	assert.NotNil(t, tk)
	assert.Equal(t, cp.usedCons.lLen(), 1)
	//sleep 一段时间 应该强制回收了
	time.Sleep(time.Second * 11)
	assert.Equal(t, cp.usedCons.lLen(), 0)
	assert.Equal(t, cp.freeCons.lLen(), 1)
	//强制回收之后 Return失败
	assert.Equal(t, false, cp.Return(tk))
	assert.Equal(t, tk, cp.Take())
}

func BenchmarkNewConnPool(b *testing.B) {
	cp := NewConnPool(func() *Holder {
		newVal := 1
		fmt.Println(fmt.Sprintf("create new val:%d", newVal))
		return NewHolder(newVal, nil, func() bool {
			fmt.Println(fmt.Sprintf("close..%d", newVal))
			return true
		})
	}, nil)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tk := cp.Take()
			if tk != nil {
				cp.Return(tk)
			}
		}
	})
}
