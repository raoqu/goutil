package shell

import "strings"

type ShellStat struct {
	psoutput []string
}

func (s *ShellStat) Check(output string) bool {
	return s.CheckOutput("ps -ef", output)
}

func (s *ShellStat) CheckOutput(command string, output string) bool {
	isWindows := OSFeature(true, false)
	if isWindows {
		return s.win_Check(command, output)
	} else {
		return s.unix_Check(command, output)
	}
}

func (s *ShellStat) win_Check(command string, output string) bool {
	panic("not implented yet.")
}

func (s *ShellStat) unix_Check(command string, output string) bool {
	if len(s.psoutput) == 0 {
		command := NewCommand("ps -ef", false)
		command.BufferLines = 200
		command.OnOutput = func(string) {}
		command.Run()
		s.psoutput = command.GetOutput()
	}

	return s.OutputContains(output)
}

func (s *ShellStat) GetOutput() []string {
	return s.psoutput
}

func (s *ShellStat) OutputContains(output string) bool {
	return s.OutputContainCount(output) > 0
}

func (s *ShellStat) OutputContainCount(output string) int {
	return outputContainCount(s.psoutput, output)
}

func outputContainCount(arr []string, sub string) int {
	count := 0
	for _, str := range arr {
		if strings.Contains(str, sub) {
			count = count + 1
		}
	}
	return count
}
