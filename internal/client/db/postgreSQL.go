package client

import (
	"CalorieGuide-db/internal/client"
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
	var goalId, dietId int
	q := `
		INSERT INTO goal DEFAULT
		VALUES
		RETURNING id`
	if err := r.client.QueryRow(ctx, q).Scan(&goalId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return err
		}
		return err
	}
	client.Goal = goalId
	q = `
		INSERT INTO diet DEFAULT
		VALUES
		RETURNING id`
	if err := r.client.QueryRow(ctx, q).Scan(&dietId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Error("Data base error", slg.PgErr(*pgErr))
			return err
		}
		return err
	}
	client.Diet = dietId
	q = `
		INSERT INTO client (user_name, surname, email, password, goal_id, diet_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	if err := r.client.QueryRow(
		ctx, q,
		client.Name, client.Surname,
		client.Email, client.Password,
		client.Goal, client.Diet,
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
	SELECT id, user_name, surname, email, password
	FROM public.client
	WHERE email = $1
	`
	rw := r.client.QueryRow(ctx, q, email)
	if err = rw.Scan(
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password,
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

func (r *repository) FindById(ctx context.Context, id int) (cl client.Client, err error) {
	q := `
	SELECT id, user_name, surname, email, password
	FROM public.client
	WHERE id = $1
	`
	rw := r.client.QueryRow(ctx, q, id)
	if err = rw.Scan(
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password,
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

func (r *repository) Update(ctx context.Context, cl client.Client) error {
	q := `
	UPDATE public.client SET
		user_name=$2, surname=$3,
		email=$4, password=$5
	WHERE id = $1;
	`
	r.client.QueryRow(ctx, q,
		&cl.Id, &cl.Name,
		&cl.Surname, &cl.Email,
		&cl.Password,
	)
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `
		DELETE FROM public.client WHERE id=$1;
	`
	r.client.QueryRow(ctx, q, id)
	return nil
}

func NewRepository(client postgreSQL.Client, log *slog.Logger) client.Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
