package shell

import (
	"os/exec"
	"strings"
	"unicode"

	"github.com/raoqu/goutil/util"
)

type CommandStatus int

const (
	INIT CommandStatus = iota
	FAIL
	COMPLETE
)

const DEFAULT_BUFFER_LINES = 50

type Command struct {
	UniqueId string
	// basic params
	Command string
	Dir     string
	Async   bool
	// callbacks
	OnOutput    func(string)
	OnClose     func()
	BufferLines int
	Status      CommandStatus

	Manager *ShellManager
	Exec    *exec.Cmd
	Err     error

	outputWriter *LineBufferWriter
}

func NewCommand(command string, async bool) Command {
	return NewCommandWithWorkDir(command, async, ".")
}

func NewCommandWithWorkDir(command string, async bool, dir string) Command {
	uniqueId := util.UUID()

	return Command{
		UniqueId: uniqueId,
		Command:  command,
		Dir:      dir,
		Async:    async,

		BufferLines: DEFAULT_BUFFER_LINES,
		Status:      INIT,
	}
}

func SplitCommand(command string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for _, char := range command {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			continue
		}

		if char == '"' {
			inQuotes = !inQuotes
			continue
		}

		if unicode.IsSpace(char) && !inQuotes {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteRune(char)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
