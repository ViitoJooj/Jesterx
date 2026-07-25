CREATE TABLE IF NOT EXISTS terms_accepted (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID NOT NULL,
    owner_uuid UUID NOT NULL,
    owner_type VARCHAR(12) NOT NULL CHECK (owner_type IN ('User', 'Organization')),
    accepted_when updated_at TIMESTAMPTZ
);