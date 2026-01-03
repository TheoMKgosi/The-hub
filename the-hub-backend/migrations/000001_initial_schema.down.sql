-- Rollback initial schema migration

-- Drop junction tables first (due to foreign key constraints)
DROP TABLE IF EXISTS task_templates;
DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS deck_users;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS tags;

-- Drop main tables
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS decks;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS budget_categories;
DROP TABLE IF EXISTS users;
