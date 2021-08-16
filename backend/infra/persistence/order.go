package persistence

import (
	"fmt"
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/cc_client"
	"hamster/lib/util"
)

type OrderRepository struct {
	client     *cc_client.Client
	db         *gorm.DB
	orderCache *model.OrderCache
}

func NewOrderRepository(
	client *cc_client.Client,
	db *gorm.DB,
	orderCache *model.OrderCache,
) repository.OrderRepository {
	return &OrderRepository{
		client:     client,
		db:         db,
		orderCache: orderCache,
	}
}

func (r *OrderRepository) SyncOpenOrders() ([]*model.OpenOrder, error) {
	res, err := r.client.GetOpenOrders()
	if err != nil {
		return nil, err
	}

	openOrders := make([]*model.OpenOrder, len(res.Orders))
	for i, o := range res.Orders {
		order, err := parseOpenOrder(o)
		if err != nil {
			return nil, err
		}
		openOrders[i] = order
	}

	r.orderCache.SetOpenOrders(openOrders)

	// TODO なんか固まる、なんで
	//err = r.db.Transaction(func(tx *gorm.DB) error {
	//	if err := tx.Where("1 = 1").Delete(&model.OpenOrder{}).Error; err != nil {
	//		return err
	//	}
	//	if err := r.db.Create(openOrders).Error; err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//clog.Logger.Error(err)

	return openOrders, nil
}

func (r *OrderRepository) GetOpenOrderCount() int {
	return r.orderCache.GetOrderCount()
}

func (r *OrderRepository) RequestOrder(orderRequest *model.OrderRequest) (*model.OrderRequest, error) {
	res, err := r.client.CreateOrder(
		orderRequest.Pair,
		orderRequest.OrderType,
		orderRequest.Rate,
		orderRequest.Amount,
		orderRequest.MarketBuyAmount,
		orderRequest.StopLossRate,
	)
	if err != nil {
		return nil, err
	}

	order := parseOrder(res)
	if order == nil {
		return nil, fmt.Errorf("failed to create order")
	}
	order.MarketBuyAmount = orderRequest.MarketBuyAmount

	return order, nil
}

func parseOpenOrder(res cc_client.OpenOrder) (*model.OpenOrder, error) {
	return &model.OpenOrder{
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

func parseOrder(res *cc_client.Order) *model.OrderRequest {
	if !res.Success {
		return nil
	}

	rate := util.ParseFloat64(res.Rate)
	amount := util.ParseFloat64(res.Amount)
	stopLossRate := util.ParseFloat64(res.StopLossRate)

	return &model.OrderRequest{
		Id:           res.Id,
		Pair:         res.Pair,
		OrderType:    res.OrderType,
		Rate:         &rate,
		Amount:       &amount,
		StopLossRate: &stopLossRate,
		CreatedAt:    res.CreatedAt,
	}
}
