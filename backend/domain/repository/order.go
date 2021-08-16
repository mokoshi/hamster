package repository

import (
	"hamster/domain/model"
)

type OrderRepository interface {
	SyncOpenOrders() ([]*model.OpenOrder, error)
	GetOpenOrderCount() int
	RequestOrder(request *model.OrderRequest) (*model.OrderRequest, error)
}
