package external

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/util"
)

type BalanceExternal struct {
	Client *cc_client.Client
}

func NewBalanceExternal(client *cc_client.Client) repository.BalanceRepository {
	return &BalanceExternal{Client: client}
}

func (obe *BalanceExternal) GetBalance() (*model.Balance, error) {
	res, err := obe.Client.GetBalance()
	if err != nil {
		return nil, err
	}

	return &model.Balance{
		Jpy:          util.ParseFloat64(res.Jpy),
		Btc:          util.ParseFloat64(res.Btc),
		JpyReserved:  util.ParseFloat64(res.JpyReserved),
		BtcReserved:  util.ParseFloat64(res.BtcReserved),
		JpyLendInUse: util.ParseFloat64(res.JpyLendInUse),
		BtcLendInUse: util.ParseFloat64(res.BtcLendInUse),
		JpyLent:      util.ParseFloat64(res.JpyLent),
		BtcLent:      util.ParseFloat64(res.BtcLent),
		JpyDebt:      util.ParseFloat64(res.JpyDebt),
		BtcDebt:      util.ParseFloat64(res.BtcDebt),
	}, nil
}
