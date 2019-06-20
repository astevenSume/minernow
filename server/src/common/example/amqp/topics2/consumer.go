package main

import (
	"common"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage : %s <exchange> <binding key>\n", os.Args[0])
		os.Exit(0)
	}

	address := "amqp://root:root@localhost:5672/"
	exchange := os.Args[1]
	bindingKey := os.Args[2]
	consumer := common.NewRabbitMQConsumer()
	err := consumer.Init(address, exchange, bindingKey, ConsumeFunc)
	if err != nil {
		return
	}

	forever := make(chan bool)

	<-forever
}

func ConsumeFunc(delivery amqp.Delivery) (err error) {
	common.LogFuncDebug("message %s %s consumed.", delivery.RoutingKey, string(delivery.Body))
	return
}
