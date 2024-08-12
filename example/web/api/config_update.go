package api

import "github.com/raoqu/goutil/example/web/process"

func APIConfigUpdate(config process.Config) (bool, error) {
	process.MANAGER.SetConfig(config)
	return true, nil
}
