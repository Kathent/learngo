package lock

import (
	"time"
	"github.com/go-redis/redis"
	"fmt"
	"crypto/rand"
	"sync"
	"encoding/base64"
	"errors"
)

const(
	REDIS_LOCK_EXPIRE = time.Second * 6
	REDIS_LOCK_KEY_FMT = "redis.lock.%s"

	REDIS_DEFAULT_DURATION = time.Second * 3
	REDIS_DEFAULT_RETRY_TIMES = 10
	REDIS_DEFAULT_FACTOR = 0.01
	REDIS_DEFAULT_SLEEP_TIME = time.Millisecond * 300

	REDIS_LUA_SCRIPT = `
		if redis.call("get", KEYS[1]) == ARGV[0] then
			return redis.call("del", KEYS[1])
		else
			return 0
	`

	REDIS_REPLY_OK = "OK"
)

var lockFailErr = errors.New("lock fail")

type Option func(l *RedisLock)

type RedisLock struct{
	id string
	cli []*redis.Client
	expire int64
	retryTimes int
	threshold int
	lock *sync.Mutex
	factor float64
	delay time.Duration
	val string
}


func NewLock(id string, cli []*redis.Client, options ...Option) (*RedisLock, error) {
	lock := &RedisLock{
		id:     fmt.Sprintf(REDIS_LOCK_KEY_FMT, id),
		cli:    cli,
		expire: time.Now().Unix() + int64(REDIS_DEFAULT_DURATION),
		threshold: len(cli) / 2 + 1,
		retryTimes: REDIS_DEFAULT_RETRY_TIMES,
		lock: new(sync.Mutex),
	}

	for _, o := range options {
		o(lock)
	}

	return lock, nil
}

func (l *RedisLock) Lock() error{
	l.lock.Lock()
	defer l.lock.Unlock()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	value := base64.StdEncoding.EncodeToString(b)
	for i:= 0; i < l.retryTimes; i++ {

		count := 0
		for _, cli := range l.cli{
			expireTime := time.Duration(float64(l.expire - time.Now().Unix()) * (1 + l.factor))
			res, err := cli.Set(l.id, value, expireTime).Result()
			if err != nil {
				continue
			}else if res != REDIS_REPLY_OK {
				continue
			}

			count++
		}

		if count >= l.threshold && time.Now().Unix() < l.expire{
			l.val = value
			return nil
		}

		for _, cli := range l.cli {
			cli.Eval(REDIS_LUA_SCRIPT, []string{l.id}, value)
		}

		time.Sleep(l.delay)
	}

	return lockFailErr
}

func (l *RedisLock)Unlock() {
	l.lock.Lock()
	defer l.lock.Unlock()

	for _, cli := range l.cli{
		cli.Eval(REDIS_LUA_SCRIPT, []string{l.id}, l.val)
	}
}