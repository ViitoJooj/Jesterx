-- 0010_CommerceProfileEnhancements.up.sql
-- Extend commerce, order, and store-user profile data

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS slug               VARCHAR(200),
    ADD COLUMN IF NOT EXISTS short_description  VARCHAR(500),
    ADD COLUMN IF NOT EXISTS brand              VARCHAR(100),
    ADD COLUMN IF NOT EXISTS model              VARCHAR(100),
    ADD COLUMN IF NOT EXISTS barcode            VARCHAR(50),
    ADD COLUMN IF NOT EXISTS condition          VARCHAR(20)
        CHECK (condition IS NULL OR condition IN ('new','used','refurbished')),
    ADD COLUMN IF NOT EXISTS weight_grams       INTEGER
        CHECK (weight_grams IS NULL OR weight_grams >= 0),
    ADD COLUMN IF NOT EXISTS width_cm           NUMERIC(10,2)
        CHECK (width_cm IS NULL OR width_cm >= 0),
    ADD COLUMN IF NOT EXISTS height_cm          NUMERIC(10,2)
        CHECK (height_cm IS NULL OR height_cm >= 0),
    ADD COLUMN IF NOT EXISTS length_cm          NUMERIC(10,2)
        CHECK (length_cm IS NULL OR length_cm >= 0),
    ADD COLUMN IF NOT EXISTS material           VARCHAR(100),
    ADD COLUMN IF NOT EXISTS color              VARCHAR(50),
    ADD COLUMN IF NOT EXISTS size               VARCHAR(50),
    ADD COLUMN IF NOT EXISTS warranty_months    INTEGER
        CHECK (warranty_months IS NULL OR warranty_months >= 0),
    ADD COLUMN IF NOT EXISTS origin_country     VARCHAR(2),
    ADD COLUMN IF NOT EXISTS tags               JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS attributes         JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS requires_shipping  BOOLEAN NOT NULL DEFAULT TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_products_slug_unique
    ON products(website_id, slug)
    WHERE slug IS NOT NULL AND slug <> '';

CREATE INDEX IF NOT EXISTS idx_products_sku
    ON products(website_id, sku)
    WHERE sku IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_products_category
    ON products(website_id, category)
    WHERE category IS NOT NULL;


ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS buyer_user_id           TEXT REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS buyer_document          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS shipping_name           VARCHAR(200),
    ADD COLUMN IF NOT EXISTS shipping_phone          VARCHAR(20),
    ADD COLUMN IF NOT EXISTS shipping_zip_code       VARCHAR(10),
    ADD COLUMN IF NOT EXISTS shipping_address_street VARCHAR(200),
    ADD COLUMN IF NOT EXISTS shipping_address_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS shipping_address_complement VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_address_district VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_address_city   VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_address_state  VARCHAR(2),
    ADD COLUMN IF NOT EXISTS shipping_address_country VARCHAR(50) DEFAULT 'BR',
    ADD COLUMN IF NOT EXISTS shipping_method         VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_cost           NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (shipping_cost >= 0),
    ADD COLUMN IF NOT EXISTS discount_total          NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
    ADD COLUMN IF NOT EXISTS tax_total               NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (tax_total >= 0),
    ADD COLUMN IF NOT EXISTS currency                VARCHAR(10) NOT NULL DEFAULT 'BRL';

CREATE INDEX IF NOT EXISTS idx_orders_buyer_user
    ON orders(buyer_user_id)
    WHERE buyer_user_id IS NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'orders_shipping_required'
    ) THEN
        ALTER TABLE orders
            ADD CONSTRAINT orders_shipping_required CHECK (
                buyer_user_id IS NULL OR (
                    shipping_zip_code IS NOT NULL AND
                    shipping_address_street IS NOT NULL AND
                    shipping_address_number IS NOT NULL AND
                    shipping_address_city IS NOT NULL AND
                    shipping_address_state IS NOT NULL AND
                    shipping_address_country IS NOT NULL
                )
            );
    END IF;
END$$;


ALTER TABLE users
    ADD COLUMN IF NOT EXISTS display_name     VARCHAR(100),
    ADD COLUMN IF NOT EXISTS birth_date       DATE,
    ADD COLUMN IF NOT EXISTS gender           VARCHAR(20)
        CHECK (gender IS NULL OR gender IN ('male','female','other','prefer_not')),
    ADD COLUMN IF NOT EXISTS bio              VARCHAR(500),
    ADD COLUMN IF NOT EXISTS instagram        VARCHAR(100),
    ADD COLUMN IF NOT EXISTS website_url      VARCHAR(200),
    ADD COLUMN IF NOT EXISTS whatsapp         VARCHAR(20),
    ADD COLUMN IF NOT EXISTS address_district VARCHAR(100);
