CREATE TABLE IF NOT EXISTS phones (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID NOT NULL,
    owner_uuid UUID NOT NULL,
    owner_type VARCHAR(12) NOT NULL CHECK (owner_type IN ('User', 'Organization')),
    label VARCHAR(250) NOT NULL,
    number VARCHAR(12) NOT NULL,
    is_default BOOLEAN NOT NULL,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE IF NOT EXISTS idx_phone_website ON phone (website_uuid);
CREATE IF NOT EXISTS idx_phone_owner ON phone (owner_uuid);
CREATE IF NOT EXISTS idx_phone_owner_type ON phone (owner_type);
CREATE IF NOT EXISTS idx_phone_label ON phone (label);
CREATE IF NOT EXISTS idx_phone_number ON phone (number);