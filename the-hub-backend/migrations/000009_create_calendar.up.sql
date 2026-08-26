CREATE TABLE IF NOT EXISTS recurrences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    frequency TEXT NOT NULL,
    interval INTEGER DEFAULT 1,
    by_day TEXT,
    by_month INTEGER,
    by_month_day INTEGER,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    "count" INTEGER,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    "start" TIMESTAMP WITH TIME ZONE NOT NULL,
    "end" TIMESTAMP WITH TIME ZONE NOT NULL,
    recurrence_rule_id UUID REFERENCES recurrences(id),
    created_by_ai BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_schedules_user_id ON scheduled_tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_recurrences_user_id ON recurrences(user_id);
CREATE INDEX IF NOT EXISTS idx_schedules_deleted_at ON scheduled_tasks(deleted_at);
CREATE INDEX IF NOT EXISTS idx_recurrences_deleted_at ON recurrences(deleted_at);
