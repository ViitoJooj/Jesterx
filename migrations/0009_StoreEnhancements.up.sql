ALTER TABLE websites ADD COLUMN IF NOT EXISTS mature_content BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE websites ADD COLUMN IF NOT EXISTS rating_avg NUMERIC(3,2) NOT NULL DEFAULT 0;
ALTER TABLE websites ADD COLUMN IF NOT EXISTS rating_count INT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS store_comments (
    id          TEXT PRIMARY KEY,
    website_id  TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content     TEXT NOT NULL CHECK (char_length(content) BETWEEN 3 AND 1000),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS store_ratings (
    id          TEXT PRIMARY KEY,
    website_id  TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stars       INT NOT NULL CHECK (stars BETWEEN 1 AND 5),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (website_id, user_id)
);

CREATE TABLE IF NOT EXISTS store_visits (
    website_id  TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    visit_date  DATE NOT NULL DEFAULT CURRENT_DATE,
    visit_count INT NOT NULL DEFAULT 1,
    PRIMARY KEY (website_id, visit_date)
);

CREATE INDEX IF NOT EXISTS idx_store_comments_website ON store_comments(website_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_store_ratings_website  ON store_ratings(website_id);
CREATE INDEX IF NOT EXISTS idx_store_visits_website   ON store_visits(website_id, visit_date DESC);
