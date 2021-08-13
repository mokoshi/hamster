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
	GetHistories(c echo.Context) error
	GetMovingAverages(c echo.Context) error
}

type OrderBooksHandlerImpl struct {
	OrderBooksUsecase usecase.OrderBooksUsecase
}

func NewOrderBooksHandler(u usecase.OrderBooksUsecase) OrderBooksHandler {
	return &OrderBooksHandlerImpl{
		OrderBooksUsecase: u,
	}
}

func (u OrderBooksHandlerImpl) GetHistories(c echo.Context) (err error) {
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	histories, err := u.OrderBooksUsecase.GetHistories(fromTime, toTime)

	res := make([]*resource.OrderBooksHistory, len(histories))
	for i, history := range histories {
		res[i] = resource.NewOrderBooksHistory(history)
	}

	return c.JSON(http.StatusOK, res)
}

func (u OrderBooksHandlerImpl) GetMovingAverages(c echo.Context) (err error) {
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	averages, err := u.OrderBooksUsecase.GetMovingAverages(fromTime, toTime)

	res := make([]*resource.OrderBooksMovingAverage, len(averages))
	for i, average := range averages {
		res[i] = resource.NewOrderBooksMovingAverage(average)
	}

	return c.JSON(http.StatusOK, res)
}
