package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sysemp_travel/model"
	"time"

	"github.com/google/uuid"
)

type AccountToPayRepository struct {
	*Repository
}

func NewAccountToPayRepository(baseRepo *Repository) AccountToPayRepository {
	return AccountToPayRepository{
		Repository: baseRepo,
	}
}

func (act *AccountToPayRepository) NewAccountToPayInsert(ctx context.Context, accountToPay model.AccountToPay, idempotencyKey string) error {
	id_account_to_pay := uuid.NewString()
	date_action, err := time.Parse("2006-01-02", accountToPay.DATE_ACTION)
	if err != nil {
		return err
	}
	date_previous, err := time.Parse("2006-01-02", accountToPay.DATE_PREVIOUS)
	if err != nil {
		return err
	}

	_, err = act.DB.ExecContext(ctx, "INSERT INTO account_to_pay "+
		"(id_account_to_pay, id_user, description, description_details, "+
		"date_action,date_previous,value_pag,value_add,value_discount,"+
		"name_pag,paid)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		id_account_to_pay,
		accountToPay.ID_USER,
		accountToPay.DESCRIPTION,
		accountToPay.DESCRIPTION_DETAILS,
		date_action,
		date_previous,
		accountToPay.VALUE_PAG,
		accountToPay.VALUE_ADD,
		accountToPay.VALUE_DISCOUNT,
		accountToPay.NAME_PAG,
		accountToPay.PAID)
	if err != nil {
		return err
	}

	_, err = act.DB.ExecContext(ctx, "UPDATE account_to_pay_idempotency SET id_action = $2 WHERE idempotency_key = $1", idempotencyKey, id_account_to_pay)
	if err != nil {
		return err
	}
	return nil
}

func (act *AccountToPayRepository) GetFrankfurterRate(ctx context.Context, coin string, coin2 string) ([]model.FrankfurterRateResponse, error) {
	url := "https://api.frankfurter.dev/v2/rate/" + coin + "/" + coin2

	cliente := http.Client{Timeout: 10 * time.Second}

	resp, err := cliente.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Frankfurter rate: %d", resp.StatusCode)
	}

	var rateResponse model.FrankfurterRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&rateResponse); err != nil {
		return nil, err
	}
	// rateResponse.Rate = math.Round(rateResponse.Rate) caso tenha a necessidade de arredondar o valor da taxa para um número inteiro

	return []model.FrankfurterRateResponse{rateResponse}, nil
}

func (act *AccountToPayRepository) CheckIdempotency(ctx context.Context, idempotencyKey string) error {
	var existingKey string
	err := act.DB.QueryRowContext(ctx, "SELECT idempotency_key FROM account_to_pay_idempotency WHERE idempotency_key = $1", idempotencyKey).Scan(&existingKey)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := act.DB.ExecContext(ctx, "INSERT INTO account_to_pay_idempotency (idempotency_key) VALUES ($1)", idempotencyKey)
			return err
		}
		return err
	}
	return fmt.Errorf("idempotency key already exists")
}
