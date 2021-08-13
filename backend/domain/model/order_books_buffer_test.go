package model

import (
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
	"time"
)

var now = time.Now()
var floatApproxOpt = cmp.Comparer(func(x, y float64) bool {
	delta := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	return delta/mean < 0.00001
})

func TestAdd_バッファが増えていく(t *testing.T) {
	buffer := OrderBooksBuffer{}

	buffer.Add(&OrderBooksHistory{Id: 1, Time: now, LowestAskPrice: 5.0, HighestBidPrice: 3.0})
	if got := len(buffer.histories); got != 1 {
		t.Fatalf("len(histories) = %v, expected %v", got, 1)
	}
	if got := len(buffer.latestHistories); got != 1 {
		t.Fatalf("len(latestHistories) = %v, expected %v", got, 1)
	}
	if got := len(buffer.movingAverages); got != 1 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 1)
	}

	buffer.Add(&OrderBooksHistory{Id: 2, Time: now.Add(time.Second), LowestAskPrice: 4.0, HighestBidPrice: 4.0})
	if got := len(buffer.histories); got != 2 {
		t.Fatalf("len(histories) = %v, expected %v", got, 2)
	}
	if got := len(buffer.latestHistories); got != 2 {
		t.Fatalf("len(latestHistories) = %v, expected %v", got, 2)
	}
	if got := len(buffer.movingAverages); got != 2 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 2)
	}
}

func TestAdd_移動平均の計算が正しい(t *testing.T) {
	buffer := OrderBooksBuffer{}

	buffer.Add(&OrderBooksHistory{Id: 1, Time: now, LowestAskPrice: 5.0, HighestBidPrice: 3.0})
	buffer.Add(&OrderBooksHistory{Id: 2, Time: now.Add(time.Second), LowestAskPrice: 4.0, HighestBidPrice: 4.0})
	buffer.Add(&OrderBooksHistory{Id: 3, Time: now.Add(time.Second * 2), LowestAskPrice: 4.2, HighestBidPrice: 2.0})
	buffer.Add(&OrderBooksHistory{Id: 4, Time: now.Add(time.Second * 10), LowestAskPrice: 4.1, HighestBidPrice: 3.6})

	expected0 := &OrderBooksMovingAverage{
		Id: 0, Time: now, Duration: time.Second * 10, AskPrice: 5.0, BidPrice: 3.0}
	if diff := cmp.Diff(expected0, buffer.movingAverages[0], floatApproxOpt); diff != "" {
		t.Errorf("buffer.movingAverages[0] differs:\n%s", diff)
	}
	expected1 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second), Duration: time.Second * 10, AskPrice: 4.5, BidPrice: 3.5}
	if diff := cmp.Diff(expected1, buffer.movingAverages[1], floatApproxOpt); diff != "" {
		t.Errorf("buffer.movingAverages[1] differs:\n%s", diff)
	}
	expected2 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second * 2), Duration: time.Second * 10, AskPrice: 4.4, BidPrice: 3.0}
	if diff := cmp.Diff(expected2, buffer.movingAverages[2], floatApproxOpt); diff != "" {
		t.Errorf("buffer.movingAverages[2] differs:\n%s", diff)
	}
	expected3 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second * 10), Duration: time.Second * 10, AskPrice: 4.1, BidPrice: 3.2}
	if diff := cmp.Diff(expected3, buffer.movingAverages[3], floatApproxOpt); diff != "" {
		t.Errorf("buffer.movingAverages[3] differs:\n%s", diff)
	}
}

func TestAdd_LatestHistoriesの先頭が順次削除されていく(t *testing.T) {
	buffer := OrderBooksBuffer{}

	buffer.Add(&OrderBooksHistory{Id: 1, Time: now})
	buffer.Add(&OrderBooksHistory{Id: 2, Time: now.Add(time.Second)})
	buffer.Add(&OrderBooksHistory{Id: 3, Time: now.Add(time.Second * 2)})
	if got := len(buffer.latestHistories); got != 3 {
		t.Fatalf("len(latestHistories) = %v, expected %v", got, 3)
	}

	buffer.Add(&OrderBooksHistory{Id: 4, Time: now.Add(time.Second * 10)})
	if got := len(buffer.latestHistories); got != 3 {
		t.Fatalf("len(latestHistories) = %v, expected %v", got, 3)
	}
	if got := buffer.latestHistories[0].Id; got != 2 {
		t.Fatalf("latestHistories[0].Id = %v, expected %v", got, 2)
	}

	buffer.Add(&OrderBooksHistory{Id: 5, Time: now.Add(time.Second * 19)})
	if got := len(buffer.latestHistories); got != 2 {
		t.Fatalf("len(latestHistories) = %v, expected %v", got, 2)
	}
	if got := buffer.latestHistories[0].Id; got != 4 {
		t.Fatalf("latestHistories[0].Id = %v, expected %v", got, 4)
	}
}

func TestRead(t *testing.T) {
	buffer := OrderBooksBuffer{}

	buffer.Add(&OrderBooksHistory{Id: 1, Time: now})
	buffer.Add(&OrderBooksHistory{Id: 2, Time: now.Add(time.Second)})

	histories, movingAverages := buffer.Read()
	if got := len(histories); got != 2 {
		t.Fatalf("len(histories) = %v, expected %v", got, 2)
	}
	if got := len(movingAverages); got != 2 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 2)
	}
	if got := len(buffer.histories); got != 0 {
		t.Fatalf("len(buffer.histories) = %v, expected %v", got, 0)
	}
	if got := len(buffer.movingAverages); got != 0 {
		t.Fatalf("len(buffer.movingAverages) = %v, expected %v", got, 0)
	}
	if got := len(buffer.latestHistories); got != 2 {
		t.Fatalf("len(buffer.latestHistories) = %v, expected %v", got, 2)
	}
}

func TestGetBufferedSize(t *testing.T) {
	buffer := OrderBooksBuffer{}
	buffer.Add(&OrderBooksHistory{Id: 1, Time: now})
	if got := buffer.BufferedSize(); got != 1 {
		t.Fatalf("BufferedSize() = %v, expected %v", got, 1)
	}
}
