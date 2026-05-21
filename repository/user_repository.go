package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sysemp_feed/config"
	"sysemp_feed/model"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type UserRepository struct {
	*Repository
}

func NewUserRepository(baseRepo *Repository) UserRepository {
	return UserRepository{
		Repository: baseRepo,
	}
}

func (r *UserRepository) FindByUsername(ctx context.Context, userOrEmail string) (*User, error) {
	var user User
	err := r.DB.QueryRowContext(ctx, "SELECT id_user, username, password FROM users WHERE username = $1 OR email = $1",
		userOrEmail).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsApprovedUser(ctx context.Context, id int64) (bool, error) {
	var id_user int64
	err := r.DB.QueryRowContext(ctx, "SELECT id_user FROM approved_users WHERE id_user = $1", id).Scan(&id_user)

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (r *UserRepository) CreateUserApprove(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(
		ctx,
		"INSERT INTO approved_users (id_user) VALUES ($1)",
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user model.User) (int, error) {
	var id_user int

	retornoUser, err := ur.FindByUsername(ctx, user.Username)
	if err != nil {
		return 0, err
	}

	if retornoUser != nil {
		return 409, fmt.Errorf("user already exists")
	}

	passwordEncrypted, err := config.HashPassword(user.Password)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = ur.DB.QueryRowContext(
		ctx,
		"INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id_user",
		user.Email,
		user.Username,
		passwordEncrypted,
	).Scan(&id_user)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = ur.CreateUserApprove(ctx, int64(id_user))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id_user, nil
}

func (ur *UserRepository) ApproveUser(ctx context.Context, id string) error {
	_, err := ur.DB.ExecContext(
		ctx,
		"DELETE FROM approved_users WHERE id_user = $1",
		id,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
