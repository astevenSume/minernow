package main

import (
	"common"
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	addr := "amqp://root:root@localhost:5672/"
	exchange := "logs_topic"
	exchangeAe := "logs_topic_ae"
	queueAe := "logs_topic_unrouted"

	//创建通道
	// Connection和Channel在进程退出的时候会自动释放，不需要显式释放
	ch, err := CreateChannel(addr)
	if err != nil {
		return
	}

	err = DeclareExchange(exchangeAe, amqp.ExchangeTopic, ch, nil)
	if err != nil {
		return
	}

	qAe, err := DeclareQueue(queueAe, ch)
	if err != nil {
		return
	}

	err = ch.QueueBind(
		qAe.Name,
		"#", //注意备份交换器，需要指定routing key为“#”，否则无法获取到所有unrouted消息
		exchangeAe,
		false,
		nil)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//备份交换器处理逻辑
	{
		msgs, err := ch.Consume(
			qAe.Name, //queue
			"",       //consumer
			false,    //注意：设置消息处理为非自动确认
			false,
			false,
			false,
			nil)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		go func() {
			for d := range msgs {
				common.LogFuncDebug("[XAE] Received : %v, %s", d, string(d.Body))
				d.Ack(false)
			}
		}()
	}

	//声明主交换器
	args := amqp.Table{}
	args["alternate-exchange"] = exchangeAe //绑定备用交换器
	err = DeclareExchange(exchange, amqp.ExchangeTopic, ch, args)
	if err != nil {
		return
	}

	for {
		var routingKey, msg string
		fmt.Println("please input message :")
		_, err := fmt.Scanf("%s %s", &routingKey, &msg)
		err = ch.Publish(
			exchange,   //exchange
			routingKey, // routing key
			false,      // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, //注意：设置消息为可持久化
				ContentType:  "text/plain",
				Body:         []byte(msg),
			},
		)
		if err != nil {
			common.LogFuncError("%v", err)
			continue
		}

		common.LogFuncDebug("[x] sent %s", msg)
	}
}
