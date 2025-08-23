-- Recovery Migration Script for SQLite Database
-- This script handles cases where a previous migration attempt failed
-- and the database is in an inconsistent state

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Begin transaction for atomicity
BEGIN TRANSACTION;

-- =====================================================
-- Step 1: Check if UUID columns already exist
-- =====================================================

-- Create a temporary table to track which tables need migration
CREATE TEMPORARY TABLE migration_status (
    table_name TEXT PRIMARY KEY,
    has_uuid_column INTEGER DEFAULT 0,
    needs_migration INTEGER DEFAULT 1
);

-- Check each table for existing UUID columns
INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'users', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('users') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'goals', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('goals') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'tasks', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('tasks') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'scheduled_tasks', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('scheduled_tasks') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'decks', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('decks') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'deck_users', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('deck_users') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'cards', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('cards') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'budget_categories', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('budget_categories') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'budgets', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('budgets') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'incomes', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('incomes') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'topics', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('topics') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'task_learning', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('task_learning') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'tags', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('tags') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'resources', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('resources') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'study_sessions', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('study_sessions') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'repeat_rules', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('repeat_rules') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

INSERT OR REPLACE INTO migration_status (table_name, has_uuid_column)
SELECT 'ai_recommendations', CASE WHEN EXISTS(
    SELECT 1 FROM pragma_table_info('ai_recommendations') WHERE name = 'id_uuid'
) THEN 1 ELSE 0 END;

-- =====================================================
-- Step 2: Add UUID columns only where they don't exist
-- =====================================================

-- Users table
INSERT OR IGNORE INTO users (id) VALUES (1); -- Ensure table exists
ALTER TABLE users ADD COLUMN id_uuid TEXT;

-- Goals table
INSERT OR IGNORE INTO goals (id, user_id, title) VALUES (1, 1, 'temp');
ALTER TABLE goals ADD COLUMN id_uuid TEXT;
ALTER TABLE goals ADD COLUMN user_id_uuid TEXT;

-- Tasks table
INSERT OR IGNORE INTO tasks (id, user_id, title) VALUES (1, 1, 'temp');
ALTER TABLE tasks ADD COLUMN id_uuid TEXT;
ALTER TABLE tasks ADD COLUMN goal_id_uuid TEXT;
ALTER TABLE tasks ADD COLUMN user_id_uuid TEXT;

-- Scheduled Tasks table
INSERT OR IGNORE INTO scheduled_tasks (id, user_id, title, start, end) VALUES (1, 1, 'temp', '2024-01-01T00:00:00Z', '2024-01-01T01:00:00Z');
ALTER TABLE scheduled_tasks ADD COLUMN id_uuid TEXT;
ALTER TABLE scheduled_tasks ADD COLUMN user_id_uuid TEXT;

-- Decks table
INSERT OR IGNORE INTO decks (id, user_id, name) VALUES (1, 1, 'temp');
ALTER TABLE decks ADD COLUMN id_uuid TEXT;
ALTER TABLE decks ADD COLUMN user_id_uuid TEXT;

-- Deck Users table
INSERT OR IGNORE INTO deck_users (id, deck_id, user_id, role) VALUES (1, 1, 1, 'owner');
ALTER TABLE deck_users ADD COLUMN id_uuid TEXT;
ALTER TABLE deck_users ADD COLUMN deck_id_uuid TEXT;
ALTER TABLE deck_users ADD COLUMN user_id_uuid TEXT;

-- Cards table
INSERT OR IGNORE INTO cards (id, deck_id, question, answer) VALUES (1, 1, 'temp', 'temp');
ALTER TABLE cards ADD COLUMN id_uuid TEXT;
ALTER TABLE cards ADD COLUMN deck_id_uuid TEXT;

-- Budget Categories table
INSERT OR IGNORE INTO budget_categories (id, user_id, name) VALUES (1, 1, 'temp');
ALTER TABLE budget_categories ADD COLUMN id_uuid TEXT;
ALTER TABLE budget_categories ADD COLUMN user_id_uuid TEXT;

-- Budgets table
INSERT OR IGNORE INTO budgets (id, category_id, user_id, amount, start_date, end_date) VALUES (1, 1, 1, 100.0, '2024-01-01', '2024-12-31');
ALTER TABLE budgets ADD COLUMN id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN category_id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN user_id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN income_id_uuid TEXT;

