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

	q, err := ch.QueueDeclare(
		"hello",
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

	msgs, err := ch.Consume(q.Name,
		"",
		false,
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
			common.LogFuncDebug("received : %v, %s", d, string(d.Body))
			d.Ack(false)
		}
	}()

	<-forever
}
