package main

import (
	"time"
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"github.com/orcaman/concurrent-map"
	"fmt"
	"github.com/alecthomas/log4go"
	"flag"
	"net"
	"context"
	"encoding/json"
)

const (
	SLEEP_TIME = time.Second * 1
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
	Publish(key string, body []byte) (PubFuture, error)
}

type publisher struct {
	Uri          string
	Config       amqp.Config
	Exchange     string
	ExchangeType string
	ReplyTo 	 string
	size 		 int

	channel *amqp.Channel
	conn    *amqp.Connection
	errChan chan *amqp.Error
	msgChan chan amqp.Delivery
	id atomic.Int64
	callBacks cmap.ConcurrentMap
}

func NewPublisher(uri, exchange, exchangeType string, replyTo string, config amqp.Config, size int) PubRpc{
	pub := &publisher{Uri: uri, Config: config, Exchange: exchange, ExchangeType: exchangeType,
		ReplyTo: replyTo, size: size}

	go pub.process()
	return pub
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

	_, err = channel.QueueDeclare(publisher.ReplyTo, false, true, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (publisher *publisher) Publish(key string, body []byte) (PubFuture, error){
	inc := publisher.id.Inc()
	corId := fmt.Sprintf("%d", inc)
	cbi := callBackInfo{id: corId, df: defaultPubFuture{c: make(chan interface{})}}
	publisher.callBacks.Set(corId, cbi)

	log4go.Info("prepared to publish.corId:%s", corId)
	pubErr := publisher.channel.Publish(publisher.Exchange, key, false, false,
		amqp.Publishing{ReplyTo: publisher.ReplyTo, Body: body, CorrelationId: corId,
			ContentType:"application/json"})
	if pubErr != nil {
		return nil, pubErr
	}

	return &cbi.df, nil
}

func (publisher *publisher) process() {
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
					log4go.Info("process needBreak")
					publisher.close()
					break
				}
			}
		}else {
			log4go.Info("process sleep")
			time.Sleep(SLEEP_TIME)
		}
	}
}

func (publisher *publisher) dealMsg() {
	for delivery := range publisher.msgChan{
		if publisher.callBacks.Has(delivery.CorrelationId) {
			if val, ok := publisher.callBacks.Get(delivery.CorrelationId); ok{
				if cbi, ok := val.(callBackInfo); ok {
					publisher.callBacks.Remove(delivery.CorrelationId)
					cbi.df.c <- delivery.Body
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

	pub := NewPublisher(*url, *exchangeKey, *exchangeType, *replyTo, amqp.Config{Heartbeat: 3 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			fmt.Println("dial", network, addr)
			return net.DialTimeout(network, addr, 3 * time.Second)
		}}, 20)
	// pub.Publish(*routeKey, []byte("Hello"))

	go customConsumer(*url, *queueName, *exchangeKey, *exchangeType, *routeKey)

	for {
		time.Sleep(time.Second * 3)
		msg := Message{MessageId: 1, Msg:"Hello"}
		msgJson, _ := json.Marshal(msg)
		pub, err := pub.Publish(*routeKey, msgJson)
		if err != nil {
			fmt.Println(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		get, err := pub.GetWithContext(ctx)
		if err != nil {
			log4go.Warn("err:%v", err)
		}
		cancel()

		log4go.Info("get is :%+v", get)
	}
}

func customConsumer(url, queueName, exchangeKey, exchangeType, routeKey string){
	fmt.Println("customConsumer", url, queueName, exchangeKey, exchangeType, routeKey)

	connection, err := amqp.DialConfig(url, amqp.Config{Dial: func(network, addr string) (net.Conn, error) {
		return net.Dial(network, addr)
	}, Heartbeat: SLEEP_TIME})

	if err != nil {
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare(exchangeKey, amqp.ExchangeTopic, true, false,
		false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(queueName, routeKey, exchangeKey, false, nil)
	if err != nil {
		panic(err)
	}

	deliveries, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	count := 0x1
	for d := range deliveries {
		log4go.Info("receive msg: %s", string(d.Body))
		count++
		err := channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{
			CorrelationId: d.CorrelationId,
			Body:          []byte{byte(count)},
		})
		if err != nil {
			log4go.Warn("publish err:%v", err)
		}

		channel.Ack(d.DeliveryTag, false)
	}

	log4go.Info("end................")
}
