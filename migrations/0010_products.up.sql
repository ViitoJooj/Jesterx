CREATE TABLE IF NOT EXISTS products (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    name VARCHAR(250) NOT NULL,
    description VARCHAR(3000),
    short_description VARCHAR(500),
    height INT,
    width INT,
    thickness INT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products (name);
CREATE INDEX IF NOT EXISTS idx_products_active ON products (active);
