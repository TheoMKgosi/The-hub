CREATE TABLE IF NOT EXISTS financial_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    target_amount DECIMAL(10,2),
    current_amount DECIMAL(10,2) DEFAULT 0,
    type TEXT NOT NULL CHECK (type IN ('savings', 'checklist')),
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled')),
    priority INTEGER CHECK (priority >= 1 AND priority <= 5),
    target_date TIMESTAMP WITH TIME ZONE,
    category_id UUID REFERENCES budget_categories(id) ON DELETE SET NULL,
    color TEXT DEFAULT '#3B82F6',
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_financial_goals_user_id ON financial_goals(user_id);
CREATE INDEX IF NOT EXISTS idx_financial_goals_deleted_at ON financial_goals(deleted_at);
CREATE INDEX IF NOT EXISTS idx_financial_goals_category_id ON financial_goals(category_id);
