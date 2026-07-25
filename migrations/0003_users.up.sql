CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_v7(),
    website_uuid UUID,
    image_url VARCHAR(1000),
    name VARCHAR(250),
    email VARCHAR(250) UNIQUE NOT NULL,
    role UUID NOT NULL,
    password NOT NULL,
    cpf VARCHAR(11) UNIQUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT
);

CREATE INDEX  IF NOT EXISTS idx_user_website ON users (website_uuid);
CREATE INDEX  IF NOT EXISTS idx_user_name ON users (name);
CREATE INDEX  IF NOT EXISTS idx_user_email ON users (email);