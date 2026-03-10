-- 0005_Reports.up.sql
-- Support ticket / report system

DO $$ BEGIN
    CREATE TYPE report_status AS ENUM ('open','in_review','resolved','rejected');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
    CREATE TYPE report_reason AS ENUM (
        'spam','inappropriate_content','fake_store','fraud','copyright','other'
    );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE TABLE IF NOT EXISTS reports (
    id              TEXT          PRIMARY KEY,
    -- Auto-incrementing human-readable ticket number
    ticket_number   SERIAL        NOT NULL,
    -- Intentionally no FK: reports about deleted stores should persist
    website_id      TEXT          NOT NULL,
    reporter_name   VARCHAR(200)  NOT NULL,
    reporter_email  VARCHAR(250)  NOT NULL,
    reason          report_reason NOT NULL,
    description     TEXT          NOT NULL CHECK (char_length(description) BETWEEN 10 AND 2000),
    status          report_status NOT NULL DEFAULT 'open',
    admin_response  TEXT,
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT reports_ticket_unique UNIQUE (ticket_number)
);

CREATE INDEX IF NOT EXISTS idx_reports_website_id ON reports(website_id);
CREATE INDEX IF NOT EXISTS idx_reports_status     ON reports(status, created_at DESC);
