# The Hub API Documentation

This document provides a comprehensive overview of The Hub's REST API endpoints, organized by feature area.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All API endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### User Management

#### Get User Profile
```http
GET /users/{id}
```
Get a user's profile information including settings.

**Response:**
```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "settings": {
    "theme": "dark",
    "language": "en",
    "notifications": true
  }
}
```

#### Update User Profile
```http
PUT /users/{id}
```
Update user profile information and settings.

**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "newemail@example.com",
  "settings": {
    "theme": "light",
    "language": "en"
  }
}
```

#### Get User Settings
```http
GET /users/{id}/settings
```
Get only the settings for a specific user.

#### Update User Settings
```http
PUT /users/{id}/settings
```
Replace all settings for a user.

**Request Body:**
```json
{
  "theme": "dark",
  "language": "en",
  "notifications": true
}
```

#### Patch User Settings
```http
PATCH /users/{id}/settings
```
Partially update user settings (merge with existing).

**Request Body:**
```json
{
  "theme": "auto"
}
```

### Task Management

#### Get All Tasks
```http
GET /tasks
```
Retrieve all tasks for the authenticated user.

#### Create Task
```http
POST /tasks
```
Create a new task.

**Request Body:**
```json
{
  "title": "Complete project",
  "description": "Finish the quarterly project",
  "due_date": "2024-12-31T23:59:59Z",
  "priority": "high",
  "status": "pending"
}
```

#### Get Task by ID
```http
GET /tasks/{id}
```
Get a specific task by its ID.

#### Update Task
```http
PUT /tasks/{id}
```
Update an existing task.

#### Delete Task
```http
DELETE /tasks/{id}
```
Delete a task.

### Goal Management

#### Get All Goals
```http
GET /goals
```
Retrieve all goals for the authenticated user.

#### Create Goal
```http
POST /goals
```
Create a new goal.

**Request Body:**
```json
{
  "title": "Learn Go",
  "description": "Complete Go programming course",
  "target_date": "2024-06-30T23:59:59Z",
  "category": "learning",
  "status": "in_progress"
}
```

#### Get Goal by ID
```http
GET /goals/{id}
```

#### Update Goal
```http
PUT /goals/{id}
```

#### Delete Goal
```http
DELETE /goals/{id}
```

### Learning Management

#### Get All Decks
```http
GET /decks
```
Get all flashcard decks.

#### Create Deck
```http
POST /decks
```
Create a new flashcard deck.

**Request Body:**
```json
{
  "name": "Go Programming",
  "description": "Learn Go fundamentals",
  "topic_id": 1
}
```

#### Get Cards in Deck
```http
GET /decks/{id}/cards
```
Get all cards in a specific deck.

#### Create Card
```http
POST /cards
```
Create a new flashcard.

**Request Body:**
```json
{
  "question": "What is a goroutine?",
  "answer": "A lightweight thread managed by the Go runtime",
  "deck_id": 1
}
```

#### Review Cards
```http
GET /decks/{id}/review
```
Get cards due for review using spaced repetition.

### Finance Management

#### Get All Transactions
```http
GET /transactions
```
Get all financial transactions.

#### Create Transaction
```http
POST /transactions
```
Create a new transaction.

**Request Body:**
```json
{
  "amount": 50.00,
  "description": "Grocery shopping",
  "category_id": 1,
  "date": "2024-01-15T10:00:00Z",
  "type": "expense"
}
```

#### Get Budgets
```http
GET /budgets
```
Get all budgets.

#### Create Budget
```http
POST /budgets
```
Create a new budget.

**Request Body:**
```json
{
  "amount": 1500.00,
  "category_id": 1,
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "income_id": 1
}
```

#### Get Income Sources
```http
GET /income
```
Get all income sources.

#### Create Income Source
```http
POST /income
```
Create a new income source.

**Request Body:**
```json
{
  "name": "Salary",
  "amount": 5000.00,
  "frequency": "monthly",
  "start_date": "2024-01-01"
}
```

### Time Management

#### Get Schedule
```http
GET /schedule
```
Get user's schedule and time blocks.

#### Create Schedule Entry
```http
POST /schedule
```
Create a new schedule entry.

**Request Body:**
```json
{
  "title": "Meeting",
  "description": "Team standup",
  "start_time": "2024-01-15T09:00:00Z",
  "end_time": "2024-01-15T10:00:00Z",
  "category": "work"
}
```

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

**Error Response Format:**
```json
{
  "error": "Error message description",
  "code": "ERROR_CODE"
}
```

## Rate Limiting

API endpoints are rate-limited to prevent abuse. Current limits:
- 1000 requests per hour for authenticated users
- 100 requests per hour for unauthenticated endpoints

## Data Formats

- All dates use ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ`
- Currency amounts are in decimal format (e.g., 50.00)
- Boolean values are `true` or `false`

## WebSocket Support

Real-time updates are available via WebSocket for:
- Task status changes
- Goal progress updates
- Notification alerts

Connect to: `ws://localhost:8080/ws`

## Additional Resources

- [Swagger UI](http://localhost:8080/swagger/index.html) - Interactive API documentation
- [User Settings API](user-settings-api.md) - Detailed user settings documentation
- [Postman Collection](https://example.com/collection) - Importable collection for testing