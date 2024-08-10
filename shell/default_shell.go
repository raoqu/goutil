package shell

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var DEFAULT_SHELL string

func getDefaultShell() string {
	if len(DEFAULT_SHELL) > 0 {
		return DEFAULT_SHELL
	}

	shell := getShellFromEnv()
	if shell != "" {
		DEFAULT_SHELL = checkShellPriority(shell)
		return DEFAULT_SHELL
	}

	shell = pickShell()
	if shell != "" {
		DEFAULT_SHELL = shell
		return DEFAULT_SHELL
	}

	// shell = getShellFromPasswd()
	// if shell != "" {
	//  DEFAULT_SHELL = checkShellPriority(shell)
	// 	return DEFAULT_SHELL
	// }

	DEFAULT_SHELL = "/bin/sh"
	return DEFAULT_SHELL
}

func getShellFromEnv() string {
	if shell, exists := getEnv("SHELL"); exists {
		return shell
	}
	return ""
}

// 从 /etc/passwd 获取 shell
func getShellFromPasswd() string {
	output, err := exec.Command("getent", "passwd").Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Split(line, ":")
				if len(parts) > 6 {
					return strings.TrimSpace(parts[6]) // 第7个字段是用户的 shell
				}
			}
		}
	}
	return ""
}

// 检查 shell 的优先级
func checkShellPriority(shell string) string {
	if strings.HasSuffix(shell, "zsh") {
		return "/bin/zsh"
	} else if strings.HasSuffix(shell, "bash") {
		return "/bin/bash"
	} else if strings.HasSuffix(shell, "sh") {
		return "/bin/sh"
	}
	return ""
}

// 获取环境变量的自定义函数
func getEnv(key string) (string, bool) {
	value, exists := syscall.Getenv(key)
	return value, exists
}

// 根据优先级检查并选择 shell
func pickShell() string {
	if shellExists("/bin/zsh") {
		return "/bin/zsh"
	} else if shellExists("/bin/bash") {
		return "/bin/bash"
	} else if shellExists("/bin/sh") {
		return "/bin/sh"
	}
	return ""
}

// 检查指定 shell 是否存在
func shellExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
