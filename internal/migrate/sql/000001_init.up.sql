-- Enable pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ---------------------------------------------------------------
-- customers
-- ---------------------------------------------------------------
CREATE TABLE IF NOT EXISTS customers (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ---------------------------------------------------------------
-- products
-- ---------------------------------------------------------------
CREATE TABLE IF NOT EXISTS products (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255)  NOT NULL,
    description    TEXT          NOT NULL DEFAULT '',
    price          NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    stock_quantity INT           NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- ---------------------------------------------------------------
-- orders
-- ---------------------------------------------------------------
CREATE TABLE IF NOT EXISTS orders (
    id           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id  UUID          NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    status       VARCHAR(50)   NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);

-- ---------------------------------------------------------------
-- order_items
-- ---------------------------------------------------------------
CREATE TABLE IF NOT EXISTS order_items (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id   UUID          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID          NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity   INT           NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL CHECK (unit_price >= 0),
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id   ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
