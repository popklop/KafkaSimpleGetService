package order

import "context"

type Repository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetAll(ctx context.Context) ([]*Order, error)
}
