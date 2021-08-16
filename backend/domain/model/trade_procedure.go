package model

import (
	"github.com/google/uuid"
	"time"
)

type TradeProcedure struct {
	Id         string
	StartedAt  *time.Time
	FinishedAt *time.Time
	Status     string
	Reason     string
	Orders     []*Order
}

func NewTradeProcedure(orderRequests []*OrderRequest, reason string) *TradeProcedure {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	procedure := &TradeProcedure{
		Id:        id.String(),
		StartedAt: &now,
		Status:    "pending",
		Reason:    reason,
	}
	for _, req := range orderRequests {
		procedure.AddNewOrder(NewOrder(req))
	}
	return procedure
}

func (p *TradeProcedure) AddNewOrder(order *Order) {
	p.Orders = append(p.Orders, order)
}

func (p *TradeProcedure) AddTransactions(transactions []*Transaction) bool {
	orderIdMap := map[uint64][]*Transaction{}
	for _, t := range transactions {
		orderIdMap[t.OrderId] = append(orderIdMap[t.OrderId], t)
	}

	procedureCompleted := true

	for _, order := range p.Orders {
		if !order.IsPending() {
			continue
		}

		orderCompleted := false
		transactions, ok := orderIdMap[order.Id]
		if ok {
			orderCompleted = order.AddTransactions(transactions)
		}

		if !orderCompleted {
			procedureCompleted = false
		}
	}

	if procedureCompleted {
		p.SetAsCompleted()
	}

	return procedureCompleted
}

func (p *TradeProcedure) IsPending() bool {
	if p.Status == "pending" {
		return true
	} else {
		return false
	}
}

func (p *TradeProcedure) SetAsCompleted() {
	p.Status = "completed"
}

func (p *TradeProcedure) Cancel(reason string) {
	p.Reason = reason
	for _, order := range p.Orders {
		if order.IsPending() {
			order.SetAsCanceled(reason)
		}
	}
	p.Status = "canceled"
}
