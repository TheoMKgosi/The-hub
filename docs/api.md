# The Hub API Documentation

## Overview

The Hub provides a comprehensive REST API for managing personal productivity data including tasks, goals, learning materials, financial information, and user settings. The API is built with Go and Gin, featuring JWT authentication, comprehensive error handling, and detailed Swagger documentation.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Authentication Endpoints

#### POST /auth/login
Authenticate user and return JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

#### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### POST /auth/refresh
Refresh JWT token (if refresh tokens are implemented).

## User Management

### GET /users/{id}
Get user profile information including settings.

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

### PUT /users/{id}
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

### GET /users/{id}/settings
Get only the settings for a specific user.

### PUT /users/{id}/settings
Replace all settings for a user.

**Request Body:**
```json
{
  "theme": "dark",
  "language": "en",
  "notifications": true
}
```

### PATCH /users/{id}/settings
Partially update user settings (merge with existing).

**Request Body:**
```json
{
  "theme": "auto"
}
```

## Task Management

### GET /tasks
Retrieve all tasks for the authenticated user with optional filtering.

**Query Parameters:**
- `status` - Filter by status (pending, in_progress, completed)
- `priority` - Filter by priority (1-5)
- `goal_id` - Filter by associated goal
- `order_by` - Sort field (created_at, due_date, priority)
- `sort` - Sort direction (asc, desc)

**Response:**
```json
[
  {
    "task_id": 1,
    "title": "Complete project proposal",
    "description": "Finish the quarterly project proposal document",
    "status": "pending",
    "priority": 3,
    "due_date": "2024-12-31T23:59:59Z",
    "goal_id": 1,
    "order": 1,
    "user_id": 1
  }
]
```

### POST /tasks
Create a new task.

**Request Body:**
```json
{
  "title": "Complete project proposal",
  "description": "Finish the quarterly project proposal document",
  "due_date": "2024-12-31T23:59:59Z",
  "priority": 3,
  "goal_id": 1,
  "order": 1
}
```

### GET /tasks/{id}
Get a specific task by its ID.

### PUT /tasks/{id}
Update an existing task.

### DELETE /tasks/{id}
Delete a task (soft delete).

### POST /tasks/reorder
Reorder tasks within a goal or list.

**Request Body:**
```json
{
  "task_orders": [
    {
      "task_id": 1,
      "order": 2
    },
    {
      "task_id": 2,
      "order": 1
    }
  ]
}
```

## Goal Management

### GET /goals
Retrieve all goals for the authenticated user.

**Query Parameters:**
- `status` - Filter by status (not_started, in_progress, completed)
- `category` - Filter by category
- `order_by` - Sort field (created_at, target_date, progress)
- `sort` - Sort direction (asc, desc)

### POST /goals
Create a new goal.

**Request Body:**
```json
{
  "title": "Learn Go Programming",
  "description": "Complete Go programming course and build projects",
  "target_date": "2024-06-30T23:59:59Z",
  "category": "learning",
  "status": "in_progress"
}
```

### GET /goals/{id}
Get a specific goal by its ID.

### PUT /goals/{id}
Update an existing goal.

### DELETE /goals/{id}
Delete a goal.

### GET /goals/{id}/tasks
Get all tasks associated with a specific goal.

## Learning Management

### GET /decks
Get all flashcard decks for the user.

### POST /decks
Create a new flashcard deck.

**Request Body:**
```json
{
  "name": "Go Programming Fundamentals",
  "description": "Basic concepts and syntax"
}
```

### GET /decks/{id}
Get a specific deck.

### PUT /decks/{id}
Update a deck.

### DELETE /decks/{id}
Delete a deck.

### GET /decks/{deckID}/cards
Get all cards in a specific deck.

### POST /cards
Create a new flashcard.

**Request Body:**
```json
{
  "question": "What is a goroutine?",
  "answer": "A lightweight thread managed by the Go runtime",
  "deck_id": 1
}
```

### GET /cards/{id}
Get a specific card.

### PUT /cards/{id}
Update a card.

### DELETE /cards/{id}
Delete a card.

### POST /cards/{id}/review
Review a card using spaced repetition algorithm.

**Request Body:**
```json
{
  "quality": 4
}
```

Quality values:
- 0: Complete blackout
- 1: Incorrect response
- 2: Incorrect response with serious difficulty
- 3: Incorrect response with difficulty
- 4: Correct response with difficulty
- 5: Perfect response

### GET /decks/{deckID}/cards/due
Get cards that are due for review in a deck.

## Finance Management

### GET /transactions
Get all financial transactions.

**Query Parameters:**
- `category_id` - Filter by category
- `type` - Filter by type (income, expense)
- `start_date` - Filter by date range
- `end_date` - Filter by date range
- `order_by` - Sort field (date, amount)
- `sort` - Sort direction (asc, desc)

### POST /transactions
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

### GET /transactions/{id}
Get a specific transaction.

### PUT /transactions/{id}
Update a transaction.

### DELETE /transactions/{id}
Delete a transaction.

### GET /categories
Get all transaction categories.

### POST /categories
Create a new category.

**Request Body:**
```json
{
  "name": "Groceries",
  "type": "expense",
  "color": "#FF5733"
}
```

### GET /budgets
Get all budgets.

### POST /budgets
Create a new budget.

**Request Body:**
```json
{
  "amount": 1500.00,
  "category_id": 1,
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}
```

