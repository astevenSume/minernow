package main

import (
	"common"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage : %s <exchange>\n", os.Args[0])
		return
	}

	address := "amqp://root:root@localhost:5672/"
	exchange := os.Args[1]
	producer := common.NewRabbitMQProducer()
	err := producer.Init(address, exchange, AEMessageFunc)
	if err != nil {
		return
	}

	for {
		var key, body string
		fmt.Println("please input message : <key> <body>")
		fmt.Scanf("%s %s", &key, &body)
		err = producer.Publish(key, []byte(body))
		if err != nil {
			return
		}
	}
}

func AEMessageFunc(delivery amqp.Delivery) (err error) {
	common.LogFuncDebug("message %s %s unrouted.", delivery.RoutingKey, string(delivery.Body))
	return
}
