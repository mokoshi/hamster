package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksSnapshot struct {
	model.OrderBooksSnapshot
	Id   uint64    `gorm:"primary_key"`
	Time time.Time `gorm:"index"`
}
