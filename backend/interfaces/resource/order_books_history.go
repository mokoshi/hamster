package resource

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksHistory struct {
	Id                 uint64    `json:"id"`
	Time               time.Time `json:"time"`
	LowestAskPrice     float64   `json:"lowestAskPrice"`
	LowestAskQuantity  float64   `json:"lowestAskQuantity"`
	HighestBidPrice    float64   `json:"highestBidPrice"`
	HighestBidQuantity float64   `json:"highestBidQuantity"`
}

func NewOrderBooksHistory(model *model.OrderBooksHistory) *OrderBooksHistory {
	return &OrderBooksHistory{
		Id:                 model.Id,
		Time:               model.Time,
		LowestAskPrice:     model.LowestAskPrice,
		LowestAskQuantity:  model.LowestAskQuantity,
		HighestBidPrice:    model.HighestBidPrice,
		HighestBidQuantity: model.HighestBidQuantity,
	}
}
