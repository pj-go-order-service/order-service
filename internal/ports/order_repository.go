package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type OrderRepository interface {
	Save(ctx context.Context, o *order.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error)
	Ping(ctx context.Context) error
}
