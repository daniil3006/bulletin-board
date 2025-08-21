package user

import (
	"bulletin-board/internal/ad"
	"context"
	"errors"
)

type Repository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (User, error)
	GetUsersAds(ctx context.Context, userId int) ([]ad.Ad, error)
	Create(ctx context.Context, newUser User) (User, error)
	Update(ctx context.Context, user User, id int) (User, error)
	Delete(ctx context.Context, id int) error
}

var ErrUserNotFound = errors.New("user not found")
