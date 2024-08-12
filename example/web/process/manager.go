package process

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/raoqu/goutil/shell"
)

const CONFIG_FILENAME = ".quprocess"

var ALL_COMMANDS = make([]string, 0) // uuid list

var MANAGER = LoadFromFile()

func (m *Manager) GetCommands() []Command {
	commands := make([]Command, 0)
	for _, item := range m.Instances {
		cmd, _ := m.Commands[item]
		commands = append(commands, cmd)
	}
	return commands
}

func (m *Manager) AddCommand(cmd Command) bool {
	m.Instances = append(m.Instances, cmd.Uuid)
	m.Commands[cmd.Uuid] = cmd
	m.Save()
	return true
}

func (m *Manager) RemoveCommand(uuid string) bool {
	instances := []string{}
	for _, item := range m.Instances {
		if item != uuid {
			instances = append(instances, item)
		}
	}
	m.Instances = instances

	delete(m.Commands, uuid)
	delete(m.Configs, uuid)
	m.Save()
	return true
}

func (m *Manager) StartCommand(uuid string) bool {
	config, exists := m.Configs[uuid]
	if exists {
		command := shell.NewCommandWithWorkDir(config.Command, true, config.Dir)
		m.ShellManager.Start(command)
		return true
	}
	return false
}

func (m *Manager) StopCommand(uuid string) {
	// m.ShellManager
}

func (m *Manager) GetConfig(uuid string) Config {
	config, exists := m.Configs[uuid]
	if exists {
		return config
	}
	return Config{
		Uuid: uuid,
	}
}

func (m *Manager) SetConfig(config Config) bool {
	uuid := config.Uuid
	m.Configs[uuid] = config

	m.Save()
	return true
}

func (m *Manager) Save() bool {
	usr, err := user.Current()
	if err != nil {
		return false
	}

	configPath := filepath.Join(usr.HomeDir, CONFIG_FILENAME)

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return false
	}

	return true
}

func initManager(manager *Manager) Manager {
	shellManager := shell.NewShellManager()
	if manager != nil {
		manager.ShellManager = &shellManager
	} else {
		manager = &Manager{
			Instances:    make([]string, 0),
			Commands:     make(map[string]Command),
			Configs:      make(map[string]Config),
			ShellManager: &shellManager,
		}
	}
	return *manager
}

func LoadFromFile() Manager {
	var manager Manager

	usr, err := user.Current()
	if err != nil {
		return initManager(nil)
	}

	configPath := filepath.Join(usr.HomeDir, CONFIG_FILENAME)

	// 尝试打开配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return initManager(nil)
	}
	defer file.Close()

	// 解析 JSON 数据
	if err := json.NewDecoder(file).Decode(&manager); err != nil {
		return initManager(nil)
	}

	return initManager(&manager)
}
