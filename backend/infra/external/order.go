package external

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/util"
)

type OrderExternal struct {
	Client *cc_client.Client
}

func NewOrderExternal(client *cc_client.Client) repository.OrderRepository {
	return &OrderExternal{Client: client}
}

func (obe *OrderExternal) GetOpenOrders() ([]*model.Order, error) {
	res, err := obe.Client.GetOpenOrders()
	if err != nil {
		return nil, err
	}

	orders := make([]*model.Order, len(res.Orders))
	for i, o := range res.Orders {
		order, err := parseOrder(o)
		if err != nil {
			return nil, err
		}
		orders[i] = order
	}

	return orders, nil
}

func parseOrder(res cc_client.Order) (*model.Order, error) {
	return &model.Order{
		Id:                     res.Id,
		OrderType:              res.OrderType,
		Rate:                   util.ParseFloat64(res.Rate),
		Pair:                   res.Pair,
		PendingAmount:          util.ParseFloat64(res.PendingAmount),
		PendingMarketBuyAmount: util.ParseFloat64(res.PendingMarketBuyAmount),
		StopLossRate:           util.ParseFloat64(res.StopLossRate),
		CreatedAt:              res.CreatedAt,
	}, nil
}
