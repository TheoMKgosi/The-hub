package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		command = flag.String("command", "", "Command to run (users, goals, all, clean)")
		help    = flag.Bool("help", false, "Show help message")
		force   = flag.Bool("force", false, "Force operation without confirmation")
	)
	flag.Parse()

	if *help || *command == "" {
		showHelp()
		return
	}

	// Initialize database connection
	if err := config.InitDBManager(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := config.GetDB()

	switch *command {
	case "users":
		if err := seedUsers(db, *force); err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		fmt.Println("✅ Users seeded successfully!")

	case "goals":
		if err := seedGoals(db, *force); err != nil {
			log.Fatalf("Failed to seed goals: %v", err)
		}
		fmt.Println("✅ Goals and tasks seeded successfully!")

	case "all":
		if err := seedAll(db, *force); err != nil {
			log.Fatalf("Failed to seed all data: %v", err)
		}
		fmt.Println("✅ All data seeded successfully!")

	case "clean":
		if !*force {
			fmt.Print("This will remove all seeded data. Are you sure? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled.")
				return
			}
		}
		if err := cleanSeededData(db); err != nil {
			log.Fatalf("Failed to clean seeded data: %v", err)
		}
		fmt.Println("✅ Seeded data cleaned successfully!")

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Database Seeder Tool")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("This tool seeds the database with sample data for development and testing.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seeder/main.go -command <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  users  - Seed sample users (admin and regular users)")
	fmt.Println("  goals  - Seed sample goals and associated tasks")
	fmt.Println("  all    - Seed all data (users, goals, tasks)")
	fmt.Println("  clean  - Remove all seeded data")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -force")
	fmt.Println("        Skip confirmation prompts")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Seed all data")
	fmt.Println("  go run cmd/seeder/main.go -command all")
	fmt.Println()
	fmt.Println("  # Seed users only")
	fmt.Println("  go run cmd/seeder/main.go -command users")
	fmt.Println()
	fmt.Println("  # Clean seeded data")
	fmt.Println("  go run cmd/seeder/main.go -command clean")
	fmt.Println()
	fmt.Println("  # Force clean without confirmation")
	fmt.Println("  go run cmd/seeder/main.go -command clean -force")
}

// seedUsers creates sample users for development
func seedUsers(db *gorm.DB, force bool) error {
	fmt.Println("Seeding users...")

	// Check if seeded users already exist
	var seededCount int64
	if err := db.Model(&models.User{}).Where("email LIKE ?", "%@thehub.com").Or("email LIKE ?", "%.example.com").Count(&seededCount).Error; err != nil {
		return fmt.Errorf("failed to check existing seeded users: %w", err)
	}

	if seededCount > 0 && !force {
		fmt.Printf("Found %d seeded users. Use -force to re-seed.\n", seededCount)
		return nil
	}

	// If force is true, clean existing seeded users first
	if force && seededCount > 0 {
		if err := db.Where("email LIKE ?", "%@thehub.com").Or("email LIKE ?", "%.example.com").Delete(&models.User{}).Error; err != nil {
			return fmt.Errorf("failed to clean existing seeded users: %w", err)
		}
		fmt.Printf("Cleaned %d existing seeded users\n", seededCount)
	}

	users := []models.User{
		{
			Name:  "Admin User",
			Email: "admin@thehub.com",
			Role:  "admin",
		},
		{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Role:  "user",
		},
		{
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
			Role:  "user",
		},
		{
			Name:  "Bob Johnson",
			Email: "bob.johnson@example.com",
			Role:  "user",
		},
	}

	for _, user := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", user.Email, err)
		}
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
		fmt.Printf("Created user: %s (%s)\n", user.Name, user.Email)
	}

	return nil
}

