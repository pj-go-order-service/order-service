package service

type CreateOrderItem struct {
	ProductID string
	Price     int64
	Quantity  int
}
type CreateOrderCommand struct {
	Items []CreateOrderItem
}
