ALTER TABLE organisations
ADD COLUMN updated_at TIMESTAMPTZ DEFAULT now(); 