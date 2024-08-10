package main

import "github.com/raoqu/goutil/shell"

func main() {
	run("ping baidu.com -t 1")
	ps()
}

func run(cmd string) {
	command := shell.NewCommand(cmd, false)
	command.OnOutput = func(line string) {
		println(line)
	}
	command.OnClose = func() {
		println("<<< OnClose >>>")
	}
	command.Run()
}

func ps() {
	tocheck := []string{
		"minio",
		"seata1.6.1/lib",
		"tomcat",
		"ping",
	}
	stat := shell.ShellStat{}
	for _, item := range tocheck {
		print(item)
		if stat.Check(item) {
			println(": true")
		} else {
			println(": false")
		}
	}
}
