package model

import (
	"time"
)

const movingAverageDuration = time.Second * 10

type OrderBooksBuffer struct {
	histories []*OrderBooksHistory

	// 移動平均を計算するため保持しておく最近の板情報
	latestHistories []*OrderBooksHistory
	movingAverages  []*OrderBooksMovingAverage
}

func (b *OrderBooksBuffer) Add(history *OrderBooksHistory) {
	b.histories = append(b.histories, history)
	b.latestHistories = append(b.latestHistories, history)

	boundaryTime := history.Time.Add(-movingAverageDuration)
	count := 0
	askPriceTotal := float64(0)
	bidPriceTotal := float64(0)
	for i := len(b.latestHistories) - 1; i >= 0; i-- {
		h := b.latestHistories[i]
		if !h.Time.After(boundaryTime) {
			break
		}

		count++
		askPriceTotal += h.LowestAskPrice
		bidPriceTotal += h.HighestBidPrice
	}
	b.movingAverages = append(b.movingAverages, &OrderBooksMovingAverage{
		Time:     history.Time,
		Duration: movingAverageDuration,
		AskPrice: askPriceTotal / float64(count),
		BidPrice: bidPriceTotal / float64(count),
	})
	b.latestHistories = b.latestHistories[len(b.latestHistories)-count:]
}

func (b *OrderBooksBuffer) Read() ([]*OrderBooksHistory, []*OrderBooksMovingAverage) {
	defer func() {
		b.histories = nil
		b.movingAverages = nil
	}()
	return b.histories, b.movingAverages
}

func (b *OrderBooksBuffer) BufferedSize() int {
	return len(b.histories)
}
