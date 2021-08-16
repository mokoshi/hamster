package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hamster/common"
	"hamster/domain/model"
	"hamster/infra/external"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/interfaces/handler"
	"hamster/lib/cc_client"
	"hamster/usecase"
)

func main() {
	common.Init()

	err := db.Init()
	if err != nil {
		panic(err)
	}

	ccClient := cc_client.NewClient(common.Env.CoincheckApiKey, common.Env.CoincheckApiSecret)

	orderBooksCache := model.NewOrderBooksCache()
	orderCache := model.NewOrderCache()

	orderBooksRepository := persistence.NewOrderBooksRepository(ccClient, db.DB, orderBooksCache, nil)
	orderRepository := persistence.NewOrderRepository(ccClient, db.DB, orderCache)
	balanceExternal := external.NewBalanceExternal(ccClient)
	rateHistoryRepository := persistence.NewRateHistoryPersistence(db.DB, ccClient)

	orderBooksHistoryUsecase := usecase.NewOrderBooksUsecase(orderBooksRepository)
	exchangeUsecase := usecase.NewExchangeUsecase(orderRepository, rateHistoryRepository)
	accountUsecase := usecase.NewAccountUsecase(balanceExternal)

	orderBooksHandler := handler.NewOrderBooksHandler(orderBooksHistoryUsecase)
	healthCheckHandler := handler.NewHealthCheckHistoryHandler()
	exchangeHandler := handler.NewExchangeHandler(exchangeUsecase)
	accountHandler := handler.NewAccountHandler(accountUsecase)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/health_check", healthCheckHandler.Check)
	e.GET("/order_books/snapshots", orderBooksHandler.GetSnapshots)
	e.GET("/order_books/moving_averages", orderBooksHandler.GetMovingAverages)
	e.GET("/exchange/open_orders", exchangeHandler.GetOpenOrders)
	e.GET("/exchange/rate_histories", exchangeHandler.GetRateHistories)
	e.GET("/account/balance", accountHandler.GetBalance)
	e.Logger.Fatal(e.Start(":4100"))
}
