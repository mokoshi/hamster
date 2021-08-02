package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthCheckHistoryHandler interface {
	Check(c echo.Context) error
}

type HealthCheckHistoryHandlerImpl struct{}

func NewHealthCheckHistoryHandler() HealthCheckHistoryHandler {
	return &HealthCheckHistoryHandlerImpl{}
}

func (u HealthCheckHistoryHandlerImpl) Check(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, &struct {
		Ok bool `json:"ok"`
	}{true})
}
