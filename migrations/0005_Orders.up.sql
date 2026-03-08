CREATE TABLE IF NOT EXISTS orders (
    id            TEXT PRIMARY KEY,
    website_id    TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    buyer_name    VARCHAR(200) NOT NULL DEFAULT '',
    buyer_email   VARCHAR(250) NOT NULL DEFAULT '',
    buyer_phone   VARCHAR(30),
    status        VARCHAR(30) NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending','paid','shipped','delivered','cancelled','refunded')),
    subtotal      NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
    platform_fee  NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (platform_fee >= 0),
    total         NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
    notes         TEXT,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id          TEXT PRIMARY KEY,
    order_id    TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id  TEXT NOT NULL,
    product_name VARCHAR(200) NOT NULL,
    unit_price  NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
    qty         INTEGER NOT NULL DEFAULT 1 CHECK (qty > 0),
    total       NUMERIC(12,2) NOT NULL CHECK (total >= 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_website_id ON orders(website_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
