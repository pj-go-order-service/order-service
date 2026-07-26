CREATE TABLE orders (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    total_amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id TEXT NOT NULL,
    price_amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);