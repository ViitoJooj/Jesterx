CREATE TABLE IF NOT EXISTS rbac (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID,
    label VARCHAR(100),
    can_read BOOLEAN NOT NULL,
    can_write BOOLEAN NOT NULL,
    can_update BOOLEAN NOT NULL,
    can_upgrade BOOLEAN NOT NULL,
    can_delete BOOLEAN NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX  IF NOT EXISTS idx_website ON rbac (website_uuid);
CREATE INDEX  IF NOT EXISTS idx_label ON rbac (label);