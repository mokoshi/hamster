package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/clog"
	"time"
)

const FlushBufferSize = 100

type OrderBooksHistoryUsecase interface {
	WriteWithBuffering(lowestAsk *model.Order, highestBid *model.Order) (*model.OrderBooksHistory, error)
}

type OrderBooksHistoryUsecaseImpl struct {
	OrderBooksHistoryRepository repository.OrderBooksHistoryRepository
}

func NewOrderBooksHistoryUsecase(orderBooksHistoryRepository repository.OrderBooksHistoryRepository) OrderBooksHistoryUsecase {
	return &OrderBooksHistoryUsecaseImpl{OrderBooksHistoryRepository: orderBooksHistoryRepository}
}

func (u *OrderBooksHistoryUsecaseImpl) WriteWithBuffering(lowestAsk *model.Order, highestBid *model.Order) (*model.OrderBooksHistory, error) {
	orderBooksHistory := &model.OrderBooksHistory{
		Time:               time.Now(),
		LowestAskPrice:     lowestAsk.Price,
		LowestAskQuantity:  lowestAsk.Quantity,
		HighestBidPrice:    highestBid.Price,
		HighestBidQuantity: highestBid.Quantity,
	}
	size, err := u.OrderBooksHistoryRepository.AddToBuffer(orderBooksHistory)
	if err != nil {
		return nil, err
	}

	if size > FlushBufferSize {
		go func() {
			clog.Logger.Debug("Flush orderBooksHistory")
			u.OrderBooksHistoryRepository.Flush()
		}()
	}

	return orderBooksHistory, err
}
