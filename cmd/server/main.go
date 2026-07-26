package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/pj-go-order-service/order-service/internal/adapters/httpadapter"
	"github.com/pj-go-order-service/order-service/internal/adapters/memory"
	"github.com/pj-go-order-service/order-service/internal/adapters/postgres"
	"github.com/pj-go-order-service/order-service/internal/app/service"
	"github.com/pj-go-order-service/order-service/internal/config"
	"github.com/pj-go-order-service/order-service/internal/ports"
)

func main() {
	cfg := config.Load()

	// выбираем репозиторий: postgres если указан DATABASE_URL, иначе memory
	var repo ports.OrderRepository = memory.NewOrderRepository()

	if cfg.DatabaseURL != "" {
		pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		defer pool.Close()
		repo = postgres.NewOrderRepository(pool)
		log.Println("connected to postgres")
	} else {
		log.Println("using in-memory repository")
	}

	orderService := service.NewOrderService(repo)
	handler := httpadapter.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.Handle("/orders", handler)
	mux.Handle("/orders/", handler)

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
		log.Printf("order-service started on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
