package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/util"
	"time"
)

type RateHistoryPersistence struct {
	Client *cc_client.Client
	Db     *gorm.DB
}

func NewRateHistoryPersistence(db *gorm.DB, client *cc_client.Client) repository.RateHistoryRepository {
	return &RateHistoryPersistence{Db: db, Client: client}
}

func (p *RateHistoryPersistence) GetRateHistories(from time.Time, to time.Time) ([]*model.RateHistory, error) {
	var histories []*model.RateHistory
	p.Db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&histories)
	return histories, nil
}

func (p *RateHistoryPersistence) SyncCurrentRate(orderType string, pair string) (*model.RateHistory, error) {
	res, err := p.Client.GetRate(orderType, pair)
	if err != nil {
		return nil, err
	}

	rate := &model.RateHistory{
		Pair:      pair,
		OrderType: orderType,
		Rate:      util.ParseFloat64(res.Rate),
		CreatedAt: time.Now(),
	}

	if err := p.Db.Create(&rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}
