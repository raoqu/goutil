package api

import "github.com/raoqu/goutil/example/web/process"

func APICommandRemove(uuid string) (bool, error) {
	return process.MANAGER.RemoveCommand(uuid), nil
}
