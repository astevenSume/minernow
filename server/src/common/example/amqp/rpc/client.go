package main

import (
	"common"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
	msgs <-chan amqp.Delivery
)

func clientInit(regionId, serverId int) (err error) {
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
		fmt.Sprintf("rpc.%d.%d.queue", regionId, serverId),
		//"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	msgs, err = ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func fibonacciRPC(n int) (res int, err error) {
	//conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}
	//
	//defer conn.Close()
	//
	//ch, err := conn.Channel()
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}
	//
	//defer ch.Close()
	//
	//q, err := ch.QueueDeclare(
	//	"",
	//	false,
	//	false,
	//	true,
	//	false,
	//	nil,
	//)
	//if err != nil {
	//	common.LogFuncError("%v", err)
	//	return
	//}
	//
	//msgs, err = ch.Consume(
	//	q.Name,
	//	"",
	//	false,
	//	false,
	//	false,
	//	false,
	//	nil)

	corrId := randomString(32)

	err = ch.Publish(
		"",
		"rpc.queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	common.LogFuncDebug("[X] requesting fib(%d) to rpc.queue", n)

	for d := range msgs {
		if d.CorrelationId == corrId {
			res, err = strconv.Atoi(string(d.Body))
			if err != nil {
				common.LogFuncError("%v", err)
				return
			}
			break
		} else {
			common.LogFuncDebug("received CorrelationId %v", d.CorrelationId)
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 3 {
		common.LogFuncError("Usage : %s <region id> <server id>", os.Args[0])
		return
	}

	var (
		regionId, serverId int
		err                error
	)

	regionId, err = strconv.Atoi(os.Args[1])
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	serverId, err = strconv.Atoi(os.Args[2])
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	err = clientInit(regionId, serverId)
	if err != nil {
		return
	}

	for {
		var n int
		fmt.Println("please input a number :")
		fmt.Scanf("%d", &n)
		res, err := fibonacciRPC(n)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
		common.LogFuncDebug("[.] Got %d", res)
	}
}

func bodyFrom(args []string) int {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		common.LogFuncError("%v", err)
		return 0
	}

	return n
}
