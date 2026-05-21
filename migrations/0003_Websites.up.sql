-- 0003_Websites.up.sql
-- Website builder: sites, routes, versions, themes.

CREATE TYPE website_type     AS ENUM ('JESTERX','ECOMMERCE','LANDING_PAGE','SOFTWARE_SELL','COURSE','VIDEO');
CREATE TYPE source_type      AS ENUM ('JXML','REACT','SVELTE','ELEMENTOR_JSON');
CREATE TYPE scan_status_type AS ENUM ('clean','warning','blocked');

CREATE TABLE websites (
    id                TEXT         PRIMARY KEY,
    website_type      website_type NOT NULL DEFAULT 'ECOMMERCE',
    -- URL to Supabase Storage (not BYTEA — never store blobs in the DB)
    image_url         TEXT,
    name              VARCHAR(50)  NOT NULL,
    short_description VARCHAR(500),
    description       VARCHAR(1500),
    -- No FK: the platform website (00000000-...0001) uses a placeholder creator_id
    creator_id        TEXT         NOT NULL,
    -- Optional: website owned by a company rather than an individual
    company_id        TEXT         REFERENCES companies(id) ON DELETE SET NULL,
    banned            BOOLEAN      NOT NULL DEFAULT FALSE,
    mature_content    BOOLEAN      NOT NULL DEFAULT FALSE,
    -- Denormalised from store_ratings for fast reads without aggregation
    rating_avg        NUMERIC(3,2) NOT NULL DEFAULT 0.00,
    rating_count      INTEGER      NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT websites_name_unique UNIQUE (name)
);

CREATE INDEX idx_websites_creator_id ON websites(creator_id);
CREATE INDEX idx_websites_company_id ON websites(company_id) WHERE company_id IS NOT NULL;
CREATE INDEX idx_websites_active     ON websites(banned) WHERE banned = FALSE;

CREATE TABLE website_routes (
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

CREATE INDEX idx_routes_website_position ON website_routes(website_id, position ASC, created_at ASC);

CREATE TABLE website_versions (
    id            TEXT             PRIMARY KEY,
    website_id    TEXT             NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    version       INTEGER          NOT NULL CHECK (version > 0),
    source_type   source_type      NOT NULL,
    source        TEXT             NOT NULL,
    compiled_html TEXT             NOT NULL,
    scan_status   scan_status_type NOT NULL DEFAULT 'clean',
    scan_score    INTEGER          NOT NULL DEFAULT 100,
    scan_findings TEXT,
    published     BOOLEAN          NOT NULL DEFAULT FALSE,
    published_at  TIMESTAMPTZ,
    -- No FK: keep version history even if the creating user is deleted
    created_by    TEXT             NOT NULL,
    created_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    CONSTRAINT website_versions_unique UNIQUE (website_id, version)
);

-- Covers FindLatestVersionByWebsiteID and ListVersionsByWebsiteID
CREATE INDEX idx_versions_website_version ON website_versions(website_id, version DESC);
-- Covers FindPublishedVersionByWebsiteID without scanning unpublished rows
CREATE INDEX idx_versions_published       ON website_versions(website_id, version DESC) WHERE published = TRUE;

CREATE TABLE themes (
    id          TEXT        PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    category    VARCHAR(50)  NOT NULL DEFAULT 'general',
    preview_url TEXT,
    source_type source_type  NOT NULL DEFAULT 'ELEMENTOR_JSON',
    source      TEXT         NOT NULL DEFAULT '{}',
    active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_themes_active_category ON themes(category) WHERE active = TRUE;
