# The Hub

## Description
The Hub is a personal productivity web app designed to automate your life so you can focus on the important stuff. It provides a centralized platform for managing various aspects of your life, including tasks, goals, time, learning, and finances.

## Features
- **Task and Goal Management:** Organize your tasks and goals in one place, set deadlines, and track your progress.
- **Time Management:** Schedule your time effectively, set reminders, and track your time spent on different activities.
- **Learning Management:** Manage your learning resources, track your progress, and discover new learning opportunities with flashcard decks and spaced repetition.
- **Finance Management:** Track your income and expenses, set budgets, and manage your financial goals with category-based budgeting.
- **AI Integration:** Leverage AI to automate tasks, get personalized recommendations, and gain insights into your productivity.

## Architecture

The Hub is built with a modern full-stack architecture:

- **Frontend:** Nuxt.js (Vue.js) with TypeScript, using composables for state management and API integration
- **Backend:** Go with Gin framework, providing RESTful APIs
- **Database:** PostgreSQL with GORM ORM (SQLite for development)
- **Authentication:** JWT-based authentication system
- **Documentation:** Swagger/OpenAPI for API documentation

## Project Structure

```
the-hub/
├── the-hub-frontend/     # Nuxt.js frontend application
├── the-hub-backend/      # Go backend API server
├── docs/                 # Project documentation
└── tools/                # Utility scripts and tools
```

## Getting Started

### Prerequisites

- Go 1.19+ (for backend)
- Node.js 16+ or Bun (for frontend)
- PostgreSQL (production) or SQLite (development)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   cd the-hub
   ```

2. **Backend Setup:**
   ```bash
   cd the-hub-backend
   go mod download
   go mod tidy
   cp .env.example .env  # Configure your environment variables
   ```

3. **Frontend Setup:**
   ```bash
   cd ../the-hub-frontend
   bun install  # or npm install or yarn install
   ```

### Configuration

Copy `.env.example` to `.env` and configure the following variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=the_hub

# JWT
JWT_SECRET=your_jwt_secret

# Server
PORT=8080
```

### Running the Application

1. **Start the backend:**
   ```bash
   cd the-hub-backend
   go run main.go
   ```

2. **Start the frontend (in a new terminal):**
   ```bash
   cd the-hub-frontend
   bun run dev  # or npm run dev
   ```

3. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html

### Development

- **Backend Tests:** `cd the-hub-backend && go test ./...`
- **Frontend Tests:** `cd the-hub-frontend && npm run test`
- **Linting:** Follow the guidelines in `AGENTS.md`

## API Documentation

The backend provides a comprehensive REST API documented with Swagger/OpenAPI:

- **Swagger UI:** Available at `/swagger/index.html` when running the backend
- **API Docs:** See `the-hub-backend/docs/` for detailed endpoint documentation
- **User Settings API:** Detailed documentation in `the-hub-backend/docs/user-settings-api.md`

Key API endpoints include:
- `/users` - User management and settings
- `/tasks` - Task management
- `/goals` - Goal tracking
- `/finance` - Financial management
- `/learning` - Learning resources and flashcards

## Contribution

We welcome contributions to The Hub! To contribute, please follow these steps:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes.
4.  Submit a pull request.

## License

[License] - Add License details here
