package api

import (
	"github.com/raoqu/goutil/example/web/process"
	"github.com/raoqu/goutil/web"
)

func APIBatchStat(req *web.Any) (map[string]int, error) {
	return process.MANAGER.BatchStat(), nil
}
