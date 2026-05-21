CREATE TABLE IF NOT EXISTS user_addresses (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label       VARCHAR(50),
    zip_code    VARCHAR(10),
    street      VARCHAR(200),
    number      VARCHAR(20),
    complement  VARCHAR(100),
    district    VARCHAR(100),
    city        VARCHAR(100),
    state       VARCHAR(2),
    country     VARCHAR(50) NOT NULL DEFAULT 'BR',
    is_default  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS user_addresses_one_default
    ON user_addresses (user_id)
    WHERE is_default = TRUE;

INSERT INTO user_addresses (id, user_id, zip_code, street, number, complement, district, city, state, country, is_default, created_at, updated_at)
SELECT
    gen_random_uuid(),
    id,
    zip_code,
    address_street,
    address_number,
    address_complement,
    address_district,
    address_city,
    address_state,
    COALESCE(address_country, 'BR'),
    TRUE,
    created_at,
    updated_at
FROM users
WHERE address_street IS NOT NULL AND address_street <> '';

ALTER TABLE users
    DROP COLUMN IF EXISTS zip_code,
    DROP COLUMN IF EXISTS address_street,
    DROP COLUMN IF EXISTS address_number,
    DROP COLUMN IF EXISTS address_complement,
    DROP COLUMN IF EXISTS address_district,
    DROP COLUMN IF EXISTS address_city,
    DROP COLUMN IF EXISTS address_state,
    DROP COLUMN IF EXISTS address_country;
