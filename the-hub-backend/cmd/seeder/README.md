# Database Seeder

This tool provides database seeding functionality for The Hub application, allowing you to populate the database with sample data for development and testing purposes.

## Features

- **Sample Users**: Create admin and regular user accounts with predefined credentials
- **Sample Goals**: Generate realistic goals with different priorities and categories
- **Sample Tasks**: Create tasks associated with goals in various states (pending, in progress, completed)
- **Safety Checks**: Prevent accidental data duplication with smart detection of seeded data
- **Clean Operations**: Remove only seeded data without affecting real user data

## Usage

### Command Line Interface

```bash
go run cmd/seeder/main.go -command <command> [options]
```

### Available Commands

- `users` - Seed sample users only
- `goals` - Seed sample goals and tasks only
- `all` - Seed all data (users, goals, and tasks)
- `clean` - Remove all seeded data

### Options

- `-force` - Skip confirmation prompts and override existing seeded data
- `-help` - Show help message

### Makefile Targets

```bash
# Seed all data
make seed

# Seed users only
make seed-users

# Seed goals and tasks only
make seed-goals

# Clean seeded data
make seed-clean
```

## Sample Data Created

### Users
- **Admin User** (`admin@thehub.com`) - Administrator account
- **John Doe** (`john.doe@example.com`) - Regular user
- **Jane Smith** (`jane.smith@example.com`) - Regular user
- **Bob Johnson** (`bob.johnson@example.com`) - Regular user

All users have the password: `password123`

### Goals and Tasks

1. **Complete Project Documentation** (Work, High Priority)
   - Write API endpoint documentation
   - Create user guide
   - Add code examples
   - Review and publish docs

2. **Learn Go Advanced Patterns** (Learning, Medium Priority)
   - Study concurrency patterns
   - Learn about generics
   - Practice with interfaces
   - Build a sample project

3. **Home Organization Project** (Personal, Low Priority)
   - Clean out desk drawers
   - Organize digital files
   - Sort through old documents
   - Set up new filing system

## Safety Features

- **Smart Detection**: Identifies seeded data by email patterns (`@thehub.com`, `.example.com`) and goal titles
- **Duplicate Prevention**: Won't create duplicate data unless `-force` flag is used
- **Selective Cleaning**: Only removes seeded data, preserving real user data
- **Confirmation Prompts**: Asks for confirmation before destructive operations (unless `-force` used)

## Examples

### Basic Usage

```bash
# Seed everything
make seed

# Seed users only
make seed-users

# Seed goals for existing users
make seed-goals
```

### Advanced Usage

```bash
# Force re-seed all data (overwrites existing seeded data)
go run cmd/seeder/main.go -command all -force

# Clean seeded data without confirmation
make seed-clean
```

### Development Workflow

```bash
# Set up fresh development database
make migrate-up
make seed

# Reset seeded data
make seed-clean
make seed
```

## Notes

- The seeder requires database connection environment variables to be set
- Seeded users have realistic but obviously fake email addresses for easy identification
- Tasks are created with varying priorities and statuses to demonstrate different states
- The seeder integrates with the existing database migration system