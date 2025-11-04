-- 000001_init.down.sql
-- Revert the initial organisations table created by 000001_init.up.sql

DROP INDEX IF EXISTS ux_organisations_name;
DROP INDEX IF EXISTS ux_organisations_slug;
DROP TABLE IF EXISTS organisations CASCADE;
