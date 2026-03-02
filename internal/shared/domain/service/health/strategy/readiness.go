package strategy

import (
	"github.com/Arheon/markus-backend/internal/shared/domain/service/health"
	"github.com/samber/do/v2"
)

type Readiness struct {
	injector *do.Injector
}

func (r *Readiness) SetContainer(injector *do.Injector) {
	r.injector = injector
}

func (r *Readiness) Execute() (*health.StrategyResult, error) {
	payload := map[string]bool{
		"readiness": true,
	}

	return &health.StrategyResult{
		Type:    health.Rediness,
		Payload: payload,
	}, nil
}

func (r *Readiness) Support(probe health.StrategyType) bool {
	return probe == health.Rediness
}
