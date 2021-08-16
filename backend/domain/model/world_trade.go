package model

import "time"

type WorldTrade struct {
	Id        uint64
	Time      time.Time
	Pair      string
	OrderType string
	Rate      float64
	Amount    float64
}
