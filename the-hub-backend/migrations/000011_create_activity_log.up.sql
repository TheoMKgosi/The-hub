-- Add activity_entries table for the personal activity log
CREATE TABLE IF NOT EXISTS activity_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    notes TEXT DEFAULT '',
    category TEXT DEFAULT '',
    tags TEXT[] DEFAULT '{}',
    activity_date DATE NOT NULL,
    duration_min INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activity_entries_user_id ON activity_entries(user_id);
CREATE INDEX idx_activity_entries_user_date ON activity_entries(user_id, activity_date DESC);
CREATE INDEX idx_activity_entries_category ON activity_entries(user_id, category);
