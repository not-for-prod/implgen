package in

import (
	"context"

	"github.com/not-for-prod/implgen/example/in/model"
)

type CreateRequest struct {
}

type Aboba interface {
	Create(ctx context.Context, req CreateRequest) (model.OrderID, error)
	Get(ctx context.Context, id model.OrderID) (model.Order, error)
}

type AbobaRepository interface {
	Create(ctx context.Context, req CreateRequest) (model.OrderID, error) // sqlx:GetContext
	Get(ctx context.Context, id model.OrderID) (model.Order, error)       // sqlx:GetContext
}
