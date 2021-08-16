package main

import (
	"context"
	"fmt"
	"hamster/common"
	"hamster/domain/model"
	"hamster/domain/service"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/lib/cc_client"
	"hamster/lib/clog"
	"hamster/usecase"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	common.Init()

	err := db.Init()
	if err != nil {
		panic(err)
	}

	ccClient := cc_client.NewClient(common.Env.CoincheckApiKey, common.Env.CoincheckApiSecret)

	orderBooksCache := model.NewOrderBooksCache()
	worldTradeCache := model.NewWorldTradeCache()
	orderCache := model.NewOrderCache()
	transactionCache := model.NewTransactionCache()

	orderBooksRepository := persistence.NewOrderBooksRepository(ccClient, db.DB, orderBooksCache, nil)
	worldTradeRepository := persistence.NewWorldTradeRepository(ccClient, db.DB, worldTradeCache, nil)
	orderRepository := persistence.NewOrderRepository(ccClient, db.DB, orderCache)
	transactionRepository := persistence.NewTransactionRepository(ccClient, db.DB, transactionCache)
	rateHistoryRepository := persistence.NewRateHistoryPersistence(db.DB, ccClient)
	tradeProcedureRepository := persistence.NewTradeProcedureRepository(db.DB)

	traderService := service.NewTraderService(orderBooksRepository, orderRepository)

	exchangeUsecase := usecase.NewExchangeUsecase(orderRepository, rateHistoryRepository)
	orderBooksUsecase := usecase.NewOrderBooksUsecase(orderBooksRepository)
	worldTradeUsecase := usecase.NewWorldTradeUsecase(worldTradeRepository)
	orderUsecase := usecase.NewTraderUsecase(
		traderService,
		orderRepository,
		orderBooksRepository,
		transactionRepository,
		tradeProcedureRepository,
	)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("signal received.")
		cancel()
	}()

	// 板情報を同期する
	orderBooksUsecase.StartSync(func(ob *model.OrderBooks) {
		err := orderUsecase.ProcessNext()
		if err != nil {
			clog.Logger.Error(err)
		}
	})

	// 取引情報を同期する
	worldTradeUsecase.StartSync()

	// レートを同期する goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if false {
					exchangeUsecase.SyncCurrentRate("buy", "btc_jpy")
					exchangeUsecase.SyncCurrentRate("sell", "btc_jpy")
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()

	wg.Wait()
}
