package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/infra/persistence"
	"hamster/interfaces/handler"
	"hamster/usecase"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>hamster neko pokemon 3<h1>")
}

func main() {
	dsn := "local:local@tcp(db:3306)/hamster?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect db...", err)
	}

	db.AutoMigrate(&model.TradingHistory{}, &model.OrderBooksHistory{})

	tradingHistoryPersistence := persistence.NewTradingHistoryPersistence(db)
	tradingHistoryUsecase := usecase.NewTradingHistoryUsecase(tradingHistoryPersistence)
	tradingHistoryHandler := handler.NewTradingHistoryHandler(tradingHistoryUsecase)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/trade", tradingHistoryHandler.HandleTrade)
	log.Fatal(http.ListenAndServe(":4100", nil))
}
