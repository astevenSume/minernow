package common

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

const (
	rabbitMQDirectReplyTo = "amq.rabbitmq.reply-to"

	RpcTimeOutDefault = time.Second * 3
)

// RabbitMQRpcMsg rabbitmq rpc message
//easyjson:json
type RabbitMQRpcMsg struct {
	Cmd  int    `json:"cmd"`  //Cmd defination , check utils/common/amqp.go
	Code int    `json:"code"` //Error Code
	Body []byte `json:"body"`
}

type rabbitMQRpc struct {
	prefix, appName    string
	regionId, serverId int64
}

func NewRabbitMQRpc(prefix, appName string, regionId, serverId int64) *rabbitMQRpc {
	return &rabbitMQRpc{
		prefix:   prefix,
		appName:  appName,
		regionId: regionId,
		serverId: serverId,
	}
}

func (r rabbitMQRpc) key() string {
	return fmt.Sprintf("%v.%v.%v.%v", r.prefix, r.appName, r.regionId, r.serverId)
}

type RpcMessageFunc func(delivery amqp.Delivery) ([]byte, error)

//
type rabbitMQRpcServer struct {
	rabbitMQBase
	q     amqp.Queue
	local rabbitMQRpc
	f     RpcMessageFunc
}

func NewRabbitMQRpcServer() *rabbitMQRpcServer {
	return &rabbitMQRpcServer{}
}

func (s *rabbitMQRpcServer) Init(address, appName string, regionId, serverId int64, f RpcMessageFunc) (err error) {
	s.address = address
	s.local.prefix = "rpc.server"
	s.local.appName = appName
	s.local.regionId = regionId
	s.local.serverId = serverId
	s.f = f

	s.conn, err = amqp.Dial(address)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	s.chErr = make(chan *amqp.Error)
	s.conn.NotifyClose(s.chErr)
	s.fReset = func() {
		s.reset()
	}

	//go routine process connection exceptions.
	go SafeRun(
		func() {
			processErr(&s.rabbitMQBase)
		})()

	err = s.reset()
	if err != nil {
		return
	}

	return
}

func (s *rabbitMQRpcServer) reset() (err error) {
	s.ch, err = s.conn.Channel()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// declare exchange
	err = s.ch.ExchangeDeclare(
		s.local.key(),
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// declare queue, should be exclusive
	s.q, err = s.ch.QueueDeclare(
		s.local.key()+".queue",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	err = s.ch.QueueBind(
		s.q.Name,
		"#",
		s.local.key(),
		false,
		nil,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	var msgs <-chan amqp.Delivery
	msgs, err = s.ch.Consume(
		s.q.Name,
		"",
		true, //auto ack,no retry
		false,
		false,
		false,
		nil)

	// main loop, process messages
	go SafeRun(func() {
		for d := range msgs {
			//LogFuncDebug("Received message from %s %s : %s", d.ReplyTo, d.CorrelationId, string(d.Body))
			go SafeRun(func() {
				resp, err := s.f(d)
				if err != nil {
					LogFuncError("%v", err)
					return
				}

				err = s.ch.Publish(
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          resp,
					})

				if err != nil {
					LogFuncError("%v", err)
					return
				}
				//LogFuncDebug("RpcServer Publish to %s : %s", d.ReplyTo, string(d.Body))
				//d.Ack(false)
			})()
		}
	})()

	return
}

func (s *rabbitMQRpcServer) Close() (err error) {
	err = s.ch.Close()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	err = s.conn.Close()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

type corrId2ChanManager struct {
	items map[string]chan []byte
	lock  sync.Mutex
}

func NewCorrId2ChanManager() *corrId2ChanManager {
	return &corrId2ChanManager{
		items: make(map[string]chan []byte),
	}
}

// Get get channel of corrId
func (m *corrId2ChanManager) Get(corrId string) chan []byte {
	m.lock.Lock()
	defer m.lock.Unlock()

	if ch, ok := m.items[corrId]; ok {
		return ch
	}

	return nil
}

// Set set corrId's channel
func (m *corrId2ChanManager) Set(corrId string, ch chan []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items[corrId] = ch
}

// Del delete channel of corrId
func (m *corrId2ChanManager) Del(corrId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.items, corrId)
}

// rabbitmq rpc client defination
type rabbitMQRpcClient struct {
	rabbitMQBase
	local          rabbitMQRpc
	server         string
	q              amqp.Queue
	msgs           <-chan amqp.Delivery
	corrId2ChanMgr *corrId2ChanManager
}

func NewRabbitMQRpcClient() *rabbitMQRpcClient {
	return &rabbitMQRpcClient{
		corrId2ChanMgr: NewCorrId2ChanManager(),
	}
}

func (c *rabbitMQRpcClient) Init(address, appName string, regionId, serverId int64) (err error) {
	c.address = address
	c.local.prefix = "rpc.client"
	c.local.appName = appName
	c.local.regionId = regionId
	c.local.serverId = serverId
	c.server = fmt.Sprintf("rpc.server.%s.%d.%d", appName, regionId, serverId)
	c.conn, err = amqp.Dial(address)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	c.chErr = make(chan *amqp.Error)
	c.conn.NotifyClose(c.chErr)
	c.fReset = func() {
		c.reset()
	}

	//go routine process connection exceptions.
	go SafeRun(
		func() {
			processErr(&c.rabbitMQBase)
		})()

	err = c.reset()
	if err != nil {
		return
	}

	return
}

func (c *rabbitMQRpcClient) reset() (err error) {
	c.ch, err = c.conn.Channel()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	//defer ch.Close()

	err = c.ch.ExchangeDeclare(
		c.server,
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	c.msgs, err = c.ch.Consume(
		rabbitMQDirectReplyTo,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// process all messages from "amq.rabbitmq.reply-to"
	go func() {
		for d := range c.msgs {
			respChan := c.corrId2ChanMgr.Get(d.CorrelationId)
			if respChan != nil {
				LogFuncDebug("get right delivery %s : %s", d.CorrelationId, string(d.Body))
				respChan <- d.Body
				continue
			} else {
				LogFuncDebug("get wrong delivery %s : %s", d.CorrelationId, d.Body)
				continue
			}
		}
	}()

	return
}

func (c *rabbitMQRpcClient) Publish(body []byte, timeOut time.Duration) (resp []byte, err error) {
	//channel to store rpc response message
	respCh := make(chan []byte)

	//generate CorrelationId
	corrId := RandomStr(32)

	// publish to rpc server
	err = c.ch.Publish(
		c.server,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       rabbitMQDirectReplyTo,
			Body:          body,
		})
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	// save respChan
	c.corrId2ChanMgr.Set(corrId, respCh)

	// time ticker
	ticker := time.After(time.Duration(timeOut) * time.Millisecond)

	// loop to get response
	for {
		select {
		case resp = <-respCh:
			{
				LogFuncDebug("Received %s : %s", corrId, string(body))
				return
			}
		case <-ticker: //timeout
			{
				c.corrId2ChanMgr.Del(corrId)
				LogFuncWarning("rpc request to %s timeout for %s", c.server, corrId)
				return
			}
		}
	}

	return
}

func (c *rabbitMQRpcClient) Close() (err error) {
	err = c.ch.Close()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	err = c.conn.Close()
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

func rpcServer(appName string, regionId, serverId int64) string {
	return fmt.Sprintf("rpc.server.%v.%v.%v", appName, regionId, serverId)
}
