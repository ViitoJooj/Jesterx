-- 0002_Builder.up.sql
-- Website builder tables: routes, versions, themes

CREATE TABLE IF NOT EXISTS website_routes (
    id            TEXT         PRIMARY KEY,
    website_id    TEXT         NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    path          VARCHAR(180) NOT NULL,
    title         VARCHAR(100) NOT NULL,
    requires_auth BOOLEAN      NOT NULL DEFAULT FALSE,
    position      INTEGER      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT website_routes_unique_path UNIQUE (website_id, path)
);

CREATE INDEX IF NOT EXISTS idx_routes_website_position ON website_routes(website_id, position ASC, created_at ASC);



CREATE TABLE IF NOT EXISTS website_versions (
    id            TEXT         PRIMARY KEY,
    website_id    TEXT         NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    version       INTEGER      NOT NULL CHECK (version > 0),
    source_type   VARCHAR(30)  NOT NULL
                  CHECK (source_type IN ('JXML','REACT','SVELTE','ELEMENTOR_JSON')),
    source        TEXT         NOT NULL,
    compiled_html TEXT         NOT NULL,
    scan_status   VARCHAR(20)  NOT NULL DEFAULT 'clean'
                  CHECK (scan_status IN ('clean','warning','blocked')),
    scan_score    INTEGER      NOT NULL DEFAULT 100,
    scan_findings TEXT,
    published     BOOLEAN      NOT NULL DEFAULT FALSE,
    published_at  TIMESTAMPTZ,
    -- No FK: keep version history even if the creating user is deleted
    created_by    TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT website_versions_unique UNIQUE (website_id, version)
);

-- Covers FindLatestVersionByWebsiteID and ListVersionsByWebsiteID
CREATE INDEX IF NOT EXISTS idx_versions_website_version ON website_versions(website_id, version DESC);
-- Covers FindPublishedVersionByWebsiteID efficiently
CREATE INDEX IF NOT EXISTS idx_versions_published       ON website_versions(website_id, version DESC) WHERE published = TRUE;



CREATE TABLE IF NOT EXISTS themes (
    id          TEXT         PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    category    VARCHAR(50)  NOT NULL DEFAULT 'general',
    preview_url TEXT,
    source_type VARCHAR(30)  NOT NULL DEFAULT 'ELEMENTOR_JSON',
    source      TEXT         NOT NULL DEFAULT '{}',
    active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_themes_active_category ON themes(category) WHERE active = TRUE;
