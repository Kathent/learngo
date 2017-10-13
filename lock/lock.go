package lock

import (
	"time"
	"github.com/go-redis/redis"
	"fmt"
	"crypto/rand"
	"sync"
	"encoding/base64"
	"errors"
	"github.com/alecthomas/log4go"
)

const(
	REDIS_LOCK_EXPIRE = time.Second * 6
	REDIS_LOCK_KEY_FMT = "redis.lock.%s"

	REDIS_DEFAULT_DURATION = 500 * time.Millisecond
	REDIS_DEFAULT_RETRY_TIMES = 10
	REDIS_DEFAULT_FACTOR = 0.01
	REDIS_DEFAULT_SLEEP_TIME = time.Millisecond * 300

	REDIS_LUA_SCRIPT = `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
)

var lockFailErr = errors.New("lock fail")

type Option func(l *RedisLock)

type RedisLock struct{
	id string
	cli []*redis.Client
	expire time.Duration
	retryTimes int
	threshold int
	lock sync.Mutex
	factor float64
	delay time.Duration
	val string
	log log4go.Logger
}


func NewLock(id string, cli []*redis.Client, log log4go.Logger, options ...Option) (*RedisLock, error) {
	lock := &RedisLock{
		id:     fmt.Sprintf(REDIS_LOCK_KEY_FMT, id),
		cli:    cli,
		expire: REDIS_DEFAULT_DURATION,
		threshold: len(cli) / 2 + 1,
		retryTimes: REDIS_DEFAULT_RETRY_TIMES,
		factor: REDIS_DEFAULT_FACTOR,
		log: log,
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
		start := time.Now()
		for _, cli := range l.cli{
			res, err := cli.SetNX(l.id, value, l.expire).Result()
			if err != nil {
				l.log.Warn("err is:%v", err)
				continue
			}else if !res {
				continue
			}

			count++
		}

		until := time.Now().Add(l.expire - time.Now().Sub(start) - time.Duration(int64(float64(l.expire)*l.factor)) +
			2*time.Millisecond)
		if count >= l.threshold && time.Now().Before(until){
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

	for i := 0; i < l.retryTimes; i++ {
		count := 0
		for _, cli := range l.cli{
			eval, err := cli.Eval(REDIS_LUA_SCRIPT, []string{l.id}, l.val).Result()
			if err != nil {
				l.log.Warn("err is:%v", err)
				continue
			}

			if eval == "1" {
				count++
			}
		}

		if count >= l.threshold {
			return
		}
	}
}