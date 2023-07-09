package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     uint64 `yaml:"Port"`
	Interval uint64 `yaml:"Interval"`
}

const (
	pathToConfig   = "./data/config.yaml"
	configDirPerm  = 0755
	configFilePerm = 0644
)

var config *Config

func GetConfig() *Config {
	if config == nil {
		err := loadConfig()
		if err != nil {
			log.Fatalf("[config-GetConfig-1] initialization failed - error: %s", err.Error())
		}
	}
	return config
}

func loadConfig() error {
	_, err := os.Stat(pathToConfig)
	if os.IsNotExist(err) {
		err = createConfig()
		if err != nil {
			return err
		}
	}
	file, err := os.Open(pathToConfig)
	if err != nil {
		return err
	}
	defer file.Close()
	instance := &Config{}
	err = yaml.NewDecoder(file).Decode(instance)
	if err != nil {
		return err
	}
	config = instance
	err = loadConfigFromEnv()
	if err != nil {
		return err
	}
	return nil
}

func createConfig() error {
	config := Config{
		Port:     80,
		Interval: 600,
	}
	err := os.MkdirAll(filepath.Dir(pathToConfig), configDirPerm)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, configFilePerm)
	if err != nil {
		return err
	}
	log.Println("[config-createConfig-1] created default configuration")
	return nil
}

func loadConfigFromEnv() error {
	port, err := parseUintEnv("APP_PORT")
	if err != nil {
		return err
	}
	if port != 0 {
		config.Port = port
	}
	interval, err := parseUintEnv("APP_INTERVAL")
	if err == nil {
		return err
	}
	if interval != 0 {
		config.Interval = interval
	}
	return nil
}

func parseUintEnv(envName string) (uint64, error) {
	valueStr := os.Getenv(envName)
	if valueStr == "" {
		return 0, nil
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func UpdateConfig(updatedConfig *Config) error {
	data, err := yaml.Marshal(updatedConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, configFilePerm)
	if err != nil {
		return err
	}
	config = updatedConfig
	return nil
}
