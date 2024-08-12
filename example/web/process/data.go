package process

import "github.com/raoqu/goutil/shell"

type Command struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Config struct {
	Uuid    string `json:"uuid"`
	Dir     string `json:"dir"`
	Command string `json:"command"`
}

type Manager struct {
	Instances    []string            `json:"instances"`
	Commands     map[string]Command  `json:"commands"`
	Configs      map[string]Config   `json:"configs"`
	ShellManager *shell.ShellManager `json:"-"`
}
