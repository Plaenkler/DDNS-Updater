package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

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

var (
	once     sync.Once
	instance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		err := initConfig()
		if err != nil {
			log.Fatalf("[config-GetConfig-1] initialization failed - error: %s", err.Error())
		}
	})
	return instance
}

func initConfig() error {
	instance = &Config{}
	if _, err := os.Stat(pathToConfig); err != nil {
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
	err = yaml.NewDecoder(file).Decode(&instance)
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
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(pathToConfig), configDirPerm)
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

func UpdateConfig(config *Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, configFilePerm)
	if err != nil {
		return err
	}
	instance = config
	return nil
}
