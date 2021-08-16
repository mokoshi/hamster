package repository

import (
	"hamster/domain/model"
	"time"
)

type WorldTradeRepository interface {
	SubscribeAndSync(listener func(average *model.WorldTrade)) error

	GetWorldTrades(from time.Time, to time.Time) ([]*model.WorldTrade, error)
	GetMovingAverages(from time.Time, to time.Time) ([]*model.WorldTradeMovingAverage, error)

	GetLatestMovingAverage(offsetFromLatest int) *model.WorldTradeMovingAverage
}
