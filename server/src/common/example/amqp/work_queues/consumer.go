package main

import (
	"bytes"
	"common"
	"github.com/streadway/amqp"
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

	q, err := ch.QueueDeclare(
		"task_queue",
		true, //注意：队列设置为可持久化的
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//保证rabbitmq每次只往消费者通道发送一个消息，直到该消息处理完了，否则不会往该通道再发消息
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, //global
	)
	if err != nil {
		common.LogFuncError("failed to set Qos : %v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer
		false,  //注意：设置消息处理为非自动确认
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
			dot_count := bytes.Count(d.Body, []byte("."))
			time.Sleep(time.Duration(dot_count) * time.Second)
			common.LogFuncDebug("Done")
			d.Ack(false) //注意，如果没有返回确认消息，可能会导致rabbitmq重复发送消息，并且造成rabbitmq内存暴涨
		}
	}()

	<-forever
}
