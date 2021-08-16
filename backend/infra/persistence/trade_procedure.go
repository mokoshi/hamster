package persistence

import (
	"gorm.io/gorm"
	"hamster/domain/model"
	"hamster/domain/repository"
)

type TradeProcedureRepository struct {
	db                    *gorm.DB
	activeTradeProcedures []*model.TradeProcedure
}

func NewTradeProcedureRepository(db *gorm.DB) repository.TradeProcedureRepository {
	return &TradeProcedureRepository{
		db: db,
	}
}

func (r *TradeProcedureRepository) CreateTradeProcedure(reason string) (*model.TradeProcedure, error) {
	procedure := model.NewTradeProcedure(nil, reason)

	if err := r.db.Create(procedure).Error; err != nil {
		return nil, err
	}

	r.activeTradeProcedures = append(r.activeTradeProcedures, procedure)

	return procedure, nil
}

func (r *TradeProcedureRepository) SaveTradeProcedure(procedure *model.TradeProcedure) (*model.TradeProcedure, error) {
	if err := r.db.Save(procedure).Error; err != nil {
		return nil, err
	}

	for i, activeProcedure := range r.activeTradeProcedures {
		if activeProcedure.Id == procedure.Id {
			if procedure.IsPending() {
				r.activeTradeProcedures[i] = procedure
			} else {
				r.activeTradeProcedures[i] = r.activeTradeProcedures[len(r.activeTradeProcedures)-1]
				r.activeTradeProcedures = r.activeTradeProcedures[:len(r.activeTradeProcedures)-1]
			}
			break
		}
	}

	return procedure, nil
}

func (r *TradeProcedureRepository) SyncByOpenOrders(openOrders []*model.OpenOrder) ([]*model.TradeProcedure, error) {
	return r.activeTradeProcedures, nil
}

func (r *TradeProcedureRepository) GetActiveTradeProcedures() []*model.TradeProcedure {
	return r.activeTradeProcedures
}
