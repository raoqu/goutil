package web

import (
	"sort"
)

type WSSHub struct {
	clients    map[*WSSClient]bool
	broadcast  chan string
	register   chan *WSSClient
	unregister chan *WSSClient
}

func NewWSSHub() *WSSHub {
	return &WSSHub{
		broadcast:  make(chan string),
		register:   make(chan *WSSClient),
		unregister: make(chan *WSSClient),
		clients:    make(map[*WSSClient]bool),
	}
}

func (h *WSSHub) Broadcast(message string) bool {
	h.broadcast <- message
	return true
}

func (h *WSSHub) Add(client *WSSClient) {
	h.register <- client
}

func (h *WSSHub) Remove(client *WSSClient) {
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
		case message := <-h.broadcast:
			print(client2string(h.clients))

			for client := range h.clients {
				doit := true
				var newMessage string
				if instance := client.Instance; instance != nil {
					newMessage, doit = instance.BeforeBroadcast(client, string(message))
				} else {
					newMessage = message
				}
				if doit {
					select {
					case client.send <- newMessage:
					default:
						client.Print("close")
						close(client.send)
						delete(h.clients, client)
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
