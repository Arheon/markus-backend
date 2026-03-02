package strategy

import (
	"errors"

	"github.com/Arheon/markus-backend/internal/shared/domain/service/health"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type Liveness struct {
	container *do.Injector
}

func (l *Liveness) SetContainer(container *do.Injector) {
	l.container = container
}

func (l *Liveness) Execute() (*health.StrategyResult, error) {
	databaseLiveness := true

	db, err := do.InvokeAs[*gorm.DB](*l.container)
	if err != nil {
		databaseLiveness = false
	}

	if err := db.Error; err != nil {
		databaseLiveness = false
	}

	payload := map[string]map[string]bool{
		"status": {
			"database": databaseLiveness,
		},
	}

	result := &health.StrategyResult{
		Type:    health.Liveness,
		Payload: payload,
	}

	var livenessError error
	if databaseLiveness == false {
		livenessError = errors.New("Not liveness")
	}

	return result, livenessError
}

func (l *Liveness) Support(probe health.StrategyType) bool {
	return probe == health.Liveness
}
