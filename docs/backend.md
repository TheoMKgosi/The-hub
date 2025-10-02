# Backend Documentation

This document provides detailed information about The Hub's backend architecture, models, handlers, and API endpoints.

## Technology Stack

- **Language:** Go 1.24+
- **Framework:** Gin (HTTP web framework)
- **ORM:** GORM (Go ORM)
- **Database:** PostgreSQL (production) / SQLite (development)
- **Authentication:** JWT (JSON Web Tokens)
- **Documentation:** Swagger/OpenAPI
- **Testing:** Go's built-in testing package + stretchr/testify
- **Load Testing:** k6
- **AI Integration:** OpenRouter API for recommendations
- **Push Notifications:** Web Push API support
- **Calendar Integration:** External calendar service integration

## Project Structure

```
the-hub-backend/
├── cmd/                   # Command-line applications
│   └── task-cleaner/      # Background task cleanup service
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
├── docs/                  # API documentation
├── main.go                # Application entry point
├── go.mod                 # Go module file
└── go.sum                 # Go module checksums
```

## Configuration

The backend uses environment-based configuration with the following key settings:

```go
type Config struct {
    Database DatabaseConfig
    JWT      JWTConfig
    Server   ServerConfig
    Logger   LoggerConfig
}
```

### Database Configuration
- **PostgreSQL:** Production database with connection pooling and native UUID support
- **UUID Primary Keys:** All models use UUID for better security and scalability
- **Migrations:** Automatic schema migrations using GORM

### JWT Configuration
- **Secret:** Environment-based JWT signing secret
- **Expiration:** Configurable token expiration times
- **Middleware:** Authentication middleware for protected routes

## Data Models

### User Model
```go
type User struct {
    ID        uuid.UUID              `json:"user_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name      string                 `json:"name"`
    Email     string                 `json:"email" gorm:"unique"`
    Password  string                 `json:"-"`
    Settings  map[string]interface{} `json:"settings" gorm:"type:jsonb"`
    CreatedAt time.Time              `json:"-"`
    UpdatedAt time.Time              `json:"-"`
    DeletedAt gorm.DeletedAt         `json:"-" gorm:"index"`
}
```

**Features:**
- JSONB settings field for flexible user preferences
- Soft delete support with DeletedAt field
- Unique email constraint
- Password hashing (not exposed in JSON)

### Task Model
```go
type Task struct {
    ID          uuid.UUID      `json:"task_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID      `json:"user_id"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description"`
    Status      string         `json:"status" gorm:"default:pending"`
    Priority    *int           `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
    DueDate     *time.Time     `json:"due_date"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
}
```

**Features:**
- User association with foreign key
- Optional due date
- Status and priority fields
- Soft delete support

### Goal Model
```go
type Goal struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    UserID      uint           `json:"user_id"`
    Title       string         `json:"title"`
    Description string         `json:"description"`
    Category    string         `json:"category"`
    Status      string         `json:"status"`
    TargetDate  *time.Time     `json:"target_date"`
    Progress    int            `json:"progress"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
}
```

**Features:**
- Progress tracking (0-100%)
- Category classification
- Target date for completion tracking

### Learning Models

#### Topic Model
```go
type Topic struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    UserID      uint           `json:"user_id"`
    Name        string         `json:"name"`
    Description string         `json:"description"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
}
```

#### Deck Model
```go
type Deck struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    UserID      uint           `json:"user_id"`
    TopicID     uint           `json:"topic_id"`
    Name        string         `json:"name"`
    Description string         `json:"description"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
    Topic       Topic          `json:"-" gorm:"foreignKey:TopicID"`
}
```

