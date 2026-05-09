package sender

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf   `yaml:"logger"`
	RabbitMQ RabbitMQConf `yaml:"rabbitmq"`
	HTTP     HTTPConf     `yaml:"http"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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

	expandedData := os.ExpandEnv(string(data))

	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
