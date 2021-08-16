package repository

import (
	"hamster/domain/model"
)

type TransactionRepository interface {
	SyncRecentTransactions() ([]*model.Transaction, error)
}