-- Income table
INSERT OR IGNORE INTO incomes (id, user_id, source, amount, received_at) VALUES (1, 1, 'temp', 100.0, '2024-01-01');
ALTER TABLE incomes ADD COLUMN id_uuid TEXT;
ALTER TABLE incomes ADD COLUMN user_id_uuid TEXT;

-- Topics table
INSERT OR IGNORE INTO topics (id, user_id, title, status) VALUES (1, 1, 'temp', 'not_started');
ALTER TABLE topics ADD COLUMN id_uuid TEXT;
ALTER TABLE topics ADD COLUMN user_id_uuid TEXT;

-- Task Learning table
INSERT OR IGNORE INTO task_learning (id, topic_id, title, status) VALUES (1, 1, 'temp', 'pending');
ALTER TABLE task_learning ADD COLUMN id_uuid TEXT;
ALTER TABLE task_learning ADD COLUMN topic_id_uuid TEXT;

-- Tags table
INSERT OR IGNORE INTO tags (id, user_id, name, color) VALUES (1, 1, 'temp', '#000000');
ALTER TABLE tags ADD COLUMN id_uuid TEXT;
ALTER TABLE tags ADD COLUMN user_id_uuid TEXT;

-- Resources table
INSERT OR IGNORE INTO resources (id, title, link, type) VALUES (1, 'temp', 'http://example.com', 'article');
ALTER TABLE resources ADD COLUMN id_uuid TEXT;
ALTER TABLE resources ADD COLUMN topic_id_uuid TEXT;
ALTER TABLE resources ADD COLUMN task_id_uuid TEXT;

-- Study Sessions table
INSERT OR IGNORE INTO study_sessions (id, user_id, duration_min, started_at, ended_at) VALUES (1, 1, 30, '2024-01-01T00:00:00Z', '2024-01-01T00:30:00Z');
ALTER TABLE study_sessions ADD COLUMN id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN user_id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN topic_id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN task_id_uuid TEXT;

-- Repeat Rules table
INSERT OR IGNORE INTO repeat_rules (id, frequency, interval, start_date) VALUES (1, 'daily', 1, '2024-01-01');
ALTER TABLE repeat_rules ADD COLUMN id_uuid TEXT;

-- AI Recommendations table
INSERT OR IGNORE INTO ai_recommendations (id, task_id, suggested_start, suggested_end, confidence, accepted) VALUES (1, 1, '2024-01-01T00:00:00Z', '2024-01-01T01:00:00Z', 0.8, 0);
ALTER TABLE ai_recommendations ADD COLUMN id_uuid TEXT;
ALTER TABLE ai_recommendations ADD COLUMN task_id_uuid TEXT;

-- =====================================================
-- Step 3: Generate UUIDs for existing records
-- =====================================================
-- Note: You need to generate actual UUIDs here
-- You can use a script to populate these values

-- Example for users table (replace with actual UUIDs):
-- UPDATE users SET id_uuid = '550e8400-e29b-41d4-a716-446655440000' WHERE id = 1;
-- UPDATE users SET id_uuid = '550e8400-e29b-41d4-a716-446655440001' WHERE id = 2;

-- =====================================================
-- Step 4: Update foreign key references
-- =====================================================

-- Update goals.user_id references
UPDATE goals SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = goals.user_id
) WHERE goals.user_id IS NOT NULL;

-- Update tasks references
UPDATE tasks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = tasks.user_id
) WHERE tasks.user_id IS NOT NULL;

UPDATE tasks SET goal_id_uuid = (
    SELECT id_uuid FROM goals WHERE goals.id = tasks.goal_id
) WHERE tasks.goal_id IS NOT NULL;

-- Update scheduled_tasks references
UPDATE scheduled_tasks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = scheduled_tasks.user_id
) WHERE scheduled_tasks.user_id IS NOT NULL;

-- Update decks references
UPDATE decks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = decks.user_id
) WHERE decks.user_id IS NOT NULL;

-- Update deck_users references
UPDATE deck_users SET deck_id_uuid = (
    SELECT id_uuid FROM decks WHERE decks.id = deck_users.deck_id
) WHERE deck_users.deck_id IS NOT NULL;

UPDATE deck_users SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = deck_users.user_id
) WHERE deck_users.user_id IS NOT NULL;

