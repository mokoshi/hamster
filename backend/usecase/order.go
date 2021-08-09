package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
)

type OrderUsecase interface {
	GetOpenOrders() ([]*model.Order, error)
}

type OrderUsecaseImpl struct {
	OrderRepository repository.OrderRepository
}

func NewOrderUsecase(myOrderRepository repository.OrderRepository) OrderUsecase {
	return &OrderUsecaseImpl{OrderRepository: myOrderRepository}
}

func (u *OrderUsecaseImpl) GetOpenOrders() ([]*model.Order, error) {
	return u.OrderRepository.GetOpenOrders()
}
