CREATE TABLE IF NOT EXISTS preparing_shipping_products (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    product_uuid UUID NOT NULL,
    address_uuid UUID NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_preparing_shipping_product ON preparing_shipping_products (product_uuid);
CREATE INDEX IF NOT EXISTS idx_preparing_shipping_address ON preparing_shipping_products (address_uuid);
