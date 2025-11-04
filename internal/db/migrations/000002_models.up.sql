-- 000002_models.up.sql
-- Creates enum types and tables derived from GORM models in internal/db/models

-- Create enum types (if they don't already exist)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'project_status') THEN
        CREATE TYPE project_status AS ENUM ('Active', 'Completed', 'Archived');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'project_role') THEN
        CREATE TYPE project_role AS ENUM ('Admin', 'Edit', 'View');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organisation_role') THEN
        CREATE TYPE organisation_role AS ENUM (
            'CEO',
            'Project Manager',
            'IT Manager',
            'Senior Software Engineer',
            'Junior Software Engineer',
            'IT Support',
            'HR Manager',
            'Recruiter',
            'Not specified'
        );
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'kanban_status') THEN
        CREATE TYPE kanban_status AS ENUM ('Planning', 'In Progress', 'Done', 'Archived');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'kanban_item_priority') THEN
        CREATE TYPE kanban_item_priority AS ENUM ('Extreme', 'High', 'Medium', 'Low', 'None');
    END IF;
END$$;

-- Users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    username VARCHAR(50) NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_username ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email ON users(email);
CREATE INDEX IF NOT EXISTS ix_users_deleted_at ON users(deleted_at);

-- User sessions (refresh tokens)
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(21) NOT NULL,
    jti VARCHAR(21) NOT NULL,
    expire_date TIMESTAMPTZ NOT NULL,
    remember_me BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_sessions_jti ON user_sessions(jti);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(user_id, jti, expire_date);

-- Projects and members
CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    name VARCHAR(255),
    organisation_id VARCHAR(21) NOT NULL,
    status project_status DEFAULT 'Active',
    CONSTRAINT fk_projects_organisation FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_members (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(21) NOT NULL,
    user_id VARCHAR(21) NOT NULL,
    role project_role NOT NULL,
    CONSTRAINT fk_project_members_project FOREIGN KEY (project_id) REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_project_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_project_member_project_user ON project_members(project_id, user_id);

-- Organisation members
CREATE TABLE IF NOT EXISTS organisation_members (
    id SERIAL PRIMARY KEY,
    organisation_id VARCHAR(21) NOT NULL,
    user_id VARCHAR(21) NOT NULL,
    role organisation_role,
    CONSTRAINT fk_org_members_organisation FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_org_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_organisation_member_org_user ON organisation_members(organisation_id, user_id);

-- Kanban boards, categories and items
CREATE TABLE IF NOT EXISTS kanbans (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    project_id VARCHAR(21) NOT NULL,
    name VARCHAR(50) NOT NULL,
    status kanban_status NOT NULL DEFAULT 'Planning',
    CONSTRAINT fk_kanbans_project FOREIGN KEY (project_id) REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_kanban_project_name ON kanbans(project_id, name);

CREATE TABLE IF NOT EXISTS kanban_categories (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    kanban_id VARCHAR(21) NOT NULL,
    name VARCHAR(50),
    CONSTRAINT fk_kanban_categories_kanban FOREIGN KEY (kanban_id) REFERENCES kanbans(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kanban_items (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    kanban_category_id VARCHAR(21) NOT NULL,
    deleted_at TIMESTAMPTZ,
    priority kanban_item_priority NOT NULL DEFAULT 'None',
    due_date TIMESTAMPTZ,
    estimated_time INTEGER,
    title VARCHAR(40) NOT NULL,
    description TEXT,
    CONSTRAINT fk_kanban_items_category FOREIGN KEY (kanban_category_id) REFERENCES kanban_categories(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Chat: rooms, members and messages
CREATE TABLE IF NOT EXISTS chat_rooms (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    organisation_id VARCHAR(21),
    name VARCHAR(255),
    CONSTRAINT fk_chat_rooms_organisation FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    room_id VARCHAR(21) NOT NULL,
    user_id VARCHAR(21) NOT NULL,
    message TEXT NOT NULL,
    CONSTRAINT fk_chat_messages_room FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_chat_messages_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS ix_chat_messages_created_at ON chat_messages(created_at);

CREATE TABLE IF NOT EXISTS chat_members (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    room_id VARCHAR(21) NOT NULL,
    user_id VARCHAR(21) NOT NULL,
    last_message_id VARCHAR(21),
    CONSTRAINT fk_chat_members_room FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_chat_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_chat_members_last_message FOREIGN KEY (last_message_id) REFERENCES chat_messages(id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- Announcements
CREATE TABLE IF NOT EXISTS announcements (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    announcment_text TEXT,
    organisation_id VARCHAR(21)
);
ALTER TABLE announcements
    ADD CONSTRAINT fk_announcements_organisation FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- Whiteboards and line data/points
CREATE TABLE IF NOT EXISTS whiteboard_rooms (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    project_id VARCHAR(21),
    name VARCHAR(255),
    CONSTRAINT fk_whiteboard_project FOREIGN KEY (project_id) REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS line_data (
    id VARCHAR(21) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    whiteboard_id VARCHAR(21) NOT NULL,
    stroke TEXT,
    stroke_width INTEGER,
    tool TEXT,
    text_content TEXT,
    CONSTRAINT fk_line_data_whiteboard FOREIGN KEY (whiteboard_id) REFERENCES whiteboard_rooms(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS line_points (
    id SERIAL PRIMARY KEY,
    point DOUBLE PRECISION NOT NULL,
    line_data_id VARCHAR(21) NOT NULL,
    CONSTRAINT fk_line_points_line_data FOREIGN KEY (line_data_id) REFERENCES line_data(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- End of migration
