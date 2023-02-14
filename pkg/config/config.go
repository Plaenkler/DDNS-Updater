package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	pathToConfig = "./data/config.yaml"
)

var (
	configOnce sync.Once
	instance   *Config
)

type Config struct {
	Port uint `yaml:"Port"`
}

func GetConfig() *Config {
	configOnce.Do(func() {
		err := initConfig()
		if err != nil {
			log.Fatalf("[get-config-1] initialization failed - error: %s", err.Error())
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
		Port: 80,
	}
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(pathToConfig), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, 0644)
	if err != nil {
		return err
	}
	log.Println("[create-config-1] created default configuration exiting...")
	os.Exit(0)
	return nil
}
