package mqtt

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"strings"
	"testing"
	"time"
)

//var host = "tcp://192.168.9.135:1883"
var host = "tcp://localhost:44765"

//var host = "tcp://192.168.9.93:44765"
//var host = "tcp://192.168.5.146:1883"

func TestConnect(t *testing.T) {
	c := &service.Client{}
	connMsg := newConnectMessage("test_subscriber")
	c.Connect(host, connMsg)

	subMsg := message.NewSubscribeMessage()
	subMsg.AddTopic([]byte("PLC"), message.QosAtMostOnce)
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
	pubMsg.SetTopic([]byte("PLC"))
	//             mac 			sta  total 	   ng        error codes
	var message = "58528af7ff84 0010 0040 0000 0001 0000 0000 1000 0008 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000"
	pubMsg.SetPayload([]byte(strings.ReplaceAll(message, " ", "")))
	c.Publish(pubMsg, nil)
}

func TestHandleMessage(t *testing.T) {
	// 92 46
	// 58528af7ff84 0011 0375 0000 0000 0000 0010 00f0 0000 0000 0000 0000000000000000000000000000000000000000
	//var message = "58528af7ff84 0020 1faa 0000 00b1 0000 0000 0000 0000 0004 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000"
	//                               1         1         1                        1    1
	var message = "58528af7ff84 0010 1e73 0000 0000 0000 0000 0000 0000 0000 0000 0000 0400 00000000000000000000000000000000"
	//handleMessage(strings.ReplaceAll(message, " ", ""))
	payload := strings.ReplaceAll(message, " ", "")
	result, err := analyzeMessage(payload[12:])
	fmt.Printf("result: %+v, err: %v\n", result, err)
}

func TestWordsToErrorIdxs(t *testing.T) {
	words, _ := hexToWords("00f0f000")
	idxs := wordsToErrorIdxs(words)
	fmt.Println(idxs)
}
