package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/usecase"
	"net/http"
)

type ExchangeHandler interface {
	GetOpenOrders(c echo.Context) error
}

type ExchangeHandlerImpl struct {
	ExchangeUsecase usecase.ExchangeUsecase
}

func NewExchangeHandler(u usecase.ExchangeUsecase) ExchangeHandler {
	return &ExchangeHandlerImpl{ExchangeUsecase: u}
}

func (u ExchangeHandlerImpl) GetOpenOrders(c echo.Context) (err error) {
	orders, err := u.ExchangeUsecase.GetOpenOrders()

	res := make([]*resource.Order, len(orders))
	for i, order := range orders {
		res[i] = resource.NewOrder(order)
	}

	return c.JSON(http.StatusOK, res)
}
