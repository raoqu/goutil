package shell

import "strings"

type ShellStat struct {
	psoutput []string
}

func (s *ShellStat) Check(str string) bool {
	isWindows := OSFeature(true, false)
	if isWindows {
		return s.win_Check(str)
	} else {
		return s.unix_Check(str)
	}
}

func (s *ShellStat) win_Check(str string) bool {
	panic("not implented yet.")
}

func (s *ShellStat) unix_Check(str string) bool {
	if len(s.psoutput) == 0 {
		command := NewCommand("ps", false)
		command.BufferLines = 200
		command.OnOutput = func(string) {}
		command.Run()
		s.psoutput = command.GetOutput()
	}

	return outputContains(s.psoutput, str)
}

func outputContains(arr []string, sub string) bool {
	for _, str := range arr {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}
