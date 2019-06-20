package main

import (
	"common"
	"fmt"
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

	err = ch.ExchangeDeclare(
		"logs_direct",
		amqp.ExchangeTopic,
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

	for {
		var routingKey, msg string
		fmt.Println("please input message :")
		_, err := fmt.Scanf("%s %s", &routingKey, &msg)
		err = ch.Publish(
			"logs_direct", //exchange
			routingKey,    // routing key
			false,         // mandatory
			false,
			amqp.Publishing{
				//DeliveryMode: amqp.Persistent, //注意：设置消息为可持久化
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			common.LogFuncError("%v", err)
			continue
		}
		common.LogFuncDebug("[x] sent %s", msg)
	}
}
