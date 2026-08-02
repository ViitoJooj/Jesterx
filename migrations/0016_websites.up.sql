CREATE TABLE IF NOT EXISTS websites (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    owner_uuid UUID NOT NULL,
    owner_type VARCHAR(12) NOT NULL CHECK (owner_type IN ('User', 'Organization')),
    label VARCHAR(250) NOT NULL,
    url VARCHAR(1000),
    write_in VARCHAR(20) NOT NULL CHECK (write_in IN ('Component', 'React', 'Svelte')),
    description VARCHAR(3000),
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_websites_owner ON websites (owner_uuid);
CREATE INDEX IF NOT EXISTS idx_websites_owner_type ON websites (owner_type);
CREATE INDEX IF NOT EXISTS idx_websites_label ON websites (label);
