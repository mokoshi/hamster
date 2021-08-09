package repository

import (
	"hamster/domain/model"
)

type BalanceRepository interface {
	GetBalance() (*model.Balance, error)
}
