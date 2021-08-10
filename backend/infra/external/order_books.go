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
	return &OrderBooksExternal{
		Client:     client,
		orderBooks: model.NewOrderBooks(nil, nil),
	}
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

	asks, bids := parseOrders(res)

	obe.orderBooks = model.NewOrderBooks(asks, bids)

	return obe.orderBooks, nil
}

func (obe *OrderBooksExternal) Subscribe(listener func(books *model.OrderBooks)) error {
	// Subscribe する前に一度APIで情報を取っておきたいが、
	// APIでの更新とSubscribeの更新タイミングの問題で古いデータが残ってしまう可能性がありそう
	// とりあえず API 呼び出しはしないようにしておく
	//if obe.orderBooks == nil {
	//	_, err := obe.Get(false)
	//	if err != nil {
	//		return err
	//	}
	//}

	err := obe.Client.SubscribeOrderBooks(func(pair string, diff *cc_client.OrderBooks) {
		asks, bids := parseOrders(diff)
		obe.orderBooks.Update(asks, bids)

		listener(obe.orderBooks)
	})

	return err
}

func parseOrders(res *cc_client.OrderBooks) (asks []*model.OrderBookItem, bids []*model.OrderBookItem) {
	ParseOrder := func(res [2]interface{}) (*model.OrderBookItem, error) {
		price, err := strconv.ParseFloat(res[0].(string), 64)
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.ParseFloat(res[1].(string), 64)
		if err != nil {
			return nil, err
		}
		return &model.OrderBookItem{Price: price, Quantity: quantity}, nil
	}

	asks = make([]*model.OrderBookItem, len(res.Asks))
	bids = make([]*model.OrderBookItem, len(res.Bids))

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
