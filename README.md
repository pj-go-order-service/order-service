# Order Service

Микросервис управления заказами на Go. Чистая архитектура + DDD.

## API

### Создать заказ
```bash
POST /orders
Content-Type: application/json

{
  "items": [
    {"product_id": "p-1", "price": 100, "quantity": 2},
    {"product_id": "p-2", "price": 200, "quantity": 1}
  ]
}
```

### Получить заказ
```bash
GET /orders/{id}
```

### Оплатить заказ
```bash
POST /orders/{id}/pay
```

### Отменить заказ
```bash
POST /orders/{id}/cancel
```

## Запуск

### Локально (in-memory)
```bash
go run ./cmd/server
```

### С Postgres
```bash
docker-compose up --build
```

## Структура проекта

```
├── cmd/server/          # точка входа
├── internal/
│   ├── domain/order/    # доменная модель (DDD)
│   ├── ports/           # интерфейсы (порты)
│   ├── app/service/     # сервисный слой
│   ├── adapters/        # адаптеры
│   │   ├── httpadapter/ # HTTP API
│   │   ├── memory/      # in-memory репозиторий
│   │   └── postgres/    # Postgres репозиторий (pgx/v5)
│   └── config/          # конфигурация
├── migrations/          # SQL миграции
├── Dockerfile
├── docker-compose.yml
└── .github/workflows/   # CI