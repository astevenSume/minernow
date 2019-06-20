package main

import (
	"common"
	"github.com/streadway/amqp"
)

func CreateChannel(address string) (ch *amqp.Channel, err error) {
	var conn *amqp.Connection
	conn, err = amqp.Dial(address)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	ch, err = conn.Channel()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func DeclareExchange(exchange, exchangeType string, ch *amqp.Channel, args amqp.Table) (err error) {
	err = ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func DeclareQueue(queue string, ch *amqp.Channel) (q amqp.Queue, err error) {
	q, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false, //勿声明为独有的(true)，否则对应的exchange只能由一个queue占有
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
