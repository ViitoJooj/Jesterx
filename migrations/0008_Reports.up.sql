DO $$ BEGIN
    CREATE TYPE report_status AS ENUM ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'DISMISSED');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE report_reason AS ENUM ('SPAM', 'FRAUD', 'SCAM', 'INAPPROPRIATE', 'COUNTERFEIT', 'OTHER');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS reports (
    id              TEXT PRIMARY KEY,
    ticket_number   SERIAL NOT NULL UNIQUE,
    website_id      TEXT NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    reporter_name   TEXT NOT NULL,
    reporter_email  TEXT NOT NULL,
    reason          report_reason NOT NULL,
    description     TEXT NOT NULL,
    status          report_status NOT NULL DEFAULT 'OPEN',
    admin_response  TEXT,
    resolved_at     TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reports_website_id ON reports(website_id);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at DESC);
