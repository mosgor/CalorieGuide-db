package food

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, food *Food) error
	FindAll(ctx context.Context, sortType string) (u []Food, err error)
	FindOne(ctx context.Context, id string) (Food, error)
	Update(ctx context.Context, food Food) error
	Delete(ctx context.Context, id string) error
}
