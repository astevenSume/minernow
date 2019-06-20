package common

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

const (
	RabbitMQClientTypeUnkown = iota
	RabbitMQClientTypeProducer
	RabbitMQClientTypeConsumer
)

const (
	RabbitMQExchangeTypeUnkown = iota
	RabbitMQExchangeTypeTopic
	RabbitMQExchangeTypeFanout
	RabbitMQExchangeTypeDirect
)

type rabbitMQClient struct {
	rabbitMQBase
	q                        amqp.Queue
	exchange                 string
	autoAck, needAe, durable bool
	isDelay                  bool
	exchangeType             int
	clientType               int
	f                        MessageFunc
}

type MessageFunc func(delivery amqp.Delivery) (err error)

// @Description init
func (r *rabbitMQClient) Init(address, exchange string, autoAck, needAe, durable, isDelay bool, f MessageFunc) (err error) {
	r.address = address
	r.exchange = exchange
	r.autoAck = autoAck
	r.needAe = needAe
	r.durable = durable
	r.isDelay = isDelay
	r.f = f

	r.fReset = func() {
		r.reset()
	}

	err = r.Dail()
	if err != nil {
		return
	}

	err = r.reset()
	if err != nil {
		return
	}

	return
}

// @Description connect to rabbitmq server
func (r *rabbitMQClient) Dail() (err error) {
	//var conn *amqp.Connection
	r.conn, err = amqp.Dial(r.address)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	r.chErr = make(chan *amqp.Error)
	r.conn.NotifyClose(r.chErr)

	//go routine process connection exceptions.
	go SafeRun(
		func() {
			processErr(&r.rabbitMQBase)
		})()

	return
}

func (r *rabbitMQClient) reset() (err error) {
	r.ch, err = r.conn.Channel()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	err = r.declareExchange()
	if err != nil {
		return
	}

	if r.clientType == RabbitMQClientTypeProducer && r.needAe {
		err = r.declareAlternateExchange()
		if err != nil {
			return
		}

		err = r.queueInit("#", r.needAe, nil)
		if err != nil {
			return
		}
	} else if r.clientType == RabbitMQClientTypeConsumer {
		// make sure unacknowledged message in consumer at once.
		err = r.ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, //global
		)
		if err != nil {
			LogFuncError("failed to set Qos : %v", err)
			return
		}

		err = r.queueInit("#", false, nil)
		if err != nil {
			return
		}
	}

	if r.clientType != RabbitMQClientTypeProducer || !r.isDelay {
		return
	}

	// 延迟队列逻辑
	args := amqp.Table{"x-dead-letter-exchange": r.exchange}
	// 替换 exchange
	r.exchange += "_delay"
	err = r.declareExchange()
	if err != nil {
		return
	}
	if r.needAe {
		err = r.declareAlternateExchange()
		if err != nil {
			return
		}
	}
	//延迟队列
	err = r.queueInit("#", false, args)
	if err != nil {
		return
	}
	if r.needAe {
		err = r.queueInit("#", true, nil)
		if err != nil {
			return
		}
	}
	return
}

func reconnect(address string) (conn *amqp.Connection) {
	var err error
	for {
		conn, err = amqp.Dial(address)
		if err == nil {
			return
		}

		LogFuncWarning("rabbitMQClient try reconnect to %s failed : %v", address, err)
		time.Sleep(time.Millisecond * 500)
	}
}

// @Description declare exchange
func (r *rabbitMQClient) declareExchange() (err error) {
	if r.needAe {
		return r.subDeclareExchange(r.exchange, amqp.Table{
			"alternate-exchange": r.exchange + "_ae",
		})
	} else {
		return r.subDeclareExchange(r.exchange, nil)
	}
}

// @Description declare exchange
func (r *rabbitMQClient) declareAlternateExchange() (err error) {
	return r.subDeclareExchange(r.exchange+"_ae", nil)
}

