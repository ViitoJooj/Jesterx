-- 0002_Plans_Payments.up.sql
-- Platform subscription plans and payment records.
-- Payments are platform-level (per user account), not per website.

CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'canceled');

CREATE TABLE plans (
    id             TEXT          PRIMARY KEY,
    name           VARCHAR(100)  NOT NULL,
    description    VARCHAR(500),
    description_md TEXT,
    features       JSONB         NOT NULL DEFAULT '[]',
    price          NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    billing_cycle  VARCHAR(20)   NOT NULL DEFAULT 'monthly',
    active         BOOLEAN       NOT NULL DEFAULT TRUE,
    max_sites      INTEGER       NOT NULL DEFAULT 1  CHECK (max_sites  >= 0),
    max_routes     INTEGER       NOT NULL DEFAULT 5  CHECK (max_routes >= 0),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT plans_name_unique UNIQUE (name)
);

CREATE TABLE payments (
    id           TEXT           PRIMARY KEY,
    user_id      TEXT           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id      TEXT           REFERENCES plans(id) ON DELETE SET NULL,
    reference_id TEXT,
    type         VARCHAR(50)    NOT NULL DEFAULT 'subscription',
    quantity     INTEGER        NOT NULL DEFAULT 1  CHECK (quantity > 0),
    amount       NUMERIC(10,2)  NOT NULL             CHECK (amount >= 0),
    currency     VARCHAR(10)    NOT NULL DEFAULT 'BRL',
    status       payment_status NOT NULL DEFAULT 'pending',
    purchased_in TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user_status  ON payments(user_id, status, purchased_in DESC);
CREATE INDEX idx_payments_reference_id ON payments(reference_id) WHERE reference_id IS NOT NULL;
