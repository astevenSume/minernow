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
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	// 不再需要声明queue，因为广播要达到的目的是消息发送到所有的消费者，并且在rabbitmq重启以后，
	// 不需要重复处理未处理的消息（消息不用持久化）
	//q, err := ch.QueueDeclare(
	//	"task_queue",
	//	true,  //注意：队列设置为可持久化
	//	false, //delete when unused
	//	false,
	//	false,
	//	nil,
	//)

	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	for {
		var msg string
		fmt.Println("please input message :")
		_, err := fmt.Scanf("%s", &msg)
		err = ch.Publish(
			"logs", //exchange
			"",     // routing key
			false,  // mandatory
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
