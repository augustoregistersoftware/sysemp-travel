package usecase

import (
	"context"
	"sysemp_travel/model"
	"sysemp_travel/publisher"
	"sysemp_travel/repository"
	"time"

	"github.com/go-redis/redis/v8"
)

type PaymentsUseCase struct {
	repository repository.PaymentsRepository
}

type PaymentsPayUsecase struct {
	repository repository.PaymentsRepository
	publisher  *publisher.PaymentPublisher
	cache      *redis.Client
	cacheTTL   time.Duration
}

func NewPaymentsUseCase(paymentsRepo repository.PaymentsRepository) PaymentsUseCase {
	return PaymentsUseCase{
		repository: paymentsRepo,
	}
}

func NewPaymentsPayUseCase(
	repo repository.PaymentsRepository,
	pub *publisher.PaymentPublisher,
	cache *redis.Client,
) PaymentsPayUsecase {
	return PaymentsPayUsecase{
		repository: repo,
		publisher:  pub,
		cache:      cache,
		cacheTTL:   5 * time.Minute,
	}
}

func (u *PaymentsUseCase) GetPayments(ctx context.Context) ([]model.Payment, error) {
	payment, err := u.repository.GetPayments(ctx)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentsPayUsecase) Pay(ctx context.Context, pay model.Pay) error {
	err := u.repository.ProcessPayment(ctx, pay)
	if err != nil {
		return err
	}

	err = u.publisher.PublishPaymentCreated(pay.AccountToPayID)
	if err != nil {
		return err
	}

	// Cache the payment result in Redis

	return nil
}
