-- Rollback initial schema migration

-- Drop junction tables first (due to foreign key constraints)
DROP TABLE IF EXISTS task_templates;
DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS deck_users;

-- Drop main tables
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS users;

-- Note: We don't drop the uuid-ossp extension as it might be used by other databases