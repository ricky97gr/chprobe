package utils

import (
	"net"
	"os"
	"runtime"
	"strings"
	"os/exec"
)

func GetHostname() (string, error) {
	return os.Hostname()
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", nil
}

func GetMachineID() (string, error) {
	// 尝试从Linux系统的machine-id文件读取
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return string(data), nil
	}

	// 尝试从Windows系统的注册表读取（这里简化处理，实际需要使用syscall）
	// 对于其他系统，可以返回空字符串或基于其他信息生成
	return "", nil
}

func GetOsType() string {
	return runtime.GOOS
}

func GetOs() string {
	if runtime.GOOS == "linux" {
		// 尝试读取 /etc/os-release 文件获取发行版信息
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				}
			}
		}
		return "Linux"
	}
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

func GetKernelVersion() string {
	if runtime.GOOS == "linux" {
		// 执行 uname -r 命令获取内核版本
		cmd := exec.Command("uname", "-r")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	return "unknown"
}
