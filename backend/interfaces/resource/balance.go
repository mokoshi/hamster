package resource

import (
	"hamster/domain/model"
)

type Balance struct {
	Jpy          float64 `json:"jpy"`
	Btc          float64 `json:"btc"`
	JpyReserved  float64 `json:"jpyReserved"`
	BtcReserved  float64 `json:"btcReserved"`
	JpyLendInUse float64 `json:"jpyLendInUse"`
	BtcLendInUse float64 `json:"btcLendInUse"`
	JpyLent      float64 `json:"jpyLent"`
	BtcLent      float64 `json:"btcLent"`
	JpyDebt      float64 `json:"jpyDebt"`
	BtcDebt      float64 `json:"btcDebt"`
}

func NewBalance(model *model.Balance) *Balance {
	return &Balance{
		Jpy:          model.Jpy,
		Btc:          model.Btc,
		JpyReserved:  model.JpyReserved,
		BtcReserved:  model.BtcReserved,
		JpyLendInUse: model.JpyLendInUse,
		BtcLendInUse: model.BtcLendInUse,
		JpyLent:      model.JpyLent,
		BtcLent:      model.BtcLent,
		JpyDebt:      model.JpyDebt,
		BtcDebt:      model.BtcDebt,
	}
}
