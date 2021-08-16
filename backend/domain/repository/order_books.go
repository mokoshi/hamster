package repository

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksRepository interface {
	SyncOrderBooks() (*model.OrderBooks, error)
	// TODO StartSync() ?
	SubscribeAndSync(listener func(*model.OrderBooks)) error

	GetSnapshots(from time.Time, to time.Time) ([]*model.OrderBooksSnapshot, error)
	GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error)

	GetLatestSnapshot(offsetFromLatest int) *model.OrderBooksSnapshot
	GetLatestMovingAverage(offsetFromLatest int) *model.OrderBooksMovingAverage
}
