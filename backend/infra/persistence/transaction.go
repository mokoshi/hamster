package persistence

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/util"
)

type TransactionRepository struct {
	client           *cc_client.Client
	db               *gorm.DB
	transactionCache *model.TransactionCache
}

func NewTransactionRepository(
	client *cc_client.Client,
	db *gorm.DB,
	transactionCache *model.TransactionCache,
) repository.TransactionRepository {
	return &TransactionRepository{
		client:           client,
		db:               db,
		transactionCache: transactionCache,
	}
}

func (r *TransactionRepository) SyncRecentTransactions() ([]*model.Transaction, error) {
	res, err := r.client.GetRecentTransactions()
	if err != nil {
		return nil, err
	}

	transactions := parseTransactions(res)

	addedTransactions := r.transactionCache.AddTransactions(transactions)

	if len(addedTransactions) > 0 {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(transactions).Error; err != nil {
				return err
			}
			return nil
		})
	}

	return addedTransactions, err
}

func parseTransactions(res *cc_client.Transactions) []*model.Transaction {
	if !res.Success {
		return nil
	}

	transactions := make([]*model.Transaction, len(res.Transactions))
	for i, t := range res.Transactions {
		fundBtc := util.ParseFloat64(t.Funds.Btc)
		fundJpy := util.ParseFloat64(t.Funds.Jpy)
		rate := util.ParseFloat64(t.Rate)
		fee := util.ParseFloat64(t.Fee)

		transactions[i] = &model.Transaction{
			Id:          t.Id,
			OrderId:     t.OrderId,
			CreatedAt:   t.CreatedAt,
			FundBtc:     fundBtc,
			FundJpy:     fundJpy,
			Pair:        t.Pair,
			Rate:        rate,
			FeeCurrency: t.FeeCurrency,
			Fee:         fee,
			Liquidity:   t.Liquidity,
			Side:        t.Side,
		}
	}

	return transactions
}
