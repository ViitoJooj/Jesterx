-- 0007_FixReports.up.sql
-- Fix report enum values (they were lowercase, domain uses UPPERCASE)
-- and add evidence_urls + reporter_user_id columns

-- Drop and recreate reports table with corrected enums and new columns
-- (safe to drop since reports rely on no FK from other tables)
DROP TABLE IF EXISTS reports;

DROP TYPE IF EXISTS report_status;
DROP TYPE IF EXISTS report_reason;

CREATE TYPE report_status AS ENUM (
    'OPEN',
    'IN_PROGRESS',
    'RESOLVED',
    'DISMISSED'
);

CREATE TYPE report_reason AS ENUM (
    'SPAM',
    'FRAUD',
    'SCAM',
    'INAPPROPRIATE',
    'COUNTERFEIT',
    'OTHER'
);

CREATE TABLE reports (
    id               TEXT          PRIMARY KEY,
    ticket_number    SERIAL        NOT NULL,
    -- No FK: reports about deleted stores should persist
    website_id       TEXT          NOT NULL,
    -- Optional: links to authenticated user who filed the report
    reporter_user_id TEXT,
    reporter_name    VARCHAR(200)  NOT NULL,
    reporter_email   VARCHAR(250)  NOT NULL,
    reason           report_reason NOT NULL,
    description      TEXT          NOT NULL CHECK (char_length(description) BETWEEN 10 AND 2000),
    -- Array of base64 data-URL strings or storage URLs (max 5)
    evidence_urls    JSONB         NOT NULL DEFAULT '[]',
    status           report_status NOT NULL DEFAULT 'OPEN',
    admin_response   TEXT,
    resolved_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT reports_ticket_unique UNIQUE (ticket_number)
);

CREATE INDEX IF NOT EXISTS idx_reports_website_id      ON reports(website_id);
CREATE INDEX IF NOT EXISTS idx_reports_status          ON reports(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reports_reporter_user   ON reports(reporter_user_id) WHERE reporter_user_id IS NOT NULL;
