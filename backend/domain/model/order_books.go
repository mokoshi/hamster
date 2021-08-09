package model

import (
	"math"
)

type OrderBooks struct {
	// Asks, Bids はそれぞれ [Price]: *OrderBookItem の map で保持
	Asks map[float64]*OrderBookItem // 売り注文
	Bids map[float64]*OrderBookItem // 買い注文
}

type OrderBookItem struct {
	Price    float64
	Quantity float64
}

func NewOrderBooks(asks []*OrderBookItem, bids []*OrderBookItem) *OrderBooks {
	asksMap := map[float64]*OrderBookItem{}
	bidsMap := map[float64]*OrderBookItem{}

	for _, ask := range asks {
		asksMap[ask.Price] = ask
	}
	for _, bid := range bids {
		bidsMap[bid.Price] = bid
	}

	return &OrderBooks{Asks: asksMap, Bids: bidsMap}
}

func (o OrderBooks) Update(asks []*OrderBookItem, bids []*OrderBookItem) {
	for _, ask := range asks {
		if ask.Quantity == 0 {
			delete(o.Asks, ask.Price)
		} else {
			o.Asks[ask.Price] = ask
		}
	}
	for _, bid := range bids {
		if bid.Quantity == 0 {
			delete(o.Bids, bid.Price)
		} else {
			o.Bids[bid.Price] = bid
		}
	}
}

// TODO 都度最小値見つけるの無駄なので、データ構造工夫したほうが良いかもしれない
func (o OrderBooks) GetLowestAsk() *OrderBookItem {
	lowest := math.MaxFloat64
	for price, _ := range o.Asks {
		if lowest > price {
			lowest = price
		}
	}

	lowestOrder, ok := o.Asks[lowest]
	if ok {
		return lowestOrder
	} else {
		return nil
	}
}

func (o OrderBooks) GetHighestBid() *OrderBookItem {
	highest := float64(0)
	for price, _ := range o.Bids {
		if highest < price {
			highest = price
		}
	}

	highestOrder, ok := o.Bids[highest]
	if ok {
		return highestOrder
	} else {
		return nil
	}
}
