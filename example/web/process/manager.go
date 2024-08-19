package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/raoqu/goutil/shell"
	"github.com/raoqu/goutil/web"
)

const CONFIG_FILENAME = ".quprocess"

var ALL_COMMANDS = make([]string, 0) // uuid list

var MANAGER = LoadFromFile()

type Manager struct {
	Instances    []string                  `json:"instances"`
	Commands     map[string]Command        `json:"commands"`
	Configs      map[string]Config         `json:"configs"`
	ShellManager *shell.ShellManager       `json:"-"`
	WSSClients   map[string]*web.WSSClient `json:"-"`
	WSSHub       *web.WSSHub               `json:"-"`
}

func LoadFromFile() Manager {
	var manager *Manager = nil
	usr, err := user.Current()
	if err == nil {
		configPath := filepath.Join(usr.HomeDir, CONFIG_FILENAME)

		file, err := os.Open(configPath)
		if err == nil {
			defer file.Close()

			// deserialize from configuration file
			var deserialized Manager
			if json.NewDecoder(file).Decode(&deserialized) == nil {
				manager = &deserialized
			}
		}
	}

	if manager == nil {
		manager = NewManager()
	}

	shellManager := shell.NewShellManager()

	manager.ShellManager = &shellManager
	return *manager
}

func NewManager() *Manager {
	return &Manager{
		Instances:  make([]string, 0),
		Commands:   make(map[string]Command),
		Configs:    make(map[string]Config),
		WSSClients: make(map[string]*web.WSSClient),
	}
}

func (m *Manager) StartWSSHub() *web.WSSHub {
	if m.WSSHub == nil {
		m.WSSHub = web.NewWSSHub()
		go m.WSSHub.Run()
	}
	return m.WSSHub
}

func (m *Manager) GetCommands() []Command {
	commands := make([]Command, 0)
	for _, item := range m.Instances {
		cmd := m.Commands[item]
		// update status
		stat := m.ShellManager.GetStatus(cmd.Uuid, false)
		status := shell.MapCommandStatus(stat)
		if cmd.Status != status {
			cmd.Status = status
			m.Commands[item] = cmd
		}

		commands = append(commands, cmd)
	}
	return commands
}

func (m *Manager) GetCommand(uuid string) (Command, bool) {
	command, exists := m.Commands[uuid]
	return command, exists
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

func (m *Manager) GetStat(uuid string) (Stat, error) {
	config, exists := m.Configs[uuid]
	if !exists {
		return Stat{
			Status: "unknown",
		}, errors.New("not configured")
	}

	shellCommand := m.ShellManager.Get(uuid)
	if shellCommand == nil {
		cmd := shell.Command{
			Uuid:     uuid,
			Command:  config.Command,
			Dir:      config.Dir,
			Attached: true,
			Async:    true,
			AliveConfig: shell.AliveCheckConfig{
				Ping: config.Ping,
			},
		}
		m.ShellManager.Attach(cmd)
		shellCommand = &cmd
	}

	// Get command alive status
	cmdStatus := m.ShellManager.GetStatus(uuid, true)
	status := shell.MapCommandStatus(cmdStatus)
	output := shellCommand.GetOutput()
	return Stat{
		Status: status,
		Output: output,
	}, nil
}

func (m *Manager) BatchStat() map[string]int {
	result := make(map[string]int)
	stat := shell.ShellStat{}
	stat.Check("trigger ps checking")

	for _, uuid := range m.Instances {
		config, exists := m.Configs[uuid]
		if exists && len(strings.TrimSpace(config.Ping)) > 0 {
			result[uuid] = stat.OutputContainCount(config.Ping)
		} else {
			result[uuid] = 0
		}
	}

	return result
}

func (m *Manager) StartCommand(uuid string) bool {
	config, exists := m.Configs[uuid]
	if exists {
		command := shell.NewCommandWithWorkDir(config.Command, true, config.Dir)
		command.OnOutput = func(message string) {
			if m.WSSHub != nil {
				m.WSSHub.Broadcast(fmt.Sprintf("%s:%s", uuid, message))
			}
		}
		m.ShellManager.Start(command)
		return true
	}
	return false
}

func (m *Manager) StopCommand(uuid string) bool {
	cmd := m.ShellManager.Get(uuid)
	if cmd == nil || cmd.Exec == nil {
		config, exists := m.Configs[uuid]
		if exists {
			return m.ShellManager.Kill(config.Ping)
		}
	} else {
		if cmd.Exec.Process != nil {
			cmd.Exec.Process.Kill()
		}
	}

	return false
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
