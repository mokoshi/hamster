package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type TradingHistoryUsecase interface {
	Create(tradedAt time.Time) (*model.TradingHistory, error)
}

type TradingHistoryUsecaseImpl struct {
	TradingHistoryRepository repository.TradingHistoryRepository
}

func NewTradingHistoryUsecase(tradingHistoryRepository repository.TradingHistoryRepository) TradingHistoryUsecase {
	return &TradingHistoryUsecaseImpl{TradingHistoryRepository: tradingHistoryRepository}
}

func (u *TradingHistoryUsecaseImpl) Create(tradedAt time.Time) (*model.TradingHistory, error) {
	return u.TradingHistoryRepository.Create(&model.TradingHistory{TradedAt: time.Now()})
}
