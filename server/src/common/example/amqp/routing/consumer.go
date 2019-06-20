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

	q, err := ch.QueueDeclare(
		"",
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

	//if len(os.Args) < 2 {
	//	common.LogFuncError("Usage : %s <routing_key_0> <routing_key_1> <routing_key_2> ...")
	//	os.Exit(0)
	//}

	//绑定所有要订阅的routing key
	//for _, s := range os.Args[1:] {
	err = ch.QueueBind(
		q.Name,
		//s,
		"#",
		"logs_direct",
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	//}

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
		common.LogFuncDebug("start receiving messages")
		for d := range msgs {
			common.LogFuncDebug("Received : %v", d)
			//d.Ack(false)
		}
	}()

	<-forever
}
