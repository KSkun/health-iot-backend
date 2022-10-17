package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const defaultConfigPath = "default.yml"

type Config struct {
	AppConfig struct {
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
	} `yaml:"app"`
	MongoConfig struct {
		Addr     string `yaml:"addr"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"mongo"`
	Debug bool `yaml:"debug"`
}

var C Config
var Debug bool

// InitConfig Read config from default.yml or specified config file
func InitConfig() {
	configPath, found := os.LookupEnv("CONFIG_FILE")
	if !found {
		log.Printf("[Config] CONFIG_FILE env not found, use default config file %s instead", defaultConfigPath)
		configPath = defaultConfigPath
	}
	log.Printf("[Config] Loading config file %s", configPath)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("[Config] Error when reading config file %s, %s", configPath, err.Error())
	}
	err = yaml.Unmarshal(configBytes, &C)
	if err != nil {
		log.Fatalf("[Config] Error when reading config file %s, %s", configPath, err.Error())
	}
	Debug = C.Debug
	log.Printf("[Config] Init done")
}
