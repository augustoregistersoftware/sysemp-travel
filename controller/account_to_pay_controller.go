package controller

import (
	"net/http"
	"sysemp_travel/usecase"

	"sysemp_travel/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type accountToPayController struct {
	AccountToPayUseCase usecase.AccountToPayUseCase
}

func NewAccountToPayController(accountToPayUseCase usecase.AccountToPayUseCase) accountToPayController {
	return accountToPayController{
		AccountToPayUseCase: accountToPayUseCase,
	}
}

func (c *accountToPayController) CreateAccountToPay(ctx *gin.Context) {
	var accountToPay model.AccountToPay

	typ := ctx.Param("type")

	err := ctx.BindJSON(&accountToPay)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(accountToPay)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idempotencyKey := ctx.Request.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key required"})
		return
	}

	err = c.AccountToPayUseCase.CheckIdempotency(ctx.Request.Context(), idempotencyKey)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Duplicate request"})
		return
	}

	err = c.AccountToPayUseCase.CreateAccountToPay(ctx.Request.Context(), typ, accountToPay, idempotencyKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Account to pay created successfully"})
}

func (c *accountToPayController) GetFrankfurterRate(ctx *gin.Context) {
	coin := ctx.Param("coin")
	if coin == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "coin query parameter is required"})
		return
	}
	coin2 := ctx.Param("coin2")
	if coin2 == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "coin2 query parameter is required"})
		return
	}

	response, err := c.AccountToPayUseCase.GetFrankfurterRate(ctx.Request.Context(), coin, coin2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"rate": response})
}
