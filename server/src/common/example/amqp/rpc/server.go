package main

import (
	"common"
	"github.com/streadway/amqp"
	"strconv"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
	msgs <-chan amqp.Delivery
)

func serverInit() (err error) {
	conn, err = amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//defer conn.Close()

	ch, err = conn.Channel()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//defer ch.Close()

	q, err = ch.QueueDeclare(
		"rpc.queue",
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

	err = ch.Qos(1, 0, false)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)

	go common.SafeRun(func() {
		common.LogFuncDebug("begin routine...")
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			if err != nil {
				common.LogFuncError("%v", err)
				continue
			}
			common.LogFuncDebug("n : %d, will reply to %s", n, d.ReplyTo)
			response := fib(n)
			err = ch.Publish(
				"",
				d.ReplyTo,
				false, false, amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			common.LogFuncDebug("Publish to %s", d.ReplyTo)
			d.Ack(false)
		}
	})()

	return
}

func main() {
	err := serverInit()
	if err != nil {
		return
	}

	forever := make(chan bool)
	<-forever
}

func fib(n int) int {

	//if n > 100 {
	//	return 0
	//}
	//
	//if n == 0 {
	//	return 0
	//} else if n == 1 {
	//	return 1
	//} else {
	//	return fib(n-1) + fib(n-2)
	//}

	return n
}
