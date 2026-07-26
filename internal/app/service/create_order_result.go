package service

import (
	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type CreateOrderResult struct {
	OrderID uuid.UUID
	Status  order.Status
}
