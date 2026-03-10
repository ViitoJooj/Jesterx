-- 0000_Reset.up.sql
-- Drops all application tables and clears migration history so the
-- new clean migrations (0001-0006) are applied from scratch.
-- schema_migrations itself is preserved (the runner recreates it).

DROP TABLE IF EXISTS store_members    CASCADE;
DROP TABLE IF EXISTS store_visits     CASCADE;
DROP TABLE IF EXISTS store_ratings    CASCADE;
DROP TABLE IF EXISTS store_comments   CASCADE;
DROP TABLE IF EXISTS reports          CASCADE;
DROP TABLE IF EXISTS order_items      CASCADE;
DROP TABLE IF EXISTS orders           CASCADE;
DROP TABLE IF EXISTS products         CASCADE;
DROP TABLE IF EXISTS website_versions CASCADE;
DROP TABLE IF EXISTS website_routes   CASCADE;
DROP TABLE IF EXISTS themes           CASCADE;
DROP TABLE IF EXISTS payments         CASCADE;
DROP TABLE IF EXISTS plans            CASCADE;
DROP TABLE IF EXISTS websites         CASCADE;
DROP TABLE IF EXISTS users            CASCADE;

-- Drop legacy ENUM types if present
DROP TYPE IF EXISTS report_status CASCADE;
DROP TYPE IF EXISTS report_reason CASCADE;

-- Clear migration history so 0001-0006 are treated as new
TRUNCATE schema_migrations;
