package kafka

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func KafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	return config
}

func NewProducer(brokers []string, logger *zap.Logger) (sarama.SyncProducer, error) {
	config := KafkaConfig()
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Error("Error in creating producer", zap.Error(err))
		return nil, err
	}
	return producer, nil
}
