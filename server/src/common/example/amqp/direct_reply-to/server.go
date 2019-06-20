package main

import (
	"common"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	defer ch.Close()

	//声明为fanout类型的Exchange，用于将消息广播到所有订阅的队列
	err = ch.ExchangeDeclare(
		"rpc",
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	q, err := ch.QueueDeclare(
		"rpc",
		false, //注意：日志消息，丢了就丢了，不做持久化
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//绑定所有要订阅的routing key
	err = ch.QueueBind(
		q.Name,
		"#",
		"rpc",
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	forever := make(chan bool)

	helloMsg := "hello world"

	for d := range msgs {
		n, err := strconv.Atoi(string(d.Body))
		if err != nil {
			common.LogFuncError("%v", err)
			continue
		}

		go common.SafeRun(func() {
			tmp := d
			// 4debug sleep for different seconds.
			time.Sleep(time.Duration(n) * time.Second)

			common.LogFuncDebug("Received : %s, %s, response to %s", tmp.CorrelationId, string(tmp.Body), tmp.ReplyTo)

			err = ch.Publish(
				"",          //exchange
				tmp.ReplyTo, // routing key
				false,       // mandatory
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					Body:          []byte(helloMsg + " after " + string(tmp.Body) + " seconds"),
					CorrelationId: tmp.CorrelationId,
				},
			)
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			common.LogFuncDebug("[x] sent %s", helloMsg)
		})()
	}

	<-forever
}
