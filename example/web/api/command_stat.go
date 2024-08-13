package api

import "github.com/raoqu/goutil/example/web/process"

func APICommandStat(uuid string) (process.Stat, error) {
	return process.MANAGER.GetStat(uuid)
}