### GET /income
Get all income sources.

### POST /income
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

## Time Management

### GET /schedule
Get user's schedule and time blocks.

### POST /schedule
Create a new schedule entry.

**Request Body:**
```json
{
  "title": "Team Meeting",
  "description": "Weekly team standup",
  "start_time": "2024-01-15T09:00:00Z",
  "end_time": "2024-01-15T10:00:00Z",
  "category": "work"
}
```

### GET /schedule/{id}
Get a specific schedule entry.

### PUT /schedule/{id}
Update a schedule entry.

### DELETE /schedule/{id}
Delete a schedule entry.

## Calendar Integration

### GET /calendar/integrations
Get calendar integration settings.

### POST /calendar/integrations
Set up calendar integration.

**Request Body:**
```json
{
  "provider": "google",
  "calendar_id": "primary",
  "sync_enabled": true
}
```

### GET /calendar/zones
Get calendar zones for time blocking.

### POST /calendar/zones
Create a calendar zone.

**Request Body:**
```json
{
  "name": "Deep Work",
  "color": "#4CAF50",
  "start_time": "09:00",
  "end_time": "12:00",
  "days": ["monday", "tuesday", "wednesday", "thursday", "friday"]
}
```

## Analytics & Statistics

### GET /stats/tasks
Get task completion statistics.

### GET /stats/goals
Get goal achievement statistics.

### GET /stats/learning
Get learning progress statistics.

### GET /stats/finance
Get financial analytics.

### GET /stats/time
Get time tracking analytics.

## Push Notifications

### POST /push/subscribe
Subscribe to push notifications.

**Request Body:**
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/...",
  "keys": {
    "p256dh": "...",
    "auth": "..."
  }
}
```

### POST /push/send
Send a push notification (admin endpoint).

## Tags Management

### GET /tags
Get all tags.

### POST /tags
Create a new tag.

**Request Body:**
```json
{
  "name": "Urgent",
  "color": "#FF0000"
}
```

### GET /tags/{id}
Get a specific tag.

### PUT /tags/{id}
Update a tag.

### DELETE /tags/{id}
Delete a tag.

## Learning Paths

### GET /learning-paths
Get all learning paths.

### POST /learning-paths
Create a new learning path.

**Request Body:**
```json
{
  "title": "Full Stack Development",
  "description": "Complete web development curriculum",
  "topics": [1, 2, 3]
}
```

### GET /learning-paths/{id}
Get a specific learning path.

### PUT /learning-paths/{id}
Update a learning path.

### DELETE /learning-paths/{id}
Delete a learning path.

## Study Sessions

### GET /study-sessions
Get study session history.

### POST /study-sessions
Record a new study session.

**Request Body:**
```json
{
  "topic_id": 1,
  "duration": 3600,
  "notes": "Covered basic algorithms"
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
- `409` - Conflict
- `422` - Unprocessable Entity
- `500` - Internal Server Error

**Error Response Format:**
```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": {}
}
```

## Rate Limiting

API endpoints are rate-limited to prevent abuse:

- **Authenticated users**: 1000 requests per hour
- **Unauthenticated endpoints**: 100 requests per hour
- **Login endpoint**: 5 attempts per 15 minutes

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Data Formats

- **Dates**: ISO 8601 format (`YYYY-MM-DDTHH:MM:SSZ`)
- **Currency amounts**: Decimal format (e.g., `50.00`)
- **Boolean values**: `true` or `false`
- **UUIDs**: Used for all primary keys
- **Colors**: Hex color codes (e.g., `#FF5733`)

## Pagination

List endpoints support pagination:

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)

**Response Format:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

## Filtering & Sorting

Most list endpoints support filtering and sorting:

**Common Query Parameters:**
- `order_by` - Field to sort by
- `sort` - Sort direction (`asc` or `desc`)
- `search` - Search term for text fields
- `start_date` - Filter by date range
- `end_date` - Filter by date range

## WebSocket Support

Real-time updates are available via WebSocket:

**Connection URL:**
```
ws://localhost:8080/ws
```

**Supported Events:**
- Task status changes
- Goal progress updates
- New notifications
- Calendar updates
- Real-time collaboration

## Versioning

The API uses URL versioning:

- Current version: `v1`
- Future versions will be available at `/api/v2`, etc.

## SDKs & Libraries

### JavaScript/TypeScript Client
```javascript
import { HubAPI } from '@the-hub/api-client';

const client = new HubAPI({
  baseURL: 'http://localhost:8080/api/v1',
  token: 'your-jwt-token'
});

// Example usage
const tasks = await client.tasks.list();
const newTask = await client.tasks.create({
  title: 'New Task',
  priority: 3
});
```

### Go Client
```go
package main

import (
    "github.com/the-hub/go-client"
)

func main() {
    client := hub.NewClient("http://localhost:8080/api/v1", "your-jwt-token")

    tasks, err := client.Tasks.List(context.Background(), &hub.TaskListOptions{})
    if err != nil {
        log.Fatal(err)
    }
}
```

## Additional Resources

- [Swagger UI](http://localhost:8080/swagger/index.html) - Interactive API documentation
- [Postman Collection](https://github.com/your-org/the-hub/tree/main/docs/postman) - Importable collection for testing
- [OpenAPI Specification](https://github.com/your-org/the-hub/tree/main/docs/openapi) - API specification files
- [WebSocket Documentation](docs/websockets.md) - Real-time features documentation