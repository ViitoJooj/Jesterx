CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL UNIQUE,
    website_id TEXT NOT NULL,
    avatar_url TEXT DEFAULT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name  VARCHAR(50) NOT NULL,
    email VARCHAR(250) NOT NULL,
    verified_email BOOLEAN NOT NULL DEFAULT false,
    CHECK (email = LOWER(email)),
    cpf_cnpj VARCHAR(20) DEFAULT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    account_type VARCHAR(10) NOT NULL DEFAULT 'personal' CHECK (account_type IN ('personal', 'business')),
    company_name VARCHAR(100) DEFAULT NULL,
    trade_name VARCHAR(100) DEFAULT NULL,
    phone VARCHAR(20) DEFAULT NULL,
    zip_code VARCHAR(10) DEFAULT NULL,
    address_street VARCHAR(200) DEFAULT NULL,
    address_number VARCHAR(20) DEFAULT NULL,
    address_complement VARCHAR(100) DEFAULT NULL,
    address_city VARCHAR(100) DEFAULT NULL,
    address_state VARCHAR(2) DEFAULT NULL,
    address_country VARCHAR(50) DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS websites (
    id TEXT PRIMARY KEY,
    website_type TEXT NOT NULL CHECK (website_type IN ('JESTERX', 'ECOMMERCE','LANDING_PAGE','SOFTWARE_SELL','COURSE','VIDEO')),
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
    max_sites  INT NOT NULL DEFAULT 1,
    max_routes INT NOT NULL DEFAULT 5,
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

CREATE TABLE IF NOT EXISTS website_routes (
    id TEXT PRIMARY KEY NOT NULL,
    website_id TEXT NOT NULL,
    path VARCHAR(180) NOT NULL,
    title VARCHAR(100) NOT NULL,
    requires_auth BOOLEAN NOT NULL DEFAULT false,
    position INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (website_id, path)
);

CREATE TABLE IF NOT EXISTS website_versions (
    id TEXT PRIMARY KEY NOT NULL,
    website_id TEXT NOT NULL,
    version INT NOT NULL,
    source_type VARCHAR(30) NOT NULL CHECK (source_type IN ('JXML','REACT','SVELTE','ELEMENTOR_JSON')),
    source TEXT NOT NULL,
    compiled_html TEXT NOT NULL,
    scan_status VARCHAR(20) NOT NULL CHECK (scan_status IN ('clean','warning','blocked')),
    scan_score INT NOT NULL DEFAULT 100,
    scan_findings TEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMP NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (website_id, version)
);

CREATE TABLE products (
    id           TEXT          PRIMARY KEY,
    website_id   TEXT          NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    name         VARCHAR(200)  NOT NULL,
    description  TEXT          NOT NULL DEFAULT '',
    price        NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    compare_price NUMERIC(12,2) CHECK (compare_price IS NULL OR compare_price >= 0),
    stock        INTEGER       NOT NULL DEFAULT 0 CHECK (stock >= 0),
    sku          VARCHAR(100),
    category     VARCHAR(100),
    images       JSONB         NOT NULL DEFAULT '[]',
    sold_count   INTEGER       NOT NULL DEFAULT 0 CHECK (sold_count >= 0),
    active       BOOLEAN       NOT NULL DEFAULT true,
    created_by   TEXT          NOT NULL REFERENCES users(id),
    updated_at   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);