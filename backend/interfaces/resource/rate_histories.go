package resource

import (
	"hamster/domain/model"
	"time"
)

type RateHistory struct {
	Id        uint64    `json:"id"`
	OrderType string    `json:"orderType"`
	Pair      string    `json:"pair"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewRateHistory(model *model.RateHistory) *RateHistory {
	return &RateHistory{
		Id:        model.Id,
		OrderType: model.OrderType,
		Pair:      model.Pair,
		Rate:      model.Rate,
		CreatedAt: model.CreatedAt,
	}
}
