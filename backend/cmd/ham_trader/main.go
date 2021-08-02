package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/infra/external"
	"hamster/infra/persistence"
	"hamster/lib/cc_client"
	"hamster/usecase"
	"log"
	"time"
)

func main() {
	dsn := "local:local@tcp(db:3306)/hamster?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect db...", err)
	}

	orderBooksHistoryPersistence := persistence.NewOrderBooksHistoryPersistence(db)
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
