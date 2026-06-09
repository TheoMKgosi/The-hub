-- Add ai_checked column to tasks table
ALTER TABLE tasks ADD COLUMN ai_checked BOOLEAN DEFAULT false;