package model

import (
	"time"
)

const (
	orderBooksMovingAverageDuration = time.Second * 10
	orderBooksRecentCacheSize       = 1000
)

type OrderBooksCache struct {
	current        *OrderBooks
	snapshots      []*OrderBooksSnapshot
	movingAverages []*OrderBooksMovingAverage

	// 直近のデータを常に保持しておくところ
	cachedRecentSnapshots      []*OrderBooksSnapshot // こちらは移動平均を計算するためにも使われる
	cachedRecentMovingAverages []*OrderBooksMovingAverage
}

func NewOrderBooksCache() *OrderBooksCache {
	return &OrderBooksCache{
		current: NewOrderBooks(nil, nil),
	}
}

func (c *OrderBooksCache) GetCurrent() *OrderBooks {
	return c.current
}

func (c *OrderBooksCache) SetCurrent(ob *OrderBooks) *OrderBooks {
	c.current = ob
	c.AddLatestSnapshot(ob)
	return c.current
}

func (c *OrderBooksCache) Update(asks []*OrderBookItem, bids []*OrderBookItem) *OrderBooks {
	current := c.current.Update(asks, bids)
	c.AddLatestSnapshot(current)
	return current
}

func (c *OrderBooksCache) AddLatestSnapshot(orderBooks *OrderBooks) {
	lowestAsk := orderBooks.GetLowestAsk()
	highestBid := orderBooks.GetHighestBid()
	snapshot := &OrderBooksSnapshot{
		Time:               time.Now(),
		LowestAskPrice:     lowestAsk.Price,
		LowestAskQuantity:  lowestAsk.Quantity,
		HighestBidPrice:    highestBid.Price,
		HighestBidQuantity: highestBid.Quantity,
	}

	c.snapshots = append(c.snapshots, snapshot)
	c.cachedRecentSnapshots = append(c.cachedRecentSnapshots, snapshot)

	boundaryTime := snapshot.Time.Add(-orderBooksMovingAverageDuration)
	count := 0
	askPriceTotal := float64(0)
	bidPriceTotal := float64(0)
	for i := len(c.cachedRecentSnapshots) - 1; i >= 0; i-- {
		h := c.cachedRecentSnapshots[i]
		if !h.Time.After(boundaryTime) {
			break
		}

		count++
		askPriceTotal += h.LowestAskPrice
		bidPriceTotal += h.HighestBidPrice
	}
	movingAverage := &OrderBooksMovingAverage{
		Time:        snapshot.Time,
		Duration:    orderBooksMovingAverageDuration,
		MiddlePrice: (askPriceTotal + bidPriceTotal) / float64(count) / 2,
		AskPrice:    askPriceTotal / float64(count),
		BidPrice:    bidPriceTotal / float64(count),
	}
	c.movingAverages = append(c.movingAverages, movingAverage)
	c.cachedRecentMovingAverages = append(c.cachedRecentMovingAverages, movingAverage)

	// キャッシュサイズをオーバーしないように、はみ出る分は削除していく
	if len(c.cachedRecentSnapshots) > orderBooksRecentCacheSize {
		c.cachedRecentSnapshots = c.cachedRecentSnapshots[len(c.cachedRecentSnapshots)-orderBooksRecentCacheSize:]
	}
	if len(c.cachedRecentMovingAverages) > orderBooksRecentCacheSize {
		c.cachedRecentMovingAverages = c.cachedRecentMovingAverages[len(c.cachedRecentMovingAverages)-orderBooksRecentCacheSize:]
	}
}

func (c *OrderBooksCache) ReadNext() ([]*OrderBooksSnapshot, []*OrderBooksMovingAverage) {
	defer func() {
		c.snapshots = nil
		c.movingAverages = nil
	}()
	return c.snapshots, c.movingAverages
}

func (c *OrderBooksCache) BufferedSize() int {
	return len(c.snapshots)
}

func (c *OrderBooksCache) GetLatestSnapshot(offsetFromLatest int) *OrderBooksSnapshot {
	l := len(c.cachedRecentSnapshots)
	if l <= offsetFromLatest {
		return nil
	}
	return c.cachedRecentSnapshots[l-offsetFromLatest-1]
}

func (c *OrderBooksCache) GetLatestMovingAverage(offsetFromLatest int) *OrderBooksMovingAverage {
	l := len(c.cachedRecentMovingAverages)
	if l <= offsetFromLatest {
		return nil
	}
	return c.cachedRecentMovingAverages[l-offsetFromLatest-1]
}
