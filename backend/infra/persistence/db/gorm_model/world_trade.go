package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type WorldTrade struct {
	model.WorldTrade
	Id   uint64    `gorm:"primary_key"`
	Time time.Time `gorm:"index"`
}
