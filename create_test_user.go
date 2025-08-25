package main

import (
	"fmt"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	config.ConnectDB()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Create test user
	testUser := models.User{
		ID:       uuid.New(),
		Name:     "Calendar Test",
		Email:    "calendar@test.com",
		Password: string(hashedPassword),
	}

	result := config.GetDB().Create(&testUser)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Printf("Created test user: %s / test123\n", testUser.Email)
}
