package handler

import (
	"github.com/labstack/echo/v4"
	"hamster/interfaces/resource"
	"hamster/usecase"
	"net/http"
)

type AccountHandler interface {
	GetBalance(c echo.Context) error
}

type AccountHandlerImpl struct {
	AccountUsecase usecase.AccountUsecase
}

func NewAccountHandler(u usecase.AccountUsecase) AccountHandler {
	return &AccountHandlerImpl{AccountUsecase: u}
}

func (u AccountHandlerImpl) GetBalance(c echo.Context) (err error) {
	balance, err := u.AccountUsecase.GetBalance()

	return c.JSON(http.StatusOK, resource.NewBalance(balance))
}
