CREATE TABLE organisations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    owner_id TEXT NOT NULL,
    website TEXT,
    logo_url TEXT,
    location TEXT,
    timezone TEXT DEFAULT 'UTC',
    is_active BOOLEAN DEFAULT true,
    archived_at TIMESTAMPTZ,
    settings JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes / constraints inferred from GORM tags
CREATE UNIQUE INDEX IF NOT EXISTS ux_organisations_name ON organisations(name);
CREATE UNIQUE INDEX IF NOT EXISTS ux_organisations_slug ON organisations(slug);