package cc_client

import "encoding/json"

type OrderBooks struct {
	Asks [][2]interface{} `json:"asks"`
	Bids [][2]interface{} `json:"bids"`
}

type OrderBooksDiffStream [2]json.RawMessage
