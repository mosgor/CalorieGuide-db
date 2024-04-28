package client

import (
	"CalorieGuide-db/internal/client"
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

func (r *repository) Create(ctx context.Context, client *client.Client) error {
	q := `INSERT INTO goal DEFAULT VALUES`
	r.client.QueryRow(ctx, q)
	q = `INSERT INTO diet DEFAULT VALUES`
	r.client.QueryRow(ctx, q)
	q = `
		INSERT INTO client (user_name, surname, email, password, picture)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	if err := r.client.QueryRow(
		ctx, q,
		client.Name, client.Surname,
		client.Email, client.Password,
		client.Picture,
	).Scan(&client.Id); err != nil {
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
	r.client.QueryRow(ctx, q,
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password, &cl.Picture,
	)
	return nil
}

func (r *repository) UpdateGoal(ctx context.Context, goal client.Goal, goalId int) error {
	q := `
	UPDATE public.goal SET
		calories_goal=$2, fats_goal=$3,
		proteins_goal=$4, carbohydrates_goal=$5
	WHERE id = $1;
	`
	r.client.QueryRow(ctx, q,
		goalId, &goal.CaloriesGoal,
		&goal.FatsGoal, &goal.ProteinsGoal,
		&goal.CarbohydratesGoal,
	)
	return nil
}

func (r *repository) UpdateDiet(ctx context.Context, diet client.Diet, dietId int) error {
	q := `UPDATE public.diet SET breakfast_id=$2 WHERE id = $1;`
	if diet.BreakfastId != 0 {
		r.client.QueryRow(ctx, q, dietId, &diet.BreakfastId)
	}
	q = `UPDATE public.diet SET lunch_id=$2 WHERE id = $1;`
	if diet.LunchId != 0 {
		r.client.QueryRow(ctx, q, dietId, &diet.LunchId)
	}
	q = `UPDATE public.diet SET dinner_id=$2 WHERE id = $1;`
	if diet.DinnerId != 0 {
		r.client.QueryRow(ctx, q, dietId, &diet.DinnerId)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int, fdRepo food.Repository) error {
	q := `DELETE FROM public.food_client WHERE user_id=$1`
	r.client.QueryRow(ctx, q, id)
	q = `DELETE FROM public.meal_client WHERE user_id=$1`
	r.client.QueryRow(ctx, q, id)
	//q = `DELETE FROM public.food WHERE author_id=$1`
	//r.client.QueryRow(ctx, q, id)
	q = `SELECT id FROM public.food WHERE author_id=$1`
	rows, err := r.client.Query(ctx, q, id)
	if err != nil {
		return err
	}
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
	q = `DELETE FROM public.client WHERE id=$1`
	r.client.QueryRow(ctx, q, id)
	q = `DELETE FROM public.goal WHERE id=$1`
	r.client.QueryRow(ctx, q, id)
	q = `DELETE FROM public.diet WHERE id=$1`
	r.client.QueryRow(ctx, q, id)
	return nil
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) client.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
