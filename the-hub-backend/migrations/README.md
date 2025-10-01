# Database Migrations

This directory contains database migration files for The Hub application using golang-migrate.

## Migration Files

Migrations follow the naming convention: `{version}_{description}.{direction}.sql`

- `up.sql` - Contains the migration to apply the changes
- `down.sql` - Contains the rollback migration to undo the changes

## Running Migrations

### Using Make Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check current migration version
make migrate-version

# Show migration status
make migrate-status
```

### Using the Migration Tool Directly

```bash
# Build the migration tool
go build -o migrate ./cmd/migrate

# Run migrations
./migrate up

# Rollback one migration
./migrate down 1

# Check version
./migrate version
```

## Development Workflow

1. **Create Migration Files**: Add new `.up.sql` and `.down.sql` files with the next version number
2. **Test Locally**: Run migrations against your development database
3. **Test Rollback**: Ensure down migrations work correctly
4. **Commit**: Include migration files in your commit
5. **Deploy**: Migrations run automatically during deployment

## Migration Safety

- **Always test** migrations on a copy of production data first
- **Create backups** before running migrations in production
- **Verify** migration success before proceeding with deployment
- **Have rollback plans** ready for failed migrations

## Current Migrations

- `000001_initial_schema` - Creates core tables (users, goals, tasks, etc.)
- Future migrations will be added as schema changes are needed

## Troubleshooting

### Migration Fails
1. Check database connectivity
2. Verify migration SQL syntax
3. Check database permissions
4. Review application logs

### Dirty Database State
If a migration fails partway through, the database may be in a "dirty" state. You may need to:
1. Manually fix the database state
2. Use `migrate force {version}` to mark a migration as applied
3. Or restore from backup and retry

### Rollback Issues
- Ensure down migrations are tested
- Check for foreign key constraints
- Verify data dependencies before rollback