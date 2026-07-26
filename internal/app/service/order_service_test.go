package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/app/service"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type mockOrderRepository struct {
	savedOrder *order.Order
	orders     map[uuid.UUID]*order.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uuid.UUID]*order.Order),
	}
}

func (m *mockOrderRepository) Save(ctx context.Context, o *order.Order) error {
	m.savedOrder = o
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, order.ErrNotFound
	}
	return o, nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	items := []service.CreateOrderItem{
		{
			ProductID: "p-1",
			Price:     100,
			Quantity:  2,
		},
		{
			ProductID: "p-2",
			Price:     200,
			Quantity:  2,
		},
	}

	cmd := service.CreateOrderCommand{
		Items: items,
	}

	result, err := svc.CreateOrder(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if result.OrderID == uuid.Nil {
		t.Fatal("expected non-zero order ID")
	}

	if result.Status != order.StatusCreated {
		t.Fatalf("expected status %s, got %s", order.StatusCreated, result.Status)
	}

	if repo.savedOrder == nil {
		t.Fatal("order was not saved")
	}

	if len(repo.savedOrder.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(repo.savedOrder.Items))
	}
}

func TestOrderService_CreateOrder_EmptyItems(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	cmd := service.CreateOrderCommand{
		Items: []service.CreateOrderItem{},
	}

	_, err := svc.CreateOrder(context.Background(), cmd)
	if !errors.Is(err, order.ErrEmptyOrder) {
		t.Fatalf("expected ErrEmptyOrder, got %v", err)
	}
}

func TestOrderService_PayOrder(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	// create order first
	items := []service.CreateOrderItem{
		{ProductID: "p-1", Price: 100, Quantity: 1},
	}
	cmd := service.CreateOrderCommand{Items: items}
	result, err := svc.CreateOrder(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// pay order
	err = svc.PayOrder(context.Background(), result.OrderID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// verify status
	o, err := repo.GetByID(context.Background(), result.OrderID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if o.Status != order.StatusPaid {
		t.Fatalf("expected status %s, got %s", order.StatusPaid, o.Status)
	}
}

func TestOrderService_PayOrder_NotFound(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	err := svc.PayOrder(context.Background(), uuid.New())
	if !errors.Is(err, order.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestOrderService_CancelOrder(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	// create order first
	items := []service.CreateOrderItem{
		{ProductID: "p-1", Price: 100, Quantity: 1},
	}
	cmd := service.CreateOrderCommand{Items: items}
	result, err := svc.CreateOrder(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// cancel order
	err = svc.CancelOrder(context.Background(), result.OrderID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// verify status
	o, err := repo.GetByID(context.Background(), result.OrderID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if o.Status != order.StatusCanceled {
		t.Fatalf("expected status %s, got %s", order.StatusCanceled, o.Status)
	}
}

func TestOrderService_CancelOrder_NotFound(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	err := svc.CancelOrder(context.Background(), uuid.New())
	if !errors.Is(err, order.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestOrderService_CancelOrder_AlreadyPaid(t *testing.T) {
	repo := newMockOrderRepository()
	svc := service.NewOrderService(repo)

	// create and pay order
	items := []service.CreateOrderItem{
		{ProductID: "p-1", Price: 100, Quantity: 1},
	}
	cmd := service.CreateOrderCommand{Items: items}
	result, err := svc.CreateOrder(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	err = svc.PayOrder(context.Background(), result.OrderID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// try to cancel paid order
	err = svc.CancelOrder(context.Background(), result.OrderID)
	if !errors.Is(err, order.ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}
}
