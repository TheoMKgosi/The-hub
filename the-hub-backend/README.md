# The Hub Backend

A robust Go-based REST API backend for The Hub productivity application, built with Gin framework and GORM ORM.

## 🚀 Features

- **RESTful API**: Clean, well-documented REST endpoints
- **JWT Authentication**: Secure token-based authentication system
- **PostgreSQL/SQLite**: Flexible database support for development and production
- **Comprehensive Models**: Task, Goal, Finance, Learning, and User management
- **AI-Powered Scheduling**: Intelligent task scheduling using OpenRouter AI
- **Natural Language Processing**: Parse tasks from natural language input
- **Swagger Documentation**: Interactive API documentation at `/swagger/index.html`
- **Middleware Support**: Authentication, CORS, logging, and error handling
- **Testing Suite**: Unit and integration tests with coverage reporting
- **Docker Support**: Containerized deployment ready

## 🏗 Architecture

### Technology Stack
- **Language**: Go 1.19+
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Go ORM)
- **Database**: PostgreSQL (production) / SQLite (development)
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go's built-in testing + stretchr/testify
- **Load Testing**: k6

### Project Structure
```
the-hub-backend/
├── internal/              # Private application code
│   ├── ai/                # AI recommendation logic
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP request handlers
│   ├── migrations/        # Database migrations
│   ├── models/            # Data models
│   ├── routes/            # Route definitions
│   └── util/              # Utility functions
├── tests/                 # Test files
│   ├── integration/       # Integration tests
│   └── unit/              # Unit tests
├── docs/                  # Documentation
├── main.go                # Application entry point
└── go.mod                 # Go module file
```

## 📋 Prerequisites

- Go 1.19 or higher
- PostgreSQL (production) or SQLite (development)
- Git

## 🛠 Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository_url>
cd the-hub-backend
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Environment Configuration
Copy the example environment file and configure your settings:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=the_hub
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRATION_HOURS=24

# Server Configuration
PORT=8080
GIN_MODE=release

# Development/Production Mode
ENV=development

# OpenRouter AI Configuration (Optional)
# Get your API key from: https://openrouter.ai/keys
# Enables AI-powered scheduling suggestions and natural language processing
OPENROUTER_API_KEY=your_openrouter_api_key_here
```

### 4. Database Setup

#### For PostgreSQL (Production):
```bash
# Create database
createdb the_hub

# Run migrations (if using migration files)
go run main.go migrate
```

#### For SQLite (Development):
The application will automatically create the SQLite database file on first run.

### 5. Run the Application
```bash
# Development mode
go run main.go

# Production mode
GIN_MODE=release go run main.go
```

The server will start on `http://localhost:8080`

## 🧪 Testing

### Run All Tests
```bash
go test ./tests/... -v
```

### Run Unit Tests Only
```bash
go test ./tests/unit/... -v
```

### Run Integration Tests
```bash
go test ./tests/integration/... -v
```

### Test Coverage
```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📚 API Documentation

### Swagger UI
Once the server is running, visit:
```
http://localhost:8080/swagger/index.html
```

### API Endpoints Overview

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh JWT token

#### User Management
- `GET /api/v1/users/{id}` - Get user profile
- `PUT /api/v1/users/{id}` - Update user profile
- `GET /api/v1/users/{id}/settings` - Get user settings
- `PUT /api/v1/users/{id}/settings` - Update user settings
- `PATCH /api/v1/users/{id}/settings` - Partial settings update

#### Task Management
- `GET /api/v1/tasks` - List user tasks
- `POST /api/v1/tasks` - Create new task
- `GET /api/v1/tasks/{id}` - Get specific task
- `PUT /api/v1/tasks/{id}` - Update task
- `DELETE /api/v1/tasks/{id}` - Delete task

#### Goal Management
- `GET /api/v1/goals` - List user goals
- `POST /api/v1/goals` - Create new goal
- `GET /api/v1/goals/{id}` - Get specific goal
- `PUT /api/v1/goals/{id}` - Update goal
- `DELETE /api/v1/goals/{id}` - Delete goal

#### Learning Management
- `GET /api/v1/decks` - List flashcard decks
- `POST /api/v1/decks` - Create new deck
- `GET /api/v1/decks/{id}/cards` - Get cards in deck
- `POST /api/v1/cards` - Create new flashcard
- `GET /api/v1/decks/{id}/review` - Get cards for review

#### Finance Management
- `GET /api/v1/transactions` - List transactions
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/budgets` - List budgets
- `POST /api/v1/budgets` - Create budget
- `GET /api/v1/income` - List income sources

#### Time Management
- `GET /api/v1/schedule` - Get user schedule
- `POST /api/v1/schedule` - Create schedule entry

## 🔧 Development

### Code Style
- Follow Go naming conventions
- Use `gofmt` for code formatting
- Write comprehensive documentation
- Use meaningful variable names

### Adding New Features
1. Create/update models in `internal/models/`
2. Implement handlers in `internal/handlers/`
3. Add routes in `internal/routes/`
4. Write tests in `tests/unit/` or `tests/integration/`
5. Update API documentation

### Database Migrations
```go
// Example migration
func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.NewModel{},
        // ... other models
    )
}
```

## 🐳 Docker Deployment

### Build Docker Image
```bash
docker build -t the-hub-backend .
```

### Run with Docker
```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e JWT_SECRET=your_secret \
  the-hub-backend
```

### Docker Compose
```yaml
version: '3.8'
services:
  the-hub-backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=your_password
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=the_hub
      - POSTGRES_USER=your_user
      - POSTGRES_PASSWORD=your_password
```

## 🔒 Security

### Authentication
- JWT tokens with configurable expiration
- bcrypt password hashing
- Secure token storage guidelines

### Best Practices
- Input validation on all endpoints
- SQL injection prevention via GORM
- XSS protection with input sanitization
- Rate limiting to prevent abuse
- HTTPS enforcement in production

## 📊 Monitoring

### Health Checks
- `GET /health` - Application health status
- Database connectivity checks
- Memory and CPU usage monitoring

### Logging
- Structured JSON logging
- Configurable log levels
- Request/response logging
- Error tracking and reporting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Workflow
- Write tests for new features
- Ensure all tests pass
- Update documentation as needed
- Follow the established code style
- Get code reviewed before merging

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Failed**
- Check your `.env` configuration
- Ensure PostgreSQL is running
- Verify database credentials

**JWT Token Issues**
- Check JWT_SECRET in environment
- Verify token expiration settings
- Ensure proper Authorization header format

**Port Already in Use**
- Change PORT in `.env` file
- Kill process using the port: `lsof -ti:8080 | xargs kill`

### Getting Help
- Check the API documentation at `/swagger/index.html`
- Review the logs for error messages
- Test endpoints with tools like Postman or curl

## 📞 Support

For support and questions:
- Create an issue on GitHub
- Check the documentation in `docs/`
- Review the main project README

---

**Happy coding! 🚀**