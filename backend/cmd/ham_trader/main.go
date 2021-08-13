package main

import (
	"context"
	"fmt"
	"hamster/common"
	"hamster/domain/model"
	"hamster/infra/external"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/lib/cc_client"
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

	orderBooksRepository := persistence.NewOrderBooksRepository(ccClient, db.DB)
	orderExternal := external.NewOrderExternal(ccClient)

	rateHistoryRepository := persistence.NewRateHistoryPersistence(db.DB, ccClient)

	exchangeUsecase := usecase.NewExchangeUsecase(orderExternal, rateHistoryRepository)
	orderBooksUsecase := usecase.NewOrderBooksUsecase(orderBooksRepository)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("signal received.")
		cancel()
	}()

	// 板情報を同期する goroutine
	orderBooksRepository.SubscribeLatest(func(orderBooks *model.OrderBooks) {
		lowestAsk := orderBooks.GetLowestAsk()
		highestBid := orderBooks.GetHighestBid()
		fmt.Printf(
			"売 %.0f : 買 %.0f\n",
			lowestAsk.Price,
			highestBid.Price,
		)

		orderBooksUsecase.WriteWithBuffering(lowestAsk, highestBid, 10)
	})

	// レートを同期する goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				exchangeUsecase.SyncCurrentRate("buy", "btc_jpy")
				exchangeUsecase.SyncCurrentRate("sell", "btc_jpy")
				time.Sleep(time.Second * 5)
			}
		}
	}()

	wg.Wait()
}
