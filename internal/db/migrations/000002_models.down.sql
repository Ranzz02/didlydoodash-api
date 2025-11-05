-- 000002_models.down.sql
-- Drops tables and types created in 000002_models.up.sql (safe order)

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS line_points CASCADE;
DROP TABLE IF EXISTS line_data CASCADE;
DROP TABLE IF EXISTS whiteboard_rooms CASCADE;
DROP TABLE IF EXISTS announcements CASCADE;
DROP TABLE IF EXISTS chat_members CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_rooms CASCADE;
DROP TABLE IF EXISTS kanban_items CASCADE;
DROP TABLE IF EXISTS kanban_categories CASCADE;
DROP TABLE IF EXISTS kanbans CASCADE;
DROP TABLE IF EXISTS project_members CASCADE;
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS organisation_members CASCADE;

-- Drop enum types (if no longer used)
DROP TYPE IF EXISTS kanban_item_priority;
DROP TYPE IF EXISTS kanban_status;
DROP TYPE IF EXISTS organisation_role;
DROP TYPE IF EXISTS project_role;
DROP TYPE IF EXISTS project_status;

-- End of down migration
