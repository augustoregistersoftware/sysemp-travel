package controller

import (
	"net/http"
	"sysemp_feed/auth"
	"sysemp_feed/usecase"

	"github.com/gin-gonic/gin"
)

type authController struct {
	authService *auth.Service
	authUseCase *usecase.AuthUseCase
}

type loginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func NewAuthController(
	authService *auth.Service,
	authUseCase *usecase.AuthUseCase,
) *authController {
	return &authController{
		authService: authService,
		authUseCase: authUseCase,
	}
}

func (a *authController) Login(ctx *gin.Context) {
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body invalid"})
		return
	}

	username, valid, err := a.authUseCase.ValidateCredentials(
		ctx.Request.Context(),
		req.User,
		req.Password,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := a.authService.GenerateToken(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   a.authService.ExpiresInSeconds(),
	})
}
