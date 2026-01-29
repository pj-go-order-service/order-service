CREATE TABLE orders (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    total_amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);