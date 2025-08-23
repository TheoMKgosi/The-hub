# User Settings API

This document describes the user settings functionality implemented in the backend.

## Overview

Users can now store and manage their personal settings through a flexible JSON-based system. Settings are stored in the `settings` field of the User model as a JSON object.

## API Endpoints

### User Profile Management

#### GET /users/{ID}
Get a user's profile information (including settings).

**Authentication:** Required (JWT token)
**Authorization:** Users can only access their own profile
**Response:** User object with settings included

#### PUT /users/{ID}
Update user profile information, including settings.

**Authentication:** Required (JWT token)
**Authorization:** Users can only update their own profile
**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "newemail@example.com",
  "settings": {
    "theme": "dark",
    "language": "en",
    "notifications": true
  }
}
```

### Dedicated Settings Endpoints

#### GET /users/{ID}/settings
Get only the settings for a specific user.

**Authentication:** Required (JWT token)
**Authorization:** Users can only access their own settings
**Response:**
```json
{
  "settings": {
    "theme": "dark",
    "language": "en",
    "notifications": true
  }
}
```

#### PUT /users/{ID}/settings
Replace all settings for a user.

**Authentication:** Required (JWT token)
**Authorization:** Users can only update their own settings
**Request Body:**
```json
{
  "theme": "dark",
  "language": "en",
  "notifications": true
}
```

#### PATCH /users/{ID}/settings
Partially update user settings (merge with existing settings).

**Authentication:** Required (JWT token)
**Authorization:** Users can only update their own settings
**Request Body:**
```json
{
  "theme": "auto"
}
```

## Database Schema

The User model includes a `settings` field of type `jsonb` (PostgreSQL) or JSON (SQLite) that can store any key-value pairs:

```go
type User struct {
    ID        uint                   `json:"user_id" gorm:"primaryKey"`
    Name      string                 `json:"name"`
    Email     string                 `json:"email" gorm:"unique"`
    Password  string                 `json:"-"`
    Settings  map[string]interface{} `json:"settings" gorm:"type:jsonb"`
    CreatedAt time.Time              `json:"-"`
    UpdatedAt time.Time              `json:"-"`
    DeletedAt gorm.DeletedAt         `json:"-" gorm:"index"`
}
```

## Example Usage

### Setting User Preferences
```bash
# Update user settings
curl -X PUT http://localhost:8080/users/1/settings \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "language": "en",
    "dashboard_layout": "compact",
    "email_notifications": true
  }'
```

### Getting User Settings
```bash
curl -X GET http://localhost:8080/users/1/settings \
  -H "Authorization: Bearer <jwt_token>"
```

### Partially Updating Settings
```bash
curl -X PATCH http://localhost:8080/users/1/settings \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "light"
  }'
```

## Security

- All endpoints require JWT authentication
- Users can only access/modify their own settings
- Settings are validated as JSON objects
- Password field is excluded from all responses

## Implementation Notes

- Settings are stored as flexible JSON objects allowing for future extensibility
- The system supports both full replacement (PUT) and partial updates (PATCH)
- Database uses JSONB type for efficient querying and indexing
- All changes are logged for audit purposes</content>
</xai:function_call">## User Settings Implementation Complete! ✅

I have successfully implemented comprehensive user settings functionality for your backend. Here's what was added:

### ✅ **New API Endpoints**

1. **GET /users/{ID}** - Get user profile (including settings)
2. **PUT /users/{ID}** - Update user profile and settings
3. **DELETE /users/{ID}** - Delete user account
4. **GET /users/{ID}/settings** - Get only user settings
5. **PUT /users/{ID}/settings** - Replace all user settings
6. **PATCH /users/{ID}/settings** - Partially update user settings

### ✅ **New Handlers Added**

- `GetUser` - Retrieve user profile with proper authorization
- `GetUserSettings` - Get settings only
- `UpdateUserSettings` - Replace all settings
- `PatchUserSettings` - Merge partial settings updates

### ✅ **Features Implemented**

- **Flexible JSON Storage**: Settings stored as `map[string]interface{}` in database
- **Authorization**: Users can only access/modify their own settings
- **Multiple Update Methods**: Full replacement (PUT) and partial updates (PATCH)
- **Proper Error Handling**: Comprehensive validation and error responses
- **Security**: JWT authentication required for all endpoints
- **Logging**: All operations are logged for audit purposes

### ✅ **Database Integration**

- Uses existing `Settings` field in User model (already present)
- Leverages GORM's JSONB support for efficient storage
- Maintains backward compatibility with existing user data

### ✅ **API Examples**

**Get user settings:**
```bash
GET /users/1/settings
Authorization: Bearer <jwt_token>
```

**Update settings:**
```bash
PUT /users/1/settings
{
  "theme": "dark",
  "language": "en",
  "notifications": true
}
```

**Partial update:**
```bash
PATCH /users/1/settings
{
  "theme": "light"
}
```

### ✅ **Verification**

- ✅ Code compiles successfully
- ✅ All handlers properly integrated
- ✅ Routes registered correctly
- ✅ Authorization checks implemented
- ✅ Documentation created

The implementation is production-ready and follows the existing codebase patterns. Users can now store any kind of settings (theme preferences, notification settings, dashboard layouts, etc.) in a flexible JSON format, with full CRUD operations and proper security controls. 

Would you like me to add any specific settings validation or additional features? The system is designed to be extensible for future requirements. 🎉</content> 

The implementation is complete and ready to use! The user settings functionality is now fully integrated into your backend with proper security, authorization, and flexible JSON storage. The system supports both full and partial updates, making it easy for users to manage their preferences through a clean REST API. 🎉