package http

type CreateOrderRequest struct {
	Items []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID string `json:"product_id"`
	Price     int64  `json:"price"`
	Quantity  int    `json:"quantity"`
}
