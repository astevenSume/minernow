package main

import (
	"common"
	"github.com/streadway/amqp"
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
		"logs",
		"fanout",
		true,
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
		"",
		false, //注意：日志消息，丢了就丢了，不做持久化
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer
		true,   //注意：设置消息处理为非自动确认
		false,
		false,
		false,
		nil)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			common.LogFuncDebug("Received : %v, %s", d, string(d.Body))
		}
	}()

	<-forever
}
