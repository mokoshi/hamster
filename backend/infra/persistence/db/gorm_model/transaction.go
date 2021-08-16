package gorm_model

import (
	"hamster/domain/model"
)

type Transaction struct {
	model.Transaction
	Id uint64 `gorm:"primary_key"`
}
