package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*order.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[uuid.UUID]*order.Order),
	}
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.ID] = o
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, order.ErrNotFound
	}

	return o, nil
}
