package model

type TransactionCache struct {
	transactionMap map[uint64]*Transaction
}

func NewTransactionCache() *TransactionCache {
	return &TransactionCache{
		transactionMap: map[uint64]*Transaction{},
	}
}

func (c *TransactionCache) AddTransactions(transactions []*Transaction) []*Transaction {
	var addedTransactions []*Transaction = nil
	for _, transaction := range transactions {
		_, ok := c.transactionMap[transaction.Id]
		if !ok {
			addedTransactions = append(addedTransactions, transaction)
		}

		c.transactionMap[transaction.Id] = transaction
	}
	return addedTransactions
}
