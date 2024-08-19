package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var WSSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSSClient struct {
	WSSHub *WSSHub
	conn   *websocket.Conn
	send   chan string

	Data     map[string]string
	Instance WSSInstance
}

var uniqueIndex = 1

type WSSCreateCallback func(*WSSClient)
type WSSMessageCallback func(*WSSClient, string)

type WSSInstance interface {
	OnCreate(client *WSSClient)
	OnMessage(client *WSSClient, message string) // message processor
	BeforeBroadcast(client *WSSClient, message string) (string, bool)
}

func NewWSSClient(conn *websocket.Conn) *WSSClient {
	client := WSSClient{
		conn: conn,
		send: make(chan string),
		Data: make(map[string]string),
	}
	client.Set("index", strconv.Itoa(uniqueIndex))
	uniqueIndex = uniqueIndex + 1
	return &client
}

func (c *WSSClient) Send(msg string) {
	c.send <- msg
}

func (c *WSSClient) Set(key string, value string) {
	c.Data[key] = value
}

func (c *WSSClient) Get(key string) string {
	return c.Data[key]
}

func (c *WSSClient) Close() {
	if c.WSSHub != nil {
		c.WSSHub.Remove(c)
	}
}

func (c *WSSClient) Print(str string) {
	fmt.Printf("%s : %s\n", c.Get("index"), str)
}

const WSS_FIRST_MESSAGE = 3 * time.Second
const WSS_KEEP_ALIVE = 3 * time.Second
const WSS_READ_LIMIT = 512
const WSS_WRITE_TIMEOUT = 3 * time.Second

func (c *WSSClient) readPump() {
	defer func() {
		c.Close()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(WSS_READ_LIMIT)
	c.conn.SetReadDeadline(time.Now().Add(WSS_FIRST_MESSAGE))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(WSS_KEEP_ALIVE)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Print("Error reading ...")
				// log.Printf("error: %v", err)
			}
			break
		}

		c.conn.SetReadDeadline(time.Now().Add(WSS_KEEP_ALIVE))
		messageStr := string(message)
		if messageStr != "<<<keePAlive>>>" {
			c.Instance.OnMessage(c, messageStr)
		}
	}
}

func (c *WSSClient) writePump() {
	ticker := time.NewTicker(WSS_KEEP_ALIVE)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				break
			}
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if !ok {
			// 	// The hub closed the channel.
			// 	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			// 	return
			// }

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message))
			// Add queued chat messages to the current websocket message.
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WSS_WRITE_TIMEOUT))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
