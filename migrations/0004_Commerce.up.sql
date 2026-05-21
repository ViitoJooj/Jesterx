-- 0004_Commerce.up.sql
-- E-commerce: products, orders, order_items.
-- All fields that were split across 0003_Commerce and 0010_Enhancements are here from the start.

CREATE TYPE order_status      AS ENUM ('pending','processing','shipped','delivered','canceled','refunded');
CREATE TYPE product_condition AS ENUM ('new','used','refurbished');

CREATE TABLE products (
    id                TEXT              PRIMARY KEY,
    website_id        TEXT              NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    name              VARCHAR(200)      NOT NULL,
    description       TEXT              NOT NULL DEFAULT '',
    short_description VARCHAR(500),
    price             NUMERIC(12,2)     NOT NULL CHECK (price >= 0),
    compare_price     NUMERIC(12,2)              CHECK (compare_price IS NULL OR compare_price >= 0),
    stock             INTEGER           NOT NULL DEFAULT 0  CHECK (stock >= 0),
    sku               VARCHAR(100),
    slug              VARCHAR(200),
    category          VARCHAR(100),
    brand             VARCHAR(100),
    model             VARCHAR(100),
    barcode           VARCHAR(50),
    condition         product_condition,
    weight_grams      INTEGER                    CHECK (weight_grams IS NULL OR weight_grams >= 0),
    width_cm          NUMERIC(10,2)              CHECK (width_cm IS NULL OR width_cm >= 0),
    height_cm         NUMERIC(10,2)              CHECK (height_cm IS NULL OR height_cm >= 0),
    length_cm         NUMERIC(10,2)              CHECK (length_cm IS NULL OR length_cm >= 0),
    material          VARCHAR(100),
    color             VARCHAR(50),
    size              VARCHAR(50),
    warranty_months   INTEGER                    CHECK (warranty_months IS NULL OR warranty_months >= 0),
    origin_country    VARCHAR(2),
    requires_shipping BOOLEAN           NOT NULL DEFAULT TRUE,
    tags              JSONB             NOT NULL DEFAULT '[]',
    attributes        JSONB             NOT NULL DEFAULT '{}',
    images            JSONB             NOT NULL DEFAULT '[]',
    active            BOOLEAN           NOT NULL DEFAULT TRUE,
    sold_count        INTEGER           NOT NULL DEFAULT 0  CHECK (sold_count >= 0),
    -- No FK: keep products even if the creator's user account is deleted
    created_by        TEXT              NOT NULL,
    created_at        TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_website_active ON products(website_id, active, created_at DESC);
CREATE INDEX idx_products_website_sold   ON products(website_id, sold_count DESC) WHERE active = TRUE;
CREATE UNIQUE INDEX idx_products_slug    ON products(website_id, slug) WHERE slug IS NOT NULL AND slug <> '';
CREATE INDEX idx_products_sku            ON products(website_id, sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_products_category       ON products(website_id, category) WHERE category IS NOT NULL;
-- GIN indexes for efficient JSONB filtering (tags and attribute facets)
CREATE INDEX idx_products_tags           ON products USING GIN (tags);
CREATE INDEX idx_products_attributes     ON products USING GIN (attributes);

CREATE TABLE orders (
    id                          TEXT         PRIMARY KEY,
    website_id                  TEXT         NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    -- NULL for guest checkouts
    buyer_user_id               TEXT         REFERENCES users(id) ON DELETE SET NULL,
    buyer_name                  VARCHAR(200) NOT NULL,
    buyer_email                 VARCHAR(250) NOT NULL,
    buyer_phone                 VARCHAR(20),
    buyer_document              VARCHAR(20),
    shipping_name               VARCHAR(200),
    shipping_phone              VARCHAR(20),
    shipping_zip_code           VARCHAR(10),
    shipping_address_street     VARCHAR(200),
    shipping_address_number     VARCHAR(20),
    shipping_address_complement VARCHAR(100),
    shipping_address_district   VARCHAR(100),
    shipping_address_city       VARCHAR(100),
    shipping_address_state      VARCHAR(2),
    shipping_address_country    VARCHAR(50)   DEFAULT 'BR',
    shipping_method             VARCHAR(100),
    shipping_cost               NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (shipping_cost >= 0),
    subtotal                    NUMERIC(12,2) NOT NULL             CHECK (subtotal >= 0),
    discount_total              NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
    tax_total                   NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (tax_total >= 0),
    platform_fee                NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (platform_fee >= 0),
    total                       NUMERIC(12,2) NOT NULL             CHECK (total >= 0),
    currency                    VARCHAR(10)   NOT NULL DEFAULT 'BRL',
    status                      order_status  NOT NULL DEFAULT 'pending',
    notes                       TEXT,
    created_at                  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_website_created ON orders(website_id, created_at DESC);
CREATE INDEX idx_orders_created_at      ON orders(created_at);
CREATE INDEX idx_orders_website_status  ON orders(website_id, status)
    WHERE status NOT IN ('delivered','canceled','refunded');
CREATE INDEX idx_orders_buyer_user      ON orders(buyer_user_id) WHERE buyer_user_id IS NOT NULL;

CREATE TABLE order_items (
    id           TEXT          PRIMARY KEY,
    order_id     TEXT          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    -- SET NULL so order history survives product deletion
    product_id   TEXT          REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(200)  NOT NULL,
    unit_price   NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
    qty          INTEGER       NOT NULL DEFAULT 1 CHECK (qty > 0),
    -- Computed automatically; never insert manually
    total        NUMERIC(12,2) GENERATED ALWAYS AS (unit_price * qty) STORED
);

CREATE INDEX idx_order_items_order_id   ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id) WHERE product_id IS NOT NULL;
