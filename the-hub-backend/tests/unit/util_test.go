package unit

import (
	"os"
	"testing"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	config.InitLogger()

	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-tokens")
	os.Setenv("GIN_MODE", "test")

	code := m.Run()
	os.Exit(code)
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "thisisaverylongpasswordthatshouldstillworkcorrectly",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := util.HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if hashedPassword == "" {
					t.Error("HashPassword() returned empty hash")
				}

				if hashedPassword == tt.password {
					t.Error("HashPassword() returned unhashed password")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expected       bool
	}{
		{
			name:     "correct password",
			password: "testpassword123",
			expected: true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the original password for testing
			hashedPassword, err := util.HashPassword("testpassword123")
			if err != nil {
				t.Fatalf("Failed to hash password for test setup: %v", err)
			}

			result := util.CheckPasswordHash(hashedPassword, tt.password)

			if tt.name == "correct password" {
				if !result {
					t.Error("CheckPasswordHash() should return true for correct password")
				}
			} else {
				if result {
					t.Error("CheckPasswordHash() should return false for incorrect password")
				}
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	testUUID1 := uuid.New()
	testUUID2 := uuid.New()
	testUUID3 := uuid.New()

	tests := []struct {
		name    string
		userID  uuid.UUID
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  testUUID1,
			wantErr: false,
		},
		{
			name:    "another valid user ID",
			userID:  testUUID2,
			wantErr: false,
		},
		{
			name:    "third valid user ID",
			userID:  testUUID3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := util.GenerateJWT(tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateJWT() returned empty token")
				}

				// Token should be longer than a typical JWT header
				if len(token) < 20 {
					t.Error("GenerateJWT() returned suspiciously short token")
				}
			}
		})
	}
}

func TestPasswordHashConsistency(t *testing.T) {
	password := "testpassword123"

	// Hash the same password multiple times
	hash1, err1 := util.HashPassword(password)
	hash2, err2 := util.HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword() failed: err1=%v, err2=%v", err1, err2)
	}

	// Hashes should be different (due to salt) but both should validate the same password
	if hash1 == hash2 {
		t.Error("HashPassword() should generate different hashes for the same password due to salt")
	}

	// Both hashes should validate against the original password
	if !util.CheckPasswordHash(hash1, password) {
		t.Error("First hash should validate against original password")
	}

	if !util.CheckPasswordHash(hash2, password) {
		t.Error("Second hash should validate against original password")
	}
}
