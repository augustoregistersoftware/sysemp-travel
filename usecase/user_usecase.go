package usecase

import (
	"context"
	"sysemp_travel/model"
	"sysemp_travel/repository"
)

type UserUseCase struct {
	repository repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return UserUseCase{
		repository: userRepo,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	_, err := u.repository.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *UserUseCase) ApproveUser(ctx context.Context, id string) error {
	return u.repository.ApproveUser(ctx, id)
}

func (u *UserUseCase) ReproveUser(ctx context.Context, id string) error {
	return u.repository.ReproveUser(ctx, id)
}
