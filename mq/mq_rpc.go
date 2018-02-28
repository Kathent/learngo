package main

import (
	"github.com/streadway/amqp"
	"flag"
	"time"
	"net"
	"icsoclib/rabbitmq"
	"fmt"
	//"encoding/json"
	"github.com/alecthomas/log4go"
	"go.uber.org/atomic"
	"context"
	"github.com/orcaman/concurrent-map"
)

const (
	SLEEP_TIME = time.Second * 2
)

type callBackInfo struct {
	id string
	df defaultPubFuture
}

type PubFuture interface {
	Get() (interface{}, error)
	GetWithContext(ctx context.Context) (interface{}, error)
}

type defaultPubFuture struct {
	c chan interface{}
}

func (df *defaultPubFuture) Get() (interface{}, error){
	val := <- df.c
	return val, nil
}

func (df *defaultPubFuture) GetWithContext(ctx context.Context) (interface{}, error) {
	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	case res := <- df.c:
		return res, nil
	}
}

type PubRpc interface {
	Publish(key string, body []byte) error
}

type publisher struct {
	Uri          string
	Config       amqp.Config
	Exchange     string
	ExchangeType string
	RouteKey     string
	ReplyTo 	 string
	size 		 int

	channel *amqp.Channel
	conn    *amqp.Connection
	errChan chan *amqp.Error
	msgChan chan amqp.Delivery
	id atomic.Int64
	callBacks cmap.ConcurrentMap
}

func NewPublisher(uri, exchange, exchangeType, key, replyTo string, config amqp.Config, size int) PubRpc{
	return &publisher{Uri: uri, Config: config, Exchange: exchange, ExchangeType: exchangeType,
		RouteKey: key, ReplyTo: replyTo, size: size}
}

type Message struct {
	MessageId int64
	Msg string
}

func (publisher *publisher) close(){
	if publisher.channel != nil {
		publisher.channel.Close()
	}

	if publisher.conn != nil {
		publisher.channel.Close()
	}

	if publisher.errChan != nil {
		close(publisher.errChan)
	}

	if publisher.msgChan != nil {
		close(publisher.msgChan)
	}

	if publisher.callBacks != nil {
		publisher.callBacks = nil
	}
}

func (publisher *publisher) init() error{
	publisher.errChan = make(chan *amqp.Error)
	publisher.callBacks = cmap.New()

	if publisher.size < 10 {
		publisher.size = 10
	}
	publisher.msgChan = make(chan amqp.Delivery, publisher.size)

	connection, err := amqp.DialConfig(publisher.Uri, publisher.Config)
	if err != nil {
		return err
	}

	publisher.conn = connection

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	publisher.channel = channel

	channel.NotifyClose(publisher.errChan)

	err = channel.ExchangeDeclare(publisher.Exchange, publisher.ExchangeType, true, false,
		false, false, nil)
	if err != nil {
		return nil
	}

	queue, err := channel.QueueDeclare(publisher.ReplyTo, true, true, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.QueueBind(queue.Name, publisher.RouteKey, publisher.Exchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (publisher *publisher) Publish(key string, body []byte) error{
	inc := publisher.id.Inc()
	corId := fmt.Sprintf("%d", inc)
	pubErr := publisher.channel.Publish(publisher.Exchange, key, false, false,
		amqp.Publishing{ReplyTo: publisher.ReplyTo, Body: body, CorrelationId: corId})
	if pubErr != nil {
		return pubErr
	}

	cbi := callBackInfo{id: corId, df: defaultPubFuture{c: make(chan interface{})}}
	publisher.callBacks.Set(corId, cbi)
	return nil
}

func (publisher *publisher) Process() {
	for {
		if err := publisher.init() ; err != nil {
			log4go.Warn("init err:%+v", err)
			time.Sleep(SLEEP_TIME)
		}else if deliveries, err := publisher.channel.Consume(publisher.ReplyTo, "",
			true, false, false, false, nil); err == nil{
			go publisher.dealMsg()

			for {
				needBreak := false
				select {
				case delivery, ok :=<- deliveries:
					if ok {
						publisher.msgChan <- delivery
					}else {
						needBreak = true
					}
				case <-publisher.errChan:
					needBreak = true
				}

				if needBreak {
					publisher.close()
					break
				}
			}
		}else {
			time.Sleep(SLEEP_TIME)
		}
	}
}
func (publisher *publisher) dealMsg() {
	for delivery := range publisher.msgChan{
		if publisher.callBacks.Has(delivery.CorrelationId) {
			if val, ok := publisher.callBacks.Get(delivery.CorrelationId); ok{
				if cbi, ok := val.(callBackInfo); ok {
					go func() {
						cbi.df.c <- delivery.Body
					}()
				}
			}
		}
	}
}

func main(){
	url := flag.String("Uri", "amqp://guest:guest@192.168.96.204:5672", "mq addr")
	exchangeKey := flag.String("Exchange", "test-ex", "Exchange name")
	exchangeType := flag.String("ExchangeType", "topic", "Exchange type")
	queueName := flag.String("queueName", "test-queue", "queue name")
	routeKey := flag.String("RouteKey", "test-route", "route key")
	replyTo := flag.String("ReplyTO", "callback-queue", "callback queue name")
	//pub := &publisher{Uri: *url, Config: amqp.Config{Heartbeat: 3 * time.Second,
	//Dial: func(network, addr string) (net.Conn, error) {
	//	fmt.Println("dial", network, addr)
	//	return net.DialTimeout(network, addr, 3 * time.Second)
	//}},
	//Exchange: *exchangeKey,
	//ExchangeType: *exchangeType,
	//RouteKey: *routeKey,
	//ReplyTo: *replyTo, }
	//
	//initErr := pub.init()
	//if initErr != nil {
	//	fmt.Println(initErr)
	//}

	pub := NewPublisher(*url, *exchangeKey, *exchangeType, *routeKey, *replyTo, amqp.Config{Heartbeat: 3 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			fmt.Println("dial", network, addr)
			return net.DialTimeout(network, addr, 3 * time.Second)
		}}, 20)
	pub.Publish(*routeKey, []byte("Hello"))

	//go func() {
	//	for {
	//		time.Sleep(time.Second * 3)
	//		msg := Message{MessageId: 1, Msg:"Hello"}
	//		msgJson, _ := json.Marshal(msg)
	//		err := pub.Publish(*routeKey, msgJson)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	//}()

	customConsumer(*url, *queueName, *exchangeKey, *exchangeType, *routeKey)
}

func customConsumer(url, queueName, exchangeKey, exchangeType, routeKey string){
	fmt.Println("customConsumer", url, queueName, exchangeKey, exchangeType, routeKey)
	consumer := rabbitmq.NewRabbitmqConsumer(url, exchangeKey, exchangeType, queueName, false,
		1024, routeKey)
	go consumer.Process()
	for v := range consumer.Consume() {
		fmt.Println("receive:" + string(v))
	}
}
