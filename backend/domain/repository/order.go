package repository

import (
	"hamster/domain/model"
)

type OrderRepository interface {
	GetOpenOrders() ([]*model.Order, error)
}
