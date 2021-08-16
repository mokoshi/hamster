package model

import (
	"time"
)

type Transaction struct {
	Id          uint64
	OrderId     uint64
	CreatedAt   time.Time
	FundBtc     float64
	FundJpy     float64
	Pair        string
	Rate        float64
	FeeCurrency string
	Fee         float64
	Liquidity   string
	Side        string
}
