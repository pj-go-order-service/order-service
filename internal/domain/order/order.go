package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	Items     []Item
	Total     Money
	Status    Status
	CreatedAt time.Time
}

func NewOrder(items []Item) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	total, err := CalculateTotal(items)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:        uuid.New(),
		Items:     items,
		Total:     total,
		Status:    StatusCreated,
		CreatedAt: time.Now(),
	}, nil
}

func (o *Order) Pay() error {
	if o.Status != StatusCreated {
		return ErrInvalidState
	}
	o.Status = StatusPaid
	return nil
}

func (o *Order) Cancel() error {
	if o.Status != StatusCreated {
		return ErrInvalidState
	}
	o.Status = StatusCanceled
	return nil
}

func CalculateTotal(items []Item) (Money, error) {
	if len(items) == 0 {
		return Money{}, ErrEmptyOrder
	}

	var total int64
	currency := items[0].Price.Currency

	for _, item := range items {
		total += item.Price.Amount * int64(item.Quantity)
	}

	return NewMoney(total, currency), nil
}
