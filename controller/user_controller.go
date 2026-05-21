package controller

import (
	"net/http"
	"sysemp_travel/model"
	"sysemp_travel/usecase"

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

func (u *UserController) ApproveUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := u.UserUseCase.ApproveUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "User approved successfully"})
}

func (u *UserController) ReproveUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := u.UserUseCase.ReproveUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "User reproved successfully"})
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	var user model.User

	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	_, err = u.UserUseCase.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "User created successfully"})
}
