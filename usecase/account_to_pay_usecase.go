package usecase

import (
	"context"
	"sysemp_travel/model"
	"sysemp_travel/repository"
)

type AccountToPayUseCase struct {
	repository repository.AccountToPayRepository
}

func NewAccountToPayUseCase(accountToPayRepo repository.AccountToPayRepository) AccountToPayUseCase {
	return AccountToPayUseCase{
		repository: accountToPayRepo,
	}
}

func (u *AccountToPayUseCase) CreateAccountToPay(ctx context.Context, typ string, accountToPay model.AccountToPay) error {
	if typ == "0" {
		accountToPay.DESCRIPTION_DETAILS = accountToPay.DESCRIPTION_DETAILS + " - foreign payment"
		println("Description details " + accountToPay.DESCRIPTION_DETAILS)
		return u.repository.NewAccountToPayInsert(ctx, accountToPay)
	} else {
		return nil
	}
}

func (u *AccountToPayUseCase) GetFrankfurterRate(ctx context.Context, coin string, coin2 string) ([]model.FrankfurterRateResponse, error) {
	response, err := u.repository.GetFrankfurterRate(ctx, coin, coin2)
	if err != nil {
		return []model.FrankfurterRateResponse{}, err
	}
	return response, nil
}
