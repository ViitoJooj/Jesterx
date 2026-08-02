CREATE TABLE IF NOT EXISTS storage_products (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    product_uuid UUID NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_storage_products_product ON storage_products (product_uuid);
