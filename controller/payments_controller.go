package controller

import (
	"net/http"
	"sysemp_travel/model"
	"sysemp_travel/usecase"

	"github.com/gin-gonic/gin"
)

type paymentsController struct {
	PaymentsUseCase usecase.PaymentsUseCase
}

type paymentspayController struct {
	PaymentsPayUseCase usecase.PaymentsPayUsecase
}

func NewPaymentsController(paymentsUseCase usecase.PaymentsUseCase) paymentsController {
	return paymentsController{
		PaymentsUseCase: paymentsUseCase,
	}
}

func NewPaymentsPayController(paymentsPayUseCase usecase.PaymentsPayUsecase) paymentspayController {
	return paymentspayController{
		PaymentsPayUseCase: paymentsPayUseCase,
	}
}

func (c *paymentsController) Payments(ctx *gin.Context) {
	payments, err := c.PaymentsUseCase.GetPayments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payments == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"payments": payments})
}

func (c *paymentspayController) Pay(ctx *gin.Context) {
	var pay model.Pay
	if err := ctx.ShouldBindJSON(&pay); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.PaymentsPayUseCase.Pay(ctx.Request.Context(), pay)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