#### Card Model
```go
type Card struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    DeckID       uint           `json:"deck_id"`
    Question     string         `json:"question"`
    Answer       string         `json:"answer"`
    Repetitions  int            `json:"repetitions"`
    Easiness     float64        `json:"easiness"`
    Interval     int            `json:"interval"`
    NextReview   time.Time      `json:"next_review"`
    CreatedAt    time.Time      `json:"-"`
    UpdatedAt    time.Time      `json:"-"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    Deck         Deck           `json:"-" gorm:"foreignKey:DeckID"`
}
```

**Features:**
- Spaced repetition algorithm implementation
- SM-2 algorithm parameters (easiness, interval, repetitions)
- Next review date calculation

### Finance Models

#### Transaction Model
```go
type Transaction struct {
    ID          uuid.UUID      `json:"transaction_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID      `json:"user_id"`
    Amount      float64        `json:"amount"`
    Description string         `json:"description"`
    CategoryID  uuid.UUID      `json:"category_id"`
    Date        time.Time      `json:"date"`
    Type        string         `json:"type"` // "income" or "expense"
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
    Category    Category       `json:"-" gorm:"foreignKey:CategoryID"`
}
```

#### Category Model
```go
type Category struct {
    ID          uuid.UUID      `json:"category_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID      `json:"user_id"`
    Name        string         `json:"name"`
    Type        string         `json:"type"` // "income" or "expense"
    Color       string         `json:"color"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
}
```

#### Budget Model
```go
type Budget struct {
    ID         uuid.UUID `json:"budget_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID     uuid.UUID `json:"user_id"`
    Amount     float64   `json:"amount"`
    CategoryID uuid.UUID `json:"category_id"`
    StartDate  time.Time `json:"start_date"`
    EndDate    time.Time `json:"end_date"`
    CreatedAt  time.Time `json:"-"`
    UpdatedAt  time.Time `json:"-"`
    User       User      `json:"-" gorm:"foreignKey:UserID"`
    Category   Category  `json:"-" gorm:"foreignKey:CategoryID"`
}
```

#### Income Model
```go
type Income struct {
    ID         uuid.UUID `json:"income_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID     uuid.UUID `json:"user_id"`
    Name       string    `json:"name"`
    Amount     float64   `json:"amount"`
    Frequency  string    `json:"frequency"` // "monthly", "weekly", etc.
    StartDate  time.Time `json:"start_date"`
    CreatedAt  time.Time `json:"-"`
    UpdatedAt  time.Time `json:"-"`
    User       User      `json:"-" gorm:"foreignKey:UserID"`
}
```

### Schedule Models

#### ScheduledTask Model
```go
type ScheduledTask struct {
    ID          uuid.UUID      `json:"scheduled_task_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID      `json:"user_id"`
    Title       string         `json:"title"`
    Description string         `json:"description"`
    StartTime   time.Time      `json:"start_time"`
    EndTime     time.Time      `json:"end_time"`
    Category    string         `json:"category"`
    CreatedAt   time.Time      `json:"-"`
    UpdatedAt   time.Time      `json:"-"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    User        User           `json:"-" gorm:"foreignKey:UserID"`
}
```

### AI Recommendation Models

#### AIRecommendation Model
```go
type AIRecommendation struct {
    ID          uuid.UUID `json:"recommendation_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID `json:"user_id"`
    Type        string    `json:"type"` // "task", "goal", "learning", etc.
    Content     string    `json:"content"`
    Priority    int       `json:"priority"`
    IsRead      bool      `json:"is_read" gorm:"default:false"`
    CreatedAt   time.Time `json:"-"`
    User        User      `json:"-" gorm:"foreignKey:UserID"`
}
```

### Calendar Integration Models

#### CalendarIntegration Model
```go
type CalendarIntegration struct {
    ID            uuid.UUID `json:"calendar_integration_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID        uuid.UUID `json:"user_id"`
    Provider      string    `json:"provider"` // "google", "outlook", etc.
    AccessToken   string    `json:"-"`
    RefreshToken  string    `json:"-"`
    TokenExpiry   time.Time `json:"-"`
    CalendarID    string    `json:"calendar_id"`
    IsActive      bool      `json:"is_active" gorm:"default:true"`
    CreatedAt     time.Time `json:"-"`
    UpdatedAt     time.Time `json:"-"`
    User          User      `json:"-" gorm:"foreignKey:UserID"`
}
```

#### CalendarZone Model
```go
type CalendarZone struct {
    ID          uuid.UUID `json:"calendar_zone_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      uuid.UUID `json:"user_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Color       string    `json:"color"`
    Timezone    string    `json:"timezone"`
    CreatedAt   time.Time `json:"-"`
    UpdatedAt   time.Time `json:"-"`
    User        User      `json:"-" gorm:"foreignKey:UserID"`
}
```

### Push Notification Models

#### PushSubscription Model
```go
type PushSubscription struct {
    ID       uuid.UUID `json:"push_subscription_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID   uuid.UUID `json:"user_id"`
    Endpoint string    `json:"endpoint"`
    P256dh   string    `json:"p256dh"`
    Auth     string    `json:"auth"`
    User     User      `json:"-" gorm:"foreignKey:UserID"`
}
```

### Authentication Models

#### PasswordResetToken Model
```go
type PasswordResetToken struct {
    ID        uuid.UUID `json:"reset_token_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID    uuid.UUID `json:"user_id"`
    Token     string    `json:"token" gorm:"unique"`
    ExpiresAt time.Time `json:"expires_at"`
    Used      bool      `json:"used" gorm:"default:false"`
    CreatedAt time.Time `json:"-"`
    User      User      `json:"-" gorm:"foreignKey:UserID"`
}
```

#### RefreshToken Model
```go
type RefreshToken struct {
    ID        uuid.UUID `json:"refresh_token_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID    uuid.UUID `json:"user_id"`
    Token     string    `json:"token" gorm:"unique"`
    ExpiresAt time.Time `json:"expires_at"`
    Revoked   bool      `json:"revoked" gorm:"default:false"`
    CreatedAt time.Time `json:"-"`
    User      User      `json:"-" gorm:"foreignKey:UserID"`
}
```

## HTTP Handlers

### User Handlers

#### Authentication
- `Register`: User registration with password hashing
- `Login`: JWT token generation and validation
- `GetProfile`: Retrieve user profile with settings
- `UpdateProfile`: Update user information
- `DeleteUser`: Soft delete user account

#### Settings Management
- `GetUserSettings`: Retrieve user settings
- `UpdateUserSettings`: Replace all user settings
- `PatchUserSettings`: Partial settings update

### Task Handlers
- `GetTasks`: List all user tasks with filtering
- `GetTask`: Retrieve specific task
- `CreateTask`: Create new task
- `UpdateTask`: Modify existing task
- `DeleteTask`: Soft delete task

### Goal Handlers
- `GetGoals`: List all user goals
- `GetGoal`: Retrieve specific goal
- `CreateGoal`: Create new goal
- `UpdateGoal`: Modify goal
- `DeleteGoal`: Soft delete goal

### Learning Handlers
- `GetTopics`: List learning topics
- `CreateTopic`: Create new topic
- `GetDecks`: List flashcard decks
- `CreateDeck`: Create new deck
- `GetCards`: List cards in deck
- `CreateCard`: Create new flashcard
- `ReviewCards`: Get cards due for review
- `UpdateCardProgress`: Update card review progress

### Finance Handlers
- `GetTransactions`: List financial transactions
- `CreateTransaction`: Create new transaction
- `GetCategories`: List transaction categories
- `CreateCategory`: Create new category
- `GetBudgets`: List budgets
- `CreateBudget`: Create new budget
- `GetIncome`: List income sources
- `CreateIncome`: Create new income source

### Schedule Handlers
- `GetSchedule`: List scheduled tasks
- `CreateScheduledTask`: Create new scheduled task
- `UpdateScheduledTask`: Modify scheduled task
- `DeleteScheduledTask`: Delete scheduled task

### AI Recommendation Handlers
- `GetRecommendations`: Get AI-generated recommendations
- `MarkRecommendationRead`: Mark recommendation as read
- `GenerateRecommendations`: Trigger recommendation generation

### Calendar Integration Handlers
- `ConnectCalendar`: Connect external calendar service
- `GetCalendarEvents`: Retrieve calendar events
- `SyncCalendarEvents`: Sync events with external calendar
- `DisconnectCalendar`: Disconnect calendar integration

### Push Notification Handlers
- `SubscribePush`: Subscribe to push notifications
- `UnsubscribePush`: Unsubscribe from push notifications
- `SendPushNotification`: Send push notification to user

### Statistics Handlers
- `GetUserStats`: Get comprehensive user statistics
- `GetTaskStats`: Get task completion statistics
- `GetGoalStats`: Get goal progress statistics
- `GetLearningStats`: Get learning progress statistics

### Learning Path Handlers
- `GetLearningPaths`: Get available learning paths
- `CreateLearningPath`: Create custom learning path
- `UpdateLearningPath`: Modify learning path
- `DeleteLearningPath`: Delete learning path

### Study Session Handlers
- `StartStudySession`: Start a study session
- `EndStudySession`: End a study session
- `GetStudySessions`: Get study session history
- `GetStudyStats`: Get study statistics

### Tag Handlers
- `GetTags`: Get all user tags
- `CreateTag`: Create new tag
- `UpdateTag`: Modify tag
- `DeleteTag`: Delete tag
- `TagResource`: Tag a resource (task, goal, etc.)

## Middleware

### Authentication Middleware
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Validate JWT token
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### CORS Middleware
Handles Cross-Origin Resource Sharing for frontend integration.

### Logging Middleware
Logs all HTTP requests with timing information.

## Database Migrations

The application uses GORM's auto-migration feature for schema management:

```go
func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Task{},
        &models.Goal{},
        &models.Topic{},
        &models.Deck{},
        &models.Card{},
        &models.Transaction{},
        &models.Category{},
        &models.Budget{},
        &models.Income{},
        &models.ScheduledTask{},
    )
}
```

## Error Handling

### Custom Error Types
```go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"status"`
}

