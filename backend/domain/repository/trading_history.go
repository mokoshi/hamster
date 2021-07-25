package repository

import (
	"hamster/domain/model"
)

type TradingHistoryRepository interface {
	Create(tradingHistory *model.TradingHistory) (*model.TradingHistory, error)
}
