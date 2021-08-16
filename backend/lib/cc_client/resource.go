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

type Trade struct {
	Id        uint64 `json:"id"`
	Pair      string `json:"pair"`
	Rate      string `json:"rate"`
	Amount    string `json:"amount"`
	OrderType string `json:"order_type"`
}
type TradeStream [5]json.RawMessage

type OpenOrder struct {
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
	Success bool        `json:"success"`
	Orders  []OpenOrder `json:"orders"`
}

type Order struct {
	Success      bool      `json:"success"`
	Id           uint64    `json:"id"`
	OrderType    string    `json:"order_type"`
	Rate         string    `json:"rate"`
	Pair         string    `json:"pair"`
	Amount       string    `json:"amount"`
	StopLossRate string    `json:"stop_loss_rate"`
	CreatedAt    time.Time `json:"created_at"`
}

type Transactions struct {
	Success      bool          `json:"success"`
	Transactions []Transaction `json:"transactions"`
}
type Transaction struct {
	Id        uint64    `json:"id"`
	OrderId   uint64    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
	Funds     struct {
		Btc string `json:"btc"`
		Jpy string `json:"jpy"`
	} `json:"funds"`
	Pair        string `json:"pair"`
	Rate        string `json:"rate"`
	FeeCurrency string `json:"fee_currency"`
	Fee         string `json:"fee"`
	Liquidity   string `json:"liquidity"`
	Side        string `json:"side"`
}

type OrderCancel struct {
	Success bool   `json:"success"`
	Id      uint64 `json:"id"`
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

type Rate struct {
	Success bool   `json:"success"`
	Rate    string `json:"rate"`
	Price   string `json:"price"`
	Amount  string `json:"amount"`
}
