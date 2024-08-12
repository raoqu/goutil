package api

import (
	"github.com/raoqu/goutil/example/web/process"
	"github.com/raoqu/goutil/web"
)

func APIConfig(req *web.Any) (map[string]process.Config, error) {
	return process.MANAGER.Configs, nil
}
