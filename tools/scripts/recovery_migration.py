#!/usr/bin/env python3
"""
Recovery Migration Script for SQLite Database

This script handles recovery from failed UUID migration attempts.
It can detect the current state of the database and apply the appropriate
recovery strategy.

Usage:
    python recovery_migration.py <database_path>

Requirements:
    pip install sqlite3 uuid

Author: AI Assistant
Date: 2025
"""

import sqlite3
import uuid
import sys
import os
from typing import List, Dict, Any, Set

class RecoveryMigrator:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self.backup_path = f"{db_path}.recovery.backup"

    def backup_database(self):
        """Create a backup of the database"""
        print(f"Creating recovery backup: {self.backup_path}")
        import shutil
        shutil.copy2(self.db_path, self.backup_path)
        print("Recovery backup created successfully")

    def analyze_database_state(self) -> Dict[str, Any]:
        """Analyze the current state of the database"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        state = {
            'tables': {},
            'has_uuid_columns': False,
            'migration_status': 'unknown'
        }

        # Get all tables
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
        tables = [row[0] for row in cursor.fetchall()]

        for table in tables:
            table_info = {'exists': True, 'columns': [], 'has_uuid_column': False}

            # Get column information
            try:
                cursor.execute(f"PRAGMA table_info({table})")
                columns = cursor.fetchall()
                table_info['columns'] = [{'name': col[1], 'type': col[2], 'pk': col[5]} for col in columns]

                # Check if UUID column exists
                for col in table_info['columns']:
                    if col['name'] == 'id_uuid':
                        table_info['has_uuid_column'] = True
                        state['has_uuid_columns'] = True
                        break

            except sqlite3.Error as e:
                print(f"Warning: Could not analyze table {table}: {e}")
                table_info['exists'] = False

            state['tables'][table] = table_info

        conn.close()

        # Determine migration status
        if state['has_uuid_columns']:
            # Check if any tables have been fully migrated (no integer primary key)
            fully_migrated = 0
            partially_migrated = 0

            for table, info in state['tables'].items():
                if info['exists'] and info['columns']:
                    has_int_pk = any(col['pk'] and col['type'].lower() in ['integer', 'int'] for col in info['columns'])
                    has_uuid_col = info['has_uuid_column']

                    if has_uuid_col and not has_int_pk:
                        fully_migrated += 1
                    elif has_uuid_col and has_int_pk:
                        partially_migrated += 1

            if fully_migrated > 0 and partially_migrated == 0:
                state['migration_status'] = 'completed'
            elif partially_migrated > 0:
                state['migration_status'] = 'partial'
            else:
                state['migration_status'] = 'uuid_columns_exist'
        else:
            state['migration_status'] = 'not_started'

        return state

    def generate_uuids_for_table(self, table: str) -> Dict[int, str]:
        """Generate UUIDs for all records in a table"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # Get primary key column name
        cursor.execute(f"PRAGMA table_info({table})")
        columns = cursor.fetchall()
        pk_column = None
        for col in columns:
            if col[5] == 1:  # is_pk
                pk_column = col[1]
                break

        if not pk_column:
            raise ValueError(f"No primary key found for table {table}")

        # Get all existing IDs
        cursor.execute(f"SELECT {pk_column} FROM {table}")
        existing_ids = cursor.fetchall()

        # Generate UUIDs for each existing record
        uuid_mapping = {}
        for (existing_id,) in existing_ids:
            uuid_mapping[existing_id] = str(uuid.uuid4())

        conn.close()
        return uuid_mapping

    def populate_uuid_columns(self):
        """Populate UUID columns with generated UUIDs"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # Get all tables that need UUID generation
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
        tables = [row[0] for row in cursor.fetchall()]

        for table in tables:
            try:
                # Check if table has id_uuid column
                cursor.execute(f"PRAGMA table_info({table})")
                columns = cursor.fetchall()
                has_uuid_column = any(col[1] == 'id_uuid' for col in columns)

                if has_uuid_column:
                    # Check if UUIDs are already populated
                    cursor.execute(f"SELECT COUNT(*) FROM {table} WHERE id_uuid IS NULL OR id_uuid = ''")
                    null_count = cursor.fetchone()[0]

                    if null_count > 0:
                        print(f"Generating UUIDs for {null_count} records in {table}")

                        # Generate UUIDs for records that don't have them
                        cursor.execute(f"SELECT id FROM {table} WHERE id_uuid IS NULL OR id_uuid = ''")
                        records_without_uuid = cursor.fetchall()

                        for (record_id,) in records_without_uuid:
                            new_uuid = str(uuid.uuid4())
                            cursor.execute(f"UPDATE {table} SET id_uuid = ? WHERE id = ?", (new_uuid, record_id))

                        print(f"Generated UUIDs for {len(records_without_uuid)} records in {table}")
                    else:
                        print(f"Table {table} already has UUIDs populated")
                else:
                    print(f"Table {table} does not have id_uuid column")

            except sqlite3.Error as e:
                print(f"Error processing table {table}: {e}")
                conn.rollback()
                raise

        conn.commit()
        conn.close()

    def recreate_tables_with_uuid_pk(self):
        """Recreate tables with UUID primary keys"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # Define table schemas with UUID primary keys
        table_schemas = {
            'users': '''
                CREATE TABLE users_new (
                    id TEXT PRIMARY KEY,
                    name TEXT,
                    email TEXT UNIQUE,
                    password TEXT,
                    settings TEXT,
                    created_at DATETIME,
                    updated_at DATETIME,
                    deleted_at DATETIME
                )
            ''',
            'goals': '''
                CREATE TABLE goals_new (
                    id TEXT PRIMARY KEY,
                    user_id TEXT,
                    title TEXT,
                    description TEXT,
                    created_at DATETIME,
                    updated_at DATETIME,
                    deleted_at DATETIME,
                    FOREIGN KEY (user_id) REFERENCES users_new(id)
                )
            ''',
            'tasks': '''
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
                )
            ''',
            'scheduled_tasks': '''
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
                )
            ''',
            'decks': '''
                CREATE TABLE decks_new (
                    id TEXT PRIMARY KEY,
                    name TEXT NOT NULL,
                    user_id TEXT NOT NULL,
                    created_at DATETIME,
                    updated_at DATETIME,
                    deleted_at DATETIME,
                    FOREIGN KEY (user_id) REFERENCES users_new(id)
                )
            ''',
            'deck_users': '''
                CREATE TABLE deck_users_new (
                    id TEXT PRIMARY KEY,
                    deck_id TEXT,
                    user_id TEXT,
                    role TEXT,
                    FOREIGN KEY (deck_id) REFERENCES decks_new(id),
                    FOREIGN KEY (user_id) REFERENCES users_new(id)
                )
            ''',
            'cards': '''
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
                )
            ''',
            'budget_categories': '''
                CREATE TABLE budget_categories_new (
                    id TEXT PRIMARY KEY,
                    name TEXT NOT NULL,
                    user_id TEXT,
                    created_at DATETIME,
                    updated_at DATETIME,
                    deleted_at DATETIME,
                    FOREIGN KEY (user_id) REFERENCES users_new(id)
                )
            ''',
            'budgets': '''
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
                )
            ''',
            'incomes': '''
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
                )
            ''',
            'topics': '''
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
                )
            ''',
            'task_learning': '''
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
                )
            ''',
            'tags': '''
                CREATE TABLE tags_new (
                    id TEXT PRIMARY KEY,
                    user_id TEXT,
                    name TEXT UNIQUE NOT NULL,
                    color TEXT,
                    FOREIGN KEY (user_id) REFERENCES users_new(id)
                )
            ''',
            'resources': '''
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
                )
            ''',
            'study_sessions': '''
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
                )
            ''',
            'repeat_rules': '''
                CREATE TABLE repeat_rules_new (
                    id TEXT PRIMARY KEY,
                    frequency TEXT,
                    interval INTEGER,
                    by_day TEXT,
                    start_date DATE,
                    end_date DATE,
                    created_at DATETIME,
                    updated_at DATETIME
                )
            ''',
            'ai_recommendations': '''
                CREATE TABLE ai_recommendations_new (
                    id TEXT PRIMARY KEY,
                    task_id TEXT,
                    suggested_start DATETIME,
                    suggested_end DATETIME,
                    confidence REAL,
                    accepted BOOLEAN,
                    created_at DATETIME,
                    FOREIGN KEY (task_id) REFERENCES tasks_new(id)
                )
            '''
        }

        # Create new tables
        for table, schema in table_schemas.items():
            print(f"Creating new table: {table}")
            cursor.execute(schema)

        # Copy data from old tables to new tables
        data_copy_queries = {
            'users': "INSERT INTO users_new SELECT id_uuid, name, email, password, settings, created_at, updated_at, deleted_at FROM users WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'goals': "INSERT INTO goals_new SELECT id_uuid, user_id_uuid, title, description, created_at, updated_at, deleted_at FROM goals WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'tasks': "INSERT INTO tasks_new SELECT id_uuid, title, description, due_date, priority, status, order_index, goal_id_uuid, user_id_uuid, created_at, updated_at, deleted_at FROM tasks WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'scheduled_tasks': "INSERT INTO scheduled_tasks_new SELECT id_uuid, title, start, end, user_id_uuid, created_by_ai, created_at, updated_at FROM scheduled_tasks WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'decks': "INSERT INTO decks_new SELECT id_uuid, name, user_id_uuid, created_at, updated_at, deleted_at FROM decks WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'deck_users': "INSERT INTO deck_users_new SELECT id_uuid, deck_id_uuid, user_id_uuid, role FROM deck_users WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'cards': "INSERT INTO cards_new SELECT id_uuid, deck_id_uuid, question, answer, easiness, interval, repetitions, last_reviewed, next_review, created_at, updated_at, deleted_at FROM cards WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'budget_categories': "INSERT INTO budget_categories_new SELECT id_uuid, name, user_id_uuid, created_at, updated_at, deleted_at FROM budget_categories WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'budgets': "INSERT INTO budgets_new SELECT id_uuid, category_id_uuid, amount, start_date, end_date, user_id_uuid, income_id_uuid, created_at, updated_at, deleted_at FROM budgets WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'incomes': "INSERT INTO incomes_new SELECT id_uuid, source, amount, user_id_uuid, received_at, created_at, updated_at, deleted_at FROM incomes WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'topics': "INSERT INTO topics_new SELECT id_uuid, user_id_uuid, title, description, status, deadline, created_at, updated_at, deleted_at FROM topics WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'task_learning': "INSERT INTO task_learning_new SELECT id_uuid, topic_id_uuid, title, notes, status, order_index, created_at, updated_at, deleted_at FROM task_learning WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'tags': "INSERT INTO tags_new SELECT id_uuid, user_id_uuid, name, color FROM tags WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'resources': "INSERT INTO resources_new SELECT id_uuid, topic_id_uuid, task_id_uuid, title, link, type, notes FROM resources WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'study_sessions': "INSERT INTO study_sessions_new SELECT id_uuid, user_id_uuid, topic_id_uuid, task_id_uuid, duration_min, started_at, ended_at FROM study_sessions WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'repeat_rules': "INSERT INTO repeat_rules_new SELECT id_uuid, frequency, interval, by_day, start_date, end_date, created_at, updated_at FROM repeat_rules WHERE id_uuid IS NOT NULL AND id_uuid != ''",
            'ai_recommendations': "INSERT INTO ai_recommendations_new SELECT id_uuid, task_id_uuid, suggested_start, suggested_end, confidence, accepted, created_at FROM ai_recommendations WHERE id_uuid IS NOT NULL AND id_uuid != ''"
        }

        # Copy data
        for table, query in data_copy_queries.items():
            try:
                cursor.execute(query)
                print(f"Copied data to {table}_new")
            except sqlite3.Error as e:
                print(f"Error copying data to {table}_new: {e}")
                # Continue with other tables

        # Drop old tables
        old_tables = list(table_schemas.keys())
        for table in old_tables:
            try:
                cursor.execute(f"DROP TABLE {table}")
                print(f"Dropped old table: {table}")
            except sqlite3.Error as e:
                print(f"Error dropping table {table}: {e}")

        # Rename new tables
        for table in old_tables:
            try:
                cursor.execute(f"ALTER TABLE {table}_new RENAME TO {table}")
                print(f"Renamed {table}_new to {table}")
            except sqlite3.Error as e:
                print(f"Error renaming {table}_new: {e}")

        conn.commit()
        conn.close()

    def migrate(self):
        """Main migration method"""
        print(f"Starting recovery migration for database: {self.db_path}")

        # Create backup
        self.backup_database()

        # Analyze current state
        print("Analyzing database state...")
        state = self.analyze_database_state()

        print(f"Migration status: {state['migration_status']}")
        print(f"Has UUID columns: {state['has_uuid_columns']}")

        if state['migration_status'] == 'completed':
            print("Migration appears to be already completed. No action needed.")
            return

        # Populate UUID columns if they exist but are empty
        if state['has_uuid_columns']:
            print("Populating UUID columns...")
            self.populate_uuid_columns()

        # Recreate tables with UUID primary keys
        print("Recreating tables with UUID primary keys...")
        self.recreate_tables_with_uuid_pk()

        print("Recovery migration completed successfully!")
        print(f"Backup saved at: {self.backup_path}")

def main():
    if len(sys.argv) != 2:
        print("Usage: python recovery_migration.py <database_path>")
        sys.exit(1)

    db_path = sys.argv[1]

    if not os.path.exists(db_path):
        print(f"Database file not found: {db_path}")
        sys.exit(1)

    migrator = RecoveryMigrator(db_path)
    migrator.migrate()

if __name__ == "__main__":
    main()