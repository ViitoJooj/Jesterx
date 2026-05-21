-- 0000_Reset.up.sql
-- Drops everything and clears migration history.
-- Run this before applying the clean 0001-0008 migrations.

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
DROP TABLE IF EXISTS companies        CASCADE;
DROP TABLE IF EXISTS users            CASCADE;

DROP TYPE IF EXISTS gender_type       CASCADE;
DROP TYPE IF EXISTS member_role       CASCADE;
DROP TYPE IF EXISTS order_status      CASCADE;
DROP TYPE IF EXISTS payment_status    CASCADE;
DROP TYPE IF EXISTS product_condition CASCADE;
DROP TYPE IF EXISTS report_status     CASCADE;
DROP TYPE IF EXISTS report_reason     CASCADE;
DROP TYPE IF EXISTS scan_status_type  CASCADE;
DROP TYPE IF EXISTS source_type       CASCADE;
DROP TYPE IF EXISTS website_type      CASCADE;

TRUNCATE schema_migrations;
