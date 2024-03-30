package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"github.com/jackc/pgconn"
	"log/slog"
)

type repository struct {
	client postgreSQL.Client
	log    *slog.Logger
}

func (r *repository) Create(ctx context.Context, food *food.Food) error {
	q := `
		INSERT INTO food (food_name, description, calories, proteins, carbohydrates, fats)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	if err := r.client.QueryRow(
		ctx, q, food.Name, food.Description, food.Calories, food.Proteins, food.Carbohydrates, food.Fats,
	).Scan(&food.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return nil
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) (u []food.Food, err error) {
	q := `
		SELECT 
		    id, food_name, description, calories, proteins, carbohydrates, fats, picture
		FROM public.food;
	`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	allFood := make([]food.Food, 0)
	for rows.Next() {
		var fd food.Food
		err = rows.Scan(
			&fd.Id, &fd.Name, &fd.Description, &fd.Calories, &fd.Proteins, &fd.Carbohydrates, &fd.Fats, &fd.Picture,
		)
		if err != nil {
			return nil, err
		}
		allFood = append(allFood, fd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return allFood, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (food.Food, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, food food.Food) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) food.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
