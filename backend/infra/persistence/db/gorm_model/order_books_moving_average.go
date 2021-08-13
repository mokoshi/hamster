package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksMovingAverage struct {
	model.OrderBooksMovingAverage
	Id   uint64    `gorm:"primary_key"`
	Time time.Time `gorm:"index"`
}
