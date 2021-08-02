package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
	"time"
)

type OrderBooksHistoryPersistence struct {
	Db     *gorm.DB
	Buffer []*model.OrderBooksHistory
}

func NewOrderBooksHistoryPersistence(db *gorm.DB) repository.OrderBooksHistoryRepository {
	return &OrderBooksHistoryPersistence{Db: db}
}

func (obhp *OrderBooksHistoryPersistence) Get(from time.Time, to time.Time) ([]*model.OrderBooksHistory, error) {
	var histories []*model.OrderBooksHistory
	obhp.Db.Where("time BETWEEN ? AND ?", from, to).Find(&histories)
	return histories, nil
}

func (obhp *OrderBooksHistoryPersistence) AddToBuffer(orderBooksHistory *model.OrderBooksHistory) (int, error) {
	obhp.Buffer = append(obhp.Buffer, orderBooksHistory)
	return len(obhp.Buffer), nil
}

func (obhp *OrderBooksHistoryPersistence) Flush() error {
	if err := obhp.Db.Create(&obhp.Buffer).Error; err != nil {
		return err
	}
	obhp.Buffer = nil
	return nil
}

func (obhp *OrderBooksHistoryPersistence) GetBufferingSize() int {
	return len(obhp.Buffer)
}
