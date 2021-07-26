package main

import (
	"fmt"
	"hamster/domain/model"
	"hamster/infra/external"
	"hamster/lib/cc_client"
	"time"
)

func main() {
	ccClient := cc_client.NewClient("DUMMY-API-KEY")

	orderBookRepository := external.NewOrderBooksExternal(ccClient)

	err := orderBookRepository.Subscribe(func(orderBooks *model.OrderBooks) {
		fmt.Println(orderBooks.GetLowestAsk(), orderBooks.GetHighestBid())
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
