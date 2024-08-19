package api

import (
	"github.com/raoqu/goutil/example/web/process"
	"github.com/raoqu/goutil/types"
	"github.com/raoqu/goutil/web"
)

type WSSOutput struct {
}

var WSS_INSTANCE_OUTPUT = &WSSOutput{}

func (*WSSOutput) OnCreate(client *web.WSSClient) {
	hub := process.MANAGER.WSSHub
	if hub != nil {
		hub.Add(client)
	}
}

func (*WSSOutput) OnMessage(client *web.WSSClient, msg string) {
	msgType, message := types.Split2(msg, ":")
	switch msgType {
	case process.MSG_MESSAGE:
		client.Print(message)
	case process.MSG_GROUP:
		client.Group = message
	default:
		// client.Print(".")
	}
}
