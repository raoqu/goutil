package main

import (
	"github.com/raoqu/goutil/example/web/api"
	"github.com/raoqu/goutil/web"
)

func main() {
	server := web.NewServer(7777, true, "assets")

	println("Listening on", server.Address())

	server.API("stat", web.NewAPI(api.APIStat))

	server.API("commands", web.NewAPI(api.APICommands))
	server.API("command_add", web.NewAPI(api.APICOmmandAdd))
	server.API("command_remove", web.NewAPI(api.APICommandRemove))
	server.API("command_start", web.NewAPI(api.APICommandStart))
	server.API("command_stop", web.NewAPI(api.APICommandStop))

	server.API("configs", web.NewAPI(api.APIConfig))
	server.API("config_update", web.NewAPI(api.APIConfigUpdate))

	server.Start()

	select {}
}
