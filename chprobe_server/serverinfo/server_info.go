package serverinfo

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
)

// GetServerInfo 获取服务器基本信息
func GetServerInfo() models.ServerInfo {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// 获取IP地址
	ip := getLocalIP()

	// 获取内核版本
	kernel := getKernelVersion()

	// 获取CPU信息
	cpu := getCPUInfo()

	// 获取内存信息
	memory := getMemoryInfo()

	// 生成产品序列号（基于主机名和IP地址）
	serial := generateProductSerial(hostname, ip)

	// 构建服务器信息
	serverInfo := models.ServerInfo{
		Hostname:    hostname,
		IP:          ip,
		Kernel:      kernel,
		CPU:         cpu,
		Memory:      memory,
		Serial:      serial,                                   // 产品序列号
		Version:     "1.0.0",                                  // 版本号，实际项目中应该从构建参数中获取
		CommitID:    "",                                       // CommitID，实际项目中应该从构建参数中获取
		BuildTime:   time.Now().Format("2006-01-02 15:04:05"), // 编译时间，实际项目中应该从构建参数中获取
		ProductName: "CHProbe",                                // 产品名称
		StartupTime: time.Now().UnixMilli(),                   // 启动时间
	}

	return serverInfo
}

// UpdateServerInfo 更新服务器信息到数据库
func UpdateServerInfo() error {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		return err
	}

	// 获取服务器信息
	serverInfo := GetServerInfo()

	// 检查是否已存在服务器信息
	var existingServerInfo models.ServerInfo
	result := db.First(&existingServerInfo)

	if result.Error == nil {
		// 已存在，更新
		existingServerInfo.Hostname = serverInfo.Hostname
		existingServerInfo.IP = serverInfo.IP
		existingServerInfo.Kernel = serverInfo.Kernel
		existingServerInfo.CPU = serverInfo.CPU
		existingServerInfo.Memory = serverInfo.Memory
		existingServerInfo.Serial = serverInfo.Serial
		existingServerInfo.Version = serverInfo.Version
		existingServerInfo.CommitID = serverInfo.CommitID
		existingServerInfo.BuildTime = serverInfo.BuildTime
		existingServerInfo.ProductName = serverInfo.ProductName
		existingServerInfo.StartupTime = serverInfo.StartupTime

		if err := db.Save(&existingServerInfo).Error; err != nil {
			return err
		}
	} else {
		// 不存在，插入
		if err := db.Create(&serverInfo).Error; err != nil {
			return err
		}
	}

	return nil
}

// getLocalIP 获取本地IP地址
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}

	return "unknown"
}

// getKernelVersion 获取内核版本
func getKernelVersion() string {
	if runtime.GOOS == "linux" {
		// 在Linux系统上，可以通过读取/proc/version文件获取内核版本
		data, err := os.ReadFile("/proc/version")
		if err == nil {
			version := string(data)
			parts := strings.Fields(version)
			if len(parts) > 2 {
				return parts[2]
			}
		}
	} else if runtime.GOOS == "windows" {
		// 在Windows系统上，可以通过系统命令获取内核版本
		// 这里简化处理，返回操作系统信息
		return runtime.GOOS + " " + runtime.GOARCH
	}

	return runtime.GOOS + " " + runtime.GOARCH
}

// getCPUInfo 获取CPU信息
func getCPUInfo() string {
	if runtime.GOOS == "linux" {
		// 在Linux系统上，可以通过读取/proc/cpuinfo文件获取CPU信息
		data, err := os.ReadFile("/proc/cpuinfo")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "model name\t") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						return strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}

	// 简化处理，返回CPU核心数
	return fmt.Sprintf("%d cores", runtime.NumCPU())
}

// getMemoryInfo 获取内存信息
func getMemoryInfo() string {
	if runtime.GOOS == "linux" {
		// 在Linux系统上，可以通过读取/proc/meminfo文件获取内存信息
		data, err := os.ReadFile("/proc/meminfo")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "MemTotal:") {
					parts := strings.Fields(line)
					if len(parts) > 2 {
						memKB, err := strconv.ParseInt(parts[1], 10, 64)
						if err == nil {
							memGB := float64(memKB) / 1024 / 1024
							return fmt.Sprintf("%.2f GB", memGB)
						}
					}
				}
			}
		}
	}

	// 简化处理，返回未知
	return "unknown"
}

// getMachineInfo 获取机器信息（包含machine-id和MAC地址）
func getMachineInfo() (string, string) {
	var machineID, macAddr string

	// 在Linux系统上，尝试从/etc/machine-id获取
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/etc/machine-id")
		if err == nil && len(data) > 0 {
			machineID = strings.TrimSpace(string(data))
		}
	}

	// 尝试获取MAC地址
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if iface.HardwareAddr != nil && len(iface.HardwareAddr) > 0 && iface.Name != "lo" {
				macAddr = iface.HardwareAddr.String()
				break
			}
		}
	}

	return machineID, macAddr
}

// generateProductSerial 生成产品序列号
func generateProductSerial(hostname, ip string) string {
	// 获取机器信息
	machineID, macAddr := getMachineInfo()

	// 产品ID
	productID := "CHProbe"

	// 使用所有信息生成稳定的序列号
	data := fmt.Sprintf("%s-%s-%s-%s-%s", productID, machineID, ip, hostname, macAddr)

	// 生成UUID格式的序列号
	// 使用MD5哈希作为基础
	hash := md5.Sum([]byte(data))

	// 按照UUID格式格式化
	// UUID格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		hash[0:4],   // 8位
		hash[4:6],   // 4位
		hash[6:8],   // 4位
		hash[8:10],  // 4位
		hash[10:16], // 12位
	)

	return uuid
}
