package food

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, food *Food) error
	FindAll(ctx context.Context, sortType string) (u []Food, err error)
	Like(ctx context.Context, prodId int, userId int) (bool, error)
	FindOne(ctx context.Context, id int) (Food, error)
	Update(ctx context.Context, fd Food) error
	Delete(ctx context.Context, id int) error
}
