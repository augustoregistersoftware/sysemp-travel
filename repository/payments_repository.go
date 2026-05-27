package repository

import (
	"context"
	"sysemp_travel/model"

	"github.com/google/uuid"
)

type PaymentsRepository struct {
	*Repository
}

func NewPaymentsRepository(baseRepo *Repository) PaymentsRepository {
	return PaymentsRepository{
		Repository: baseRepo,
	}
}

func (r *PaymentsRepository) GetPayments(ctx context.Context) ([]model.Payment, error) {
	var payment []model.Payment
	rows, err := r.DB.QueryContext(ctx, "SELECT id_payments, name FROM payments")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		payment = append(payment, p)
	}

	return payment, nil
}

func (r *PaymentsRepository) ProcessPayment(ctx context.Context, pay model.Pay) error {
	uuid := uuid.NewString()
	_, err := r.DB.ExecContext(ctx, "INSERT INTO account_to_pay_payments (id_account_to_pay, id_payments, id_account_to_pay_payments) VALUES ($1, $2, $3)", pay.AccountToPayID, pay.PaymentID, uuid)

	return err
}
