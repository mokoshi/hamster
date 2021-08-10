package model

import "time"

type RateHistory struct {
	Id        uint64
	OrderType string
	Pair      string
	Rate      float64
	CreatedAt time.Time
}
