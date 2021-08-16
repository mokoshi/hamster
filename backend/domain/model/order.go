package model

import (
	"encoding/json"
	"math"
)

type status int

const (
	instantiated = "instantiated"
	requested    = "requested"
	pending      = "pending"
	completed    = "completed"
	canceled     = "canceled"
)

type Order struct {
	*OrderRequest
	TradeProcedureId         string
	Status                   string
	Reason                   string
	CompletedAmount          float64
	CompletedMarketBuyAmount float64
	PendingAmount            float64
	PendingMarketBuyAmount   float64
	Transactions             []*Transaction
}

func NewOrder(orderRequest *OrderRequest) *Order {
	return &Order{
		OrderRequest: orderRequest,
		Status:       requested,
	}
}

func (o *Order) String() string {
	s, _ := json.MarshalIndent(o, "", "\t")
	return string(s)
}

func (o *Order) IsPending() bool {
	if o.Status == requested || o.Status == pending {
		return true
	} else {
		return false
	}
}

func (o *Order) AddTransactions(transactions []*Transaction) bool {
	o.Transactions = append(o.Transactions, transactions...)

	btcAmount := float64(0)
	jpyAmount := float64(0)
	for _, transaction := range o.Transactions {
		btcAmount += transaction.FundBtc
		jpyAmount += transaction.FundJpy
	}

	completed := false
	if o.OrderType == "market_buy" {
		// 成行買いのときだけ、JPY指定での取引
		if math.Abs(*o.MarketBuyAmount+jpyAmount) < 0.1 {
			completed = true
		}
	} else {
		if math.Abs(*o.Amount+btcAmount) < 0.000000001 {
			completed = true
		}
	}

	if completed {
		o.SetAsCompleted()
	}

	return completed
}

func (o *Order) SetAsCompleted() {
	o.Status = completed

	if o.OrderType == "market_buy" && o.MarketBuyAmount != nil {
		o.CompletedMarketBuyAmount = *o.MarketBuyAmount
		o.PendingMarketBuyAmount = 0
	} else {
		if o.Amount == nil {
			//	このケースは無いはず
			o.CompletedAmount = 0
		} else {
			o.CompletedAmount = *o.Amount
		}
		o.PendingAmount = 0
	}
}

func (o *Order) SetAsCanceled(reason string) {
	o.Status = canceled
	o.Reason = reason
}
