package httpadapter

import "github.com/google/uuid"

// GetOrderResponse — ответ с данными заказа
type GetOrderResponse struct {
	ID        uuid.UUID      `json:"id"`
	Items     []GetOrderItem `json:"items"`
	Total     GetOrderMoney  `json:"total"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
}

// GetOrderItem — позиция заказа
type GetOrderItem struct {
	ProductID string        `json:"product_id"`
	Name      string        `json:"name"`
	Price     GetOrderMoney `json:"price"`
	Quantity  int           `json:"quantity"`
}

// GetOrderMoney — денежная сумма
type GetOrderMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// CreateOrderRequest — запрос на создание заказа
type CreateOrderRequest struct {
	Items []CreateOrderItem `json:"items"`
}

// CreateOrderItem — позиция заказа
type CreateOrderItem struct {
	ProductID string `json:"product_id"`
	Price     int64  `json:"price"`
	Quantity  int    `json:"quantity"`
}

// CreateOrderResponse — ответ на создание заказа
type CreateOrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}

// StatusResponse — ответ со статусом операции
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
