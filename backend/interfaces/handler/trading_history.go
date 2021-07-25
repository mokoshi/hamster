package handler

import (
	"fmt"
	"hamster/usecase"
	"net/http"
	"time"
)

type TradingHistoryHandler interface {
	HandleTrade(http.ResponseWriter, *http.Request)
}

type TradingHistoryHandlerImpl struct {
	TradingHistoryUsecase usecase.TradingHistoryUsecase
}

func NewTradingHistoryHandler(u usecase.TradingHistoryUsecase) TradingHistoryHandler {
	return &TradingHistoryHandlerImpl{
		TradingHistoryUsecase: u,
	}
}

func (u TradingHistoryHandlerImpl) HandleTrade(w http.ResponseWriter, r *http.Request) {
	tradingHistory, err := u.TradingHistoryUsecase.Create(time.Now())
	if err != nil {
		fmt.Fprintf(w, "<h1>hamster neko pokemon 3<h1>")
		return
	}

	fmt.Fprintf(w, "<h1>id: %d<h1>", tradingHistory.ID)
}
