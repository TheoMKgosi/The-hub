-- Remove ai_checked column from tasks table
ALTER TABLE tasks DROP COLUMN IF EXISTS ai_checked;