package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type RateHistory struct {
	model.RateHistory
	Id        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"index"`
}
