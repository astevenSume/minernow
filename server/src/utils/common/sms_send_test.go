package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func TestSMS_Send(t *testing.T) {
	service := sns.New(session.New(aws.NewConfig()), aws.NewConfig().WithRegion("us-east-1").WithCredentials(credentials.NewStaticCredentials("AKIAIRSMVUQN6TAEPOQA", "B2Pgokp2z6Z+ga2sSnIWAwL3oYAkLZKNjTonqeZQ", "")))
	//service := sns.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))
	sms := &SMS{
		Service:  service,
		Type:     Transactional,
		MaxPrice: 0.50,
	}

	err := sms.Send("Hello from SMS.Send()", "+86xxxxx")
	if err != nil {
		fmt.Printf("sms send fail err:%v", err)
	}
}

func TestSend(t *testing.T) {
	err := Send("Hello from Send()", "+86xxxxx")
	if err != nil {
		fmt.Printf("sms send fail err:%v", err)
	}
}
