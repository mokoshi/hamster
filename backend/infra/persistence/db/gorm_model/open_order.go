package gorm_model

import (
	"hamster/domain/model"
	"time"
)

type OpenOrder struct {
	model.OpenOrder
	Id        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"index"`
}
