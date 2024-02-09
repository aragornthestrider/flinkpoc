package application

import (
	"context"
	"encoding/json"
	"time"
	"transactiongenerator/kafka"
	"transactiongenerator/models"

	"github.com/IBM/sarama"
	"github.com/bxcodec/faker/v3"
	"go.uber.org/zap"
)

type Application struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func (app *Application) Init(brokers []string, logger *zap.Logger) error {
	producer, err := kafka.NewProducer(brokers, logger)
	if err != nil {
		return err
	}
	app.producer = producer
	app.logger = logger
	return nil
}

func (app *Application) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			app.ProduceMsg()
		case <-ctx.Done():
			app.logger.Info("Context cancelled, closing producer")
			ticker.Stop()
			return
		}
	}
}

func (app *Application) ProduceMsg() error {
	var txn models.Transaction
	faker.FakeData(&txn)
	txnJSON, err := json.Marshal(txn)
	if err != nil {
		app.logger.Error("Error in marshalling", zap.Error(err))
		return err
	}
	msg := &sarama.ProducerMessage{Topic: "topic", Key: nil, Value: sarama.StringEncoder(txnJSON)}
	partition, offset, err := app.producer.SendMessage(msg)
	if err != nil {
		app.logger.Error("Error in producing msg", zap.Error(err))
		return err
	}
	app.logger.Info("Message produced", zap.Int32("Partition", partition), zap.Int64("Offset", offset))
	return nil
}
