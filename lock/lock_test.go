package lock

import (
	"testing"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/alecthomas/log4go"
)

func TestNewLock(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Addr: "192.168.96.140:6379",

	})
	l1, err := NewLock("1", []*redis.Client{cli}, log4go.NewDefaultLogger(log4go.INFO))
	assert.Nil(t, err)
	lockErr := l1.Lock()
	assert.Nil(t, lockErr)

	l2, err := NewLock("1", []*redis.Client{cli}, log4go.NewDefaultLogger(log4go.INFO))
	assert.Nil(t, err)
	assert.NotNil(t, l2.Lock())

	l1.Unlock()
	assert.Nil(t, l2.Lock())
	l2.Unlock()
}

func BenchmarkNewLock(b *testing.B) {
	b.StopTimer()
	cli := redis.NewClient(&redis.Options{
		Addr: "192.168.96.140:6379",

	})

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l1, err := NewLock("1", []*redis.Client{cli}, log4go.NewDefaultLogger(log4go.INFO))
			if err != nil{
				b.Log(err)
				continue
			}

			l1.Lock()
			l1.Unlock()
		}
	})
}
