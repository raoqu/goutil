package main

import (
	"flag"

	"github.com/raoqu/goutil/main/api"
	"github.com/raoqu/goutil/shell"
	"github.com/raoqu/goutil/web"
)

type ApplicationArgs struct {
	Local bool
	Port  int
}

var AppConfig = ApplicationArgs{
	Local: true,
	Port:  8080,
}

func parseArgs() {

	flag.BoolVar(&AppConfig.Local, "local", AppConfig.Local, "Start web server for localhost only.")
	flag.IntVar(&AppConfig.Port, "port", AppConfig.Port, "Web server port (default to 8080).")

	flag.Parse()
}

func main() {
	parseArgs()
	println("local", AppConfig.Local)
	println("port", AppConfig.Port)

	var address = "0.0.0.0"
	if AppConfig.Local {
		address = "127.0.0.1"
	}

	var server = web.WebServer(address, AppConfig.Port, "./assets")
	web.WebAPIRegister(&server, "start", api.StartCommand)

	var parts = shell.SplitCommand(`command "param 1" "param 2" -d target "the last \"param\""`)
	for _, item := range parts {
		println(item)
	}

	// select {}
}
