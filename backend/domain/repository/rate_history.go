package repository

import (
	"hamster/domain/model"
	"time"
)

type RateHistoryRepository interface {
	GetRateHistories(from time.Time, to time.Time) ([]*model.RateHistory, error)
	SyncCurrentRate(orderType string, pair string) (*model.RateHistory, error)
}
