package calendar

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	HTTP    HTTPConf    `yaml:"http"`
	GRPC    GRPCConf    `yaml:"grpc"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Type          string `yaml:"type"`
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrationsDir"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func NewConfig(path string) (Config, error) {
	config := Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
