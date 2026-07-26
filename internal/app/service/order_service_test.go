package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/app/service"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type mocOrderRepository struct {
	savedOrder *order.Order
}

func (m *mocOrderRepository) Save(ctx context.Context, o *order.Order) error {
	m.savedOrder = o
	return nil
}

func (m *mocOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return nil, nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	repo := &mocOrderRepository{}
	svc := service.NewOrderService(repo)

	items := []order.Item{
		{
			ProductID: "p-1",
			Name:      "Keyboard",
			Price:     order.NewMoney(100, "$"),
			Quantity:  2,
		},
		{
			ProductID: "p-2",
			Name:      "Mouse",
			Price:     order.NewMoney(200, "$"),
			Quantity:  2,
		},
	}

	_, err := svc.CreateOrder(context.Background(), items)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if repo.savedOrder == nil {
		t.Fatal("order was not saved")
	}

	if len(repo.savedOrder.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(repo.savedOrder.Items))
	}
}
