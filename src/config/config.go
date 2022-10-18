package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

const defaultConfigPath = "default.yml"

type IntRangeConfig struct {
	Low  int `yaml:"low"`
	High int `yaml:"high"`
}

type Config struct {
	AppConfig struct {
		Addr                  string         `yaml:"addr"`
		Port                  int            `yaml:"port"`
		OnlineTimeoutSeconds  int            `yaml:"online_timeout_seconds"`
		OnlineTimeoutDuration time.Duration  `yaml:"-"`
		HeartRateThreshold    IntRangeConfig `yaml:"heart_rate_threshold"`
		BloodOxygenThreshold  int            `yaml:"blood_oxygen_threshold"`
		BatteryThreshold      int            `yaml:"battery_threshold"`
	} `yaml:"app"`
	MongoConfig struct {
		Addr     string `yaml:"addr"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"mongo"`
	JWTConfig struct {
		Secret         string        `yaml:"secret"`
		SecretBytes    []byte        `yaml:"-"`
		ExpireMinutes  int           `yaml:"expire_minutes"`
		ExpireDuration time.Duration `yaml:"-"`
	} `yaml:"jwt"`
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
	if err := C.Validate(); err != nil {
		log.Fatalf("[Config] Config validation error: %s", err.Error())
	}
	C.Compile()
	Debug = C.Debug
	log.Printf("[Config] Init done")
}

func (c IntRangeConfig) Validate() error {
	if c.Low >= c.High {
		return errors.New("low threshold is not less than high threshold")
	}
	return nil
}

func (c Config) Validate() error {
	if err := c.AppConfig.HeartRateThreshold.Validate(); err != nil {
		return errors.New("heart_rate_threshold: " + err.Error())
	}
	return nil
}

func (c Config) Compile() {
	C.AppConfig.OnlineTimeoutDuration = time.Second * time.Duration(C.AppConfig.OnlineTimeoutSeconds)
	C.JWTConfig.SecretBytes = []byte(C.JWTConfig.Secret)
	C.JWTConfig.ExpireDuration = time.Minute * time.Duration(C.JWTConfig.ExpireMinutes)
}
