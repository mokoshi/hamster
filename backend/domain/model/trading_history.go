package model

import (
	"time"
)

type TradingHistory struct {
	ID        uint
	TradedAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
