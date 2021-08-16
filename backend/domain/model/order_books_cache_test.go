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
	manager := OrderBooksCache{}

	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 1, Time: now, LowestAskPrice: 5.0, HighestBidPrice: 3.0})
	if got := len(manager.snapshots); got != 1 {
		t.Fatalf("len(snapshots) = %v, expected %v", got, 1)
	}
	if got := len(manager.cachedRecentSnapshots); got != 1 {
		t.Fatalf("len(cachedRecentSnapshots) = %v, expected %v", got, 1)
	}
	if got := len(manager.movingAverages); got != 1 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 1)
	}

	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 2, Time: now.Add(time.Second), LowestAskPrice: 4.0, HighestBidPrice: 4.0})
	if got := len(manager.snapshots); got != 2 {
		t.Fatalf("len(snapshots) = %v, expected %v", got, 2)
	}
	if got := len(manager.cachedRecentSnapshots); got != 2 {
		t.Fatalf("len(cachedRecentSnapshots) = %v, expected %v", got, 2)
	}
	if got := len(manager.movingAverages); got != 2 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 2)
	}
}

func TestAdd_移動平均の計算が正しい(t *testing.T) {
	manager := OrderBooksCache{}

	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 1, Time: now, LowestAskPrice: 5.0, HighestBidPrice: 3.0})
	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 2, Time: now.Add(time.Second), LowestAskPrice: 4.0, HighestBidPrice: 4.0})
	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 3, Time: now.Add(time.Second * 2), LowestAskPrice: 4.2, HighestBidPrice: 2.0})
	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 4, Time: now.Add(time.Second * 10), LowestAskPrice: 4.1, HighestBidPrice: 3.6})

	expected0 := &OrderBooksMovingAverage{
		Id: 0, Time: now, Duration: time.Second * 10, AskPrice: 5.0, BidPrice: 3.0}
	if diff := cmp.Diff(expected0, manager.movingAverages[0], floatApproxOpt); diff != "" {
		t.Errorf("manager.movingAverages[0] differs:\n%s", diff)
	}
	expected1 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second), Duration: time.Second * 10, AskPrice: 4.5, BidPrice: 3.5}
	if diff := cmp.Diff(expected1, manager.movingAverages[1], floatApproxOpt); diff != "" {
		t.Errorf("manager.movingAverages[1] differs:\n%s", diff)
	}
	expected2 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second * 2), Duration: time.Second * 10, AskPrice: 4.4, BidPrice: 3.0}
	if diff := cmp.Diff(expected2, manager.movingAverages[2], floatApproxOpt); diff != "" {
		t.Errorf("manager.movingAverages[2] differs:\n%s", diff)
	}
	expected3 := &OrderBooksMovingAverage{
		Id: 0, Time: now.Add(time.Second * 10), Duration: time.Second * 10, AskPrice: 4.1, BidPrice: 3.2}
	if diff := cmp.Diff(expected3, manager.movingAverages[3], floatApproxOpt); diff != "" {
		t.Errorf("manager.movingAverages[3] differs:\n%s", diff)
	}
}

func TestAdd_キャッシュが正常にパージされていく(t *testing.T) {
	manager := OrderBooksCache{}

	for i := 0; i < 1000; i++ {
		manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: uint64(i)})
	}
	if got := len(manager.cachedRecentSnapshots); got != 1000 {
		t.Fatalf("len(cachedRecentSnapshots) = %v, expected %v", got, 1000)
	}
	if got := len(manager.cachedRecentMovingAverages); got != 1000 {
		t.Fatalf("len(cachedRecentMovingAverages) = %v, expected %v", got, 1000)
	}
	if got := manager.cachedRecentSnapshots[0].Id; got != 0 {
		t.Fatalf("cachedRecentSnapshots[0].Id = %v, expected %v", got, 0)
	}

	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: uint64(1000)})
	if got := len(manager.cachedRecentSnapshots); got != 1000 {
		t.Fatalf("len(cachedRecentSnapshots) = %v, expected %v", got, 1000)
	}
	if got := len(manager.cachedRecentMovingAverages); got != 1000 {
		t.Fatalf("len(cachedRecentMovingAverages) = %v, expected %v", got, 1000)
	}
	if got := manager.cachedRecentSnapshots[0].Id; got != 1 {
		t.Fatalf("cachedRecentSnapshots[0].Id = %v, expected %v", got, 1)
	}
}

func TestRead(t *testing.T) {
	manager := OrderBooksCache{}

	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 1, Time: now})
	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 2, Time: now.Add(time.Second)})

	histories, movingAverages := manager.ReadNext()
	if got := len(histories); got != 2 {
		t.Fatalf("len(snapshots) = %v, expected %v", got, 2)
	}
	if got := len(movingAverages); got != 2 {
		t.Fatalf("len(movingAverages) = %v, expected %v", got, 2)
	}
	if got := len(manager.snapshots); got != 0 {
		t.Fatalf("len(manager.snapshots) = %v, expected %v", got, 0)
	}
	if got := len(manager.movingAverages); got != 0 {
		t.Fatalf("len(manager.movingAverages) = %v, expected %v", got, 0)
	}
	if got := len(manager.cachedRecentSnapshots); got != 2 {
		t.Fatalf("len(manager.cachedRecentSnapshots) = %v, expected %v", got, 2)
	}
}

func TestGetBufferedSize(t *testing.T) {
	manager := OrderBooksCache{}
	manager.AddLatestSnapshot(&OrderBooksSnapshot{Id: 1, Time: now})
	if got := manager.BufferedSize(); got != 1 {
		t.Fatalf("BufferedSize() = %v, expected %v", got, 1)
	}
}
