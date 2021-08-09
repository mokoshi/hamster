package resource

import (
	"hamster/domain/model"
	"time"
)

type Order struct {
	Id                     uint64    `json:"id"`
	OrderType              string    `json:"orderType"`
	Rate                   float64   `json:"rate"`
	Pair                   string    `json:"pair"`
	PendingAmount          float64   `json:"pendingAmount"`
	PendingMarketBuyAmount float64   `json:"pendingMarketBuyAmount"`
	StopLossRate           float64   `json:"stopLossRate"`
	CreatedAt              time.Time `json:"createdAt"`
}

func NewOrder(model *model.Order) *Order {
	return &Order{
		Id:                     model.Id,
		OrderType:              model.OrderType,
		Rate:                   model.Rate,
		Pair:                   model.Pair,
		PendingAmount:          model.PendingAmount,
		PendingMarketBuyAmount: model.PendingMarketBuyAmount,
		StopLossRate:           model.StopLossRate,
		CreatedAt:              model.CreatedAt,
	}
}
