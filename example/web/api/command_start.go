package api

import "github.com/raoqu/goutil/example/web/process"

func APICommandStart(uuid string) (bool, error) {
	return process.MANAGER.StartCommand(uuid), nil
}
