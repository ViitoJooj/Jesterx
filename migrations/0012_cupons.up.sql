CREATE TABLE IF NOT EXISTS cupons (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    tag_uuid UUID NOT NULL,
    label VARCHAR(250) NOT NULL,
    description VARCHAR(3000),
    value VARCHAR(100) NOT NULL,
    value_type VARCHAR(10) NOT NULL CHECK (value_type IN ('percentage', 'Value'))
);

CREATE INDEX IF NOT EXISTS idx_cupons_tag ON cupons (tag_uuid);
CREATE INDEX IF NOT EXISTS idx_cupons_label ON cupons (label);
