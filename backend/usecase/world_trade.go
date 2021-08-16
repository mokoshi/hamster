package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type WorldTradeUsecase interface {
	StartSync() error
	GetWorldTrades(from time.Time, to time.Time) ([]*model.WorldTrade, error)
	GetMovingAverages(from time.Time, to time.Time) ([]*model.WorldTradeMovingAverage, error)
}

type worldTradeUsecase struct {
	worldTradeRepository repository.WorldTradeRepository
}

func NewWorldTradeUsecase(
	worldTradeRepository repository.WorldTradeRepository,
) WorldTradeUsecase {
	return &worldTradeUsecase{
		worldTradeRepository: worldTradeRepository,
	}
}

func (u *worldTradeUsecase) StartSync() error {
	return u.worldTradeRepository.SubscribeAndSync(func(worldTrade *model.WorldTrade) {
		//clog.Logger.Debug(worldTrade)
	})
}

func (u *worldTradeUsecase) GetWorldTrades(from time.Time, to time.Time) ([]*model.WorldTrade, error) {
	histories, err := u.worldTradeRepository.GetWorldTrades(from, to)
	return histories, err
}

func (u *worldTradeUsecase) GetMovingAverages(from time.Time, to time.Time) ([]*model.WorldTradeMovingAverage, error) {
	histories, err := u.worldTradeRepository.GetMovingAverages(from, to)
	return histories, err
}
