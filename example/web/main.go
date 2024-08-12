package main

import (
	"github.com/raoqu/goutil/example/web/api"
	"github.com/raoqu/goutil/web"
)

func main() {
	server := web.NewServer(7777, true, "assets")

	println("Listening on", server.Address())

	web.RegisterAPI(server, "stat", api.APIStat)

	web.RegisterAPI(server, "commands", api.APICommands)
	web.RegisterAPI(server, "command_add", api.APICOmmandAdd)
	web.RegisterAPI(server, "command_remove", api.APICommandRemove)
	web.RegisterAPI(server, "command_start", api.APICommandStart)
	web.RegisterAPI(server, "command_stop", api.APICommandStop)

	web.RegisterAPI(server, "configs", api.APIConfig)
	web.RegisterAPI(server, "config_update", api.APIConfigUpdate)

	server.Start()

	select {}
}
