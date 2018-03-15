package go_micro_learn

import (
	"testing"
	"google.golang.org/grpc"
	"learngo/proto1"
	"context"
	"github.com/alecthomas/log4go"
	"time"
	"sync"
	"fmt"
	"runtime"
	"icsoclib/pool"
)

func BenchmarkHello(b *testing.B) {
	b.StopTimer()
	conn, err := grpc.DialContext(context.Background(),"127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	client := proto.NewGreeterSClient(conn)

	b.StartTimer()
	count := 0
	N := 100
	runtime.GOMAXPROCS(2)
	for i := 0; i < b.N; i++ {
		group := sync.WaitGroup{}
		group.Add(N)
		for j := 0; j < N; j++ {
			go func(val int) {
				response, e := client.Hello(context.Background(), &proto.HelloRequestS{Name: fmt.Sprintf("%d", val)})
				if e != nil {
					log4go.Info("err:%v", e)
					count++
				}else {
					log4go.Info("result:%s", response.Greeting)
				}

				group.Done()
			}(j)
		}

		group.Wait()
	}
	log4go.Info("count:%d", count)
}

func BenchmarkHello2(b *testing.B) {
	b.StopTimer()
	gen := func() *pool.Holder {
		conn, err := grpc.DialContext(context.Background(),"127.0.0.1:8888", grpc.WithInsecure())
		if err != nil {
			return nil
		}
		client := proto.NewGreeterSClient(conn)
		return pool.NewHolder(client, nil, func() bool {
			e := conn.Close()
			return e != nil
		})
	}
	sp := pool.NewScalablePool(10, gen)

	b.StartTimer()
	count := 0
	N := 100
	runtime.GOMAXPROCS(2)
	for i := 0; i < b.N; i++ {
		group := sync.WaitGroup{}
		group.Add(N)
		for j := 0; j < N; j++ {
			go func(val int) {
				cl, _ := sp.GetHolder(context.Background(), time.Microsecond * 20)
				response, e := cl.GetObj().(proto.GreeterSClient).Hello(context.Background(), &proto.HelloRequestS{Name: fmt.Sprintf("%d", val)})
				if e != nil {
					log4go.Info("err:%v", e)
					count++
				}else {
					log4go.Info("result:%s", response.Greeting)
				}

				group.Done()
			}(j)
		}

		group.Wait()
	}

	log4go.Info("count:%d", count)
}
