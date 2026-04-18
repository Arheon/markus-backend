package dependencyinjection

import (
	"github.com/Arheon/markus-backend/internal/message/domain/repository"
	"github.com/Arheon/markus-backend/internal/message/infrastructure/repository/message"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func InitModule(injector do.Injector) {
	do.Provide[repository.MessageRepository](injector, func(i do.Injector) (repository.MessageRepository, error) {
		db, err := do.InvokeAs[*gorm.DB](i)
		if err != nil {
			return nil, err
		}

		return message.NewRepository(db), nil
	})
}
