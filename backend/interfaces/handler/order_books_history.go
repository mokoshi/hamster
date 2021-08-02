package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/usecase"
	"net/http"
	"time"
)

type OrderBooksHistoryHandler interface {
	GetHistories(c echo.Context) error
}

type OrderBooksHistoryHandlerImpl struct {
	OrderBooksHistoryUsecase usecase.OrderBooksHistoryUsecase
}

func NewOrderBooksHistoryHandler(u usecase.OrderBooksHistoryUsecase) OrderBooksHistoryHandler {
	return &OrderBooksHistoryHandlerImpl{
		OrderBooksHistoryUsecase: u,
	}
}

func (u OrderBooksHistoryHandlerImpl) GetHistories(c echo.Context) (err error) {
	var fromTime, toTime time.Time
	fromTime, err = time.Parse("2006-01-02", c.QueryParam("from"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	toTime, err = time.Parse("2006-01-02", c.QueryParam("to"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	histories, err := u.OrderBooksHistoryUsecase.GetHistories(fromTime, toTime)

	res := make([]*resource.OrderBooksHistory, len(histories))
	for i, history := range histories {
		res[i] = resource.NewOrderBooksHistory(history)
	}

	return c.JSON(http.StatusOK, res)
}
