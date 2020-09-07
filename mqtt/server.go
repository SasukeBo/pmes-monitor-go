package mqtt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/SasukeBo/pmes-device-monitor/orm"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"math"
	"time"
)

const (
	healthCheckTopic     = "HEALTH_CHECK"
	globalSubscriber     = "GLOBAL_SUBSCRIBER"
	healthCheckPublisher = "HEALTH_CHECK_PUBLISHER"
)

var healthCheckError = errors.New("health check error")

var uri string

func Serve() {
	svr := &service.Server{
		ConnectTimeout:   10, // seconds
		AckTimeout:       0,
		TimeoutRetries:   0,
		Authenticator:    "mockSuccess", // always succeed
		SessionsProvider: "mem",         // keeps sessions in memory
		TopicsProvider:   "mem",         // keeps topic subscriptions in memory
	}

	go globalSubscribe()
	// Listen and serve connections at localhost:1883
	log.Info("MQTT server listening on %s", configer.GetString("mqtt_port"))
	svr.ListenAndServe(uri)
}

func globalSubscribe() {
	defer func() {
		if err := recover(); err != nil {
			go globalSubscribe()
		}
	}()

	c := &service.Client{}

	for {
		err := connect(c)
		if err == nil {
			log.Success("global subscribe connect successful")
			break
		}
		log.Error("global subscribe connect failed: %v", err)
		log.Info("try connect again after 1 second.")
		<-time.After(1 * time.Second)
	}

	err := setHeartBeat()
	if err != nil {
		panic(err)
	}

	for { // Loop to connect and check health
		subscribeAndCheckHealth(c)
		connect(c)
	}
}

func connect(c *service.Client) error {
	connMsg := newConnectMessage(globalSubscriber)
	return c.Connect(uri, connMsg)
}

func newConnectMessage(clientName string) *message.ConnectMessage {
	msg := message.NewConnectMessage()
	msg.SetWillQos(message.QosExactlyOnce)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(clientName))
	msg.SetKeepAlive(60 * 60)
	msg.SetWillTopic([]byte("will"))
	return msg
}

func subscribeAndCheckHealth(c *service.Client) error {
	var healthCheckChan = make(chan struct{})
	var err error
	subMsg := message.NewSubscribeMessage()
	err = subMsg.AddTopic([]byte("#"), message.QosAtMostOnce) // subscribe to all
	if err != nil {
		return err
	}

	err = c.Subscribe(subMsg, nil, func(msg *message.PublishMessage) error {
		if topic := string(msg.Topic()); topic == healthCheckTopic {
			healthCheckChan <- struct{}{}
		} else {
			payload := string(msg.Payload())
			go handleMessage(payload)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for {
		select {
		case <-healthCheckChan:
			continue
		case <-time.After(2 * time.Second):
			return healthCheckError
		}
	}
}

func setHeartBeat() error {
	c := &service.Client{}
	connMsg := newConnectMessage(healthCheckPublisher)
	err := c.Connect(uri, connMsg)
	if err != nil {
		return err
	}
	pubMsg := message.NewPublishMessage()
	err = pubMsg.SetTopic([]byte(healthCheckTopic))
	if err != nil {
		return err
	}
	pubMsg.SetPayload([]byte("ping"))
	go func() {
		for {
			<-time.After(1 * time.Second)
			err := c.Publish(pubMsg, nil)
			if err != nil {
				c.Connect(uri, connMsg)
			}
		}
	}()

	return nil
}

func init() {
	uri = fmt.Sprintf("tcp://0.0.0.0:%s", configer.GetString("mqtt_port"))
}

const (
	deviceStatusRunning          = 16
	deviceStatusStopped          = 17
	deviceStatusRunningWithError = 18
	deviceStatusOffline          = 32
	deviceStatusStoppedWithError = 33
)

func handleMessage(payload string) {
	words, err := hexToWords(payload)
	if err != nil {
		log.Errorln(err)
		return
	}
	fmt.Printf("words: %v\n", words)

	if len(words) < 8 {
		log.Errorln("Illegal payload length.")
		return
	}

	// mac
	mac := wordsToMAC(words[0:3])
	fmt.Printf("mac: %v\n", mac)
	// 状态
	status := wordToStatus(words[3])
	fmt.Println(status, words[3])
	// 产量
	total := wordsToAmount(words[4:6])
	fmt.Printf("产量：%v\n", total)
	// 不良
	ng := wordsToAmount(words[6:8])
	fmt.Printf("不良：%v\n", ng)

	var errorIndex []int
	if len(words) > 8 {
		errorIndex = wordsToErrorIdxs(words[8:])
	}
	fmt.Printf("故障信息编号：%v\n", errorIndex)
}

func wordsToErrorIdxs(words [][]byte) []int {
	var idxs []int
	fmt.Println(words)
	for i, word := range words {
		j := bytesToInt(word)
		if j == 0 {
			continue
		}
		for k := 0; k < 16; k++ {
			compare := int(math.Pow(2, float64(k)))
			if j&compare == compare {
				idxs = append(idxs, i*16+k)
			}
		}
	}

	return idxs
}

func wordsToAmount(words [][]byte) int {
	if len(words) != 2 {
		return 0
	}
	var amountBytes []byte
	amountBytes = append(amountBytes, words[1]...)
	amountBytes = append(amountBytes, words[0]...)
	return bytesToInt(amountBytes)
}

func wordToStatus(word []byte) int {
	statusCode := bytesToInt(word)
	var status int
	switch statusCode {
	case deviceStatusRunning:
		status = orm.DeviceStatusRunning
		fmt.Println("status: Running")
	case deviceStatusStopped:
		status = orm.DeviceStatusStopped
		fmt.Println("status: Stopped")
	case deviceStatusStoppedWithError, deviceStatusRunningWithError:
		status = orm.DeviceStatusError
		fmt.Println("status: Error")
	case deviceStatusOffline:
		status = orm.DeviceStatusShutdown
		fmt.Println("status: Offline")
	}

	return status
}

func wordsToMAC(words [][]byte) string {
	if len(words) < 3 {
		return ""
	}
	return fmt.Sprintf("%s%s%s", hex.EncodeToString(words[0]), hex.EncodeToString(words[1]), hex.EncodeToString(words[2]))
}

func hexToWords(hexStr string) ([][]byte, error) {
	var length = len(hexStr)
	var words [][]byte
	for i := 0; i < length; i = i + 4 {
		word, err := hex.DecodeString(hexStr[i : i+4])
		if err != nil {
			return words, err
		}

		words = append(words, word)
	}

	return words, nil
}

func bytesToInt(bytes []byte) int {
	var result int
	var length = len(bytes)
	for i := 0; i < length; i++ {
		result = result + int(bytes[i])*int(math.Pow(16*16, float64(length-i-1)))
	}

	return result
}
