package ad

import (
	"context"
	"errors"
)

type Repository interface {
	GetAll(context.Context) ([]Ad, error)
	GetByID(ctx context.Context, ID int) (Ad, error)
	Create(ctx context.Context, ad Ad) (Ad, error)
	Update(ctx context.Context, ad Ad, id int) (Ad, error)
	Delete(ctx context.Context, id int) error
}

var ErrNotFound = errors.New("ad not found")

var ErrInvalidAd = errors.New("invalid ad")
