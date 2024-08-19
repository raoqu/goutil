package web

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSSInitCallback func(*WSSClient)
type WSSAPI func(*WSSClient, string)

func NewWebSocketCallback(callback WSSAPI) WSSAPI {
	return callback
}

func RegisterWebSocket(s *Server, endpoint string, instance WSSInstance) {
	s.WSSEndpointCreate[endpoint] = instance
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
}

func wssHandler(instance WSSInstance, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewWSSClient(conn)
	client.Instance = instance

	instance.OnCreate(client)

	go client.writePump()
	go client.readPump()
}
