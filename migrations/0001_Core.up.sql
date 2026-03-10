-- 0001_Core.up.sql
-- Core platform tables: users, websites, plans, payments

CREATE TABLE IF NOT EXISTS users (
    id                  TEXT         PRIMARY KEY,
    website_id          TEXT         NOT NULL,
    first_name          VARCHAR(50)  NOT NULL DEFAULT '',
    last_name           VARCHAR(50)  NOT NULL DEFAULT '',
    email               VARCHAR(250) NOT NULL,
    verified_email      BOOLEAN      NOT NULL DEFAULT FALSE,
    CHECK (email = LOWER(email)),
    password            TEXT         NOT NULL,
    role                VARCHAR(20)  NOT NULL DEFAULT 'user',
    avatar_url          TEXT,
    account_type        VARCHAR(10)  NOT NULL DEFAULT 'personal'
                        CHECK (account_type IN ('personal', 'business')),
    company_name        VARCHAR(100),
    trade_name          VARCHAR(100),
    cpf_cnpj            VARCHAR(20),
    phone               VARCHAR(20),
    zip_code            VARCHAR(10),
    address_street      VARCHAR(200),
    address_number      VARCHAR(20),
    address_complement  VARCHAR(100),
    address_city        VARCHAR(100),
    address_state       VARCHAR(2),
    address_country     VARCHAR(50)  DEFAULT 'BR',
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_website_unique UNIQUE (email, website_id)
);

CREATE INDEX IF NOT EXISTS idx_users_website_id  ON users(website_id);
CREATE INDEX IF NOT EXISTS idx_users_email        ON users(email);
-- Partial index for fast cleanup of expired unverified users
CREATE INDEX IF NOT EXISTS idx_users_unverified   ON users(created_at) WHERE verified_email = FALSE;

CREATE TABLE IF NOT EXISTS websites (
    id                TEXT         PRIMARY KEY,
    website_type      TEXT         NOT NULL DEFAULT 'ECOMMERCE'
                      CHECK (website_type IN ('JESTERX','ECOMMERCE','LANDING_PAGE','SOFTWARE_SELL','COURSE','VIDEO')),
    image             BYTEA        CHECK (octet_length(image) <= 2097152),
    name              VARCHAR(50)  NOT NULL,
    short_description VARCHAR(500),
    description       VARCHAR(1500),
    -- creator_id intentionally has no FK: platform website uses a placeholder ID
    creator_id        TEXT         NOT NULL,
    banned            BOOLEAN      NOT NULL DEFAULT FALSE,
    mature_content    BOOLEAN      NOT NULL DEFAULT FALSE,
    -- Denormalised from store_ratings for fast reads
    rating_avg        NUMERIC(3,2) NOT NULL DEFAULT 0.00,
    rating_count      INTEGER      NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT websites_name_unique UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_websites_creator_id ON websites(creator_id);
CREATE INDEX IF NOT EXISTS idx_websites_active     ON websites(banned) WHERE banned = FALSE;

CREATE TABLE IF NOT EXISTS plans (
    id             TEXT          PRIMARY KEY,
    name           VARCHAR(100)  NOT NULL,
    description    VARCHAR(500),
    description_md TEXT,
    price          NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    billing_cycle  VARCHAR(20)   NOT NULL DEFAULT 'monthly',
    active         BOOLEAN       NOT NULL DEFAULT TRUE,
    max_sites      INTEGER       NOT NULL DEFAULT 1  CHECK (max_sites  >= 0),
    max_routes     INTEGER       NOT NULL DEFAULT 5  CHECK (max_routes >= 0),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT plans_name_unique UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS payments (
    id           TEXT          PRIMARY KEY,
    user_id      TEXT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    website_id   TEXT          NOT NULL,
    plan_id      TEXT          REFERENCES plans(id) ON DELETE SET NULL,
    reference_id TEXT,
    type         VARCHAR(50)   NOT NULL DEFAULT 'subscription',
    quantity     INTEGER       NOT NULL DEFAULT 1  CHECK (quantity > 0),
    amount       NUMERIC(10,2) NOT NULL             CHECK (amount >= 0),
    currency     VARCHAR(10)   NOT NULL DEFAULT 'BRL',
    status       VARCHAR(20)   NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending','completed','canceled')),
    purchased_in TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_user_status  ON payments(user_id, status, purchased_in DESC);
CREATE INDEX IF NOT EXISTS idx_payments_reference_id ON payments(reference_id) WHERE reference_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_payments_website_id   ON payments(website_id);
