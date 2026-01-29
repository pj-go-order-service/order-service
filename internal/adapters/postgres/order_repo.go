package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (id, status, total_amount, currency, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		o.ID,
		o.Status,
		o.Total.Amount,
		o.Total.Currency,
		o.CreatedAt,
	)
	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, status, total_amount, currency, created_at
        FROM orders WHERE id = $1`,
		id,
	)

	var o order.Order
	var amount int64
	var currency string

	err := row.Scan(
		&o.ID,
		&o.Status,
		&amount,
		&currency,
		&o.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	o.Total = order.NewMoney(amount, currency)
	return &o, nil
}

func OpenDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
