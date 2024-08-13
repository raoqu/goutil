package api

import "github.com/raoqu/goutil/example/web/process"

func APICommandStop(uuid string) (bool, error) {
	stopped := process.MANAGER.StopCommand(uuid)
	return stopped, nil
}
