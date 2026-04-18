package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/Arheon/markus-backend/configs"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/Arheon/markus-backend/pkg/outbox/broker/kafka"
	storeGorm "github.com/Arheon/markus-backend/pkg/outbox/store/gorm"
)

func InitDispatcher(env string) error {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)

	cfg, err := config.NewConfig(env)
	if err != nil {
		logger.Error("Unexpected config error", err)
		return err
	}

	dbConfig := cfg.Database
	db, err := gorm.Open(postgres.Open(dbConfig.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	db.AutoMigrate(
		&outbox.Message{},
	)

	if err != nil {
		return err
	}

	store := storeGorm.NewStore(db)
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Net.MaxOpenRequests = 1

	logger.Debug("Try to connect to kafka...")
	producer, err := sarama.NewSyncProducer([]string{cfg.Broker.KafkaDSN}, saramaConfig)
	if err != nil {
		return err
	}

	broker := kafka.NewBroker(producer, kafka.WithTopics(map[string]string{
		"event.server.create_server":   cfg.Broker.BrokerTopics.Server,
		"event.server.create_room":     cfg.Broker.BrokerTopics.Server,
		"event.message.create_message": cfg.Broker.BrokerTopics.Message,
	}))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	dispatcher := outbox.NewDispatcher(
		store,
		broker,
		outbox.WithInterval(5*time.Second),
		outbox.WithReadBatchSize(200),
		outbox.WithDeleteBatchSize(50),
		outbox.WithMaxAttempts(300),
		outbox.WithContext(ctx),
		outbox.WithLogger(logger),
		outbox.WithLimit(10),
	)

	dispatcher.Start()
	logger.Debug("Dispatch started...")

	go func() {
		for err := range dispatcher.GetErrorChan() {
			logger.Errorf("Error received: %v\n", err)
		}
	}()

	go func() {
		for msg := range dispatcher.GetDiscardedMesgChan() {
			logger.WithFields(logrus.Fields{
				"message": msg,
			}).Debugf("Message received: %s", msg.ID)
		}
	}()

	<-ctx.Done()
	logger.Debug("Dispatcher shutdown...")
	return nil
}
