package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler отвечает за health-check эндпоинты.
type HealthHandler struct {
	ping func(context.Context) error
}

// NewHealthHandler создаёт обработчик health-check.
// Если ping == nil, сервис считается всегда готовым (например, in-memory репозиторий).
func NewHealthHandler(ping func(context.Context) error) *HealthHandler {
	return &HealthHandler{ping: ping}
}

// ServeHTTP маршрутизирует запросы к /health и /health/ready.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		h.liveness(w)
	case "/health/ready":
		h.readiness(w, r)
	default:
		http.NotFound(w, r)
	}
}

// liveness — сервис жив, если процесс отвечает.
func (h *HealthHandler) liveness(w http.ResponseWriter) {
	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readiness — сервис готов принимать трафик, если доступны внешние зависимости.
func (h *HealthHandler) readiness(w http.ResponseWriter, r *http.Request) {
	if h.ping == nil {
		h.writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.ping(ctx); err != nil {
		h.writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h *HealthHandler) writeJSON(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
