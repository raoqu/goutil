package api

import (
	"github.com/raoqu/goutil/example/web/process"
	"github.com/raoqu/goutil/web"
)

func APICommands(req *web.Any) ([]process.Command, error) {
	return process.MANAGER.GetCommands(), nil
}
