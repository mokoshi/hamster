package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type ExchangeUsecase interface {
	GetOpenOrders() ([]*model.OpenOrder, error)

	GetRateHistories(from time.Time, to time.Time) ([]*model.RateHistory, error)
	SyncCurrentRate(orderType string, pair string) (*model.RateHistory, error)
}

type ExchangeUsecaseImpl struct {
	OrderRepository       repository.OrderRepository
	RateHistoryRepository repository.RateHistoryRepository
}

func NewExchangeUsecase(
	orderRepository repository.OrderRepository,
	rateHistoryRepository repository.RateHistoryRepository,
) ExchangeUsecase {
	return &ExchangeUsecaseImpl{
		OrderRepository:       orderRepository,
		RateHistoryRepository: rateHistoryRepository,
	}
}

func (u *ExchangeUsecaseImpl) GetOpenOrders() ([]*model.OpenOrder, error) {
	return u.OrderRepository.SyncOpenOrders()
}

func (u *ExchangeUsecaseImpl) GetRateHistories(from time.Time, to time.Time) ([]*model.RateHistory, error) {
	return u.RateHistoryRepository.GetRateHistories(from, to)
}

func (u *ExchangeUsecaseImpl) SyncCurrentRate(orderType string, pair string) (*model.RateHistory, error) {
	return u.RateHistoryRepository.SyncCurrentRate(orderType, pair)
}
