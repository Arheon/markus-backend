package health

import (
	"github.com/Arheon/markus-backend/internal/shared/domain/service/health"
	"github.com/Arheon/markus-backend/internal/shared/domain/service/health/strategy"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

type HealthQuery struct {
	injector do.Injector
}

func NewQuery(injector do.Injector) *HealthQuery {
	return &HealthQuery{
		injector: injector,
	}
}

func (h *HealthQuery) Handle(probeName string) (*health.StrategyResult, error) {
	logger := do.MustInvokeAs[*logrus.Logger](h.injector)

	healthprobe := health.NewHealthProbe(h.injector)
	healthprobe.AddStrategy(&strategy.Liveness{})
	healthprobe.AddStrategy(&strategy.Readiness{})

	var result *health.StrategyResult
	var err error
	switch probeName {
	case "liveness":
		result, err = healthprobe.ExecuteStrategy(health.Liveness)
	case "readiness":
		result, err = healthprobe.ExecuteStrategy(health.Rediness)
	}

	if err != nil {
		logger.WithError(err).Error("Undefined error in health strategy")
		return nil, err
	}

	return result, err
}
