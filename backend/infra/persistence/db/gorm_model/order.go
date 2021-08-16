package gorm_model

import (
	"hamster/domain/model"
)

type Order struct {
	model.Order
	Id uint64 `gorm:"primary_key"`
}
