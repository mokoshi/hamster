package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"strconv"
	"time"
)

type OrderBooksRepository struct {
	client     *cc_client.Client
	db         *gorm.DB
	orderBooks *model.OrderBooks
}

func NewOrderBooksRepository(client *cc_client.Client, db *gorm.DB) repository.OrderBooksRepository {
	return &OrderBooksRepository{
		client:     client,
		db:         db,
		orderBooks: model.NewOrderBooks(nil, nil),
	}
}

func (obe *OrderBooksRepository) FetchCurrent(refresh bool) (*model.OrderBooks, error) {
	// キャッシュがある場合はそのまま返す
	if !refresh && obe.orderBooks != nil {
		return obe.orderBooks, nil
	}

	res, err := obe.client.GetOrderBooks()
	if err != nil {
		return nil, err
	}

	asks, bids := parseOrders(res)

	obe.orderBooks = model.NewOrderBooks(asks, bids)

	return obe.orderBooks, nil
}

func (obe *OrderBooksRepository) SubscribeLatest(listener func(books *model.OrderBooks)) error {
	// Subscribe する前に一度APIで情報を取っておきたいが、
	// APIでの更新とSubscribeの更新タイミングの問題で古いデータが残ってしまう可能性がありそう
	// とりあえず API 呼び出しはしないようにしておく
	//if obe.orderBooks == nil {
	//	_, err := obe.Get(false)
	//	if err != nil {
	//		return err
	//	}
	//}

	err := obe.client.SubscribeOrderBooks(func(pair string, diff *cc_client.OrderBooks) {
		asks, bids := parseOrders(diff)
		obe.orderBooks.Update(asks, bids)

		listener(obe.orderBooks)
	})

	return err
}

func (obhp *OrderBooksRepository) GetHistories(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error) {
	var histories []*model.OrderBooksHistory
	obhp.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&histories)
	return histories, nil
}

func (obhp *OrderBooksRepository) GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error) {
	var averages []*model.OrderBooksMovingAverage
	obhp.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&averages)
	return averages, nil
}

func (obhp *OrderBooksRepository) CreateHistories(histories []*model.OrderBooksHistory) error {
	if err := obhp.db.Create(histories).Error; err != nil {
		return err
	}
	return nil
}

func (obhp *OrderBooksRepository) CreateMovingAverages(averages []*model.OrderBooksMovingAverage) error {
	if err := obhp.db.Create(averages).Error; err != nil {
		return err
	}
	return nil
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
