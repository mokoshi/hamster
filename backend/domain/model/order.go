package model

import "time"

type Order struct {
	Id                     uint64
	OrderType              string // "buy" or "sell"
	Rate                   float64
	Pair                   string
	PendingAmount          float64
	PendingMarketBuyAmount float64
	StopLossRate           float64
	CreatedAt              time.Time
}
