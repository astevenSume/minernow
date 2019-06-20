package main

import (
	"common"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"sync"
	"time"
)

var (
	ch *amqp.Channel
)

func main() {

	rand.Seed(time.Now().UnixNano())

	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	defer conn.Close()

	ch, err = conn.Channel()

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		"rpc",
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)

	msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to", //queue
		"",                      //consumer
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// process all messages from "amq.rabbitmq.reply-to"
	go func() {
		for d := range msgs {
			respChan := corrId2ChanMgr.Get(d.CorrelationId)
			if respChan != nil {
				common.LogFuncDebug("get right delivery %s : %s", d.CorrelationId, string(d.Body))
				respChan <- d.Body
				continue
			} else {
				common.LogFuncDebug("get wrong delivery %s : %s", d.CorrelationId, d.Body)
				continue
			}
		}
	}()

	for {
		var secs int
		fmt.Println("please input message :")
		_, err := fmt.Scanf("%d", &secs)
		if err != nil {
			common.LogFuncError("%v", err)
			continue
		}

		// one publish message per goroutine
		go func() {
			//channel to store rpc response message
			respCh := make(chan []byte)

			//generate CorrelationId
			corrId := common.RandomStr(32)

			// publish to rpc server
			err = ch.Publish(
				"rpc", //exchange
				"",    // routing key
				false, // mandatory
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					Body:          []byte(fmt.Sprint(secs)),
					ReplyTo:       "amq.rabbitmq.reply-to",
					CorrelationId: corrId,
				},
			)
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			common.LogFuncDebug("[x] sent %s : body %s.", corrId, fmt.Sprint(secs))

			// save respChan
			corrId2ChanMgr.Set(corrId, respCh)

			// time ticker
			ticker := time.After(time.Second * 5)

			// loop to get response
			for {
				select {
				case body := <-respCh:
					{
						common.LogFuncDebug("Received %s : %s", corrId, string(body))
						return
					}
				case <-ticker: //timeout
					{
						//close(respCh)
						corrId2ChanMgr.Del(corrId)
						common.LogFuncDebug("timeout for %s", corrId)
						return
					}
				}
			}
		}()
	}

	forever := make(chan bool)
	<-forever
}

var corrId2ChanMgr = NewCorrId2ChanManager()

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
