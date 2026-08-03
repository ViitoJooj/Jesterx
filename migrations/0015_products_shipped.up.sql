CREATE TABLE IF NOT EXISTS products_shipped (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    product_uuid UUID NOT NULL,
    address_uuid UUID NOT NULL,
    status VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_products_shipped_product ON products_shipped (product_uuid);
CREATE INDEX IF NOT EXISTS idx_products_shipped_address ON products_shipped (address_uuid);
CREATE INDEX IF NOT EXISTS idx_products_shipped_status ON products_shipped (status);
