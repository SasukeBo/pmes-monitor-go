package mqtt

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"testing"
	"time"
)

//var host = "tcp://192.168.9.135:1883"
//var host = "tcp://localhost:1883"
var host = "tcp://192.168.9.93:44765"

func TestConnect(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test_subscriber")
	c.Connect(host, connMsg)

	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("abc1"), message.QosAtMostOnce)
	c.Subscribe(subMsg, nil, func(msg *message.PublishMessage) error {
		fmt.Printf("%s: %s\n", string(msg.Topic()), string(msg.Payload()))
		return nil
	})

	fmt.Println("waiting message coming")
	<-time.After(1 * time.Hour)
}

func TestMQTTSend(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test_publisher")
	if err := c.Connect(host, connMsg); err != nil {
		t.Fatal(err)
	}

	pubMsg := message.NewPublishMessage()
	pubMsg.SetTopic([]byte("abc1"))
	pubMsg.SetPayload([]byte("hello world\n"))
	for i := 0; i < 30; i++ {
		c.Publish(pubMsg, nil)
		<-time.After(time.Second)
	}
}
