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

func (u *AccountToPayUseCase) CreateAccountToPay(ctx context.Context, typ string, accountToPay model.AccountToPay, idempotencyKey string) error {
	if typ == "0" {
		return u.repository.NewAccountToPayInsert(ctx, accountToPay, idempotencyKey)
	} else {
		accountToPay.DESCRIPTION_DETAILS = accountToPay.DESCRIPTION_DETAILS + " - foreign payment"
		return u.repository.NewAccountToPayInsert(ctx, accountToPay, idempotencyKey)
	}
}

func (u *AccountToPayUseCase) GetFrankfurterRate(ctx context.Context, coin string, coin2 string) ([]model.FrankfurterRateResponse, error) {
	response, err := u.repository.GetFrankfurterRate(ctx, coin, coin2)
	if err != nil {
		return []model.FrankfurterRateResponse{}, err
	}
	return response, nil
}

func (u *AccountToPayUseCase) CheckIdempotency(ctx context.Context, idempotencyKey string) error {
	return u.repository.CheckIdempotency(ctx, idempotencyKey)
}
