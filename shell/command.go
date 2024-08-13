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
	START
	FAIL
	COMPLETE
	UNKNOWN
)

const DEFAULT_BUFFER_LINES = 50

type Command struct {
	Uuid string
	// basic params
	Command  string
	Dir      string
	Async    bool
	Attached bool
	// callbacks
	OnOutput    func(string)
	OnClose     func()
	BufferLines int
	Status      CommandStatus
	// alive check
	AliveConfig AliveCheckConfig

	Manager *ShellManager
	Exec    *exec.Cmd
	Err     error

	outputWriter *LineBufferWriter
	extBuffer    []string
}

func NewCommand(command string, async bool) Command {
	return NewCommandWithWorkDir(command, async, ".")
}

func NewCommandWithWorkDir(command string, async bool, dir string) Command {
	uuid := util.UUID()

	return Command{
		Uuid:    uuid,
		Command: command,
		Dir:     dir,
		Async:   async,

		BufferLines: DEFAULT_BUFFER_LINES,
		Status:      INIT,

		extBuffer: make([]string, 0),
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
