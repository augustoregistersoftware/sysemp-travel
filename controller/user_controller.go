package controller

import (
	"fmt"
	"net/http"
	"sysemp_feed/model"
	"sysemp_feed/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUseCase usecase.UserUseCase
}

func NewUserController(userUseCase usecase.UserUseCase) UserController {
	return UserController{
		UserUseCase: userUseCase,
	}
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	var user model.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err = u.UserUseCase.CreateUser(ctx.Request.Context(), user)
	fmt.Println("Erro", err)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "User created successfully"})
}
