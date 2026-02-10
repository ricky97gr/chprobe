package conf

import (
	"os"
	"path/filepath"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
	"go.yaml.in/yaml/v3"
)

const (
	ConfigFileName = "config.yaml"
)

type Config struct {
	UUID       string `yaml:"uuid"`
	ServerIP   string `yaml:"serverIP"`
	ServerPort string `yaml:"serverPort"`
}

func GetServerAddr() string {
	ip := AppConfig.ServerIP
	if ip == "" {
		ip = "127.0.0.1"
	}
	port := AppConfig.ServerPort
	if port == "" {
		port = "32000"
	}
	return ip + ":" + port
}

var AppConfig *Config

func init() {
	AppConfig = &Config{}
	loadConfig()
}

func loadConfig() {
	configPath := filepath.Join(getConfigDir(), ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		utils.Logger.Infof("config file not exist, will create new one\n")
		saveConfig()
		return
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		utils.Logger.Errorf("read config file failed, err: %v\n", err)
		return
	}

	err = yaml.Unmarshal(content, AppConfig)
	if err != nil {
		utils.Logger.Errorf("unmarshal config file failed, err: %v\n", err)
		return
	}
}

func saveConfig() {
	configPath := filepath.Join(getConfigDir(), ConfigFileName)
	content, err := yaml.Marshal(AppConfig)
	if err != nil {
		utils.Logger.Errorf("marshal config failed, err: %v\n", err)
		return
	}

	err = os.WriteFile(configPath, content, 0644)
	if err != nil {
		utils.Logger.Errorf("write config file failed, err: %v\n", err)
		return
	}
}

func getConfigDir() string {
	// 首先尝试从当前工作目录的conf子目录读取
	configDir := filepath.Join(getCurrentDir(), "conf")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 如果conf目录不存在，尝试从可执行文件目录读取
		configDir = filepath.Dir(os.Args[0])
		if configDir == "." {
			configDir, _ = os.Getwd()
		}
	}
	return configDir
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func SetUUID(uuid string) {
	AppConfig.UUID = uuid
	saveConfig()
}

func GetUUID() string {
	return AppConfig.UUID
}
