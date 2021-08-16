package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/lib/util"
	"hamster/usecase"
	"net/http"
	"time"
)

type OrderBooksHandler interface {
	GetSnapshots(c echo.Context) error
	GetMovingAverages(c echo.Context) error
}

type orderBooksHandler struct {
	OrderBooksUsecase usecase.OrderBooksUsecase
}

func NewOrderBooksHandler(u usecase.OrderBooksUsecase) OrderBooksHandler {
	return &orderBooksHandler{
		OrderBooksUsecase: u,
	}
}

func (u orderBooksHandler) GetSnapshots(c echo.Context) (err error) {
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	histories, err := u.OrderBooksUsecase.GetSnapshots(fromTime, toTime)

	res := make([]*resource.OrderBooksSnapshot, len(histories))
	for i, history := range histories {
		res[i] = resource.NewOrderBooksHistory(history)
	}

	return c.JSON(http.StatusOK, res)
}

func (u orderBooksHandler) GetMovingAverages(c echo.Context) (err error) {
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	averages, err := u.OrderBooksUsecase.GetMovingAverages(fromTime, toTime)

	res := make([]*resource.OrderBooksMovingAverage, len(averages))
	for i, average := range averages {
		res[i] = resource.NewOrderBooksMovingAverage(average)
	}

	return c.JSON(http.StatusOK, res)
}
