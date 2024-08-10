package shell

import (
	"sync"

	"github.com/raoqu/goutil/types"
)

type ShellManager struct {
	Processes map[string]*Command
	Mutex     sync.Mutex
}

func NewShellManager() ShellManager {
	return ShellManager{
		Processes: make(map[string]*Command),
	}
}

func (g *ShellManager) Start(command Command) {
	command.Manager = g
	command.Run()
}

func (g *ShellManager) IsAlive(uniqueId string) bool {
	_, exists := g.Processes[uniqueId]
	return exists
}

func (g *ShellManager) List() []*Command {
	return types.Map2Array(g.Processes)
}
