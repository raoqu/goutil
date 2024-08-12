package main

import (
	"github.com/raoqu/goutil/example/web/api"
	"github.com/raoqu/goutil/web"
)

func main() {
	server := web.NewServer(7777, true, "assets")

	println("Listening on", server.Address())

	server.API("stat", web.NewAPI(api.APIStat))
	server.API("config", web.NewAPI(api.APIConfig))
	server.Start()

	select {}
}
