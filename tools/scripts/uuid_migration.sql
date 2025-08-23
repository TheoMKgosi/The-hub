-- UUID Migration Script for SQLite Database
-- This script migrates from uint primary keys to UUID primary keys
-- Run this script on your production database after creating a backup

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Begin transaction for atomicity
BEGIN TRANSACTION;

-- =====================================================
-- Step 1: Add UUID columns to all tables
-- =====================================================

-- Users table
ALTER TABLE users ADD COLUMN id_uuid TEXT;

-- Goals table
ALTER TABLE goals ADD COLUMN id_uuid TEXT;
ALTER TABLE goals ADD COLUMN user_id_uuid TEXT;

-- Tasks table
ALTER TABLE tasks ADD COLUMN id_uuid TEXT;
ALTER TABLE tasks ADD COLUMN goal_id_uuid TEXT;
ALTER TABLE tasks ADD COLUMN user_id_uuid TEXT;

-- Scheduled Tasks table
ALTER TABLE scheduled_tasks ADD COLUMN id_uuid TEXT;
ALTER TABLE scheduled_tasks ADD COLUMN user_id_uuid TEXT;

-- Decks table
ALTER TABLE decks ADD COLUMN id_uuid TEXT;
ALTER TABLE decks ADD COLUMN user_id_uuid TEXT;

-- Deck Users table
ALTER TABLE deck_users ADD COLUMN id_uuid TEXT;
ALTER TABLE deck_users ADD COLUMN deck_id_uuid TEXT;
ALTER TABLE deck_users ADD COLUMN user_id_uuid TEXT;

-- Cards table
ALTER TABLE cards ADD COLUMN id_uuid TEXT;
ALTER TABLE cards ADD COLUMN deck_id_uuid TEXT;

-- Budget Categories table
ALTER TABLE budget_categories ADD COLUMN id_uuid TEXT;
ALTER TABLE budget_categories ADD COLUMN user_id_uuid TEXT;

-- Budgets table
ALTER TABLE budgets ADD COLUMN id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN category_id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN user_id_uuid TEXT;
ALTER TABLE budgets ADD COLUMN income_id_uuid TEXT;

-- Income table
ALTER TABLE incomes ADD COLUMN id_uuid TEXT;
ALTER TABLE incomes ADD COLUMN user_id_uuid TEXT;

-- Topics table
ALTER TABLE topics ADD COLUMN id_uuid TEXT;
ALTER TABLE topics ADD COLUMN user_id_uuid TEXT;

-- Task Learning table
ALTER TABLE task_learning ADD COLUMN id_uuid TEXT;
ALTER TABLE task_learning ADD COLUMN topic_id_uuid TEXT;

-- Tags table
ALTER TABLE tags ADD COLUMN id_uuid TEXT;
ALTER TABLE tags ADD COLUMN user_id_uuid TEXT;

-- Resources table
ALTER TABLE resources ADD COLUMN id_uuid TEXT;
ALTER TABLE resources ADD COLUMN topic_id_uuid TEXT;
ALTER TABLE resources ADD COLUMN task_id_uuid TEXT;

-- Study Sessions table
ALTER TABLE study_sessions ADD COLUMN id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN user_id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN topic_id_uuid TEXT;
ALTER TABLE study_sessions ADD COLUMN task_id_uuid TEXT;

-- Repeat Rules table
ALTER TABLE repeat_rules ADD COLUMN id_uuid TEXT;

-- AI Recommendations table
ALTER TABLE ai_recommendations ADD COLUMN id_uuid TEXT;
ALTER TABLE ai_recommendations ADD COLUMN task_id_uuid TEXT;

-- =====================================================
-- Step 2: Generate UUIDs for existing records
-- =====================================================
-- Note: You need to generate actual UUIDs here
-- You can use a script to populate these values

-- Example for users table (replace with actual UUIDs):
-- UPDATE users SET id_uuid = '550e8400-e29b-41d4-a716-446655440000' WHERE id = 1;
-- UPDATE users SET id_uuid = '550e8400-e29b-41d4-a716-446655440001' WHERE id = 2;

-- =====================================================
-- Step 3: Update foreign key references
-- =====================================================

-- Update goals.user_id references
UPDATE goals SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = goals.user_id
);

-- Update tasks references
UPDATE tasks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = tasks.user_id
);
UPDATE tasks SET goal_id_uuid = (
    SELECT id_uuid FROM goals WHERE goals.id = tasks.goal_id
) WHERE goal_id IS NOT NULL;