// seedGoals creates sample goals and tasks
func seedGoals(db *gorm.DB, force bool) error {
	fmt.Println("Seeding goals and tasks...")

	// Get existing users (prefer seeded users)
	var users []models.User
	if err := db.Where("email LIKE ?", "%@thehub.com").Or("email LIKE ?", "%.example.com").Find(&users).Error; err != nil {
		return fmt.Errorf("failed to get seeded users: %w", err)
	}

	// If no seeded users, get all users
	if len(users) == 0 {
		if err := db.Find(&users).Error; err != nil {
			return fmt.Errorf("failed to get users: %w", err)
		}
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found. Run 'users' command first")
	}

	// Check if seeded goals already exist
	var seededGoalsCount int64
	if err := db.Model(&models.Goal{}).Where("title LIKE ?", "Complete Project%").Or("title LIKE ?", "Learn Go%").Or("title LIKE ?", "Home Organization%").Count(&seededGoalsCount).Error; err != nil {
		return fmt.Errorf("failed to check existing seeded goals: %w", err)
	}

	if seededGoalsCount > 0 && !force {
		fmt.Printf("Found %d seeded goals. Use -force to re-seed.\n", seededGoalsCount)
		return nil
	}

	// If force is true, clean existing seeded goals and tasks first
	if force && seededGoalsCount > 0 {
		// Get seeded goals
		var seededGoals []models.Goal
		if err := db.Where("title LIKE ?", "Complete Project%").Or("title LIKE ?", "Learn Go%").Or("title LIKE ?", "Home Organization%").Find(&seededGoals).Error; err != nil {
			return fmt.Errorf("failed to find seeded goals: %w", err)
		}

		// Delete associated tasks first
		for _, goal := range seededGoals {
			if err := db.Where("goal_id = ?", goal.ID).Delete(&models.Task{}).Error; err != nil {
				return fmt.Errorf("failed to delete tasks for goal %s: %w", goal.Title, err)
			}
		}

		// Delete goals
		if err := db.Where("title LIKE ?", "Complete Project%").Or("title LIKE ?", "Learn Go%").Or("title LIKE ?", "Home Organization%").Delete(&models.Goal{}).Error; err != nil {
			return fmt.Errorf("failed to clean existing seeded goals: %w", err)
		}
		fmt.Printf("Cleaned %d existing seeded goals and their tasks\n", seededGoalsCount)
	}

	// Sample goals data
	goalsData := []struct {
		title       string
		description string
		priority    int
		category    string
		tasks       []string
	}{
		{
			title:       "Complete Project Documentation",
			description: "Write comprehensive documentation for the new API endpoints",
			priority:    4,
			category:    "Work",
			tasks: []string{
				"Write API endpoint documentation",
				"Create user guide",
				"Add code examples",
				"Review and publish docs",
			},
		},
		{
			title:       "Learn Go Advanced Patterns",
			description: "Master advanced Go programming concepts and patterns",
			priority:    3,
			category:    "Learning",
			tasks: []string{
				"Study concurrency patterns",
				"Learn about generics",
				"Practice with interfaces",
				"Build a sample project",
			},
		},
		{
			title:       "Home Organization Project",
			description: "Organize and declutter the home office and living spaces",
			priority:    2,
			category:    "Personal",
			tasks: []string{
				"Clean out desk drawers",
				"Organize digital files",
				"Sort through old documents",
				"Set up new filing system",
			},
		},
	}

	for i, goalData := range goalsData {
		user := users[i%len(users)] // Cycle through users

		// Create goal
		goal := models.Goal{
			UserID:      user.ID,
			Title:       goalData.title,
			Description: goalData.description,
			Priority:    &goalData.priority,
			Category:    goalData.category,
			Status:      "active",
		}

		if err := db.Create(&goal).Error; err != nil {
			return fmt.Errorf("failed to create goal %s: %w", goalData.title, err)
		}

		fmt.Printf("Created goal: %s for user %s\n", goal.Title, user.Name)

		// Create tasks for this goal
		for j, taskTitle := range goalData.tasks {
			priority := (j % 5) + 1 // Vary priorities 1-5
			status := "pending"
			if j == 0 {
				status = "in_progress" // Make first task in progress
			} else if j == len(goalData.tasks)-1 {
				status = "completed" // Make last task completed
			}

			task := models.Task{
				UserID:     user.ID,
				GoalID:     &goal.ID,
				Title:      taskTitle,
				Status:     status,
				Priority:   &priority,
				OrderIndex: j,
			}

			if err := db.Create(&task).Error; err != nil {
				return fmt.Errorf("failed to create task %s: %w", taskTitle, err)
			}
		}

		fmt.Printf("Created %d tasks for goal: %s\n", len(goalData.tasks), goal.Title)
	}

	return nil
}

// seedAll runs all seeding operations
func seedAll(db *gorm.DB, force bool) error {
	fmt.Println("Seeding all data...")

	if err := seedUsers(db, force); err != nil {
		return err
	}

	if err := seedGoals(db, force); err != nil {
		return err
	}

	return nil
}

// cleanSeededData removes all seeded data
func cleanSeededData(db *gorm.DB) error {
	fmt.Println("Cleaning seeded data...")

	// Clean tasks associated with seeded goals first
	seededGoalTitles := []string{"Complete Project%", "Learn Go%", "Home Organization%"}
	for _, titlePattern := range seededGoalTitles {
		var goals []models.Goal
		if err := db.Where("title LIKE ?", titlePattern).Find(&goals).Error; err != nil {
			return fmt.Errorf("failed to find goals with pattern %s: %w", titlePattern, err)
		}

		for _, goal := range goals {
			if err := db.Where("goal_id = ?", goal.ID).Delete(&models.Task{}).Error; err != nil {
				return fmt.Errorf("failed to delete tasks for goal %s: %w", goal.Title, err)
			}
		}
	}

	// Clean seeded goals
	var goalsDeleted int64
	for _, titlePattern := range seededGoalTitles {
		result := db.Where("title LIKE ?", titlePattern).Delete(&models.Goal{})
		if result.Error != nil {
			return fmt.Errorf("failed to clean goals with pattern %s: %w", titlePattern, result.Error)
		}
		goalsDeleted += result.RowsAffected
	}
	fmt.Printf("Cleaned %d seeded goals\n", goalsDeleted)

	// Clean seeded users
	var usersDeleted int64
	result := db.Where("email LIKE ?", "%@thehub.com").Or("email LIKE ?", "%.example.com").Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean seeded users: %w", result.Error)
	}
	usersDeleted = result.RowsAffected
	fmt.Printf("Cleaned %d seeded users\n", usersDeleted)

	return nil
}
