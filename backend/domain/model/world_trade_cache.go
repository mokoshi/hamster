package model

import (
	"time"
)

const (
	worldTradeMovingAverageDuration = time.Second * 10
	worldTradeRecentCacheSize       = 1000
)

type WorldTradeCache struct {
	worldTrades    []*WorldTrade
	movingAverages []*WorldTradeMovingAverage

	// 直近のデータを常に保持しておくところ
	cachedRecentTrades         []*WorldTrade // こちらは移動平均を計算するためにも使われる
	cachedRecentMovingAverages []*WorldTradeMovingAverage
}

func NewWorldTradeCache() *WorldTradeCache {
	return &WorldTradeCache{}
}

func (c *WorldTradeCache) AddLatestTrade(trade *WorldTrade) {
	c.worldTrades = append(c.worldTrades, trade)
	c.cachedRecentTrades = append(c.cachedRecentTrades, trade)

	boundaryTime := trade.Time.Add(-worldTradeMovingAverageDuration)
	count := 0
	rateTotal := float64(0)
	for i := len(c.cachedRecentTrades) - 1; i >= 0; i-- {
		t := c.cachedRecentTrades[i]
		if !t.Time.After(boundaryTime) {
			break
		}

		count++
		rateTotal += t.Rate
	}
	movingAverage := &WorldTradeMovingAverage{
		Time:      trade.Time,
		Duration:  worldTradeMovingAverageDuration,
		Pair:      trade.Pair,
		OrderType: trade.OrderType,
		Rate:      rateTotal / float64(count),
	}
	c.movingAverages = append(c.movingAverages, movingAverage)
	c.cachedRecentMovingAverages = append(c.cachedRecentMovingAverages, movingAverage)

	// キャッシュサイズをオーバーしないように、はみ出る分は削除していく
	if len(c.cachedRecentTrades) > worldTradeRecentCacheSize {
		c.cachedRecentTrades = c.cachedRecentTrades[len(c.cachedRecentTrades)-worldTradeRecentCacheSize:]
	}
	if len(c.cachedRecentMovingAverages) > worldTradeRecentCacheSize {
		c.cachedRecentMovingAverages = c.cachedRecentMovingAverages[len(c.cachedRecentMovingAverages)-worldTradeRecentCacheSize:]
	}
}

func (c *WorldTradeCache) ReadNext() ([]*WorldTrade, []*WorldTradeMovingAverage) {
	defer func() {
		c.worldTrades = nil
		c.movingAverages = nil
	}()
	return c.worldTrades, c.movingAverages
}

func (c *WorldTradeCache) BufferedSize() int {
	return len(c.worldTrades)
}

func (c *WorldTradeCache) GetLatestMovingAverage(offsetFromLatest int) *WorldTradeMovingAverage {
	l := len(c.cachedRecentMovingAverages)
	if l <= offsetFromLatest {
		return nil
	}
	return c.cachedRecentMovingAverages[l-offsetFromLatest-1]
}