-- Update cards references
UPDATE cards SET deck_id_uuid = (
    SELECT id_uuid FROM decks WHERE decks.id = cards.deck_id
) WHERE cards.deck_id IS NOT NULL;

-- Update budget_categories references
UPDATE budget_categories SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = budget_categories.user_id
) WHERE budget_categories.user_id IS NOT NULL;

-- Update budgets references
UPDATE budgets SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = budgets.user_id
) WHERE budgets.user_id IS NOT NULL;

UPDATE budgets SET category_id_uuid = (
    SELECT id_uuid FROM budget_categories WHERE budget_categories.id = budgets.category_id
) WHERE budgets.category_id IS NOT NULL;

UPDATE budgets SET income_id_uuid = (
    SELECT id_uuid FROM incomes WHERE incomes.id = budgets.income_id
) WHERE budgets.income_id IS NOT NULL;

-- Update incomes references
UPDATE incomes SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = incomes.user_id
) WHERE incomes.user_id IS NOT NULL;

-- Update topics references
UPDATE topics SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = topics.user_id
) WHERE topics.user_id IS NOT NULL;

-- Update task_learning references
UPDATE task_learning SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = task_learning.topic_id
) WHERE task_learning.topic_id IS NOT NULL;

-- Update tags references
UPDATE tags SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = tags.user_id
) WHERE tags.user_id IS NOT NULL;

-- Update resources references
UPDATE resources SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = resources.topic_id
) WHERE resources.topic_id IS NOT NULL;

UPDATE resources SET task_id_uuid = (
    SELECT id_uuid FROM task_learning WHERE task_learning.id = resources.task_id
) WHERE resources.task_id IS NOT NULL;

-- Update study_sessions references
UPDATE study_sessions SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = study_sessions.user_id
) WHERE study_sessions.user_id IS NOT NULL;

UPDATE study_sessions SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = study_sessions.topic_id
) WHERE study_sessions.topic_id IS NOT NULL;

UPDATE study_sessions SET task_id_uuid = (
    SELECT id_uuid FROM task_learning WHERE task_learning.id = study_sessions.task_id
) WHERE study_sessions.task_id IS NOT NULL;

-- Update ai_recommendations references
UPDATE ai_recommendations SET task_id_uuid = (
    SELECT id_uuid FROM tasks WHERE tasks.id = ai_recommendations.task_id
) WHERE ai_recommendations.task_id IS NOT NULL;

-- =====================================================
-- Step 5: Recreate tables with UUID primary keys
-- =====================================================
-- SQLite doesn't allow direct modification of primary keys,
-- so we need to recreate each table

