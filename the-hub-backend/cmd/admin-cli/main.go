package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	var (
		email    = flag.String("email", "", "Email address for the admin user")
		name     = flag.String("name", "", "Name for the admin user")
		password = flag.String("password", "", "Password for the admin user")
		help     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Initialize database connection
	if err := config.InitDBManager("postgres"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := config.GetDB()

	// Get user input if not provided via flags
	if *email == "" || *name == "" || *password == "" {
		fmt.Println("Admin User Creation Tool")
		fmt.Println("========================")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)

		if *email == "" {
			*email = prompt(reader, "Enter email address: ")
		}
		if *name == "" {
			*name = prompt(reader, "Enter name: ")
		}
		if *password == "" {
			*password = promptPassword(reader)
		}
	}

	// Validate input
	if *email == "" || *name == "" || *password == "" {
		fmt.Println("Error: Email, name, and password are required")
		os.Exit(1)
	}

	// Create admin user
	userID, err := createAdminUser(db, *email, *name, *password)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("✅ Admin user created successfully!\n")
	fmt.Printf("   ID: %s\n", userID)
	fmt.Printf("   Email: %s\n", *email)
	fmt.Printf("   Name: %s\n", *name)
	fmt.Printf("   Role: admin\n")
}

func showHelp() {
	fmt.Println("Admin User Creation Tool")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("This tool creates admin users for The Hub application.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/admin-cli/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -email string")
	fmt.Println("        Email address for the admin user")
	fmt.Println("  -name string")
	fmt.Println("        Name for the admin user")
	fmt.Println("  -password string")
	fmt.Println("        Password for the admin user")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Interactive mode (will prompt for input)")
	fmt.Println("  go run cmd/admin-cli/main.go")
	fmt.Println()
	fmt.Println("  # Command line mode")
	fmt.Println("  go run cmd/admin-cli/main.go -email admin@example.com -name \"Admin User\" -password securepassword")
	fmt.Println()
	fmt.Println("  # Show help")
	fmt.Println("  go run cmd/admin-cli/main.go -help")
}

func prompt(reader *bufio.Reader, text string) string {
	fmt.Print(text)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptPassword(reader *bufio.Reader) string {
	fmt.Print("Enter password: ")
	// Note: In a real production environment, you might want to use a library
	// that hides password input like golang.org/x/term
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password)
}

func createAdminUser(db *gorm.DB, email, name, password string) (string, error) {
	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return "", fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	config.Logger.Infof("Admin user created: ID=%s, Email=%s", user.ID, user.Email)
	return user.ID.String(), nil
}
