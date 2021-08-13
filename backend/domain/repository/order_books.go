package repository

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksRepository interface {
	FetchCurrent(refresh bool) (*model.OrderBooks, error)
	SubscribeLatest(listener func(*model.OrderBooks)) error

	GetHistories(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error)
	CreateHistories([]*model.OrderBooksHistory) error

	GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error)
	CreateMovingAverages([]*model.OrderBooksMovingAverage) error
}
