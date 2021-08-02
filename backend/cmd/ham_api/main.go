package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/interfaces/handler"
	"hamster/usecase"
	"time"
)

const location = "Asia/Tokyo"

func init() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}

	//tradingHistoryPersistence := persistence.NewTradingHistoryPersistence(db.DB)
	//tradingHistoryUsecase := usecase.NewTradingHistoryUsecase(tradingHistoryPersistence)
	//tradingHistoryHandler := handler.NewTradingHistoryHandler(tradingHistoryUsecase)

	orderBooksHistoryPersistence := persistence.NewOrderBooksHistoryPersistence(db.DB)
	orderBooksHistoryUsecase := usecase.NewOrderBooksHistoryUsecase(orderBooksHistoryPersistence)
	orderBooksHistoryHandler := handler.NewOrderBooksHistoryHandler(orderBooksHistoryUsecase)
	healthCheckHandler := handler.NewHealthCheckHistoryHandler()

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/health_check", healthCheckHandler.Check)
	e.GET("/order_books_histories", orderBooksHistoryHandler.GetHistories)
	e.Logger.Fatal(e.Start(":4100"))
}
