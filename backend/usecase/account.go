package usecase

import (
	"hamster/domain/model"
	"hamster/domain/repository"
)

type AccountUsecase interface {
	GetBalance() (*model.Balance, error)
}

type AccountUsecaseImpl struct {
	BalanceRepository repository.BalanceRepository
}

func NewAccountUsecase(balanceRepo repository.BalanceRepository) AccountUsecase {
	return &AccountUsecaseImpl{BalanceRepository: balanceRepo}
}

func (u *AccountUsecaseImpl) GetBalance() (*model.Balance, error) {
	return u.BalanceRepository.GetBalance()
}
