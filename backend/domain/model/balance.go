package model

type Balance struct {
	Jpy          float64
	Btc          float64
	JpyReserved  float64
	BtcReserved  float64
	JpyLendInUse float64
	BtcLendInUse float64
	JpyLent      float64
	BtcLent      float64
	JpyDebt      float64
	BtcDebt      float64
}
