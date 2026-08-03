CREATE TABLE IF NOT EXISTS websites_components (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID NOT NULL,
    logo_url VARCHAR(1000),
    tittle VARCHAR(250),
    description VARCHAR(3000),
    path VARCHAR(500),
    content JSONB,
    visits INT NOT NULL DEFAULT 0,
    updated_by UUID,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_websites_components_website ON websites_components (website_uuid);
CREATE INDEX IF NOT EXISTS idx_websites_components_path ON websites_components (path);
