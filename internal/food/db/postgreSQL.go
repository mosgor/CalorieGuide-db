package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"log/slog"
)

type repository struct {
	client postgreSQL.Client
	log    *slog.Logger
}

func (r *repository) Create(ctx context.Context, food *food.Food) error {
	q := `
		INSERT INTO food (food_name, description, calories, proteins, carbohydrates, fats, likes, author_id, picture)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	if err := r.client.QueryRow(
		ctx, q, food.Name,
		food.Description, food.Calories,
		food.Proteins, food.Carbohydrates,
		food.Fats, food.Likes,
		food.AuthorId, food.Picture,
	).Scan(&food.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return err
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, sortType string) (u []food.Food, err error) {
	q := `
		SELECT 
		    id, food_name, description, calories, proteins, carbohydrates, fats, author_id, likes, picture
		FROM public.food
	`
	switch sortType {
	case "likesAsc":
		q += `ORDER BY likes ASC;`
	case "likesDesc":
		q += `ORDER BY likes DESC;`
	case "fromOldest":
		q += `ORDER BY id ASC;`
	case "fromNewest":
		fallthrough
	default:
		q += `ORDER BY id DESC;`
	}
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	allFood := make([]food.Food, 0)
	for rows.Next() {
		var fd food.Food
		err = rows.Scan(
			&fd.Id, &fd.Name, &fd.Description,
			&fd.Calories, &fd.Proteins,
			&fd.Carbohydrates, &fd.Fats,
			&fd.AuthorId, &fd.Likes, &fd.Picture,
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

func (r *repository) FindOne(ctx context.Context, id int) (fd food.Food, err error) {
	q := `
	SELECT 
	    id, food_name, description, calories, proteins, carbohydrates, fats, author_id, likes, picture 
	FROM public.food WHERE id = $1;
	`
	rw := r.client.QueryRow(ctx, q, id)
	if err = rw.Scan(
		&fd.Id, &fd.Name,
		&fd.Description, &fd.Calories,
		&fd.Proteins, &fd.Carbohydrates,
		&fd.Fats, &fd.AuthorId,
		&fd.Likes, &fd.Picture,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return fd, err
		}
		return fd, err
	}
	return
}

func (r *repository) Update(ctx context.Context, fd food.Food) error {
	q := `
	UPDATE public.food SET
		food_name=$2, description=$3, 
		calories=$4, proteins = $5, 
		carbohydrates = $6, fats = $7,
		picture = $8
	WHERE id = $1;
	`
	r.client.QueryRow(ctx, q,
		&fd.Id, &fd.Name,
		&fd.Description, &fd.Calories,
		&fd.Proteins, &fd.Carbohydrates,
		&fd.Fats, &fd.Picture,
	)
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) (err error) {
	q := `DELETE FROM public.food_client WHERE food_id = $1`
	r.client.QueryRow(ctx, q, id)
	q = `DELETE FROM public.meal_food WHERE food_id = $1`
	r.client.QueryRow(ctx, q, id)
	q = `DELETE FROM public.food WHERE id = $1;`
	r.client.QueryRow(ctx, q, id)
	return
}

func (r *repository) Like(ctx context.Context, prodId int, userId int) (liked bool, err error) {
	q := `
	SELECT EXISTS (SELECT 1 FROM public.food_client WHERE food_id = $1 AND user_id = $2)
	`
	var exists bool
	rw := r.client.QueryRow(ctx, q, prodId, userId)
	err = rw.Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		q = `
		DELETE FROM public.food_client WHERE food_id = $1 AND user_id = $2
		`
		r.client.QueryRow(ctx, q, prodId, userId)
		liked = false
		q = `
		UPDATE public.food SET likes = likes - 1 WHERE id = $1
		`
		r.client.QueryRow(ctx, q, prodId)
	} else {
		q = `
		INSERT INTO food_client (food_id, user_id) VALUES ($1, $2)
		`
		r.client.QueryRow(ctx, q, prodId, userId)
		liked = true
		q = `
		UPDATE public.food SET likes = likes + 1 WHERE id = $1
		`
		r.client.QueryRow(ctx, q, prodId)
	}
	return
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) food.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
