package resource

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksSnapshot struct {
	Id                 uint64    `json:"id"`
	Time               time.Time `json:"time"`
	LowestAskPrice     float64   `json:"lowestAskPrice"`
	LowestAskQuantity  float64   `json:"lowestAskQuantity"`
	HighestBidPrice    float64   `json:"highestBidPrice"`
	HighestBidQuantity float64   `json:"highestBidQuantity"`
}

func NewOrderBooksHistory(model *model.OrderBooksSnapshot) *OrderBooksSnapshot {
	return &OrderBooksSnapshot{
		Id:                 model.Id,
		Time:               model.Time,
		LowestAskPrice:     model.LowestAskPrice,
		LowestAskQuantity:  model.LowestAskQuantity,
		HighestBidPrice:    model.HighestBidPrice,
		HighestBidQuantity: model.HighestBidQuantity,
	}
}