// @Description declare exchange
func (r *rabbitMQClient) subDeclareExchange(exchange string, args amqp.Table) (err error) {

	var kind string

	switch r.exchangeType {
	case RabbitMQExchangeTypeTopic:
		kind = amqp.ExchangeTopic
	case RabbitMQExchangeTypeFanout:
		kind = amqp.ExchangeFanout
	case RabbitMQExchangeTypeDirect:
		kind = amqp.ExchangeDirect
	default:
		err = errors.New("unsupported exchange type " + fmt.Sprint(r.exchangeType))
		return
	}

	err = r.ch.ExchangeDeclare(
		exchange,
		kind,
		r.durable, //always be durable
		false,     //acknowledge always be required
		false,
		false,
		args,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

// @Description declare queue
// queue is the name of queue.
func (r *rabbitMQClient) declareQueue(isAe bool, args amqp.Table) (err error) {
	var queue string
	if isAe {
		queue = r.exchange + ".ae.queue"
	} else { //no need to declare queue
		queue = r.exchange + ".queue"
	}

	err = r.declareQueueBase(queue, args)
	return
}

func (r *rabbitMQClient) declareQueueBase(queue string, args amqp.Table) (err error) {
	r.q, err = r.ch.QueueDeclare(
		queue,
		true, //always be durable
		false,
		false, //always be false
		false,
		args,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

// @Description queue bind
func (r *rabbitMQClient) queueBind(bindingKey string, isAe bool) (err error) {
	var exchange string
	exchange = r.exchange
	if isAe {
		exchange = r.exchange + "_ae"
	}

	err = r.ch.QueueBind(
		r.q.Name,
		bindingKey,
		exchange,
		false,
		nil,
	)

	if err != nil {
		LogFuncError("%v", err)
	}

	return
}

func (r *rabbitMQClient) queueInit(bindingKey string, isAe bool, args amqp.Table) (err error) {
	// init queue
	err = r.declareQueue(isAe, args)
	if err != nil {
		return
	}

	// bind queue
	err = r.queueBind(bindingKey, isAe)
	if err != nil {
		return
	}

	// 延迟队列不创建consume
	if args != nil {
		return
	}
	// 创建消费者
	msgs, err := r.ch.Consume(
		r.q.Name, //queue
		"",       //consumer
		r.autoAck,
		false,
		false,
		false,
		nil)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// Should be recovered from crash, don't crash the whole application.
	go SafeRun(func() {
		for d := range msgs {
			r.f(d)
			if !r.autoAck {
				d.Ack(false)
			}
		}
	})()

	return
}

// Producer
type rabbitMQProducer struct {
	rabbitMQClient
}

func NewRabbitMQProducer(exchangeType int) *rabbitMQProducer {
	return &rabbitMQProducer{
		rabbitMQClient: rabbitMQClient{
			clientType:   RabbitMQClientTypeProducer,
			exchangeType: exchangeType,
		},
	}
}

// @Description init producer
func (r *rabbitMQProducer) Init(address, exchange string, autoAck, needAe, durable, isDelay bool, f MessageFunc) (err error) {
	err = r.rabbitMQClient.Init(address, exchange, autoAck, needAe, durable, isDelay, f)
	if err != nil {
		return
	}
	return
}

// @Description publish message
// key represents routing key
func (r *rabbitMQProducer) Publish(key string, msg []byte) (err error) {
	return r.publish(key, msg, "")
}

//推送延迟信息   t="5000" 5秒
func (r *rabbitMQProducer) PublishDelay(key string, msg []byte, t string) (err error) {
	return r.publish(key+"_delay", msg, t)
}

//推送延迟信息   t="5000" 5秒
func (r *rabbitMQProducer) publish(key string, msg []byte, t string) (err error) {
	err = r.ch.Publish(
		r.exchange,
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //always be durable
			ContentType:  "text/plain",
			Body:         []byte(msg),
			Expiration:   t, // 设置五秒的过期时间
		},
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

// Consumer
type rabbitMQConsumer struct {
	rabbitMQClient
}

func NewRabbitMQConsumer(exchangeType int) *rabbitMQConsumer {
	return &rabbitMQConsumer{
		rabbitMQClient: rabbitMQClient{
			clientType:   RabbitMQClientTypeConsumer,
			exchangeType: exchangeType,
		},
	}
}

// @Description init consumer
func (r *rabbitMQConsumer) Init(address, exchange string, autoAck, needAe, durable bool, f MessageFunc) (err error) {
	err = r.rabbitMQClient.Init(address, exchange, autoAck, needAe, durable, false, f)
	if err != nil {
		return
	}

	return
}
