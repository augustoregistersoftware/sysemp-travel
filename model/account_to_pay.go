package model

type AccountToPay struct {
	ID_USER int `json:"id_user" validate:"required,gt=0"`

	DESCRIPTION string `json:"description" validate:"required,min=3,max=255"`

	DESCRIPTION_DETAILS string `json:"description_details" validate:"max=1000"`

	DATE_ACTION string `json:"date_action" validate:"required,datetime=2006-01-02"`

	DATE_PREVIOUS string `json:"date_previous" validate:"required,datetime=2006-01-02"`

	VALUE_PAG float64 `json:"value_pag" validate:"required,gte=0"`

	VALUE_ADD float64 `json:"value_add" validate:"gte=0"`

	VALUE_DISCOUNT float64 `json:"value_discount" validate:"gte=0"`

	NAME_PAG string `json:"name_pag" validate:"required,min=3,max=255"`

	PAID bool `json:"paid"`
}

type FrankfurterRateResponse struct { //api de contação de moedas https://www.frankfurter.app/docs/
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Date  string  `json:"date"`
	Rate  float64 `json:"rate"`
}
