package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	KafkaBrokers []string `yaml:"kafkaBrokers"`
	ServicePort  int      `yaml:"servicePort"`
}

func ParseConfig(configFilePath string, logger *zap.Logger) *Config {
	var config Config
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		logger.Error("Error in reading file ", zap.Error(err))
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Error("Error in unmarshiling config object", zap.Error(err))
	}
	return &config
}
