-- Theme days for tasks
-- A user defines themes (e.g. "Deep Work", "Admin") and assigns weekdays to them.
-- On a given day, only tasks belonging to that day's theme are displayed when
-- the theme days mode is enabled.

CREATE TABLE IF NOT EXISTS themes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  color VARCHAR(20) DEFAULT '#3B82F6',
  icon VARCHAR(100),
  days TEXT[] NOT NULL DEFAULT '{}',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_themes_user_id ON themes(user_id);
CREATE INDEX IF NOT EXISTS idx_themes_deleted_at ON themes(deleted_at);

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS theme_id UUID REFERENCES themes(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_theme_id ON tasks(theme_id);
