-- 0004_Social.up.sql
-- Store social tables: comments, ratings, visits, members

CREATE TABLE IF NOT EXISTS store_comments (
    id                TEXT        PRIMARY KEY,
    website_id        TEXT        NOT NULL REFERENCES websites(id)      ON DELETE CASCADE,
    user_id           TEXT        NOT NULL REFERENCES users(id)         ON DELETE CASCADE,
    content           TEXT        NOT NULL CHECK (char_length(content) BETWEEN 1 AND 1000),
    -- NULL for replies from the store team; 1-5 for customer reviews
    stars             INTEGER              CHECK (stars BETWEEN 1 AND 5),
    -- NULL for top-level comments; set for replies (max 1 level deep)
    parent_comment_id TEXT                 REFERENCES store_comments(id) ON DELETE CASCADE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Covers ListComments: ordered by parent NULLS FIRST then chronologically
CREATE INDEX IF NOT EXISTS idx_comments_website_thread ON store_comments(website_id, parent_comment_id NULLS FIRST, created_at ASC);
CREATE INDEX IF NOT EXISTS idx_comments_user_id        ON store_comments(user_id);



CREATE TABLE IF NOT EXISTS store_ratings (
    id         TEXT        PRIMARY KEY,
    website_id TEXT        NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    user_id    TEXT        NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    stars      INTEGER     NOT NULL CHECK (stars BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT store_ratings_unique UNIQUE (website_id, user_id)
);

-- Covers RecalcRating aggregate queries
CREATE INDEX IF NOT EXISTS idx_ratings_website_id ON store_ratings(website_id);



-- High-write table: one row per (store, day), updated via ON CONFLICT
CREATE TABLE IF NOT EXISTS store_visits (
    website_id  TEXT    NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    visit_date  DATE    NOT NULL DEFAULT CURRENT_DATE,
    visit_count INTEGER NOT NULL DEFAULT 0 CHECK (visit_count >= 0),
    PRIMARY KEY (website_id, visit_date)
);



CREATE TABLE IF NOT EXISTS store_members (
    id         TEXT        PRIMARY KEY,
    website_id TEXT        NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    user_id    TEXT        NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    role       VARCHAR(30) NOT NULL
               CHECK (role IN ('manager','catalog_manager','support','logistics')),
    -- SET NULL if the inviter is deleted
    invited_by TEXT                 REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT store_members_unique UNIQUE (website_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_members_website_id ON store_members(website_id);
CREATE INDEX IF NOT EXISTS idx_members_user_id    ON store_members(user_id);
