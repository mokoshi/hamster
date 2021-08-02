package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hamster/domain/model"
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
	err = DB.AutoMigrate(&model.TradingHistory{}, &model.OrderBooksHistory{})
	if err != nil {
		clog.Logger.Fatal("Failed to run auto-migration...", err)
		return err
	}

	return nil
}
