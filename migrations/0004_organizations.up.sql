CREATE TABLE IF NOT EXISTS organizations (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID NOT NULL,
    owner_uuid UUID NOT NULL,
    image_url VARCHAR(1000),
    name VARCHAR(250),
    trade_name VARCHAR(250),
    cnpj VARCHAR(14),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX  IF NOT EXISTS idx_organization_website ON organizations (website_uuid);
CREATE INDEX  IF NOT EXISTS idx_organization_name ON organizations (name);
CREATE INDEX  IF NOT EXISTS idx_organization_owner ON organizations (owner_uuid);
CREATE INDEX  IF NOT EXISTS idx_organization_cnpj ON organizations (cnpj);