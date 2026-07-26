package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pj-go-order-service/order-service/internal/app/service"
	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/orders")
	path = strings.TrimPrefix(path, "/")

	switch {
	case path == "" && r.Method == http.MethodPost:
		h.createOrder(w, r)
	case path != "" && r.Method == http.MethodGet:
		h.getOrder(w, r, path)
	case strings.HasSuffix(path, "/pay") && r.Method == http.MethodPost:
		h.payOrder(w, r, strings.TrimSuffix(path, "/pay"))
	case strings.HasSuffix(path, "/cancel") && r.Method == http.MethodPost:
		h.cancelOrder(w, r, strings.TrimSuffix(path, "/cancel"))
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		writeError(w, "order must contain at least one item", http.StatusBadRequest)
		return
	}

	for _, it := range req.Items {
		if it.Price <= 0 {
			writeError(w, "price must be positive", http.StatusBadRequest)
			return
		}
		if it.Quantity <= 0 {
			writeError(w, "quantity must be positive", http.StatusBadRequest)
			return
		}
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

	result, err := h.service.CreateOrder(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, order.ErrEmptyOrder) {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, "invalid order id", http.StatusBadRequest)
		return
	}

	o, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, order.ErrNotFound) {
			writeError(w, "order not found", http.StatusNotFound)
			return
		}
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(o)
}

func (h *OrderHandler) payOrder(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(strings.TrimRight(idStr, "/"))
	if err != nil {
		writeError(w, "invalid order id", http.StatusBadRequest)
		return
	}

	if err := h.service.PayOrder(r.Context(), id); err != nil {
		if errors.Is(err, order.ErrNotFound) {
			writeError(w, "order not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, order.ErrInvalidState) {
			writeError(w, err.Error(), http.StatusConflict)
			return
		}
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "paid"})
}

func (h *OrderHandler) cancelOrder(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(strings.TrimRight(idStr, "/"))
	if err != nil {
		writeError(w, "invalid order id", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelOrder(r.Context(), id); err != nil {
		if errors.Is(err, order.ErrNotFound) {
			writeError(w, "order not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, order.ErrInvalidState) {
			writeError(w, err.Error(), http.StatusConflict)
			return
		}
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "canceled"})
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
