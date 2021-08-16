package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/lib/util"
	"hamster/usecase"
	"net/http"
	"time"
)

type ExchangeHandler interface {
	GetOpenOrders(c echo.Context) error
	GetRateHistories(c echo.Context) error
}

type ExchangeHandlerImpl struct {
	ExchangeUsecase usecase.ExchangeUsecase
}

func NewExchangeHandler(u usecase.ExchangeUsecase) ExchangeHandler {
	return &ExchangeHandlerImpl{ExchangeUsecase: u}
}

func (u ExchangeHandlerImpl) GetOpenOrders(c echo.Context) (err error) {
	orders, err := u.ExchangeUsecase.GetOpenOrders()

	res := make([]*resource.OpenOrder, len(orders))
	for i, order := range orders {
		res[i] = resource.NewOpenOrder(order)
	}

	return c.JSON(http.StatusOK, res)
}

func (u ExchangeHandlerImpl) GetRateHistories(c echo.Context) (err error) {
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	histories, err := u.ExchangeUsecase.GetRateHistories(fromTime, toTime)

	res := make([]*resource.RateHistory, len(histories))
	for i, history := range histories {
		res[i] = resource.NewRateHistory(history)
	}

	return c.JSON(http.StatusOK, res)
}
