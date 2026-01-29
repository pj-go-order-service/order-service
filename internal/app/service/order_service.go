package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
	"github.com/pj-go-order-service/order-service/internal/ports"
)

// OrderService — основной сервис для работы с заказами
type OrderService struct {
	repo ports.OrderRepository
}

// NewOrderService создаёт сервис
func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, items []order.Item) (*order.Order, error) {
	if len(items) == 0 {
		return nil, order.ErrEmptyOrder
	}

	total := order.CalculateTotal(items) // экспортированная функция из domain

	newOrder := &order.Order{
		ID:        uuid.New(),
		Items:     items,
		Total:     total,
		Status:    order.StatusCreated,
		CreatedAt: time.Now(), // можно оставить time.Now() напрямую
	}

	// сохраняем заказ в репозитории
	if err := s.repo.Save(ctx, newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}

// PayOrder меняет статус заказа на Paid
func (s *OrderService) PayOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if o.Status != order.StatusCreated {
		return order.ErrInvalidState
	}

	o.Status = order.StatusPaid

	return s.repo.Save(ctx, o)
}

// CancelOrder меняет статус заказа на Cancelled
func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if o.Status != order.StatusCreated {
		return errors.New("order cannot be cancelled")
	}

	o.Status = order.StatusCanceled

	return s.repo.Save(ctx, o)
}
