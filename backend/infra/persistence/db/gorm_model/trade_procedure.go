package gorm_model

import (
	"hamster/domain/model"
)

type TradeProcedure struct {
	model.TradeProcedure
	Id        string    `gorm:"primary_key"`
}
