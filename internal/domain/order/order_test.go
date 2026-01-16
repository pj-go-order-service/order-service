package order

import (
	"testing"
)

func TestNewOrder(t *testing.T) {
	items := []Item{
		{Name: "Item1", Price: Money{Amount: 100, Currency: "USD"}, Quantity: 2},
		{Name: "Item2", Price: Money{Amount: 50, Currency: "USD"}, Quantity: 1},
	}

	o, err := NewOrder(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTotal := int64(250)
	if o.Total.Amount != expectedTotal {
		t.Errorf("expected total %d, got %d", expectedTotal, o.Total.Amount)
	}

	if o.Status != StatusCreated {
		t.Errorf("expected status %s, got %s", StatusCreated, o.Status)
	}
}

func TestNewOrder_EmptyItems(t *testing.T) {
	_, err := NewOrder([]Item{})
	if err != ErrEmptyOrder {
		t.Fatalf("expected ErrEmptyOrder, got %v", err)
	}
}
func TestOrder_Pay(t *testing.T) {
	items := []Item{
		{Name: "Item1", Price: Money{Amount: 100, Currency: "USD"}, Quantity: 1},
	}

	o, err := NewOrder(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Successful payment
	err = o.Pay()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if o.Status != StatusPaid {
		t.Errorf("expected status %s, got %s", StatusPaid, o.Status)
	}

	// 2. Repeat payment attempt -> should return ErrInvalidState
	err = o.Pay()
	if err != ErrInvalidState {
		t.Errorf("expected ErrInvalidState, got %v", err)
	}
}
func TestCalculateTotalMultipleItems(t *testing.T) {
	items := []Item{
		{Name: "A", Price: Money{Amount: 100, Currency: "USD"}, Quantity: 3}, // 300
		{Name: "B", Price: Money{Amount: 200, Currency: "USD"}, Quantity: 2}, // 400
	}

	o, err := NewOrder(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := int64(700)
	if o.Total.Amount != expected {
		t.Errorf("expected total %d, got %d", expected, o.Total.Amount)
	}
}
