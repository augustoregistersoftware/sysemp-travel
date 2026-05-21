package usecase

import (
	"context"
	"fmt"
	"strings"
	"sysemp_feed/config"
	"sysemp_feed/repository"
)

type AuthUseCase struct {
	userRepo *repository.UserRepository
}

func NewAuthUsecase(userRepo *repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

func (a *AuthUseCase) ValidateCredentials(ctx context.Context, userOrEmail, password string) (string, bool, error) {
	userOrEmail = strings.TrimSpace(userOrEmail)

	user, err := a.userRepo.FindByUsername(ctx, userOrEmail)
	if err != nil {
		return "", false, err
	}

	if user == nil {
		return "", false, nil
	}

	validApprovedUser, err := a.userRepo.IsApprovedUser(ctx, user.ID)
	if err != nil {
		return "", false, err
	}

	if !validApprovedUser {
		return "", false, fmt.Errorf("Usuario não aprovado")
	}

	validPassword, err := config.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return "", false, err
	}

	if !validPassword {
		return "", false, nil
	}

	return user.Username, true, nil
}
