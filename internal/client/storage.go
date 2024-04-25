package client

import "context"

type Repository interface {
	Create(ctx context.Context, client *Client) error
	FindByEmail(ctx context.Context, email string) (Client, error)
	FindById(ctx context.Context, id int) (Client, error)
	UpdateClient(ctx context.Context, cl Client) error
	Delete(ctx context.Context, id int) error
	FindGoalById(ctx context.Context, id int) (Goal, error)
	FindDietById(ctx context.Context, id int) (Diet, error)
	UpdateDiet(ctx context.Context, diet Diet, id int) error
	UpdateGoal(ctx context.Context, goal Goal, id int) error
}
