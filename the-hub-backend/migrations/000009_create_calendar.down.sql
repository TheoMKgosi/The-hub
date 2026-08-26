DROP INDEX IF EXISTS idx_schedule_user_id;
DROP INDEX IF EXISTS idx_recurrences_user_id;
DROP INDEX IF EXISTS idx_schedule_deleted_at;
DROP INDEX IF EXISTS idx_recurrences_deleted_at;
DROP TABLE IF EXISTS scheduled_tasks;
DROP TABLE IF EXISTS recurrences;
