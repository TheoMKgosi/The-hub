package unit

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestUser is a simplified user model for testing with SQLite
type TestUser struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	Settings  string         `json:"settings"` // Use string instead of JSON for SQLite
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

var modelTestDB *gorm.DB

func setupModelTestDB() {
	if modelTestDB == nil {
		dsn := ":memory:"
		var err error
		modelTestDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic("Failed to connect to test database")
		}

		// Auto-migrate the schema
		err = modelTestDB.AutoMigrate(&TestUser{})
		if err != nil {
			panic("Failed to migrate test database")
		}
	}
}

func TestUserModel(t *testing.T) {
	setupModelTestDB()

	tests := []struct {
		name        string
		user        *TestUser
		wantErr     bool
		description string
	}{
		{
			name: "valid user",
			user: &TestUser{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "hashedpassword",
				Settings: "{}", // Use JSON string for SQLite
			},
			wantErr:     false,
			description: "should create user with valid data",
		},
		{
			name: "user with empty email",
			user: &TestUser{
				Email:    "",
				Name:     "Test User",
				Password: "hashedpassword",
				Settings: "{}",
			},
			wantErr:     false, // SQLite allows empty strings
			description: "should handle empty email",
		},
		{
			name: "user with empty name",
			user: &TestUser{
				Email:    "test2@example.com",
				Name:     "",
				Password: "hashedpassword",
				Settings: "{}",
			},
			wantErr:     false,
			description: "should handle empty name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before each test
			modelTestDB.Unscoped().Where("1=1").Delete(&TestUser{})

			// Attempt to create user
			result := modelTestDB.Create(tt.user)

			if (result.Error != nil) != tt.wantErr {
				t.Errorf("User creation error = %v, wantErr %v", result.Error, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify user was created
				if tt.user.ID == 0 {
					t.Error("User should have been assigned an ID")
				}

				// Verify user can be retrieved
				var retrievedUser TestUser
				err := modelTestDB.First(&retrievedUser, tt.user.ID).Error
				if err != nil {
					t.Errorf("Failed to retrieve created user: %v", err)
				}

				// Verify data integrity
				if retrievedUser.Email != tt.user.Email {
					t.Errorf("Retrieved email %s doesn't match original %s", retrievedUser.Email, tt.user.Email)
				}

				if retrievedUser.Name != tt.user.Name {
					t.Errorf("Retrieved name %s doesn't match original %s", retrievedUser.Name, tt.user.Name)
				}
			}
		})
	}
}

func TestUserUniqueEmail(t *testing.T) {
	setupModelTestDB()

	// Clean up before test
	modelTestDB.Unscoped().Where("1=1").Delete(&TestUser{})

	// Create first user
	user1 := &TestUser{
		Email:    "duplicate@example.com",
		Name:     "User 1",
		Password: "password1",
		Settings: "{}",
	}

	result1 := modelTestDB.Create(user1)
	if result1.Error != nil {
		t.Fatalf("Failed to create first user: %v", result1.Error)
	}

	// Attempt to create second user with same email
	user2 := &TestUser{
		Email:    "duplicate@example.com",
		Name:     "User 2",
		Password: "password2",
		Settings: "{}",
	}

	result2 := modelTestDB.Create(user2)

	// Should fail due to unique constraint
	if result2.Error == nil {
		t.Error("Expected error when creating user with duplicate email, but got none")
	}

	// Verify only one user exists
	var users []TestUser
	modelTestDB.Find(&users)
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

func TestUserCRUDOperations(t *testing.T) {
	setupModelTestDB()

	// Clean up before test
	modelTestDB.Unscoped().Where("1=1").Delete(&TestUser{})

	// CREATE
	user := &TestUser{
		Email:    "crud@example.com",
		Name:     "CRUD User",
		Password: "password",
		Settings: "{\"theme\": \"light\"}",
	}

	result := modelTestDB.Create(user)
	if result.Error != nil {
		t.Fatalf("Failed to create user: %v", result.Error)
	}

	// READ
	var retrievedUser TestUser
	err := modelTestDB.First(&retrievedUser, user.ID).Error
	if err != nil {
		t.Fatalf("Failed to read user: %v", err)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Retrieved email %s doesn't match original %s", retrievedUser.Email, user.Email)
	}

	// UPDATE
	updates := map[string]interface{}{
		"name":  "Updated Name",
		"email": "updated@example.com",
	}

	err = modelTestDB.Model(&retrievedUser).Updates(updates).Error
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	var updatedUser TestUser
	err = modelTestDB.First(&updatedUser, user.ID).Error
	if err != nil {
		t.Fatalf("Failed to read updated user: %v", err)
	}

	if updatedUser.Name != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got '%s'", updatedUser.Name)
	}

	if updatedUser.Email != "updated@example.com" {
		t.Errorf("Expected updated email 'updated@example.com', got '%s'", updatedUser.Email)
	}

	// DELETE
	err = modelTestDB.Delete(&updatedUser).Error
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	var deletedUser TestUser
	err = modelTestDB.First(&deletedUser, user.ID).Error
	if err == nil {
		t.Error("User should have been deleted")
	}
}

func TestUserSoftDelete(t *testing.T) {
	setupModelTestDB()

	// Clean up before test
	modelTestDB.Unscoped().Where("1=1").Delete(&TestUser{})

	// Create user
	user := &TestUser{
		Email:    "softdelete@example.com",
		Name:     "Soft Delete User",
		Password: "password",
		Settings: "{}",
	}

	result := modelTestDB.Create(user)
	if result.Error != nil {
		t.Fatalf("Failed to create user: %v", result.Error)
	}

	// Soft delete
	err := modelTestDB.Delete(&user).Error
	if err != nil {
		t.Fatalf("Failed to soft delete user: %v", err)
	}

	// Should not find with normal query
	var normalUser TestUser
	err = modelTestDB.First(&normalUser, user.ID).Error
	if err == nil {
		t.Error("Should not find soft deleted user with normal query")
	}

	// Should find with unscoped query
	var softDeletedUser TestUser
	err = modelTestDB.Unscoped().First(&softDeletedUser, user.ID).Error
	if err != nil {
		t.Fatalf("Should find soft deleted user with unscoped query: %v", err)
	}

	if softDeletedUser.DeletedAt.Time.IsZero() {
		t.Error("Soft deleted user should have DeletedAt set")
	}
}
