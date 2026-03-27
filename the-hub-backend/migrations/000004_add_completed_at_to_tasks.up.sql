-- Add completed_at timestamp to tasks

ALTER TABLE tasks
  ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE;

-- Backfill existing completed tasks (best effort)
UPDATE tasks
SET completed_at = updated_at
WHERE completed_at IS NULL
  AND status IN ('completed', 'complete');

-- Normalize legacy status values
UPDATE tasks
SET status = 'completed'
WHERE status = 'complete';

CREATE INDEX IF NOT EXISTS idx_tasks_completed_at
  ON tasks(user_id, completed_at);
