package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hamster/common"
	"hamster/infra/external"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/interfaces/handler"
	"hamster/lib/cc_client"
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
	common.Init()

	err := db.Init()
	if err != nil {
		panic(err)
	}

	ccClient := cc_client.NewClient(common.Env.CoincheckApiKey, common.Env.CoincheckApiSecret)

	//tradingHistoryPersistence := persistence.NewTradingHistoryPersistence(db.DB)
	//tradingHistoryUsecase := usecase.NewTradingHistoryUsecase(tradingHistoryPersistence)
	//tradingHistoryHandler := handler.NewTradingHistoryHandler(tradingHistoryUsecase)

	orderBooksHistoryPersistence := persistence.NewOrderBooksHistoryPersistence(db.DB)
	orderExternal := external.NewOrderExternal(ccClient)

	orderBooksHistoryUsecase := usecase.NewOrderBooksHistoryUsecase(orderBooksHistoryPersistence)
	orderUsecase := usecase.NewOrderUsecase(orderExternal)

	orderBooksHistoryHandler := handler.NewOrderBooksHistoryHandler(orderBooksHistoryUsecase)
	healthCheckHandler := handler.NewHealthCheckHistoryHandler()
	orderHandler := handler.NewOrderHandler(orderUsecase)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/health_check", healthCheckHandler.Check)
	e.GET("/order_books_histories", orderBooksHistoryHandler.GetHistories)
	e.GET("/orders/open", orderHandler.GetOpenOrders)
	e.Logger.Fatal(e.Start(":4100"))
}
