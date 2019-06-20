package main

import (
	"common"
	"github.com/streadway/amqp"
	"os"
)

func main() {
	addr := "amqp://root:root@localhost:5672/"
	exchange := "logs_topic"
	exchangeAe := "logs_topic_ae"
	queue := "logs_topic_routed"
	//queueAe := "logs_topic_unrouted"

	//创建通道
	// Connection和Channel在进程退出的时候会自动释放，不需要显式释放
	ch, err := CreateChannel(addr)
	if err != nil {
		return
	}

	//声明主交换器
	args := amqp.Table{}
	args["alternate-exchange"] = exchangeAe //绑定备用交换器

	// 尝试获取一个已经存在的交换器
	err = DeclareExchange(exchange, amqp.ExchangeTopic, ch, args)
	if err != nil {
		return
	}

	//声明备用交换器的队列
	q, err := DeclareQueue(queue, ch)
	if err != nil {
		return
	}

	if len(os.Args) < 2 {
		common.LogFuncError("Usage : %s <routing_key_0> <routing_key_1> <routing_key_2> ...", os.Args[0])
		os.Exit(0)
	}

	//绑定所有要订阅的routing key
	for _, s := range os.Args[1:] {
		err = ch.QueueBind(
			q.Name,
			s,
			exchange,
			false,
			nil,
		)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		common.LogFuncDebug("binding key %s", s)
	}

	forever := make(chan bool)

	//主交换器
	{
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

		go func() {
			for d := range msgs {
				common.LogFuncDebug("[X] Received : %v, %s", d, string(d.Body))
				d.Ack(false)
			}
		}()
	}

	<-forever
}
