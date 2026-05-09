package scheduler

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger    LoggerConf    `yaml:"logger"`
	Storage   StorageConf   `yaml:"storage"`
	RabbitMQ  RabbitMQConf  `yaml:"rabbitmq"`
	Scheduler SchedulerConf `yaml:"scheduler"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Type          string `yaml:"type"`
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrationsDir"`
}

type RabbitMQConf struct {
	URL   string `yaml:"url"`
	Queue string `yaml:"queue"`
}

type SchedulerConf struct {
	ScanInterval  time.Duration `yaml:"scan_interval"`
	CleanInterval time.Duration `yaml:"clean_interval"`
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
