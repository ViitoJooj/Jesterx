-- extensão de dependencia
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Função para retornar data formatada
CREATE OR REPLACE FUNCTION format_datetime(value TIMESTAMPTZ)
RETURNS TEXT
LANGUAGE SQL
IMMUTABLE
AS $$
    SELECT TO_CHAR(value, 'DD/MM/YYYY:HH24:MI:SS');
$$;

-- Criando uuidv7
CREATE OR REPLACE FUNCTION uuid_v7()
RETURNS UUID
LANGUAGE PLPGSQL
VOLATILE
AS $$
DECLARE
    unix_ts_ms BIGINT;
    uuid_bytes BYTEA;
BEGIN
    unix_ts_ms := FLOOR(EXTRACT(EPOCH FROM CLOCK_TIMESTAMP()) * 1000);
    uuid_bytes := gen_random_bytes(16);
    uuid_bytes := set_byte(uuid_bytes, 0, ((unix_ts_ms >> 40) & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 1, ((unix_ts_ms >> 32) & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 2, ((unix_ts_ms >> 24) & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 3, ((unix_ts_ms >> 16) & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 4, ((unix_ts_ms >> 8) & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 5, (unix_ts_ms & 255)::INT);
    uuid_bytes := set_byte(uuid_bytes, 6, (get_byte(uuid_bytes, 6) & 15) | 112);
    uuid_bytes := set_byte(uuid_bytes, 8, (get_byte(uuid_bytes, 8) & 63) | 128);
    RETURN encode(uuid_bytes, 'hex')::UUID;
END;
$$;