package api

import "github.com/raoqu/goutil/web"

type ShellConfig struct {
	Name    string `json:"name"`
	Dir     string `json:"dir"`
	Command string `json:"command"`
}

var configs = []ShellConfig{
	{
		Name:    "ping",
		Dir:     ".",
		Command: "ping baidu.com -t 5",
	},
	{
		Name:    "tomcat",
		Dir:     "/Users/raoqu/data/tomcat/tomcat1",
		Command: "sh bin/start-tomcat.sh",
	},
}

func APIConfig(req *web.Any) ([]ShellConfig, error) {
	return configs, nil
}
