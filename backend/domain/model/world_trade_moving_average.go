package model

import "time"

type WorldTradeMovingAverage struct {
	Id        uint64
	Time      time.Time
	Duration  time.Duration
	Pair      string
	OrderType string
	Rate      float64
}
