package client

import "context"

type Repository interface {
	Create(ctx context.Context, client *Client) error
	FindByEmail(ctx context.Context, email string) (Client, error)
	Update(ctx context.Context, cl Client) error
	Delete(ctx context.Context, id int) (Client, error)
}
