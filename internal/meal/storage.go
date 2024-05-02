package meal

import (
	"CalorieGuide-db/internal/food"
	"context"
)

type Repository interface {
	Create(ctx context.Context, meal *Meal) error
	FindAll(ctx context.Context, sortType string, twoDecade int) (u []Meal, err error)
	Like(ctx context.Context, mealId int, userId int) (bool, error)
	FindOne(ctx context.Context, id int) (Meal, error)
	Update(ctx context.Context, fd *Meal) error
	Delete(ctx context.Context, id int) error
	AddProduct(ctx context.Context, id int, product *food.Food, quantity int) error
}
