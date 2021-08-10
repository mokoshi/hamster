package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type OrderBooksHistory struct {
	model.OrderBooksHistory
	Id   uint64    `gorm:"primary_key"`
	Time time.Time `gorm:"index"`
}
