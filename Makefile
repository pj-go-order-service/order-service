# Order Service Makefile
# Стандартные команды для разработки и CI

# Переменные
BINARY      := bin/server
GO          := go
GOLANGCI    := golangci-lint
SWAG        := swag
DOCKER      := docker
COMPOSE     := docker-compose
# -race требует cgo; на Windows можно отключить: make test RACE=
RACE        ?= -race

# Кросс-платформенное создание директории для бинарника
ifeq ($(OS),Windows_NT)
MKDIR       = if not exist "$(dir $(BINARY))" mkdir "$(dir $(BINARY))"
RM          = if exist bin rmdir /s /q bin
else
MKDIR       = mkdir -p "$(dir $(BINARY))"
RM          = rm -rf bin
endif

# Параметры подключения к БД (по умолчанию — как в docker-compose.yml)
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_USER     ?= postgres
DB_PASSWORD ?= password
DB_NAME     ?= orders
DB_SSLMODE  ?= disable
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_DIR := migrations

.PHONY: help build run test vet lint fmt swagger migrate migrate-up migrate-down \
        docker-up docker-down docker-build clean install-tools

help: ## Показать список доступных команд
	@echo "Available commands:"
	@echo "  build          - build binary to bin/server"
	@echo "  run            - run server (in-memory)"
	@echo "  test           - run all tests with -race"
	@echo "  vet            - run go vet"
	@echo "  lint           - run golangci-lint"
	@echo "  fmt            - format code"
	@echo "  swagger        - regenerate Swagger docs"
	@echo "  migrate        - apply migrations (alias for migrate-up)"
	@echo "  migrate-up     - apply migrations to DB"
	@echo "  migrate-down   - rollback migrations"
	@echo "  docker-up      - start services via docker-compose"
	@echo "  docker-down    - stop docker-compose services"
	@echo "  docker-build   - build docker image"
	@echo "  clean          - remove build artifacts"
	@echo "  install-tools  - install golangci-lint and swag"

build: ## Собрать бинарник
	@$(MKDIR)
	$(GO) build -o $(BINARY) ./cmd/server

run: ## Запустить сервер (in-memory)
	$(GO) run ./cmd/server

test: ## Запустить все тесты
	$(GO) test $(RACE) -count=1 ./...

vet: ## Проверить код статическим анализатором
	$(GO) vet ./...

lint: ## Запустить golangci-lint
	$(GOLANGCI) run

fmt: ## Отформатировать код
	$(GO) fmt ./...

swagger: ## Перегенерировать Swagger-документацию
	$(SWAG) init -g cmd/server/main.go -o docs

migrate: migrate-up ## Применить все миграции (алиас для migrate-up)

migrate-up: ## Применить миграции к БД
	@for f in $(MIGRATIONS_DIR)/*.sql; do \
		echo "Applying $$f..."; \
		psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f "$$f" || exit 1; \
	done
	@echo "Migrations applied."

migrate-down: ## Откатить миграции (удалить таблицы orders и order_items)
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -c "DROP TABLE IF EXISTS order_items; DROP TABLE IF EXISTS orders;"
	@echo "Migrations rolled back."

docker-up: ## Запустить сервис и БД через docker-compose
	$(COMPOSE) up --build

docker-down: ## Остановить docker-compose сервисы
	$(COMPOSE) down

docker-build: ## Собрать docker-образ
	$(DOCKER) build -t order-service .

clean: ## Удалить артефакты сборки
	@$(RM)

install-tools: ## Установить инструменты разработки (golangci-lint, swag)
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest