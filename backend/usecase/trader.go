package usecase

import (
	"fmt"
	"hamster/domain/model"
	"hamster/domain/repository"
	"hamster/domain/service"
	"hamster/lib/clog"
)

var cnt = 0

type TraderUsecase interface {
	ProcessNext() error
}

type traderUsecase struct {
	traderService            service.TraderService
	orderRepository          repository.OrderRepository
	orderBooksRepository     repository.OrderBooksRepository
	transactionRepository    repository.TransactionRepository
	tradeProcedureRepository repository.TradeProcedureRepository
}

func NewTraderUsecase(
	traderService service.TraderService,
	orderRepository repository.OrderRepository,
	orderBooksRepository repository.OrderBooksRepository,
	transactionRepository repository.TransactionRepository,
	tradeProcedureRepository repository.TradeProcedureRepository,
) TraderUsecase {
	return &traderUsecase{
		traderService:            traderService,
		orderRepository:          orderRepository,
		orderBooksRepository:     orderBooksRepository,
		transactionRepository:    transactionRepository,
		tradeProcedureRepository: tradeProcedureRepository,
	}
}

func (u *traderUsecase) ProcessNext() error {
	procedures := u.tradeProcedureRepository.GetActiveTradeProcedures()
	if len(procedures) > 0 {
		transactions, err := u.transactionRepository.SyncRecentTransactions()
		if err != nil {
			clog.Logger.Warn("Failed to sync recent transactions.")
			return err
		}

		if len(transactions) == 0 {
			return nil
		}

		for _, procedure := range procedures {
			completed := procedure.AddTransactions(transactions)
			if completed {
				clog.Logger.Info("procedure completed!")
			} else {
				//clog.Logger.Info("procedure still pending!")
			}
			_, err = u.tradeProcedureRepository.SaveTradeProcedure(procedure)
			if err != nil {
				return err
			}
		}
	} else {
		//if cnt > 0 {
		//	clog.Logger.Info("１回取引終わってます")
		//	return nil
		//}
		orderRequests := u.traderService.ShouldBuy()
		if orderRequests != nil {
			cnt++

			procedure, err := u.tradeProcedureRepository.CreateTradeProcedure(
				// TODO この段階で reason 入れておきたい
				fmt.Sprintf("price 上がりそう！"),
			)
			if err != nil {
				return err
			}

			for _, orderRequest := range orderRequests {
				requestedOrder, err := u.orderRepository.RequestOrder(orderRequest)
				if err != nil {
					procedure.Cancel(fmt.Sprintf("request failed: %s", err))
					u.tradeProcedureRepository.SaveTradeProcedure(procedure)
					return err
				}
				procedure.AddNewOrder(model.NewOrder(requestedOrder))
			}

			_, err = u.tradeProcedureRepository.SaveTradeProcedure(procedure)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
