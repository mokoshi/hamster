package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/usecase"
	"net/http"
)

type OrderHandler interface {
	GetOpenOrders(c echo.Context) error
}

type OrderHandlerImpl struct {
	OrderUsecase usecase.OrderUsecase
}

func NewOrderHandler(u usecase.OrderUsecase) OrderHandler {
	return &OrderHandlerImpl{OrderUsecase: u}
}

func (u OrderHandlerImpl) GetOpenOrders(c echo.Context) (err error) {
	orders, err := u.OrderUsecase.GetOpenOrders()

	res := make([]*resource.Order, len(orders))
	for i, order := range orders {
		res[i] = resource.NewOrder(order)
	}

	return c.JSON(http.StatusOK, res)
}
