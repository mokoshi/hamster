package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type OrderBooksUsecase interface {
	StartSync(listener func(orderBooks *model.OrderBooks)) error
	GetSnapshots(from time.Time, to time.Time) ([]*model.OrderBooksSnapshot, error)
	GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error)
}

type orderBooksUsecase struct {
	orderBooksRepository repository.OrderBooksRepository
}

func NewOrderBooksUsecase(
	orderBooksRepository repository.OrderBooksRepository,
) OrderBooksUsecase {
	return &orderBooksUsecase{
		orderBooksRepository: orderBooksRepository,
	}
}

func (u *orderBooksUsecase) StartSync(listener func(orderBooks *model.OrderBooks)) error {
	return u.orderBooksRepository.SubscribeAndSync(listener)
}

func (u *orderBooksUsecase) GetSnapshots(from time.Time, to time.Time) ([]*model.OrderBooksSnapshot, error) {
	histories, err := u.orderBooksRepository.GetSnapshots(from, to)
	return histories, err
}

func (u *orderBooksUsecase) GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error) {
	histories, err := u.orderBooksRepository.GetMovingAverages(from, to)
	return histories, err
}
