package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

func (r *repository) FindAll(ctx context.Context, sortType string, twoDecade int, userId int) (u []food.WithLike, err error) {
	q := `
		SELECT 
		    id, food_name, description, calories, proteins, carbohydrates, fats, author_id, likes, picture
		FROM public.food
	`
	switch sortType {
	case "likesAsc":
		q += `ORDER BY likes ASC`
	case "likesDesc":
		q += `ORDER BY likes DESC`
	case "fromOldest":
		q += `ORDER BY id ASC`
	case "fromNewest":
		fallthrough
	default:
		q += `ORDER BY id DESC`
	}
	q += ` OFFSET $1 LIMIT 20;`
	rows, err := r.client.Query(ctx, q, (twoDecade-1)*20)
	if err != nil {
		return nil, err
	}
	allFood := make([]food.WithLike, 0)
	for rows.Next() {
		var fd food.WithLike
		err = rows.Scan(
			&fd.Id, &fd.Name, &fd.Description,
			&fd.Calories, &fd.Proteins,
			&fd.Carbohydrates, &fd.Fats,
			&fd.AuthorId, &fd.Likes, &fd.Picture,
		)
		if err != nil {
			return nil, err
		}
		if userId != 0 {
			q = `SELECT EXISTS (SELECT 1 FROM food_client WHERE user_id = $1 AND food_id  = $2)`
			rw, err := r.client.Query(ctx, q, userId, fd.Id)
			if err != nil {
				return nil, err
			}
			rw.Next()
			err = rw.Scan(&fd.Like)
			if err != nil {
				return nil, err
			}
			rw.Close()
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
	rw, err := r.client.Query(ctx, q,
		&fd.Id, &fd.Name,
		&fd.Description, &fd.Calories,
		&fd.Proteins, &fd.Carbohydrates,
		&fd.Fats, &fd.Picture,
	)
	if err != nil {
		return err
	}
	defer rw.Close()
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) (err error) {
	q := `DELETE FROM public.food_client WHERE food_id = $1`
	rw, err := r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM public.meal_food WHERE food_id = $1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM public.food WHERE id = $1;`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	return
}

func (r *repository) Like(ctx context.Context, prodId int, userId int) (liked bool, err error) {
	q := `
	SELECT EXISTS (SELECT 1 FROM public.food_client WHERE food_id = $1 AND user_id = $2)
	`
	var exists bool
	rw, err := r.client.Query(ctx, q, prodId, userId)
	if err != nil {
		return false, err
	}
	rw.Next()
	defer rw.Close()
	err = rw.Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		q = `
		DELETE FROM public.food_client WHERE food_id = $1 AND user_id = $2
		`
		rw, err = r.client.Query(ctx, q, prodId, userId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
		liked = false
		q = `
		UPDATE public.food SET likes = likes - 1 WHERE id = $1
		`
		rw, err = r.client.Query(ctx, q, prodId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
	} else {
		q = `
		INSERT INTO food_client (food_id, user_id) VALUES ($1, $2)
		`
		rw, err = r.client.Query(ctx, q, prodId, userId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
		liked = true
		q = `
		UPDATE public.food SET likes = likes + 1 WHERE id = $1
		`
		rw, err = r.client.Query(ctx, q, prodId)
		if err != nil {
			return false, err
		}
		defer rw.Close()
	}
	return
}

func (r *repository) Search(ctx context.Context, word string, userId int) (fd []food.WithLike, err error) {
	includedIds := make(map[int]struct{})
	q := `
		SELECT * FROM public.food WHERE food_name ILIKE CONCAT('%',$1::text,'%') ORDER BY likes DESC;
	`
	rows, err := r.client.Query(ctx, q, word)
	if err != nil {
		r.log.Error("Error searching food by name")
		return nil, err
	}
	defer rows.Close()
	q = `
		SELECT * FROM public.food WHERE description ILIKE CONCAT('%',$1::text,'%') ORDER BY likes DESC;
	`
	rows1, err := r.client.Query(ctx, q, word)
	if err != nil {
		r.log.Error("Error searching food by description")
		return nil, err
	}
	defer rows1.Close()
	for _, rw := range []pgx.Rows{rows, rows1} {
		for rw.Next() {
			var product food.WithLike
			err = rw.Scan(
				&product.Id, &product.Name, &product.Description,
				&product.Calories, &product.Proteins,
				&product.Carbohydrates, &product.Fats,
				&product.AuthorId, &product.Likes, &product.Picture,
			)
			if err != nil {
				return nil, err
			}
			_, ok := includedIds[product.Id]
			if !ok {
				if userId != 0 {
					q = `SELECT EXISTS (SELECT 1 FROM food_client WHERE user_id = $1 AND food_id  = $2)`
					rw, err := r.client.Query(ctx, q, userId, product.Id)
					if err != nil {
						return nil, err
					}
					rw.Next()
					err = rw.Scan(&product.Like)
					if err != nil {
						return nil, err
					}
					rw.Close()
				}
				fd = append(fd, product)
				includedIds[product.Id] = struct{}{}
			}
		}
	}
	return fd, nil
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) food.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
