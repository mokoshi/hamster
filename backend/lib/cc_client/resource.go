package cc_client

import (
	"encoding/json"
	"time"
)

type OrderBooks struct {
	Asks [][2]interface{} `json:"asks"`
	Bids [][2]interface{} `json:"bids"`
}

type OrderBooksDiffStream [2]json.RawMessage

type Order struct {
	Id                     uint64    `json:"id"`
	OrderType              string    `json:"order_type"`
	Rate                   string    `json:"rate"`
	Pair                   string    `json:"pair"`
	PendingAmount          string    `json:"pending_amount"`
	PendingMarketBuyAmount string    `json:"pending_market_buy_amount"`
	StopLossRate           string    `json:"stop_loss_rate"`
	CreatedAt              time.Time `json:"created_at"`
}
type OpenOrders struct {
	Success bool    `json:"success"`
	Orders  []Order `json:"orders"`
}

type Balance struct {
	Success      bool   `json:"success"`
	Jpy          string `json:"jpy"`
	Btc          string `json:"btc"`
	JpyReserved  string `json:"jpy_reserved"`
	BtcReserved  string `json:"btc_reserved""`
	JpyLendInUse string `json:"jpy_lend_in_use"`
	BtcLendInUse string `json:"btc_lend_in_use"`
	JpyLent      string `json:"jpy_lent"`
	BtcLent      string `json:"btc_lent"`
	JpyDebt      string `json:"jpy_debt"`
	BtcDebt      string `json:"btc_debt"`
}
