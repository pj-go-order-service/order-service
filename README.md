# Order Service

Микросервис управления заказами на Go. Чистая архитектура + DDD.

## API

Swagger UI доступен по адресу: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

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

## Health Check

### Liveness — сервис жив
```bash
GET /health
# {"status":"ok"}
```

### Readiness — сервис готов принимать трафик (проверяет доступность БД)
```bash
GET /health/ready
# {"status":"ready"}
```
При недоступности БД вернёт HTTP 503 и `{"status":"unavailable"}`.

## Генерация Swagger-документации
```bash
swag init -g cmd/server/main.go -o docs
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

## Makefile

Стандартные команды для разработки и CI:

| Команда | Описание |
|---------|----------|
| `make build` | Собрать бинарник в `bin/server` |
| `make run` | Запустить сервер (in-memory) |
| `make test` | Запустить все тесты с `-race` (на Windows: `make test RACE=`) |
| `make vet` | Статический анализ `go vet` |
| `make lint` | Запустить `golangci-lint` |
| `make fmt` | Отформатировать код |
| `make swagger` | Перегенерировать Swagger-документацию |
| `make migrate` / `make migrate-up` | Применить миграции к БД |
| `make migrate-down` | Откатить миграции |
| `make docker-up` | Запустить сервис и БД через docker-compose |
| `make docker-down` | Остановить docker-compose сервисы |
| `make docker-build` | Собрать docker-образ |
| `make clean` | Удалить артефакты сборки |
| `make install-tools` | Установить golangci-lint и swag |
| `make help` | Показать список всех команд |

Параметры подключения к БД для миграций можно переопределить:
```bash
make migrate DB_HOST=myhost DB_PORT=5433 DB_USER=admin DB_PASSWORD=secret DB_NAME=orders
```

## Структура проекта

```
├── cmd/server/          # точка входа
├── docs/                # сгенерированная Swagger-документация
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
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── .github/workflows/   # CI