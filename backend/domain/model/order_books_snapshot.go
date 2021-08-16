package model

import "time"

type OrderBooksSnapshot struct {
	Id                 uint64
	Time               time.Time
	LowestAskPrice     float64
	LowestAskQuantity  float64
	HighestBidPrice    float64
	HighestBidQuantity float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
