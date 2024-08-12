goutil is a common tool module for Golang development. It is designed to not rely on other third-party modules to avoid nested dependencies of modules, thereby reducing the size of the compiled target application.

1. The first step to using this module is to import it:

```shell
go get github.com/raoqu/goutil@latest
```

2. Then use the various tools provided by the module in your code.

```go
package main

import "github.com/raoqu/goutil/shell"

func main() {
	run("ping google.com -t 1")
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
```
