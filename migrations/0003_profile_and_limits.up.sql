-- Migration 0003: user profile fields + plan limits

-- User profile additions
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS cpf_cnpj   VARCHAR(20)  DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS avatar_url TEXT         DEFAULT NULL;

-- Plan limits
ALTER TABLE plans
  ADD COLUMN IF NOT EXISTS max_sites  INT NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS max_routes INT NOT NULL DEFAULT 5;

-- Seed sensible defaults for existing plans (update by name pattern)
UPDATE plans SET max_sites = 1,  max_routes = 5   WHERE LOWER(name) LIKE '%starter%'  OR LOWER(name) LIKE '%basic%'    OR LOWER(name) LIKE '%essencial%';
UPDATE plans SET max_sites = 5,  max_routes = 30  WHERE LOWER(name) LIKE '%pro%'      OR LOWER(name) LIKE '%business%';
UPDATE plans SET max_sites = 20, max_routes = 100 WHERE LOWER(name) LIKE '%enterprise%' OR LOWER(name) LIKE '%ultra%'  OR LOWER(name) LIKE '%scale%';
