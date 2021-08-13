package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type OrderBooksUsecase interface {
	GetHistories(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error)
	GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error)
	WriteWithBuffering(lowestAsk *model.OrderBookItem, highestBid *model.OrderBookItem, flushSize int) error
}

type orderBooksUsecase struct {
	buffer               *model.OrderBooksBuffer
	orderBooksRepository repository.OrderBooksRepository
}

func NewOrderBooksUsecase(
	orderBooksHistoryRepository repository.OrderBooksRepository,
) OrderBooksUsecase {
	return &orderBooksUsecase{
		buffer:               &model.OrderBooksBuffer{},
		orderBooksRepository: orderBooksHistoryRepository,
	}
}

func (u *orderBooksUsecase) GetHistories(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error) {
	histories, err := u.orderBooksRepository.GetHistories(from, to)
	return histories, err
}

func (u *orderBooksUsecase) GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error) {
	histories, err := u.orderBooksRepository.GetMovingAverages(from, to)
	return histories, err
}

func (u *orderBooksUsecase) WriteWithBuffering(
	lowestAsk *model.OrderBookItem,
	highestBid *model.OrderBookItem,
	flushSize int,
) error {
	u.buffer.Add(&model.OrderBooksHistory{
		Time:               time.Now(),
		LowestAskPrice:     lowestAsk.Price,
		LowestAskQuantity:  lowestAsk.Quantity,
		HighestBidPrice:    highestBid.Price,
		HighestBidQuantity: highestBid.Quantity,
	})

	if u.buffer.BufferedSize() > flushSize {
		histories, averages := u.buffer.Read()
		if err := u.orderBooksRepository.CreateHistories(histories); err != nil {
			return err
		}
		if err := u.orderBooksRepository.CreateMovingAverages(averages); err != nil {
			return err
		}
	}

	return nil
}
