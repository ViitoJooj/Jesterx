CREATE TABLE IF NOT EXISTS website_routes (
    id TEXT PRIMARY KEY NOT NULL,
    website_id TEXT NOT NULL,
    path VARCHAR(180) NOT NULL,
    title VARCHAR(100) NOT NULL,
    requires_auth BOOLEAN NOT NULL DEFAULT false,
    position INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (website_id, path)
);

CREATE TABLE IF NOT EXISTS website_versions (
    id TEXT PRIMARY KEY NOT NULL,
    website_id TEXT NOT NULL,
    version INT NOT NULL,
    source_type VARCHAR(30) NOT NULL CHECK (source_type IN ('JXML','REACT','SVELTE','ELEMENTOR_JSON')),
    source TEXT NOT NULL,
    compiled_html TEXT NOT NULL,
    scan_status VARCHAR(20) NOT NULL CHECK (scan_status IN ('clean','warning','blocked')),
    scan_score INT NOT NULL DEFAULT 100,
    scan_findings TEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMP NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (website_id, version)
);
