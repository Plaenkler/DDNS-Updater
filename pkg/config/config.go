package config

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	log "github.com/plaenkler/ddns-updater/pkg/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval uint32 `yaml:"Interval"`
	UseTOTP  bool   `yaml:"TOTP"`
	Port     uint16 `yaml:"Port"`
	Resolver string `yaml:"Resolver"`
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
		log.Fatalf("initialization failed: %s", err.Error())
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
	defer log.ErrorClose(file)
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
	config := Config{
		Interval: 600,
		UseTOTP:  false,
		Port:     80,
		Resolver: "",
	}
	err := os.MkdirAll(filepath.Dir(pathToConfig), dirPerm)
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
	log.Infof("created default configuration")
	return nil
}

func loadFromEnv() error {
	interval, err := parseUintEnv("DDNS_INTERVAL")
	if err != nil {
		return err
	}
	if interval > uint64(math.MaxUint32) {
		return fmt.Errorf("interval value exceeds uint32 maximum")
	}
	if interval != 0 {
		config.Interval = uint32(interval)
	}
	useTOTP, err := parseBoolEnv("DDNS_TOTP")
	if err != nil && err.Error() != "not set" {
		return err
	}
	if err == nil {
		config.UseTOTP = useTOTP
	}
	port, err := parseUintEnv("DDNS_PORT")
	if err != nil {
		return err
	}
	if port > uint64(math.MaxUint16) {
		return fmt.Errorf("port value exceeds uint16 maximum")
	}
	if port != 0 {
		config.Port = uint16(port)
	}
	resolver, err := parseURLEnv("DDNS_RESOLVER")
	if err != nil {
		return err
	}
	if resolver != "" {
		config.Resolver = resolver
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
