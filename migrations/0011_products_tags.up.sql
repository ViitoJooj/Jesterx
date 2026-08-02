CREATE TABLE IF NOT EXISTS products_tags (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    product_uuid UUID NOT NULL,
    label VARCHAR(250) NOT NULL,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_tags_product ON products_tags (product_uuid);
CREATE INDEX IF NOT EXISTS idx_products_tags_label ON products_tags (label);