func (e *AppError) Error() string {
    return e.Message
}
```

### Global Error Handler
```go
func ErrorHandler() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if err, ok := recovered.(error); ok {
            log.Printf("Panic recovered: %v", err)
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
        c.AbortWithStatus(500)
    })
}
```

## Testing

### Unit Tests
```go
func TestCreateTask(t *testing.T) {
    // Test task creation logic
    task := &models.Task{
        Title:       "Test Task",
        Description: "Test Description",
        UserID:      1,
    }

    err := task.Validate()
    assert.NoError(t, err)
}
```

### Integration Tests
```go
func TestUserAPI(t *testing.T) {
    router := setupTestRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## Security Features

### Password Hashing
Uses bcrypt for secure password hashing:
```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### JWT Authentication
- Tokens expire after configurable time period
- Refresh token mechanism for extended sessions
- Secure token storage recommendations

### Input Validation
- Request payload validation using struct tags
- SQL injection prevention with GORM
- XSS protection with input sanitization

### Rate Limiting
- API rate limiting to prevent abuse
- Configurable limits per endpoint
- IP-based and user-based restrictions

## Performance Optimizations

### Database
- Connection pooling with PostgreSQL
- Database indexing on frequently queried fields
- Query optimization with eager loading

### Caching
- Redis integration for session storage
- In-memory caching for frequently accessed data
- Cache invalidation strategies

### API
- Pagination for large datasets
- Efficient JSON serialization
- Gzip compression for responses

## Deployment

### Build Process
```bash
# Build the application
go mod download
go build -o the-hub-backend

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o the-hub-backend-linux
GOOS=windows GOARCH=amd64 go build -o the-hub-backend.exe
```

### Environment Configuration
```bash
# Production environment variables
DB_HOST=prod-db-host
DB_PORT=5432
DB_USER=prod-user
DB_PASSWORD=prod-password
DB_NAME=the_hub_prod

JWT_SECRET=your-production-jwt-secret
SERVER_PORT=8080
```

### Health Checks
- `/health` endpoint for load balancer health checks
- Database connectivity checks
- Memory and CPU usage monitoring

## Monitoring and Logging

### Structured Logging
```go
logger := logrus.New()
logger.SetFormatter(&logrus.JSONFormatter{})
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action": "task_created",
    "task_id": taskID,
}).Info("Task created successfully")
```

### Metrics
- Request/response metrics
- Database query performance
- Error rate monitoring
- User activity tracking

## Development Guidelines

### Code Style
- Follow Go naming conventions
- Use `gofmt` for code formatting
- Write comprehensive documentation
- Use meaningful variable names

### API Design
- RESTful endpoint design
- Consistent JSON response format
- Proper HTTP status codes
- Versioned API endpoints

### Database Design
- Use foreign keys for relationships
- Implement soft deletes where appropriate
- Add database indexes for performance
- Use transactions for complex operations

### Security Best Practices
- Validate all user inputs
- Use parameterized queries
- Implement proper authentication
- Log security events
- Regular security updates