package http

import (
	"encoding/json"
	"net/http"

	"github.com/pj-go-order-service/order-service/internal/app/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid Request", http.StatusBadRequest)
		return
	}

	cmd := service.CreateOrderCommand{
		Items: make([]service.CreateOrderItem, 0, len(req.Items)),
	}

	for _, it := range req.Items {
		cmd.Items = append(cmd.Items, service.CreateOrderItem{
			ProductID: it.ProductID,
			Price:     it.Price,
			Quantity:  it.Quantity,
		})
	}

	o, err := h.service.CreateOrder(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(o)
}
