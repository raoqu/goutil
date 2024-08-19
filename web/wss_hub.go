package web

import (
	"sort"
)

type WSSHub struct {
	clients    map[*WSSClient]bool
	broadcast  chan BroadcastMessage
	register   chan *WSSClient
	unregister chan *WSSClient
}

type BroadcastMessage struct {
	Group   string
	Message string
}

func NewWSSHub() *WSSHub {
	return &WSSHub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *WSSClient),
		unregister: make(chan *WSSClient),
		clients:    make(map[*WSSClient]bool),
	}
}

func (h *WSSHub) Broadcast(message string, groupId string) bool {
	if len(groupId) < 1 {
		groupId = "*"
	}
	h.broadcast <- BroadcastMessage{
		Group:   groupId,
		Message: message,
	}
	return true
}

func (h *WSSHub) Add(client *WSSClient) {
	h.register <- client
}

func (h *WSSHub) Remove(client *WSSClient) {
	client.Print("remove")
	if h.clients[client] {
		h.unregister <- client
	}
}

func (h *WSSHub) Run() {
	for {
		select {
		case client := <-h.register:
			client.Print("register")
			h.clients[client] = true
			client.WSSHub = h
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.Print("delete")
				delete(h.clients, client)
				close(client.send)
			}
		case broadMessage := <-h.broadcast:
			print(client2string(h.clients))
			for client := range h.clients {
				if broadMessage.Group == "*" || broadMessage.Group == client.Group {
					select {
					case client.send <- broadMessage.Message:

					}
				}
			}
		}
	}
}

func client2string(clients map[*WSSClient]bool) string {
	arr := []string{}
	for c := range clients {
		arr = append(arr, c.Get("index"))
	}
	str := Array2String(arr, true)
	return str
}

func Array2String(arr []string, asc bool) string {
	if asc {
		sort.Strings(arr)
	}
	str := "["
	for i, item := range arr {
		if i > 0 {
			str = str + ", "
		}
		str = str + item
	}
	str = str + "] "
	if str == "[] " {
		return ""
	}
	return str
}
