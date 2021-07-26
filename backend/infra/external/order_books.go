package external

import (
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"strconv"
)

type OrderBooksExternal struct {
	Client *cc_client.Client

	// キャッシュ
	orderBooks *model.OrderBooks
}

func NewOrderBooksExternal(client *cc_client.Client) repository.OrderBooksRepository {
	return &OrderBooksExternal{Client: client}
}

func (obe *OrderBooksExternal) Get(refresh bool) (*model.OrderBooks, error) {
	// キャッシュがある場合はそのまま返す
	if !refresh && obe.orderBooks != nil {
		return obe.orderBooks, nil
	}

	res, err := obe.Client.GetOrderBooks()
	if err != nil {
		return nil, err
	}

	asks, bids := ParseOrders(res)

	obe.orderBooks = model.NewOrderBooks(asks, bids)

	return obe.orderBooks, nil
}

func (obe *OrderBooksExternal) Subscribe(listener func(books *model.OrderBooks)) error {
	// まだキャッシュがない場合は、一度 API 叩いて取得する
	if obe.orderBooks == nil {
		_, err := obe.Get(false)
		if err != nil {
			return err
		}
	}

	err := obe.Client.SubscribeOrderBooks(func(pair string, diff *cc_client.OrderBooks) {
		asks, bids := ParseOrders(diff)
		obe.orderBooks.Update(asks, bids)

		listener(obe.orderBooks)
	})

	return err
}

func ParseOrders(res *cc_client.OrderBooks) (asks []*model.Order, bids []*model.Order) {
	ParseOrder := func(res [2]interface{}) (*model.Order, error) {
		price, err := strconv.ParseFloat(res[0].(string), 64)
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.ParseFloat(res[1].(string), 64)
		if err != nil {
			return nil, err
		}
		return &model.Order{Price: price, Quantity: quantity}, nil
	}

	asks = make([]*model.Order, len(res.Asks))
	bids = make([]*model.Order, len(res.Bids))

	for i, ask := range res.Asks {
		order, err := ParseOrder(ask)
		if err != nil {
			return nil, nil
		}
		asks[i] = order
	}
	for i, bid := range res.Bids {
		order, err := ParseOrder(bid)
		if err != nil {
			return nil, nil
		}
		bids[i] = order
	}

	return asks, bids
}
