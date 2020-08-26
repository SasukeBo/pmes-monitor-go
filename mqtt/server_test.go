package mqtt

import (
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test")
	c.Connect(uri, connMsg)
	pubMsg := message.NewPublishMessage()
	pubMsg.SetTopic([]byte("abc1"))
	pubMsg.SetPayload([]byte("hello world"))
	for i := 0; i < 30; i++ {
		c.Publish(pubMsg, nil)
		<-time.After(time.Second)
	}
}
