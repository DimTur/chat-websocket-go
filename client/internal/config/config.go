package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/pelletier/go-toml"
)

type Config struct {
	AuthorizationToken string
}

var (
	configInstance *Config
	once           sync.Once
)

func LoadConfig(filepath string) (*Config, error) {
	var err error
	once.Do(func() {
		file, e := os.Open(filepath)
		if e != nil {
			err = fmt.Errorf("error opening configuration file: %w", e)
			return
		}
		defer file.Close()

		configInstance = &Config{}
		if e := toml.NewDecoder(file).Decode(configInstance); e != nil {
			err = fmt.Errorf("configuration decoding error: %w", e)
		}
	})
	return configInstance, err
}

func GetConfig() *Config {
	return configInstance
}