-- Update scheduled_tasks references
UPDATE scheduled_tasks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = scheduled_tasks.user_id
);

-- Update decks references
UPDATE decks SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = decks.user_id
);

-- Update deck_users references
UPDATE deck_users SET deck_id_uuid = (
    SELECT id_uuid FROM decks WHERE decks.id = deck_users.deck_id
);
UPDATE deck_users SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = deck_users.user_id
);

-- Update cards references
UPDATE cards SET deck_id_uuid = (
    SELECT id_uuid FROM decks WHERE decks.id = cards.deck_id
);

-- Update budget_categories references
UPDATE budget_categories SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = budget_categories.user_id
);

-- Update budgets references
UPDATE budgets SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = budgets.user_id
);
UPDATE budgets SET category_id_uuid = (
    SELECT id_uuid FROM budget_categories WHERE budget_categories.id = budgets.category_id
);
UPDATE budgets SET income_id_uuid = (
    SELECT id_uuid FROM incomes WHERE incomes.id = budgets.income_id
) WHERE income_id IS NOT NULL;

-- Update incomes references
UPDATE incomes SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = incomes.user_id
);

-- Update topics references
UPDATE topics SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = topics.user_id
);

-- Update task_learning references
UPDATE task_learning SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = task_learning.topic_id
);

-- Update tags references
UPDATE tags SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = tags.user_id
);

-- Update resources references
UPDATE resources SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = resources.topic_id
) WHERE topic_id IS NOT NULL;
UPDATE resources SET task_id_uuid = (
    SELECT id_uuid FROM task_learning WHERE task_learning.id = resources.task_id
) WHERE task_id IS NOT NULL;

-- Update study_sessions references
UPDATE study_sessions SET user_id_uuid = (
    SELECT id_uuid FROM users WHERE users.id = study_sessions.user_id
) WHERE user_id IS NOT NULL;
UPDATE study_sessions SET topic_id_uuid = (
    SELECT id_uuid FROM topics WHERE topics.id = study_sessions.topic_id
) WHERE topic_id IS NOT NULL;
UPDATE study_sessions SET task_id_uuid = (
    SELECT id_uuid FROM task_learning WHERE task_learning.id = study_sessions.task_id
) WHERE task_id IS NOT NULL;

-- Update ai_recommendations references
UPDATE ai_recommendations SET task_id_uuid = (
    SELECT id_uuid FROM tasks WHERE tasks.id = ai_recommendations.task_id
);

-- =====================================================
-- Step 4: Drop old columns and rename new ones
-- =====================================================

-- Users table
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_uuid TO id;

-- Goals table
ALTER TABLE goals DROP COLUMN id;
ALTER TABLE goals DROP COLUMN user_id;
ALTER TABLE goals RENAME COLUMN id_uuid TO id;
ALTER TABLE goals RENAME COLUMN user_id_uuid TO user_id;

-- Tasks table
ALTER TABLE tasks DROP COLUMN id;
ALTER TABLE tasks DROP COLUMN goal_id;
ALTER TABLE tasks DROP COLUMN user_id;
ALTER TABLE tasks RENAME COLUMN id_uuid TO id;
ALTER TABLE tasks RENAME COLUMN goal_id_uuid TO goal_id;
ALTER TABLE tasks RENAME COLUMN user_id_uuid TO user_id;

-- Scheduled Tasks table
ALTER TABLE scheduled_tasks DROP COLUMN id;
ALTER TABLE scheduled_tasks DROP COLUMN user_id;
ALTER TABLE scheduled_tasks RENAME COLUMN id_uuid TO id;
ALTER TABLE scheduled_tasks RENAME COLUMN user_id_uuid TO user_id;

-- Decks table
ALTER TABLE decks DROP COLUMN id;
ALTER TABLE decks DROP COLUMN user_id;
ALTER TABLE decks RENAME COLUMN id_uuid TO id;
ALTER TABLE decks RENAME COLUMN user_id_uuid TO user_id;

-- Deck Users table
ALTER TABLE deck_users DROP COLUMN id;
ALTER TABLE deck_users DROP COLUMN deck_id;
ALTER TABLE deck_users DROP COLUMN user_id;
ALTER TABLE deck_users RENAME COLUMN id_uuid TO id;
ALTER TABLE deck_users RENAME COLUMN deck_id_uuid TO deck_id;
ALTER TABLE deck_users RENAME COLUMN user_id_uuid TO user_id;

