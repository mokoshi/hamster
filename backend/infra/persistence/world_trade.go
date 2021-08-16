package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/clog"
	"hamster/lib/util"
	"time"
)

type WorldTradeRepository struct {
	client          *cc_client.Client
	db              *gorm.DB
	worldTradeCache *model.WorldTradeCache
	flushSize       int
}

func NewWorldTradeRepository(
	client *cc_client.Client,
	db *gorm.DB,
	worldTradeCache *model.WorldTradeCache,
	flushSize *int,
) repository.WorldTradeRepository {
	var fs int
	if flushSize == nil {
		fs = 10
	} else {
		fs = *flushSize
	}

	return &WorldTradeRepository{
		client:          client,
		db:              db,
		worldTradeCache: worldTradeCache,
		flushSize:       fs,
	}
}

func (r *WorldTradeRepository) SubscribeAndSync(listener func(trade *model.WorldTrade)) error {
	err := r.client.SubscribeTrades(func(trade *cc_client.Trade) {
		worldTrade := r.parseTrade(trade)

		r.worldTradeCache.AddLatestTrade(worldTrade)

		if r.worldTradeCache.BufferedSize() > r.flushSize {
			trades, averages := r.worldTradeCache.ReadNext()
			if err := r.db.Create(trades).Error; err != nil {
				clog.Logger.Error(err)
			}
			if err := r.db.Create(averages).Error; err != nil {
				clog.Logger.Error(err)
			}
		}

		listener(worldTrade)
	})

	return err
}

func (r *WorldTradeRepository) GetWorldTrades(from time.Time, to time.Time) ([]*model.WorldTrade, error) {
	var trades []*model.WorldTrade
	r.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&trades)
	return trades, nil
}

func (r *WorldTradeRepository) GetMovingAverages(from time.Time, to time.Time) ([]*model.WorldTradeMovingAverage, error) {
	var averages []*model.WorldTradeMovingAverage
	r.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&averages)
	return averages, nil
}

func (r *WorldTradeRepository) GetLatestMovingAverage(offsetFromLatest int) *model.WorldTradeMovingAverage {
	return r.worldTradeCache.GetLatestMovingAverage(offsetFromLatest)
}

func (r *WorldTradeRepository) parseTrade(res *cc_client.Trade) *model.WorldTrade {
	return &model.WorldTrade{
		Id:        res.Id,
		Time:      time.Now(),
		Pair:      res.Pair,
		OrderType: res.OrderType,
		Rate:      util.ParseFloat64(res.Rate),
		Amount:    util.ParseFloat64(res.Amount),
	}
}
