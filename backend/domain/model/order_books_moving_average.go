package model

import "time"

type OrderBooksMovingAverage struct {
	Id       uint64
	Time     time.Time
	Duration time.Duration
	AskPrice float64
	BidPrice float64
}
