package api

import (
	"github.com/raoqu/goutil/example/web/process"
)

func APICOmmandAdd(req process.Command) (bool, error) {
	return process.MANAGER.AddCommand(req), nil
}
