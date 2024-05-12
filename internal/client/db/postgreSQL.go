package client

import (
	"CalorieGuide-db/internal/client"
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/meal"
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

func (r *repository) Create(ctx context.Context, client *client.Client) error {
	q := `INSERT INTO goal DEFAULT VALUES`
	rw, err := r.client.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `INSERT INTO diet DEFAULT VALUES`
	rw, err = r.client.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `
		INSERT INTO client (user_name, surname, email, password, picture)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	rw, err = r.client.Query(ctx, q,
		client.Name, client.Surname,
		client.Email, client.Password,
		client.Picture,
	)
	if err != nil {
		return err
	}
	defer rw.Close()
	rw.Next()
	if err = rw.Scan(&client.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return err
		}
		return err
	}
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (cl client.Client, err error) {
	q := `
	SELECT id, user_name, surname, email, password, picture
	FROM public.client
	WHERE email = $1
	`
	rw := r.client.QueryRow(ctx, q, email)
	if err = rw.Scan(
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password, &cl.Picture,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error in email", slg.PgErr(*pgErr))
			return cl, err
		}
		return cl, err
	}
	return
}

func (r *repository) FindGoalById(ctx context.Context, goalId int) (goal client.Goal, err error) {
	q := `
	SELECT calories_goal, fats_goal, proteins_goal, carbohydrates_goal FROM public.goal
	WHERE id = $1
	`
	rw := r.client.QueryRow(ctx, q, goalId)
	if err = rw.Scan(
		&goal.CaloriesGoal, &goal.FatsGoal,
		&goal.ProteinsGoal, &goal.CarbohydratesGoal,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error in goal", slg.PgErr(*pgErr))
			return goal, err
		}
		return goal, err
	}
	return
}

func (r *repository) FindDietById(ctx context.Context, dietId int) (diet client.Diet, err error) {
	q := `SELECT CASE WHEN breakfast_id IS NULL THEN 0 ELSE breakfast_id END FROM public.diet WHERE id = $1`
	rw := r.client.QueryRow(ctx, q, dietId)
	if err = rw.Scan(&diet.BreakfastId); err != nil {
		return diet, err
	}
	q = `SELECT CASE WHEN lunch_id IS NULL THEN 0 ELSE lunch_id END FROM public.diet WHERE id = $1`
	rw = r.client.QueryRow(ctx, q, dietId)
	if err = rw.Scan(&diet.LunchId); err != nil {
		return
	}
	q = `SELECT CASE WHEN dinner_id IS NULL THEN 0 ELSE dinner_id END FROM public.diet WHERE id = $1`
	rw = r.client.QueryRow(ctx, q, dietId)
	if err = rw.Scan(&diet.DinnerId); err != nil {
		return
	}
	return
}

func (r *repository) FindById(ctx context.Context, id int) (cl client.Client, err error) {
	q := `
	SELECT id, user_name, surname, email, password, picture
	FROM public.client
	WHERE id = $1
	`
	rw := r.client.QueryRow(ctx, q, id)
	if err = rw.Scan(
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password, &cl.Picture,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return cl, err
		}
		return cl, err
	}
	return
}

func (r *repository) UpdateClient(ctx context.Context, cl client.Client) error {
	q := `
	UPDATE public.client SET
		user_name=$2, surname=$3,
		email=$4, password=$5, picture=$6
	WHERE id = $1;
	`
	rw, err := r.client.Query(ctx, q,
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password, &cl.Picture,
	)
	if err != nil {
		return err
	}
	defer rw.Close()
	return nil
}

func (r *repository) UpdateGoal(ctx context.Context, goal client.Goal, goalId int) error {
	q := `
	UPDATE public.goal SET
		calories_goal=$2, fats_goal=$3,
		proteins_goal=$4, carbohydrates_goal=$5
	WHERE id = $1;
	`
	rw, err := r.client.Query(ctx, q,
		goalId, &goal.CaloriesGoal,
		&goal.FatsGoal, &goal.ProteinsGoal,
		&goal.CarbohydratesGoal,
	)
	if err != nil {
		return err
	}
	defer rw.Close()
	return nil
}

func (r *repository) UpdateDiet(ctx context.Context, diet client.Diet, dietId int) error {
	q := `UPDATE public.diet SET breakfast_id=$2 WHERE id = $1;`
	if diet.BreakfastId != 0 {
		rw, err := r.client.Query(ctx, q, dietId, &diet.BreakfastId)
		if err != nil {
			return err
		}
		defer rw.Close()
	} else {
		rw, err := r.client.Query(ctx, q, dietId, nil)
		if err != nil {
			return err
		}
		defer rw.Close()
	}
	q = `UPDATE public.diet SET lunch_id=$2 WHERE id = $1;`
	if diet.LunchId != 0 {
		rw, err := r.client.Query(ctx, q, dietId, &diet.LunchId)
		if err != nil {
			return err
		}
		defer rw.Close()
	} else {
		rw, err := r.client.Query(ctx, q, dietId, nil)
		if err != nil {
			return err
		}
		defer rw.Close()
	}
	q = `UPDATE public.diet SET dinner_id=$2 WHERE id = $1;`
	if diet.DinnerId != 0 {
		rw, err := r.client.Query(ctx, q, dietId, &diet.DinnerId)
		if err != nil {
			return err
		}
		defer rw.Close()
	} else {
		rw, err := r.client.Query(ctx, q, dietId, nil)
		if err != nil {
			return err
		}
		defer rw.Close()
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int, fdRepo food.Repository, mealRepo meal.Repository) error {
	q := `DELETE FROM public.food_client WHERE user_id=$1`
	rw, err := r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM public.meal_client WHERE user_id=$1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `SELECT id FROM public.food WHERE author_id=$1`
	rows, err := r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		foodId := 0
		err = rows.Scan(&foodId)
		if err != nil {
			return err
		}
		err = fdRepo.Delete(ctx, foodId)
		if err != nil {
			return err
		}
	}
	q = `SELECT id FROM public.meal WHERE author_id=$1`
	rows, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		mealId := 0
		err = rows.Scan(&mealId)
		if err != nil {
			return err
		}
		err = mealRepo.Delete(ctx, mealId)
		if err != nil {
			return err
		}
	}
	q = `DELETE FROM public.goal WHERE id=$1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM public.diet WHERE id=$1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	q = `DELETE FROM public.client WHERE id=$1`
	rw, err = r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
	defer rw.Close()
	return nil
}

func (r *repository) FindMealLikes(ctx context.Context, id int) (meals []meal.WithLike, err error) {
	q := `SELECT meal_id FROM public.meal_client WHERE user_id=$1`
	rows, err := r.client.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mealId int
		err = rows.Scan(&mealId)
		if err != nil {
			return nil, err
		}
		q = `SELECT * FROM public.meal WHERE id=$1`
		var ml meal.WithLike
		err = r.client.QueryRow(ctx, q, mealId).Scan(
			&ml.Id, &ml.Name,
			&ml.TotalCalories, &ml.TotalProteins,
			&ml.TotalFats, &ml.TotalCarbs,
			&ml.AuthorId, &ml.Description, &ml.Likes,
			&ml.Picture,
		)
		ml.Like = true
		q = `SELECT food_id, quantity FROM meal_food WHERE meal_id = $1`
		rw, ferr := r.client.Query(ctx, q, ml.Id)
		if ferr != nil {
			r.log.Error("Error getting foods")
			return nil, ferr
		}
		var foods []meal.Product
		if rw != nil {
			for rw.Next() {
				var f meal.Product
				if err = rw.Scan(&f.ProductId, &f.Quantity); err != nil {
					r.log.Error("Error scanning foods")
					return nil, err
				}
				foods = append(foods, f)
			}
		}
		ml.Products = foods
		meals = append(meals, ml)
	}
	return
}

func (r *repository) FindFoodLikes(ctx context.Context, id int) (foods []food.WithLike, err error) {
	q := `SELECT food_id FROM public.food_client WHERE user_id=$1`
	rows, err := r.client.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var foodId int
		err = rows.Scan(&foodId)
		if err != nil {
			return nil, err
		}
		q = `SELECT * FROM public.food WHERE id=$1`
		var fd food.WithLike
		err = r.client.QueryRow(ctx, q, foodId).Scan(
			&fd.Id, &fd.Name, &fd.Description,
			&fd.Calories, &fd.Proteins,
			&fd.Carbohydrates, &fd.Fats,
			&fd.AuthorId, &fd.Likes, &fd.Picture,
		)
		if err != nil {
			r.log.Error("Failed to scan food")
			return nil, err
		}
		fd.Like = true
		foods = append(foods, fd)
	}
	return
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) client.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
