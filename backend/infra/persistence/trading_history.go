package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
)

type TradingHistoryPersistence struct {
	Db *gorm.DB
}

func NewTradingHistoryPersistence(db *gorm.DB) repository.TradingHistoryRepository {
	return &TradingHistoryPersistence{Db: db}
}

func (thp TradingHistoryPersistence) Create(tradingHistory *model.TradingHistory) (*model.TradingHistory, error) {
	if err := thp.Db.Create(&tradingHistory).Error; err != nil {
		return nil, err
	}
	return tradingHistory, nil
}
