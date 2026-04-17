package internal

import (
	"log"
	"os"
	"time"

	config "github.com/Arheon/markus-backend/configs"
	messageDomainEntity "github.com/Arheon/markus-backend/internal/message/domain/entity"
	messageDI "github.com/Arheon/markus-backend/internal/message/infrastructure/dependency_injection"
	serverDomainEntity "github.com/Arheon/markus-backend/internal/server/domain/entity"
	serverDI "github.com/Arheon/markus-backend/internal/server/infrastructure/dependency_injection"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/servers/http"
	userDomainEntity "github.com/Arheon/markus-backend/internal/user/domain/entity"
	userDI "github.com/Arheon/markus-backend/internal/user/infrastructure/dependency_injection"
	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/IBM/sarama"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	authDomainEntity "github.com/Arheon/markus-backend/internal/auth/domain/entity"
	authDI "github.com/Arheon/markus-backend/internal/auth/infrastructure/dependency_injection"
	"github.com/Arheon/markus-backend/pkg/outbox/broker/kafka"
	storeGorm "github.com/Arheon/markus-backend/pkg/outbox/store/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitApp(env string) error {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)

	injector := do.New()
	cfg, err := config.NewConfig(env)
	if err != nil {
		logger.Error("Undefined config error")
		return err
	}

	dbConf := cfg.Database

	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		newLogger := gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             time.Second,     // Медленный SQL (1 сек)
				LogLevel:                  gormLogger.Info, // Уровень логов (Info логирует все)
				IgnoreRecordNotFoundError: false,           // Не игнорировать ошибки "запись не найдена"
				Colorful:                  true,            // Цветная печать
			},
		)
		db, err := gorm.Open(postgres.Open(dbConf.DSN), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
		})

		db.AutoMigrate(
			&authDomainEntity.User{},
			&userDomainEntity.User{},
			&serverDomainEntity.Member{},
			&messageDomainEntity.Message{},
			&serverDomainEntity.Room{},
			&serverDomainEntity.RoomCategory{},
			&serverDomainEntity.Server{},
		)

		if nil != err {
			log.Fatalln("Connection to database failed", err)
		} else {
			log.Println("Connected to database successfully", nil)
		}

		return db, err
	})

	do.Provide(injector, func(i do.Injector) (*logrus.Logger, error) {
		return logger, nil
	})

	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return cfg, nil
	})

	do.Provide(injector, func(i do.Injector) (*outbox.Publisher, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		store := storeGorm.NewStore(db)
		publisher := outbox.NewPublisher(store)

		return publisher, nil
	})

	do.Provide(injector, func(i do.Injector) (*outbox.Dispatcher, error) {
		config, err := do.InvokeAs[*config.Config](i)
		if err != nil {
			return nil, err
		}

		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		store := storeGorm.NewStore(db)
		saramaConfig := sarama.Config{}
		saramaConfig.Producer.Return.Successes = true

		producer, err := sarama.NewSyncProducer([]string{config.Broker.KafkaDSN}, &saramaConfig)
		if err != nil {
			return nil, err
		}

		broker := kafka.NewBroker(producer, kafka.WithTopics(map[string]string{
			"ServerCreate":         config.Broker.BrokerTopics.Server,
			"ServerInfoCartUpdate": config.Broker.BrokerTopics.Server,
		}))

		dispatcher := outbox.NewDispatcher(
			store,
			broker,
			outbox.WithInterval(5*time.Second),
			outbox.WithReadBatchSize(200),
			outbox.WithDeleteBatchSize(50),
			outbox.WithMaxAttempts(300),
		)

		return dispatcher, nil
	})

	authDI.InitModule(injector)
	messageDI.InitModule(injector)
	serverDI.InitModule(injector)
	userDI.InitModule(injector)

	http.Init(env, injector)

	return nil
}
