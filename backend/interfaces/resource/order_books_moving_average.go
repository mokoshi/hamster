package resource

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksMovingAverage struct {
	Id       uint64        `json:"id"`
	Time     time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
	AskPrice float64       `json:"askPrice"`
	BidPrice float64       `json:"bidPrice"`
}

func NewOrderBooksMovingAverage(model *model.OrderBooksMovingAverage) *OrderBooksMovingAverage {
	return &OrderBooksMovingAverage{
		Id:       model.Id,
		Time:     model.Time,
		Duration: model.Duration,
		AskPrice: model.AskPrice,
		BidPrice: model.BidPrice,
	}
}
