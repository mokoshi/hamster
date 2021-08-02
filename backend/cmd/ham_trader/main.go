package main

import (
	"fmt"
	"hamster/domain/model"
	"hamster/infra/external"
	"hamster/infra/persistence"
	"hamster/infra/persistence/db"
	"hamster/lib/cc_client"
	"hamster/usecase"
	"time"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}

	orderBooksHistoryPersistence := persistence.NewOrderBooksHistoryPersistence(db.DB)
	orderBooksHistoryUsecase := usecase.NewOrderBooksHistoryUsecase(orderBooksHistoryPersistence)

	ccClient := cc_client.NewClient("DUMMY-API-KEY")

	orderBookRepository := external.NewOrderBooksExternal(ccClient)

	err = orderBookRepository.Subscribe(func(orderBooks *model.OrderBooks) {
		lowestAsk := orderBooks.GetLowestAsk()
		highestBid := orderBooks.GetHighestBid()
		fmt.Printf(
			"売 %.0f : 買 %.0f\n",
			lowestAsk.Price,
			highestBid.Price,
		)

		orderBooksHistoryUsecase.WriteWithBuffering(lowestAsk, highestBid)
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("--tick--")
		time.Sleep(time.Second * 5)
	}
}
