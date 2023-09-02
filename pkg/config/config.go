package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/plaenkler/ddns-updater/pkg/cipher"
	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval uint64 `yaml:"Interval"`
	UseTOTP  bool   `yaml:"TOTP"`
	Port     uint64 `yaml:"Port"`
	Resolver string `yaml:"Resolver"`
	Cryptor  string `yaml:"Cryptor"`
}

const (
	pathToConfig = "./data/config.yaml"
	dirPerm      = 0755
	filePerm     = 0644
)

var (
	config *Config
	mutex  = &sync.RWMutex{}
)

func init() {
	err := load()
	if err != nil {
		log.Fatalf("[config-init-1] initialization failed: %s", err.Error())
	}
}

func load() error {
	_, err := os.Stat(pathToConfig)
	if os.IsNotExist(err) {
		err = create()
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
	err = loadFromEnv()
	if err != nil {
		return err
	}
	return nil
}

func create() error {
	Cryptor, err := cipher.GenerateRandomKey(16)
	if err != nil {
		return err
	}
	config := Config{
		Interval: 600,
		UseTOTP:  false,
		Port:     80,
		Resolver: "",
		Cryptor:  Cryptor,
	}
	err = os.MkdirAll(filepath.Dir(pathToConfig), dirPerm)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, filePerm)
	if err != nil {
		return err
	}
	log.Infof("[config-create-1] created default configuration")
	return nil
}

func loadFromEnv() error {
	interval, err := parseUintEnv("DDNS_INTERVAL")
	if err != nil {
		return err
	}
	if interval != 0 {
		config.Interval = interval
	}
	useTOTP, err := parseBoolEnv("DDNS_TOTP")
	if err == nil {
		config.UseTOTP = useTOTP
	}
	if err != nil && err.Error() != "not set" {
		return err
	}
	port, err := parseUintEnv("DDNS_PORT")
	if err != nil {
		return err
	}
	if port != 0 {
		config.Port = port
	}
	resolver, err := parseURLEnv("DDNS_RESOLVER")
	if err != nil {
		return err
	}
	if resolver != "" {
		config.Resolver = resolver
	}
	Cryptor, ok := os.LookupEnv("DDNS_Cryptor")
	if ok && Cryptor != "" {
		config.Cryptor = Cryptor
	}
	return nil
}

func parseUintEnv(envName string) (uint64, error) {
	valueStr, ok := os.LookupEnv(envName)
	if !ok {
		return 0, nil
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func parseBoolEnv(envName string) (bool, error) {
	valueStr, ok := os.LookupEnv(envName)
	if !ok {
		return false, fmt.Errorf("not set")
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}
	return value, nil
}

func parseURLEnv(envName string) (string, error) {
	value, ok := os.LookupEnv(envName)
	if !ok {
		return "", nil
	}
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func Update(updatedConfig *Config) error {
	data, err := yaml.Marshal(updatedConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(pathToConfig, data, filePerm)
	if err != nil {
		return err
	}
	mutex.Lock()
	defer mutex.Unlock()
	config = updatedConfig
	return nil
}

func Get() *Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}
