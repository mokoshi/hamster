package service

import (
	"github.com/golang/geo/r2"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/lib/clog"
	"math"
	"time"
)

const maxConcurrentOrders = 1

type TraderService interface {
	ShouldBuy() []*model.OrderRequest
}

type traderService struct {
	orderBooksRepository repository.OrderBooksRepository
}

func NewTraderService(orderBooksRepository repository.OrderBooksRepository, orderRepository repository.OrderRepository) *traderService {
	return &traderService{
		orderBooksRepository: orderBooksRepository,
	}
}

func (s *traderService) ShouldBuy() []*model.OrderRequest {
	average0 := s.orderBooksRepository.GetLatestMovingAverage(0)
	average1 := s.orderBooksRepository.GetLatestMovingAverage(1)
	average4 := s.orderBooksRepository.GetLatestMovingAverage(4)

	if average0 == nil || average1 == nil || average4 == nil {
		return nil
	}

	//currentGradient := calcGradient(
	//	&pricePoint{time: average1.Time, price: average1.MiddlePrice},
	//	&pricePoint{time: average0.Time, price: average0.MiddlePrice},
	//)
	//previousGradient := calcGradient(
	//	&pricePoint{time: average4.Time, price: average4.MiddlePrice},
	//	&pricePoint{time: average1.Time, price: average1.MiddlePrice},
	//)
	//start := currentGradient > 0 && previousGradient < 0
	v1 := newVector(
		&pricePoint{time: average4.Time, price: average4.MiddlePrice},
		&pricePoint{time: average1.Time, price: average1.MiddlePrice},
	)
	v2 := newVector(
		&pricePoint{time: average1.Time, price: average1.MiddlePrice},
		&pricePoint{time: average0.Time, price: average0.MiddlePrice},
	)
	clog.Logger.Infof("%0.2f", calcAngle(v1, v2))

	start := calcAngle(v1, v2) > 10.0

	if start {
		buffer := 1.00005
		rate := s.orderBooksRepository.GetLatestSnapshot(0).LowestAskPrice * buffer

		return []*model.OrderRequest{
			model.NewMarketBuyOrder(0.005, rate),
			model.NewLimitSellOrder(0.005, rate*1.0002),
		}
	}

	return nil

}

type pricePoint struct {
	time  time.Time
	price float64
}

func calcGradient(start *pricePoint, end *pricePoint) float64 {
	w := end.time.UnixNano() - start.time.UnixNano()
	return (end.price - start.price) / float64(w)
}
func newVector(start *pricePoint, end *pricePoint) r2.Point {
	return r2.Point{
		X: float64(end.time.UnixNano()-start.time.UnixNano()) / float64(time.Millisecond),
		Y: end.price - start.price,
	}
}
func calcAngle(v1 r2.Point, v2 r2.Point) float64 {
	theta := math.Acos(v1.Dot(v2) / v1.Norm() / v2.Norm())

	if v1.Cross(v2) < 0 {
		theta *= -1
	}

	return theta * 180 / math.Pi
}
