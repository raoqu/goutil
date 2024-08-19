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

func (*WSSOutput) OnMessage(client *web.WSSClient, message string) {
	if message != "<<<keePAlive>>>" {
		client.Set("uuid", message)
		client.Print("UUID: " + client.Get("uuid"))
	}
}

func (*WSSOutput) BeforeBroadcast(client *web.WSSClient, str string) (string, bool) {
	target, message := types.Split2(str, ":")
	uuid := client.Get("uuid")

	return message, target == uuid
}
