package model

import (
	"encoding/json"
	"time"
)

type OrderRequest struct {
	Id              uint64
	Pair            string
	OrderType       string // "buy" | "sell" | "market_buy" | "market_sell"
	Rate            *float64
	Amount          *float64
	MarketBuyAmount *float64
	StopLossRate    *float64
	CreatedAt       time.Time
}

func (r *OrderRequest) String() string {
	s, _ := json.MarshalIndent(r, "", "\t")
	return string(s)
}

func NewMarketBuyOrder(btcAmount float64, expectedRate float64) *OrderRequest {
	jpyAmount := btcAmount * expectedRate
	return &OrderRequest{
		Pair:            "btc_jpy",
		OrderType:       "market_buy",
		MarketBuyAmount: &jpyAmount,
	}
}

func NewMarketSellOrder(btcAmount float64) *OrderRequest {
	return &OrderRequest{
		Pair:      "btc_jpy",
		OrderType: "market_sell",
		Amount:    &btcAmount,
	}
}

func NewLimitBuyOrder(btcAmount float64, rate float64) *OrderRequest {
	return &OrderRequest{
		Pair:      "btc_jpy",
		OrderType: "buy",
		Rate:      &rate,
		Amount:    &btcAmount,
	}
}

func NewLimitSellOrder(btcAmount float64, rate float64) *OrderRequest {
	return &OrderRequest{
		Pair:      "btc_jpy",
		OrderType: "sell",
		Rate:      &rate,
		Amount:    &btcAmount,
	}
}
