package model

type OrderCache struct {
	openOrderMap map[uint64]*OpenOrder
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		openOrderMap: map[uint64]*OpenOrder{},
	}
}

func (c *OrderCache) SetOpenOrders(openOrders []*OpenOrder) {
	c.openOrderMap = map[uint64]*OpenOrder{}
	for _, openOrder := range openOrders {
		c.openOrderMap[openOrder.Id] = openOrder
	}
}

func (c *OrderCache) AddOrder(order *OrderRequest) {
	c.openOrderMap[order.Id] = &OpenOrder{
		Id:                     order.Id,
		OrderType:              order.OrderType,
		Rate:                   *order.Rate,
		Pair:                   order.Pair,
		PendingAmount:          *order.Amount,
		PendingMarketBuyAmount: 0,
		StopLossRate:           0,
		CreatedAt:              order.CreatedAt,
	}
}

func (c *OrderCache) GetOrderCount() int {
	return len(c.openOrderMap)
}
