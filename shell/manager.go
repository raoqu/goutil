package shell

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

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

func (g *ShellManager) ParseProcessIDs(ping string) []int {
	stat := ShellStat{}
	processIds := []int{}
	if stat.Check(ping) {
		outputs := stat.psoutput
		if len(outputs) > 0 {
			re := regexp.MustCompile(`^\s*[^\s]+\s+(\d+)\s+.*`)
			for _, line := range outputs {
				if strings.Contains(line, ping) {
					println("line", line)
					matches := re.FindStringSubmatch(line)
					if len(matches) > 1 {
						numStr := matches[1]
						println("pid", numStr)
						processId, err := strconv.Atoi(numStr)
						if err == nil {
							processIds = append(processIds, processId)
						}
					}
				}
			}
		}
	}
	return processIds
}

func (g *ShellManager) Kill(ping string) bool {
	ids := g.ParseProcessIDs(ping)
	if len(ids) != 1 {
		return false
	}
	pid := ids[0]

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.SIGTERM)
	return err == nil
}

func (g *ShellManager) IsAlive(uuid string) bool {
	status := g.GetStatus(uuid, true)
	return IsCommandStatusAlive(status)
}

func (g *ShellManager) GetStatus(uuid string, extStatus bool) CommandStatus {
	cmd, exists := g.Processes[uuid]
	if !exists {
		return UNKNOWN
	}

	if !extStatus {
		return cmd.Status
	}

	return cmd.UpdateStatus()
}

func (g *ShellManager) List() []*Command {
	return types.Map2Array(g.Processes)
}

func (g *ShellManager) Get(uuid string) *Command {
	cmd, exists := g.Processes[uuid]
	if !exists {
		return nil
	}
	return cmd
}

func (g *ShellManager) Attach(cmd Command) {
	cmd.Manager = g
	g.Processes[cmd.Uuid] = &cmd
}
