package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pj-go-order-service/order-service/internal/domain/order"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// сохраняем заказ
	_, err = tx.Exec(ctx,
		`INSERT INTO orders (id, status, total_amount, currency, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET status = $2, total_amount = $3, currency = $4`,
		o.ID,
		o.Status,
		o.Total.Amount,
		o.Total.Currency,
		o.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	// удаляем старые товары и вставляем новые
	_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, o.ID)
	if err != nil {
		return fmt.Errorf("delete order_items: %w", err)
	}

	for _, item := range o.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, price_amount, currency, quantity)
			VALUES ($1, $2, $3, $4, $5)`,
			o.ID,
			item.ProductID,
			item.Price.Amount,
			item.Price.Currency,
			item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("insert order_item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	o := &order.Order{}

	// загружаем заказ
	err := r.pool.QueryRow(ctx,
		`SELECT id, status, total_amount, currency, created_at
		FROM orders WHERE id = $1`,
		id,
	).Scan(
		&o.ID,
		&o.Status,
		&o.Total.Amount,
		&o.Total.Currency,
		&o.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, order.ErrNotFound
		}
		return nil, fmt.Errorf("query order: %w", err)
	}

	// загружаем товары
	rows, err := r.pool.Query(ctx,
		`SELECT product_id, price_amount, currency, quantity
		FROM order_items WHERE order_id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("query order_items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item order.Item
		err := rows.Scan(
			&item.ProductID,
			&item.Price.Amount,
			&item.Price.Currency,
			&item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan order_item: %w", err)
		}
		o.Items = append(o.Items, item)
	}

	return o, nil
}

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}
