-- 0006_Reports.up.sql
-- Support ticket / abuse report system.

CREATE TYPE report_status AS ENUM ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'DISMISSED');
CREATE TYPE report_reason AS ENUM ('SPAM', 'FRAUD', 'SCAM', 'INAPPROPRIATE', 'COUNTERFEIT', 'OTHER');

CREATE TABLE reports (
    id               TEXT          PRIMARY KEY,
    -- Auto-incrementing human-readable ticket number
    ticket_number    INTEGER       NOT NULL GENERATED ALWAYS AS IDENTITY,
    -- No FK: reports about deleted stores must persist
    website_id       TEXT          NOT NULL,
    -- NULL for anonymous reports
    reporter_user_id TEXT,
    reporter_name    VARCHAR(200)  NOT NULL,
    reporter_email   VARCHAR(250)  NOT NULL,
    reason           report_reason NOT NULL,
    description      TEXT          NOT NULL CHECK (char_length(description) BETWEEN 10 AND 2000),
    -- Array of storage URLs (max 5 evidence items)
    evidence_urls    JSONB         NOT NULL DEFAULT '[]',
    status           report_status NOT NULL DEFAULT 'OPEN',
    admin_response   TEXT,
    resolved_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT reports_ticket_unique UNIQUE (ticket_number)
);

CREATE INDEX idx_reports_website_id    ON reports(website_id);
CREATE INDEX idx_reports_status        ON reports(status, created_at DESC);
CREATE INDEX idx_reports_reporter_user ON reports(reporter_user_id) WHERE reporter_user_id IS NOT NULL;
