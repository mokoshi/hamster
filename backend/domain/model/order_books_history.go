package model

import "time"

type OrderBooksHistory struct {
	ID                 uint
	Time               time.Time
	LowestAskPrice     float64
	LowestAskQuantity  float64
	HighestBidPrice    float64
	HighestBidQuantity float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
