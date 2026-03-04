package health

import (
	"errors"

	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

type HealthProbe struct {
	injector     do.Injector
	strategyList []Strategy
}

func NewHealthProbe(injector do.Injector) *HealthProbe {
	return &HealthProbe{
		injector: injector,
	}
}

func (hp *HealthProbe) AddStrategy(strategy Strategy) {
	strategy.SetContainer(hp.injector)
	hp.strategyList = append(hp.strategyList, strategy)
}

func (hp *HealthProbe) ExecuteStrategy(probe StrategyType) (*StrategyResult, error) {
	logger := do.MustInvoke[*logrus.Logger](hp.injector)

	for _, strategy := range hp.strategyList {
		if !strategy.Support(probe) {
			continue
		}

		logger.Debug("Try to execute prope strategy")
		result, err := strategy.Execute()
		if err != nil {
			logger.WithError(err).Error("Unexpected error in strategy")
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("Unexpected strategy type")
}
