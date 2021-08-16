package repository

import (
	"hamster/domain/model"
)

type TradeProcedureRepository interface {
	CreateTradeProcedure(reason string) (*model.TradeProcedure, error)
	SaveTradeProcedure(procedure *model.TradeProcedure) (*model.TradeProcedure, error)
	SyncByOpenOrders(openOrders []*model.OpenOrder) ([]*model.TradeProcedure, error)
	GetActiveTradeProcedures() []*model.TradeProcedure
}
