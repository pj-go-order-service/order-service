package main

import (
	"log"
	"net/http"

	httpadapter "github.com/pj-go-order-service/order-service/internal/adapters/httpadapter"
	"github.com/pj-go-order-service/order-service/internal/adapters/memory"
	"github.com/pj-go-order-service/order-service/internal/app/service"
)

func main() {
	repo := memory.NewOrderRepository()
	orderService := service.NewOrderService(repo)

	handler := httpadapter.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", handler.CreateOrder)

	_ = orderService // на потом

	log.Println("order-service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
