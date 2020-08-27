package mqtt

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test")
	c.Connect(uri, connMsg)

	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("abc1"), message.QosAtMostOnce)
	c.Subscribe(subMsg, nil, func(msg *message.PublishMessage) error {
		fmt.Printf("%s: %s", string(msg.Topic()), string(msg.Payload()))
		return nil
	})

	fmt.Println("waiting message coming")
	<-time.After(1 * time.Hour)
}

func TestMQTTSend(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test")
	c.Connect("tcp://192.168.9.93:44765", connMsg)

	pubMsg := message.NewPublishMessage()
	pubMsg.SetTopic([]byte("abc1"))
	pubMsg.SetPayload([]byte("hello world"))
	for i := 0; i < 30; i++ {
		c.Publish(pubMsg, nil)
		<-time.After(time.Second)
	}
}
