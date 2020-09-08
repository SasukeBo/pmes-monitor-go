package mqtt

import (
	"errors"
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
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
