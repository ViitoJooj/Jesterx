-- 0001_Users_Companies.up.sql
-- Platform users and business companies.
-- Users are multi-tenant (each website has its own user base).
-- Companies are separate entities linked to a platform user (owner).

CREATE TYPE gender_type AS ENUM ('male', 'female', 'other', 'prefer_not');

CREATE TABLE users (
    id                  TEXT         PRIMARY KEY,
    website_id          TEXT         NOT NULL,
    first_name          VARCHAR(50)  NOT NULL DEFAULT '',
    last_name           VARCHAR(50)  NOT NULL DEFAULT '',
    email               VARCHAR(250) NOT NULL CHECK (email = LOWER(email)),
    verified_email      BOOLEAN      NOT NULL DEFAULT FALSE,
    password            TEXT         NOT NULL,
    role                VARCHAR(20)  NOT NULL DEFAULT 'user',
    avatar_url          TEXT,
    cpf                 VARCHAR(14),
    phone               VARCHAR(20),
    display_name        VARCHAR(100),
    birth_date          DATE,
    gender              gender_type,
    bio                 VARCHAR(500),
    instagram           VARCHAR(100),
    website_url         VARCHAR(200),
    whatsapp            VARCHAR(20),
    zip_code            VARCHAR(10),
    address_street      VARCHAR(200),
    address_number      VARCHAR(20),
    address_complement  VARCHAR(100),
    address_district    VARCHAR(100),
    address_city        VARCHAR(100),
    address_state       VARCHAR(2),
    address_country     VARCHAR(50)  DEFAULT 'BR',
    is_active           BOOLEAN      NOT NULL DEFAULT TRUE,
    deactivated_at      TIMESTAMPTZ,
    delete_after        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_website_unique UNIQUE (email, website_id)
);

CREATE INDEX idx_users_website_id   ON users(website_id);
CREATE INDEX idx_users_email        ON users(email);
-- Fast lookup of unverified users for periodic cleanup
CREATE INDEX idx_users_unverified   ON users(created_at) WHERE verified_email = FALSE;
-- Fast lookup of accounts scheduled for deletion
CREATE INDEX idx_users_delete_after ON users(delete_after) WHERE is_active = FALSE;

-- A company is a business entity created by a platform user.
-- When registering as a business, both a user record and a company record are created.
-- If the email already exists in users, only the company record is created.
CREATE TABLE companies (
    id                 TEXT         PRIMARY KEY,
    owner_user_id      TEXT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name       VARCHAR(100) NOT NULL,
    trade_name         VARCHAR(100),
    cnpj               VARCHAR(18),
    phone              VARCHAR(20),
    zip_code           VARCHAR(10),
    address_street     VARCHAR(200),
    address_number     VARCHAR(20),
    address_complement VARCHAR(100),
    address_district   VARCHAR(100),
    address_city       VARCHAR(100),
    address_state      VARCHAR(2),
    address_country    VARCHAR(50)  DEFAULT 'BR',
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_companies_owner ON companies(owner_user_id);
CREATE UNIQUE INDEX idx_companies_cnpj ON companies(cnpj) WHERE cnpj IS NOT NULL AND cnpj <> '';
