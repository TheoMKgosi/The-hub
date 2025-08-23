# UUID Migration Guide

This guide explains how to migrate your SQLite database from `uint` primary keys to `UUID` primary keys while preserving all existing data.

## Overview

Your application has been updated to use UUID primary keys instead of integer primary keys. This provides better security, scalability, and compatibility across different database systems. However, your existing database still uses `uint` primary keys, so you need to migrate the data.

## Migration Scripts

Four migration approaches are provided:

### 1. Easy Shell Script (`migrate_to_uuid.sh`) - RECOMMENDED
- **Pros**: Fully automated, handles all scenarios including recovery
- **Cons**: Requires bash and Python 3
- **Best for**: All environments, especially production

### 2. Python Script (`uuid_migration.py`)
- **Pros**: Automated, handles UUID generation and foreign key updates
- **Cons**: Requires Python and sqlite3 dependencies
- **Best for**: Development/staging environments

### 3. SQL Script (`uuid_migration.sql`)
- **Pros**: Direct SQL approach, more reliable for production
- **Cons**: Requires manual UUID generation
- **Best for**: Production environments with custom requirements

### 4. Recovery Script (`recovery_migration.py`)
- **Pros**: Handles failed migration attempts and inconsistent states
- **Cons**: Requires Python 3
- **Best for**: Recovery from failed migrations

## Prerequisites

1. **Backup your database**:
   ```bash
   cp yourdb.db yourdb.db.backup
   ```

2. **Install dependencies** (for Python script):
   ```bash
   pip install sqlite3 uuid
   ```

3. **Verify your database structure**:
   ```sql
   .tables
   PRAGMA table_info(users);
   ```

## Migration Steps

### Step 1: Choose Your Migration Method

#### Option A: Easy Shell Script (Recommended for all cases)
```bash
cd tools/scripts
./migrate_to_uuid.sh yourdb.db python
```

#### Option B: Recovery Migration (If previous attempts failed)
```bash
cd tools/scripts
./migrate_to_uuid.sh yourdb.db recovery
```

#### Option C: Python Script (Direct method)
```bash
python uuid_migration.py yourdb.db
```

#### Option D: SQL Script (Manual method)
1. **Generate UUIDs for existing data**:
   ```python
   import sqlite3
   import uuid

   conn = sqlite3.connect('yourdb.db')
   cursor = conn.cursor()

   # Generate UUIDs for each table
   tables = ['users', 'tasks', 'goals', ...]  # All your tables
   for table in tables:
       cursor.execute(f"SELECT id FROM {table}")
       rows = cursor.fetchall()
       for (old_id,) in rows:
           new_uuid = str(uuid.uuid4())
           cursor.execute(f"UPDATE {table} SET id_uuid = ? WHERE id = ?", (new_uuid, old_id))

   conn.commit()
   conn.close()
   ```

2. **Run the SQL migration**:
   ```bash
   sqlite3 yourdb.db < uuid_migration.sql
   ```

### Step 2: Verify Migration Success

1. **Check record counts**:
   ```sql
   SELECT 'users' as table, COUNT(*) FROM users
   UNION ALL
   SELECT 'tasks', COUNT(*) FROM tasks
   UNION ALL
   SELECT 'goals', COUNT(*) FROM goals;
   ```

2. **Verify foreign key constraints**:
   ```sql
   PRAGMA foreign_key_check;
   ```

3. **Test a few queries**:
   ```sql
   SELECT * FROM users LIMIT 5;
   SELECT * FROM tasks WHERE user_id IN (SELECT id FROM users LIMIT 3);
   ```

### Step 3: Update Application Configuration

After successful migration:

1. **Remove the old migration** from your migration manager
2. **Update any hardcoded queries** that might reference old column names
3. **Test your application** thoroughly

## Troubleshooting

### Common Issues

1. **Foreign Key Constraint Errors**:
   - Ensure all foreign key references are updated before dropping columns
   - Check `PRAGMA foreign_key_check;` for violations

2. **UUID Generation Issues**:
   - Ensure all `id_uuid` columns are populated before running the migration
   - Verify UUID format (should be like `550e8400-e29b-41d4-a716-446655440000`)

3. **Data Loss**:
   - Always work with a backup first
   - If migration fails, restore from backup and try again

4. **Failed Migration Recovery**:
   - If you see errors like "duplicate column name: id_uuid" or "cannot drop PRIMARY KEY column"
   - Use the recovery migration: `./migrate_to_uuid.sh yourdb.db recovery`
   - The recovery script can handle partially migrated databases

5. **SQLite Limitations**:
   - SQLite doesn't support direct primary key modifications
   - The migration recreates tables to work around this limitation
   - This is normal and expected behavior

### Rollback Procedure

If you need to rollback:

1. **Restore from backup**:
   ```bash
   cp yourdb.db.backup yourdb.db
   ```

2. **Re-run your application migrations** to get back to the original state

## Post-Migration Tasks

1. **Update any data seeding scripts** to use UUIDs
2. **Review and update any direct SQL queries** in your application
3. **Test all CRUD operations** for each entity
4. **Monitor performance** - UUIDs are larger than integers but provide better distribution

## Security Considerations

- UUIDs provide better security than sequential integers
- Consider implementing UUID v4 (random) for production use
- Review any external APIs that might expose IDs

## Support

If you encounter issues during migration:

1. Check the backup integrity
2. Verify all steps were completed in order
3. Review the error messages carefully
4. Consider testing with a smaller dataset first

## Recovery Migration

If you've already attempted to migrate your database and it failed, you may see errors like:
- "duplicate column name: id_uuid"
- "cannot drop PRIMARY KEY column"
- "unknown column" errors
- Foreign key constraint violations

In these cases, use the **recovery migration**:

```bash
cd tools/scripts
./migrate_to_uuid.sh yourdb.db recovery
```

The recovery migration:
1. Analyzes the current state of your database
2. Detects which tables have been partially migrated
3. Handles existing UUID columns appropriately
4. Recreates tables with proper UUID primary keys
5. Preserves all your data

## Files in this Directory

- `migrate_to_uuid.sh` - Easy-to-use shell script (recommended)
- `uuid_migration.py` - Python script for automated migration
- `uuid_migration.sql` - SQL script for manual migration
- `recovery_migration.py` - Recovery script for failed migrations
- `recovery_migration.sql` - SQL recovery script
- `README_UUID_MIGRATION.md` - This documentation

---

**Important**: Always test this migration process in a development environment before applying to production data.