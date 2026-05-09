package sender

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf   `yaml:"logger"`
	RabbitMQ RabbitMQConf `yaml:"rabbitmq"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type RabbitMQConf struct {
	URL   string `yaml:"url"`
	Queue string `yaml:"queue"`
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
