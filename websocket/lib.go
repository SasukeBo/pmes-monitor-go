package websocket

import (
	"github.com/SasukeBo/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"
)

var wsUpGrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsConn struct {
	*websocket.Conn
	topics map[string]struct{}
	token  string
}

func (c *WsConn) AddTopic(topic string) {
	log.Info("Add Topic %s for connection(%s)", topic, c.token)
	c.topics[topic] = struct{}{}
}

func (c *WsConn) IsSubscriber(topic string) bool {
	_, ok := c.topics[topic]
	return ok
}

func (c *WsConn) RemoveTopic(topic string) {
	delete(c.topics, topic)
}

var connPool map[string]*WsConn

const (
	SubscribePrefix   = "SUBSCRIBE:"
	UnSubscribePrefix = "UNSUBSCRIBE:"
	CloseWebSocket    = "CLOSE"
)

func (c *WsConn) Close() {
	c.Conn.Close()
	delete(connPool, c.token)
}

func (c *WsConn) Receive() bool {
	_, content, err := c.ReadMessage()
	if err != nil {
		return true
	}

	message := string(content)

	if strings.HasPrefix(message, SubscribePrefix) {
		c.AddTopic(strings.Replace(message, SubscribePrefix, "", 1))
		return false
	}

	if strings.HasPrefix(message, UnSubscribePrefix) {
		c.RemoveTopic(strings.Replace(message, UnSubscribePrefix, "", 1))
		return false
	}

	if message == CloseWebSocket {
		return true
	}

	return false
}

func (c *WsConn) Send(message []byte) {
	if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
		c.Close()
	}
}

func NewWsConn(c *gin.Context) (*WsConn, error) {
	ws, err := wsUpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewRandom()

	conn := WsConn{
		Conn:   ws,
		topics: make(map[string]struct{}),
		token:  uid.String(),
	}

	connPool[conn.token] = &conn
	return &conn, nil
}

type PublishMessage struct {
	Topic   string
	Message []byte
}

var messageChannel chan PublishMessage

// 发布消息接口
func Publish(topic string, message []byte) {
	messageChannel <- PublishMessage{
		Topic:   topic,
		Message: message,
	}
}

func messageDeliver() {
	for {
		select {
		case message := <-messageChannel:
			go handlePublish(message)
		}
	}
}

func handlePublish(pm PublishMessage) {
	for _, conn := range connPool {
		if conn.IsSubscriber(pm.Topic) {
			conn.Send(pm.Message)
		}
	}
}

func init() {
	connPool = make(map[string]*WsConn)
	messageChannel = make(chan PublishMessage, 10)
	go messageDeliver()
}
