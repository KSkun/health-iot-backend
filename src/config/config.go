package config

import (
	"github.com/KSkun/health-iot-backend/global"
	"gopkg.in/yaml.v3"
	"os"
)

const defaultConfigPath = "default.yml"

type Config struct {
	AppConfig struct {
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
	} `yaml:"app"`
	Debug bool `yaml:"debug"`
}

var C Config
var Debug bool

// InitConfig Read config from default.yml or specified config file
func InitConfig() {
	configPath, found := os.LookupEnv("CONFIG_FILE")
	if !found {
		global.Logger.Infof("[Config] CONFIG_FILE env not found, use default config file %s instead", defaultConfigPath)
		configPath = defaultConfigPath
	}
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		global.Logger.Fatalf("[Config] Error when reading config file %s, %s", configPath, err.Error())
	}
	err = yaml.Unmarshal(configBytes, &C)
	if err != nil {
		global.Logger.Fatalf("[Config] Error when reading config file %s, %s", configPath, err.Error())
	}
	Debug = C.Debug
}
