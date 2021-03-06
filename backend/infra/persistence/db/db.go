package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hamster/infra/persistence/db/gorm_model"
	"hamster/lib/clog"
)

var DB *gorm.DB

func Init() (err error) {
	// TODO 環境によって切り替える
	dsn := "local:local@tcp(db:3306)/hamster?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		clog.Logger.Fatal("Failed to connect db...", err)
		return err
	}

	// マイグレーションも走らせちゃう
	err = DB.AutoMigrate(
		&gorm_model.Order{},
		&gorm_model.OpenOrder{},
		&gorm_model.OrderBooksSnapshot{},
		&gorm_model.OrderBooksMovingAverage{},
		&gorm_model.RateHistory{},
		&gorm_model.TradeProcedure{},
		&gorm_model.Transaction{},
		&gorm_model.WorldTrade{},
		&gorm_model.WorldTradeMovingAverage{},
	)
	if err != nil {
		clog.Logger.Fatal("Failed to run auto-migration...", err)
		return err
	}

	return nil
}
