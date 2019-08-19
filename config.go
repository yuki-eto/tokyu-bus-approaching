package tokyu_bus_approaching

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	URLs        []string `yaml:"urls"`
	RefreshTime int      `yaml:"refresh_time"`
}

var config *Config

func LoadConfig(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return config
}

func init() {
	config = &Config{
		URLs:        []string{},
		RefreshTime: 60,
	}
}
