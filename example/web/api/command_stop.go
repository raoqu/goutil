package api

import "github.com/raoqu/goutil/example/web/process"

func APICommandStop(uuid string) (bool, error) {
	process.MANAGER.StopCommand(uuid)
	return true, nil
}
