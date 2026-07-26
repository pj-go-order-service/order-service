package service

import (
	"context"

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
func (s *OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*CreateOrderResult, error) {
	if len(cmd.Items) == 0 {
		return nil, order.ErrEmptyOrder
	}

	items := make([]order.Item, 0, len(cmd.Items))
	for _, it := range cmd.Items {
		price := order.NewMoney(it.Price, "RUB")
		item := order.Item{
			ProductID: order.ProductID(it.ProductID),
			Price:     price,
			Quantity:  it.Quantity,
		}
		items = append(items, item)
	}

	newOrder, err := order.NewOrder(items)
	if err != nil {
		return nil, err
	}

	// сохраняем заказ в репозитории
	if err := s.repo.Save(ctx, newOrder); err != nil {
		return nil, err
	}

	return &CreateOrderResult{
		OrderID: newOrder.ID,
		Status:  newOrder.Status,
	}, nil
}

// PayOrder меняет статус заказа на Paid
func (s *OrderService) PayOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if o == nil {
		return order.ErrNotFound
	}

	if err := o.Pay(); err != nil {
		return err
	}

	return s.repo.Save(ctx, o)
}

// GetOrder возвращает заказ по ID
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// CancelOrder меняет статус заказа на Cancelled
func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if o == nil {
		return order.ErrNotFound
	}

	if err := o.Cancel(); err != nil {
		return err
	}

	return s.repo.Save(ctx, o)
}
