-- Store team members with role-based access
CREATE TABLE IF NOT EXISTS store_members (
    id          TEXT PRIMARY KEY,
    website_id  TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role        TEXT NOT NULL CHECK (role IN ('manager', 'catalog_manager', 'support', 'logistics')),
    invited_by  TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (website_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_store_members_website ON store_members (website_id);
CREATE INDEX IF NOT EXISTS idx_store_members_user    ON store_members (user_id);

-- Allow comment replies (max 1 level deep; parent must have no parent itself)
ALTER TABLE store_comments
    ADD COLUMN IF NOT EXISTS parent_comment_id TEXT REFERENCES store_comments(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_store_comments_parent ON store_comments (parent_comment_id);
