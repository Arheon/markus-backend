package health

import "github.com/samber/do/v2"

type StrategyType int

const (
	Liveness StrategyType = iota
	Rediness
)

type StrategyResult struct {
	Type    StrategyType `json:"type"`
	Payload any          `json:"payload"`
}

type Strategy interface {
	SetContainer(container do.Injector)
	Execute() (*StrategyResult, error)
	Support(probe StrategyType) bool
}
