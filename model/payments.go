package model

type Payment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Pay struct {
	AccountToPayID string `json:"account_to_pay_id" binding:"required"`
	PaymentID      int    `json:"payment_id" binding:"required"`
}
