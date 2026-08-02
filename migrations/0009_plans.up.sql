CREATE TABLE IF NOT EXISTS plans (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    name VARCHAR(250) NOT NULL,
    description VARCHAR(3000),
    max_websites INT NOT NULL DEFAULT 0,
    max_routers INT NOT NULL DEFAULT 0,
    max_products INT NOT NULL DEFAULT 0,
    cost_per_sale_rate INT NOT NULL DEFAULT 0,
    coin VARCHAR(3) NOT NULL CHECK (coin IN ('BRL', 'USD', 'EUR')),
    price INT NOT NULL,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_plans_name ON plans (name);
