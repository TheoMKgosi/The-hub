# The Hub

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

The Hub is a comprehensive personal productivity platform designed to streamline your daily life through intelligent task management, learning systems, financial tracking, and time optimization. Built with modern web technologies, it provides a seamless experience across devices with offline capabilities and real-time synchronization.

Our mission is to create an intelligent productivity ecosystem that adapts to your workflow, learns from your habits, and provides actionable insights to help you achieve your goals more efficiently.

## Key Features

### Task & Goal Management
- **Task Organization**: Create, categorize, and prioritize tasks with advanced filtering and search
- **Goal Tracking**: Set long-term goals with progress monitoring, milestones, and achievement tracking
- **Time Tracking**: Built-in time tracking with detailed analytics and productivity insights

### Intelligent Time Management
- **Calendar Integration**: Sync with Google Calendar, Outlook, and other calendar services

### Advanced Learning Management
- **Spaced Repetition System**: Scientific flashcard system using SM-2 algorithm for optimal learning
- **Knowledge Base**: Personal wiki for notes, documentation, and reference materials
- **Study Sessions**: Timed study sessions with progress tracking and performance metrics

### Comprehensive Finance Management
- **Expense Tracking**: Categorize and monitor spending patterns with smart categorization
- **Budget Planning**: Create budgets with intelligent alerts and spending predictions
- **Income Management**: Track multiple income sources with projections and analytics
- **Financial Goals**: Set savings targets with progress visualization and milestone tracking
- **Transaction Analytics**: Detailed spending analysis with trends and insights

### Security & Privacy
- **End-to-End Encryption**: Secure data transmission and storage
- **Privacy-First Design**: Your data belongs to you, with transparent data handling
- **Data Export**: Full data export capabilities for data portability

## Architecture

The Hub is built with a modern full-stack architecture designed for scalability, maintainability, and performance:

### Backend Architecture
- **Language**: Go 1.24+ with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT-based authentication with refresh tokens
- **API**: RESTful API with comprehensive Swagger documentation

### Frontend Architecture
- **Framework**: Nuxt.js 3 (Vue.js 3) with TypeScript
- **State Management**: Pinia stores with persistence
- **Styling**: Tailwind CSS with custom design system
- **PWA**: Progressive Web App with offline capabilities
- **Build**: Vite for fast development and optimized production builds

### Infrastructure
- **CI/CD**: GitHub Actions with automated deployment
- **Monitoring**: Application monitoring and error tracking
- **Backup**: Automated database backups and recovery

## Quick Start

### Prerequisites
- Go 1.24+ (backend)
- Node.js 18+ or Bun (frontend)
- PostgreSQL 17+ (production) 

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/TheoMKgosi/The-hub.git
   cd The-hub
   ```

2. **Backend Setup:**
   ```bash
   cd the-hub-backend
   go mod download
   cp .env.example .env
   # Configure your environment variables
   go run main.go
   ```

3. **Frontend Setup:**
   ```bash
   cd the-hub-frontend
   bun install  # or npm install
   bun run dev  # or npm run dev
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html



## Documentation

- **[API Documentation](docs/api.md)** - Complete REST API reference
- **[Architecture Overview](docs/architecture.md)** - System design and architecture
- **[Backend Documentation](the-hub-backend/README.md)** - Backend setup and development
- **[Frontend Documentation](the-hub-frontend/README.md)** - Frontend setup and development
- **[Deployment Guide](docs/deployment.md)** - Production deployment instructions
- **[User Settings API](the-hub-backend/docs/user-settings-api.md)** - User settings management

## Development

### Testing
```bash
# Backend tests
cd the-hub-backend && go test ./...

# Frontend tests
cd the-hub-frontend && bun run test
```

### Code Quality
```bash
# Backend linting
cd the-hub-backend && go vet ./...

# Frontend linting
cd the-hub-frontend && bun run lint
```

### Building
```bash
# Backend build
cd the-hub-backend && go build

# Frontend build
cd the-hub-frontend && bun run build
```

## Contributing

No contribution will be taken

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### License Summary
- ✅ **Commercial Use**: You can use this project for commercial purposes
- ✅ **Modification**: You can modify the source code
- ✅ **Distribution**: You can distribute the project
- ✅ **Private Use**: You can use privately
- ❌ **Liability**: No liability for damages
- ❌ **Warranty**: No warranty provided

## Acknowledgments

- **Vue.js Team** for the amazing framework ecosystem
- **Go Community** for the powerful language and ecosystem
- **Our Users** for their feedback and support
