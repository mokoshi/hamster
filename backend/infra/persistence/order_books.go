package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/clog"
	"strconv"
	"time"
)

type OrderBooksRepository struct {
	client          *cc_client.Client
	db              *gorm.DB
	orderBooksCache *model.OrderBooksCache
	flushSize       int
}

func NewOrderBooksRepository(
	client *cc_client.Client,
	db *gorm.DB,
	orderBooksCache *model.OrderBooksCache,
	flushSize *int,
) repository.OrderBooksRepository {
	var fs int
	if flushSize == nil {
		fs = 10
	} else {
		fs = *flushSize
	}

	return &OrderBooksRepository{
		client:          client,
		db:              db,
		orderBooksCache: orderBooksCache,
		flushSize:       fs,
	}
}

func (r *OrderBooksRepository) SyncOrderBooks() (*model.OrderBooks, error) {
	res, err := r.client.GetOrderBooks()
	if err != nil {
		return nil, err
	}

	asks, bids := r.parseOrders(res)

	return r.orderBooksCache.SetCurrent(model.NewOrderBooks(asks, bids)), nil
}

func (r *OrderBooksRepository) SubscribeAndSync(listener func(books *model.OrderBooks)) error {
	// Subscribe する前に一度APIで情報を取っておきたいが、
	// APIでの更新とSubscribeの更新タイミングの問題で古いデータが残ってしまう可能性がありそう
	// とりあえず API 呼び出しはしないようにしておく
	//if r.orderBooks == nil {
	//	_, err := r.Get(false)
	//	if err != nil {
	//		return err
	//	}
	//}

	err := r.client.SubscribeOrderBooks(func(pair string, diff *cc_client.OrderBooks) {
		asks, bids := r.parseOrders(diff)

		c := r.orderBooksCache.Update(asks, bids)

		if r.orderBooksCache.BufferedSize() > r.flushSize {
			snapshots, averages := r.orderBooksCache.ReadNext()
			if err := r.db.Create(snapshots).Error; err != nil {
				clog.Logger.Error(err)
			}
			if err := r.db.Create(averages).Error; err != nil {
				clog.Logger.Error(err)
			}
		}

		listener(c)
	})

	return err
}

func (r *OrderBooksRepository) GetSnapshots(from time.Time, to time.Time) ([]*model.OrderBooksSnapshot, error) {
	var snapshots []*model.OrderBooksSnapshot
	r.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&snapshots)
	return snapshots, nil
}

func (r *OrderBooksRepository) GetMovingAverages(from time.Time, to time.Time) ([]*model.OrderBooksMovingAverage, error) {
	var averages []*model.OrderBooksMovingAverage
	r.db.Where("time BETWEEN ? AND ?", from, to).Order("time").Find(&averages)
	return averages, nil
}

func (r *OrderBooksRepository) GetLatestSnapshot(offsetFromLatest int) *model.OrderBooksSnapshot {
	return r.orderBooksCache.GetLatestSnapshot(offsetFromLatest)
}

func (r *OrderBooksRepository) GetLatestMovingAverage(offsetFromLatest int) *model.OrderBooksMovingAverage {
	return r.orderBooksCache.GetLatestMovingAverage(offsetFromLatest)
}

func (r *OrderBooksRepository) parseOrders(res *cc_client.OrderBooks) (asks []*model.OrderBookItem, bids []*model.OrderBookItem) {
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
