package repository

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksHistoryRepository interface {
	Get(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error)
	AddToBuffer(orderBooksHistory *model.OrderBooksHistory) (int, error)
	GetBufferingSize() int
	Flush() error
}
