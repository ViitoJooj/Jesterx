CREATE TABLE IF NOT EXISTS themes (
    id          TEXT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    category    VARCHAR(50) NOT NULL DEFAULT 'general',
    preview_url TEXT,
    source_type VARCHAR(30) NOT NULL DEFAULT 'ELEMENTOR_JSON',
    source      TEXT NOT NULL DEFAULT '{}',
    active      BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_themes_category ON themes(category);
