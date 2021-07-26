package repository

import (
	"hamster/domain/model"
)

type OrderBooksRepository interface {
	Get(refresh bool) (*model.OrderBooks, error)
	Subscribe(listener func(*model.OrderBooks)) error
}
