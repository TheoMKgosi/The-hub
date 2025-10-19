# Admin CLI Tool

A command-line tool for creating admin users in The Hub application.

## Features

- Create admin users with secure password hashing
- Interactive mode with prompts
- Command-line argument support
- Duplicate email detection
- Comprehensive help documentation

## Usage

### Interactive Mode (Recommended)

Run the tool without arguments to be prompted for user information:

```bash
make create-admin
```

Or directly:

```bash
go run cmd/admin-cli/main.go
```

### Command Line Mode

Provide all required information via command-line arguments:

```bash
make create-admin-args EMAIL=admin@example.com NAME="Admin User" PASSWORD=securepassword
```

Or directly:

```bash
go run cmd/admin-cli/main.go -email admin@example.com -name "Admin User" -password securepassword
```

### Help

Show the help message:

```bash
go run cmd/admin-cli/main.go -help
```

## Environment Variables

The tool requires the following environment variables to be set for database connection:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name

## Command Line Options

- `-email string` - Email address for the admin user
- `-name string` - Name for the admin user  
- `-password string` - Password for the admin user
- `-help` - Show help message

## Security Features

- Passwords are hashed using bcrypt with default cost
- Duplicate email detection prevents account conflicts
- Admin role is automatically assigned
- Database connection uses secure PostgreSQL connection

## Examples

### Create an admin user interactively:

```bash
$ make create-admin
Admin User Creation Tool
========================

Enter email address: admin@example.com
Enter name: System Administrator
Enter password: securepassword

✅ Admin user created successfully!
   ID: 123e4567-e89b-12d3-a456-426614174000
   Email: admin@example.com
   Name: System Administrator
   Role: admin
```

### Create an admin user with arguments:

```bash
$ make create-admin-args EMAIL=admin@example.com NAME="Admin User" PASSWORD=mypassword123

✅ Admin user created successfully!
   ID: 123e4567-e89b-12d3-a456-426614174001
   Email: admin@example.com
   Name: Admin User
   Role: admin
```

## Error Handling

The tool will exit with error codes for:

- Missing required fields (email, name, password)
- Database connection failures
- Duplicate email addresses
- Invalid input formats

## Building

Build the CLI tool as a standalone binary:

```bash
make build-admin-cli
```

This creates an `admin-cli` executable in the current directory.