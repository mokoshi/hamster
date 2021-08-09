package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
)

type ExchangeUsecase interface {
	GetOpenOrders() ([]*model.Order, error)
}

type ExchangeUsecaseImpl struct {
	OrderRepository repository.OrderRepository
}

func NewExchangeUsecase(myOrderRepository repository.OrderRepository) ExchangeUsecase {
	return &ExchangeUsecaseImpl{OrderRepository: myOrderRepository}
}

func (u *ExchangeUsecaseImpl) GetOpenOrders() ([]*model.Order, error) {
	return u.OrderRepository.GetOpenOrders()
}
