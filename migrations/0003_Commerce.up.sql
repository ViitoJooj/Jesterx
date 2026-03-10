-- 0003_Commerce.up.sql
-- E-commerce tables: products, orders, order_items

CREATE TABLE IF NOT EXISTS products (
    id            TEXT          PRIMARY KEY,
    website_id    TEXT          NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    name          VARCHAR(200)  NOT NULL,
    description   TEXT          NOT NULL DEFAULT '',
    price         NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    compare_price NUMERIC(12,2)          CHECK (compare_price IS NULL OR compare_price >= 0),
    stock         INTEGER       NOT NULL DEFAULT 0   CHECK (stock >= 0),
    sku           VARCHAR(100),
    category      VARCHAR(100),
    images        JSONB         NOT NULL DEFAULT '[]',
    active        BOOLEAN       NOT NULL DEFAULT TRUE,
    sold_count    INTEGER       NOT NULL DEFAULT 0   CHECK (sold_count >= 0),
    -- No FK: keep products even if the creator user is deleted
    created_by    TEXT          NOT NULL,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Covers ListProductsByWebsiteID (active filter + created_at ordering)
CREATE INDEX IF NOT EXISTS idx_products_website_active ON products(website_id, active, created_at DESC);
-- Covers ListProductsByWebsiteID when ordered by sold_count
CREATE INDEX IF NOT EXISTS idx_products_website_sold   ON products(website_id, sold_count DESC) WHERE active = TRUE;



CREATE TABLE IF NOT EXISTS orders (
    id           TEXT          PRIMARY KEY,
    website_id   TEXT          NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    buyer_name   VARCHAR(200)  NOT NULL,
    buyer_email  VARCHAR(250)  NOT NULL,
    buyer_phone  VARCHAR(20),
    status       VARCHAR(30)   NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending','processing','shipped','delivered','canceled','refunded')),
    subtotal     NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
    platform_fee NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (platform_fee >= 0),
    total        NUMERIC(12,2) NOT NULL CHECK (total >= 0),
    notes        TEXT,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

-- Covers ListBySite and default ordering
CREATE INDEX IF NOT EXISTS idx_orders_website_created  ON orders(website_id, created_at DESC);
-- Covers ListSince (date-range queries across all sites)
CREATE INDEX IF NOT EXISTS idx_orders_created_at       ON orders(created_at);
-- Covers status-filtered queries for active orders
CREATE INDEX IF NOT EXISTS idx_orders_website_status   ON orders(website_id, status)
    WHERE status NOT IN ('delivered','canceled','refunded');



CREATE TABLE IF NOT EXISTS order_items (
    id           TEXT          PRIMARY KEY,
    order_id     TEXT          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    -- SET NULL so order history is preserved even if product is deleted
    product_id   TEXT          REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(200)  NOT NULL,
    unit_price   NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
    qty          INTEGER       NOT NULL DEFAULT 1 CHECK (qty > 0),
    total        NUMERIC(12,2) NOT NULL CHECK (total >= 0)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id   ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id) WHERE product_id IS NOT NULL;