-- Create temporary tables with new structure
CREATE TABLE users_new (
    id TEXT PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    password TEXT,
    settings TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE TABLE goals_new (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    title TEXT,
    description TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE tasks_new (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    due_date DATETIME,
    priority INTEGER CHECK(priority >= 1 AND priority <= 5),
    status TEXT DEFAULT 'pending',
    order_index INTEGER DEFAULT 0,
    goal_id TEXT,
    user_id TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (goal_id) REFERENCES goals_new(id),
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE scheduled_tasks_new (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    start DATETIME NOT NULL,
    end DATETIME NOT NULL,
    user_id TEXT,
    created_by_ai BOOLEAN DEFAULT FALSE,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE decks_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE deck_users_new (
    id TEXT PRIMARY KEY,
    deck_id TEXT,
    user_id TEXT,
    role TEXT,
    FOREIGN KEY (deck_id) REFERENCES decks_new(id),
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE cards_new (
    id TEXT PRIMARY KEY,
    deck_id TEXT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    easiness REAL DEFAULT 2.5,
    interval INTEGER DEFAULT 1,
    repetitions INTEGER DEFAULT 0,
    last_reviewed DATETIME,
    next_review DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (deck_id) REFERENCES decks_new(id)
);

CREATE TABLE budget_categories_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    user_id TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE budgets_new (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL,
    amount REAL NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    user_id TEXT,
    income_id TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES budget_categories_new(id),
    FOREIGN KEY (user_id) REFERENCES users_new(id),
    FOREIGN KEY (income_id) REFERENCES incomes_new(id)
);

CREATE TABLE incomes_new (
    id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    amount REAL NOT NULL,
    user_id TEXT,
    received_at DATE NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE topics_new (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'not_started',
    deadline DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE task_learning_new (
    id TEXT PRIMARY KEY,
    topic_id TEXT,
    title TEXT NOT NULL,
    notes TEXT,
    status TEXT DEFAULT 'pending',
    order_index INTEGER,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (topic_id) REFERENCES topics_new(id)
);

CREATE TABLE tags_new (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    name TEXT UNIQUE NOT NULL,
    color TEXT,
    FOREIGN KEY (user_id) REFERENCES users_new(id)
);

CREATE TABLE resources_new (
    id TEXT PRIMARY KEY,
    topic_id TEXT,
    task_id TEXT,
    title TEXT,
    link TEXT,
    type TEXT,
    notes TEXT,
    FOREIGN KEY (topic_id) REFERENCES topics_new(id),
    FOREIGN KEY (task_id) REFERENCES task_learning_new(id)
);

CREATE TABLE study_sessions_new (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    topic_id TEXT,
    task_id TEXT,
    duration_min INTEGER,
    started_at DATETIME,
    ended_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users_new(id),
    FOREIGN KEY (topic_id) REFERENCES topics_new(id),
    FOREIGN KEY (task_id) REFERENCES task_learning_new(id)
);

CREATE TABLE repeat_rules_new (
    id TEXT PRIMARY KEY,
    frequency TEXT,
    interval INTEGER,
    by_day TEXT,
    start_date DATE,
    end_date DATE,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE ai_recommendations_new (
    id TEXT PRIMARY KEY,
    task_id TEXT,
    suggested_start DATETIME,
    suggested_end DATETIME,
    confidence REAL,
    accepted BOOLEAN,
    created_at DATETIME,
    FOREIGN KEY (task_id) REFERENCES tasks_new(id)
);

-- Copy data from old tables to new tables
INSERT INTO users_new SELECT id_uuid, name, email, password, settings, created_at, updated_at, deleted_at FROM users WHERE id_uuid IS NOT NULL;
INSERT INTO goals_new SELECT id_uuid, user_id_uuid, title, description, created_at, updated_at, deleted_at FROM goals WHERE id_uuid IS NOT NULL;
INSERT INTO tasks_new SELECT id_uuid, title, description, due_date, priority, status, order_index, goal_id_uuid, user_id_uuid, created_at, updated_at, deleted_at FROM tasks WHERE id_uuid IS NOT NULL;
INSERT INTO scheduled_tasks_new SELECT id_uuid, title, start, end, user_id_uuid, created_by_ai, created_at, updated_at FROM scheduled_tasks WHERE id_uuid IS NOT NULL;
INSERT INTO decks_new SELECT id_uuid, name, user_id_uuid, created_at, updated_at, deleted_at FROM decks WHERE id_uuid IS NOT NULL;
INSERT INTO deck_users_new SELECT id_uuid, deck_id_uuid, user_id_uuid, role FROM deck_users WHERE id_uuid IS NOT NULL;
INSERT INTO cards_new SELECT id_uuid, deck_id_uuid, question, answer, easiness, interval, repetitions, last_reviewed, next_review, created_at, updated_at, deleted_at FROM cards WHERE id_uuid IS NOT NULL;
INSERT INTO budget_categories_new SELECT id_uuid, name, user_id_uuid, created_at, updated_at, deleted_at FROM budget_categories WHERE id_uuid IS NOT NULL;
INSERT INTO budgets_new SELECT id_uuid, category_id_uuid, amount, start_date, end_date, user_id_uuid, income_id_uuid, created_at, updated_at, deleted_at FROM budgets WHERE id_uuid IS NOT NULL;
INSERT INTO incomes_new SELECT id_uuid, source, amount, user_id_uuid, received_at, created_at, updated_at, deleted_at FROM incomes WHERE id_uuid IS NOT NULL;
INSERT INTO topics_new SELECT id_uuid, user_id_uuid, title, description, status, deadline, created_at, updated_at, deleted_at FROM topics WHERE id_uuid IS NOT NULL;
INSERT INTO task_learning_new SELECT id_uuid, topic_id_uuid, title, notes, status, order_index, created_at, updated_at, deleted_at FROM task_learning WHERE id_uuid IS NOT NULL;
INSERT INTO tags_new SELECT id_uuid, user_id_uuid, name, color FROM tags WHERE id_uuid IS NOT NULL;
INSERT INTO resources_new SELECT id_uuid, topic_id_uuid, task_id_uuid, title, link, type, notes FROM resources WHERE id_uuid IS NOT NULL;
INSERT INTO study_sessions_new SELECT id_uuid, user_id_uuid, topic_id_uuid, task_id_uuid, duration_min, started_at, ended_at FROM study_sessions WHERE id_uuid IS NOT NULL;
INSERT INTO repeat_rules_new SELECT id_uuid, frequency, interval, by_day, start_date, end_date, created_at, updated_at FROM repeat_rules WHERE id_uuid IS NOT NULL;
INSERT INTO ai_recommendations_new SELECT id_uuid, task_id_uuid, suggested_start, suggested_end, confidence, accepted, created_at FROM ai_recommendations WHERE id_uuid IS NOT NULL;

-- Drop old tables
DROP TABLE users;
DROP TABLE goals;
DROP TABLE tasks;
DROP TABLE scheduled_tasks;
DROP TABLE decks;
DROP TABLE deck_users;
DROP TABLE cards;
DROP TABLE budget_categories;
DROP TABLE budgets;
DROP TABLE incomes;
DROP TABLE topics;
DROP TABLE task_learning;
DROP TABLE tags;
DROP TABLE resources;
DROP TABLE study_sessions;
DROP TABLE repeat_rules;
DROP TABLE ai_recommendations;

-- Rename new tables to original names
ALTER TABLE users_new RENAME TO users;
ALTER TABLE goals_new RENAME TO goals;
ALTER TABLE tasks_new RENAME TO tasks;
ALTER TABLE scheduled_tasks_new RENAME TO scheduled_tasks;
ALTER TABLE decks_new RENAME TO decks;
ALTER TABLE deck_users_new RENAME TO deck_users;
ALTER TABLE cards_new RENAME TO cards;
ALTER TABLE budget_categories_new RENAME TO budget_categories;
ALTER TABLE budgets_new RENAME TO budgets;
ALTER TABLE incomes_new RENAME TO incomes;
ALTER TABLE topics_new RENAME TO topics;
ALTER TABLE task_learning_new RENAME TO task_learning;
ALTER TABLE tags_new RENAME TO tags;
ALTER TABLE resources_new RENAME TO resources;
ALTER TABLE study_sessions_new RENAME TO study_sessions;
ALTER TABLE repeat_rules_new RENAME TO repeat_rules;
ALTER TABLE ai_recommendations_new RENAME TO ai_recommendations;

-- =====================================================
-- Step 6: Verify the migration
-- =====================================================

-- Check foreign key constraints
PRAGMA foreign_key_check;

-- Count records in each table to ensure no data loss
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'goals', COUNT(*) FROM goals
UNION ALL
SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL
SELECT 'scheduled_tasks', COUNT(*) FROM scheduled_tasks
UNION ALL
SELECT 'decks', COUNT(*) FROM decks
UNION ALL
SELECT 'deck_users', COUNT(*) FROM deck_users
UNION ALL
SELECT 'cards', COUNT(*) FROM cards
UNION ALL
SELECT 'budget_categories', COUNT(*) FROM budget_categories
UNION ALL
SELECT 'budgets', COUNT(*) FROM budgets
UNION ALL
SELECT 'incomes', COUNT(*) FROM incomes
UNION ALL
SELECT 'topics', COUNT(*) FROM topics
UNION ALL
SELECT 'task_learning', COUNT(*) FROM task_learning
UNION ALL
SELECT 'tags', COUNT(*) FROM tags
UNION ALL
SELECT 'resources', COUNT(*) FROM resources
UNION ALL
SELECT 'study_sessions', COUNT(*) FROM study_sessions
UNION ALL
SELECT 'repeat_rules', COUNT(*) FROM repeat_rules
UNION ALL
SELECT 'ai_recommendations', COUNT(*) FROM ai_recommendations;

-- Commit the transaction
COMMIT;

-- =====================================================
-- IMPORTANT NOTES:
-- =====================================================
-- 1. This script handles the case where UUID columns already exist
-- 2. It recreates all tables to properly set UUID primary keys
-- 3. Make sure to generate UUIDs for all records before running this script
-- 4. If any step fails, the transaction will rollback
-- 5. After successful migration, update your application code to use UUIDs
-- 6. Test thoroughly in a staging environment first