-- Cards table
ALTER TABLE cards DROP COLUMN id;
ALTER TABLE cards DROP COLUMN deck_id;
ALTER TABLE cards RENAME COLUMN id_uuid TO id;
ALTER TABLE cards RENAME COLUMN deck_id_uuid TO deck_id;

-- Budget Categories table
ALTER TABLE budget_categories DROP COLUMN id;
ALTER TABLE budget_categories DROP COLUMN user_id;
ALTER TABLE budget_categories RENAME COLUMN id_uuid TO id;
ALTER TABLE budget_categories RENAME COLUMN user_id_uuid TO user_id;

-- Budgets table
ALTER TABLE budgets DROP COLUMN id;
ALTER TABLE budgets DROP COLUMN category_id;
ALTER TABLE budgets DROP COLUMN user_id;
ALTER TABLE budgets DROP COLUMN income_id;
ALTER TABLE budgets RENAME COLUMN id_uuid TO id;
ALTER TABLE budgets RENAME COLUMN category_id_uuid TO category_id;
ALTER TABLE budgets RENAME COLUMN user_id_uuid TO user_id;
ALTER TABLE budgets RENAME COLUMN income_id_uuid TO income_id;

-- Income table
ALTER TABLE incomes DROP COLUMN id;
ALTER TABLE incomes DROP COLUMN user_id;
ALTER TABLE incomes RENAME COLUMN id_uuid TO id;
ALTER TABLE incomes RENAME COLUMN user_id_uuid TO user_id;

-- Topics table
ALTER TABLE topics DROP COLUMN id;
ALTER TABLE topics DROP COLUMN user_id;
ALTER TABLE topics RENAME COLUMN id_uuid TO id;
ALTER TABLE topics RENAME COLUMN user_id_uuid TO user_id;

-- Task Learning table
ALTER TABLE task_learning DROP COLUMN id;
ALTER TABLE task_learning DROP COLUMN topic_id;
ALTER TABLE task_learning RENAME COLUMN id_uuid TO id;
ALTER TABLE task_learning RENAME COLUMN topic_id_uuid TO topic_id;

-- Tags table
ALTER TABLE tags DROP COLUMN id;
ALTER TABLE tags DROP COLUMN user_id;
ALTER TABLE tags RENAME COLUMN id_uuid TO id;
ALTER TABLE tags RENAME COLUMN user_id_uuid TO user_id;

-- Resources table
ALTER TABLE resources DROP COLUMN id;
ALTER TABLE resources DROP COLUMN topic_id;
ALTER TABLE resources DROP COLUMN task_id;
ALTER TABLE resources RENAME COLUMN id_uuid TO id;
ALTER TABLE resources RENAME COLUMN topic_id_uuid TO topic_id;
ALTER TABLE resources RENAME COLUMN task_id_uuid TO task_id;

-- Study Sessions table
ALTER TABLE study_sessions DROP COLUMN id;
ALTER TABLE study_sessions DROP COLUMN user_id;
ALTER TABLE study_sessions DROP COLUMN topic_id;
ALTER TABLE study_sessions DROP COLUMN task_id;
ALTER TABLE study_sessions RENAME COLUMN id_uuid TO id;
ALTER TABLE study_sessions RENAME COLUMN user_id_uuid TO user_id;
ALTER TABLE study_sessions RENAME COLUMN topic_id_uuid TO topic_id;
ALTER TABLE study_sessions RENAME COLUMN task_id_uuid TO task_id;

-- Repeat Rules table
ALTER TABLE repeat_rules DROP COLUMN id;
ALTER TABLE repeat_rules RENAME COLUMN id_uuid TO id;

-- AI Recommendations table
ALTER TABLE ai_recommendations DROP COLUMN id;
ALTER TABLE ai_recommendations DROP COLUMN task_id;
ALTER TABLE ai_recommendations RENAME COLUMN id_uuid TO id;
ALTER TABLE ai_recommendations RENAME COLUMN task_id_uuid TO task_id;

-- =====================================================
-- Step 5: Verify the migration
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
-- 1. This script assumes you have already populated the id_uuid columns with actual UUIDs
-- 2. Make sure to create a backup before running this script
-- 3. If any step fails, the transaction will rollback
-- 4. After successful migration, update your application code to use UUIDs
-- 5. Test thoroughly in a staging environment first