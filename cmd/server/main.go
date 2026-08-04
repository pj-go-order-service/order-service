package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/pj-go-order-service/order-service/docs"
	httpadapter "github.com/pj-go-order-service/order-service/internal/adapters/httpadapter"
	"github.com/pj-go-order-service/order-service/internal/adapters/memory"
	"github.com/pj-go-order-service/order-service/internal/adapters/postgres"
	"github.com/pj-go-order-service/order-service/internal/app/service"
	"github.com/pj-go-order-service/order-service/internal/config"
	"github.com/pj-go-order-service/order-service/internal/ports"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Order Service API
// @version 1.0
// @description Микросервис управления заказами на Go. Чистая архитектура + DDD.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	// выбираем репозиторий: postgres если указан DATABASE_URL, иначе memory
	var repo ports.OrderRepository = memory.NewOrderRepository()

	if cfg.DatabaseURL != "" {
		pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL)
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		repo = postgres.NewOrderRepository(pool)
		slog.Info("connected to postgres")
	} else {
		slog.Info("using in-memory repository")
	}

	orderService := service.NewOrderService(repo)
	handler := httpadapter.NewOrderHandler(orderService)
	healthHandler := httpadapter.NewHealthHandler(repo.Ping)

	mux := http.NewServeMux()
	mux.Handle("/orders", handler)
	mux.Handle("/orders/", handler)
	mux.Handle("/health", healthHandler)
	mux.Handle("/health/ready", healthHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// middleware
	var h http.Handler = mux
	h = httpadapter.LoggingMiddleware(h)
	h = httpadapter.RecoveryMiddleware(h)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h,
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("order-service started", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}

// setupLogger настраивает slog: текстовый вывод в dev, JSON в prod.
func setupLogger(level string) {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)
}
