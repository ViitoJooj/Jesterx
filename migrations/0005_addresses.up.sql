CREATE TABLE IF NOT EXISTS addresses (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID NOT NULL,
    owner_uuid UUID NOT NULL,
    owner_type VARCHAR(12) NOT NULL CHECK (owner_type IN ('User', 'Organization')),
    label VARCHAR(250) NOT NULL,
    address_line1 VARCHAR(500),
    address_line2 VARCHAR(500),
    neighborhood VARCHAR(250),
    city VARCHAR(250) NOT NULL,
    state VARCHAR(100) NOT NULL,
    state_code CHAR(2) NOT NULL,
    postal_code VARCHAR(9) NOT NULL,
    reference_point VARCHAR(500),
    delivery_notes VARCHAR(1000),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX  IF NOT EXISTS idx_addresses_website ON addresses (website_uuid);
CREATE INDEX  IF NOT EXISTS idx_addresses_owner ON addresses (owner_uuid);
CREATE INDEX  IF NOT EXISTS idx_addresses_owner_type ON addresses (owner_type);
CREATE INDEX  IF NOT EXISTS idx_addresses_label ON addresses (label);
