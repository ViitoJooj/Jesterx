CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    website_id TEXT NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name  VARCHAR(50) NOT NULL,
    email VARCHAR(250) NOT NULL,
    verified_email BOOLEAN NOT NULL DEFAULT false,
    CHECK (email = LOWER(email)),
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS websites (
    id TEXT PRIMARY KEY,
    website_type TEXT NOT NULL CHECK (website_type IN ('ECOMMERCE','LANDING_PAGE','SOFTWARE_SELL','COURSE')),
    image BYTEA CHECK (octet_length(image) <= 2097152),
    name VARCHAR(50) NOT NULL UNIQUE,
    short_description VARCHAR(500),
    description VARCHAR(1500),
    creator_id TEXT NOT NULL,
    banned BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(500),
    description_md TEXT,
    price NUMERIC(10,2) NOT NULL,
    billing_cycle VARCHAR(20) DEFAULT 'monthly',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    website_id TEXT NOT NULL,
    reference_id TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    quantity INT DEFAULT 1,
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'BRL',
    status VARCHAR(50) NOT NULL CHECK (status IN ('completed','pending','canceled')),
    purchased_in TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);