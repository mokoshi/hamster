package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type WorldTradeMovingAverage struct {
	model.WorldTradeMovingAverage
	Id   uint64    `gorm:"primary_key"`
	Time time.Time `gorm:"index"`
}
