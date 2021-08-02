package repository

import (
	"hamster/domain/model"
)

type OrderBooksHistoryRepository interface {
	AddToBuffer(orderBooksHistory *model.OrderBooksHistory) (int, error)
	GetBufferingSize() int
	Flush() error
}
