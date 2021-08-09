package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/lib/util"
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
	fromTime := time.Unix(util.ParseInt64(c.QueryParam("from")), 0)
	toTime := time.Unix(util.ParseInt64(c.QueryParam("to")), 0)

	histories, err := u.OrderBooksHistoryUsecase.GetHistories(fromTime, toTime)

	res := make([]*resource.OrderBooksHistory, len(histories))
	for i, history := range histories {
		res[i] = resource.NewOrderBooksHistory(history)
	}

	return c.JSON(http.StatusOK, res)
}
