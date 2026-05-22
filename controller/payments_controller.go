package controller

import (
	"net/http"
	"sysemp_travel/usecase"

	"github.com/gin-gonic/gin"
)

type paymentsController struct {
	PaymentsUseCase usecase.PaymentsUseCase
}

func NewPaymentsController(paymentsUseCase usecase.PaymentsUseCase) paymentsController {
	return paymentsController{
		PaymentsUseCase: paymentsUseCase,
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
