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
)

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

	orderBooksRepository := persistence.NewOrderBooksRepository(ccClient, db.DB)
	orderExternal := external.NewOrderExternal(ccClient)
	balanceExternal := external.NewBalanceExternal(ccClient)
	rateHistoryRepository := persistence.NewRateHistoryPersistence(db.DB, ccClient)

	orderBooksHistoryUsecase := usecase.NewOrderBooksUsecase(orderBooksRepository)
	exchangeUsecase := usecase.NewExchangeUsecase(orderExternal, rateHistoryRepository)
	accountUsecase := usecase.NewAccountUsecase(balanceExternal)

	orderBooksHandler := handler.NewOrderBooksHandler(orderBooksHistoryUsecase)
	healthCheckHandler := handler.NewHealthCheckHistoryHandler()
	exchangeHandler := handler.NewExchangeHandler(exchangeUsecase)
	accountHandler := handler.NewAccountHandler(accountUsecase)

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/health_check", healthCheckHandler.Check)
	e.GET("/order_books/histories", orderBooksHandler.GetHistories)
	e.GET("/order_books/moving_averages", orderBooksHandler.GetMovingAverages)
	e.GET("/exchange/open_orders", exchangeHandler.GetOpenOrders)
	e.GET("/exchange/rate_histories", exchangeHandler.GetRateHistories)
	e.GET("/account/balance", accountHandler.GetBalance)
	e.Logger.Fatal(e.Start(":4100"))
}
