#!/usr/bin/env python3
"""
UUID Migration Script for SQLite Database

This script migrates a SQLite database from uint primary keys to UUID primary keys
while preserving all existing data and relationships.

Usage:
    python uuid_migration.py <database_path>

Requirements:
    pip install sqlite3 uuid

Author: AI Assistant
Date: 2025
"""

import sqlite3
import uuid
import sys
import os
from typing import List, Dict, Any

class UUIDMigrator:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self.backup_path = f"{db_path}.backup"

    def backup_database(self):
        """Create a backup of the database"""
        print(f"Creating backup: {self.backup_path}")
        import shutil
        shutil.copy2(self.db_path, self.backup_path)
        print("Backup created successfully")

    def get_tables_with_uint_keys(self) -> List[str]:
        """Get all tables that have uint primary keys"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        # Get all tables
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
        tables = [row[0] for row in cursor.fetchall()]

        tables_with_uint_keys = []

        for table in tables:
            # Check if table has an integer primary key
            cursor.execute(f"PRAGMA table_info({table})")
            columns = cursor.fetchall()

            for col in columns:
                col_name, col_type, not_null, default_value, is_pk = col[1], col[2], col[3], col[4], col[5]
                if is_pk and col_type.lower() in ['integer', 'int']:
                    tables_with_uint_keys.append(table)
                    break

        conn.close()
        return tables_with_uint_keys

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

    def migrate_table(self, table: str, uuid_mapping: Dict[int, str]):
        """Migrate a single table to use UUIDs"""
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

        print(f"Migrating table: {table} (PK: {pk_column})")

        # Step 1: Add new UUID column
        new_pk_column = f"{pk_column}_uuid"
        cursor.execute(f"ALTER TABLE {table} ADD COLUMN {new_pk_column} TEXT")

        # Step 2: Populate UUID column with generated UUIDs
        for old_id, new_uuid in uuid_mapping.items():
            cursor.execute(f"UPDATE {table} SET {new_pk_column} = ? WHERE {pk_column} = ?", (new_uuid, old_id))

        # Step 3: Update foreign key references in other tables
        if pk_column:
            self.update_foreign_key_references(cursor, table, pk_column, new_pk_column, uuid_mapping)

        # Step 4: Drop old primary key column
        cursor.execute(f"ALTER TABLE {table} DROP COLUMN {pk_column}")

        # Step 5: Rename new column to original name
        cursor.execute(f"ALTER TABLE {table} RENAME COLUMN {new_pk_column} TO {pk_column}")

        conn.commit()
        conn.close()

        print(f"Successfully migrated table: {table}")

    def update_foreign_key_references(self, cursor, table: str, pk_column: str, new_pk_column: str, uuid_mapping: Dict[int, str]):
        """Update foreign key references to the migrated table"""
        # Get all tables that reference this table
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
        all_tables = [row[0] for row in cursor.fetchall()]

        for other_table in all_tables:
            if other_table == table:
                continue

            # Check if this table has foreign key references to our table
            cursor.execute(f"PRAGMA foreign_key_list({other_table})")
            foreign_keys = cursor.fetchall()

            for fk in foreign_keys:
                if fk[2] == table and fk[3] == pk_column:  # table and from column
                    fk_column = fk[3]  # from column in the referencing table
                    print(f"  Updating foreign key reference: {other_table}.{fk_column}")

                    # Update the foreign key values
                    for old_id, new_uuid in uuid_mapping.items():
                        cursor.execute(f"UPDATE {other_table} SET {fk_column} = ? WHERE {fk_column} = ?", (new_uuid, old_id))

    def migrate(self):
        """Main migration method"""
        print(f"Starting UUID migration for database: {self.db_path}")

        # Create backup
        self.backup_database()

        # Get tables to migrate
        tables_to_migrate = self.get_tables_with_uint_keys()
        print(f"Found {len(tables_to_migrate)} tables to migrate: {tables_to_migrate}")

        if not tables_to_migrate:
            print("No tables with uint primary keys found. Migration complete.")
            return

        # Generate UUIDs for all tables first
        uuid_mappings = {}
        for table in tables_to_migrate:
            print(f"Generating UUIDs for table: {table}")
            uuid_mappings[table] = self.generate_uuids_for_table(table)

        # Migrate each table
        for table in tables_to_migrate:
            try:
                self.migrate_table(table, uuid_mappings[table])
            except Exception as e:
                print(f"Error migrating table {table}: {e}")
                print("Rolling back to backup...")
                import shutil
                shutil.copy2(self.backup_path, self.db_path)
                raise

        print("Migration completed successfully!")
        print(f"Backup saved at: {self.backup_path}")

def main():
    if len(sys.argv) != 2:
        print("Usage: python uuid_migration.py <database_path>")
        sys.exit(1)

    db_path = sys.argv[1]

    if not os.path.exists(db_path):
        print(f"Database file not found: {db_path}")
        sys.exit(1)

    migrator = UUIDMigrator(db_path)
    migrator.migrate()

if __name__ == "__main__":
    main()