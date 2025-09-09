# The Hub

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-4.9+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🌟 Description

The Hub is a comprehensive personal productivity web application designed to automate and streamline your daily life, helping you focus on what truly matters. It provides a centralized platform for managing various aspects of your personal and professional life through an intuitive, modern interface.

Our mission is to create a seamless productivity ecosystem that adapts to your workflow, learns from your habits, and provides intelligent insights to help you achieve your goals more efficiently.

## ✨ Features

### 🎯 Task & Goal Management
- **Smart Task Organization**: Create, categorize, and prioritize tasks with advanced filtering
- **Goal Tracking**: Set long-term goals with progress monitoring and milestone tracking
- **Deadline Management**: Set due dates with intelligent reminders and notifications
- **Progress Analytics**: Visualize your productivity trends and completion rates

### ⏰ Time Management
- **Intelligent Scheduling**: AI-powered time blocking and calendar integration
- **Time Tracking**: Automatic time logging with detailed analytics
- **Focus Sessions**: Pomodoro-style work sessions with distraction blocking
- **Calendar Integration**: Sync with Google Calendar, Outlook, and other calendar services

### 📚 Learning Management
- **Spaced Repetition**: Scientific flashcard system using SM-2 algorithm
- **Course Tracking**: Organize learning resources and track progress
- **Knowledge Base**: Personal wiki for notes and documentation
- **Progress Insights**: Learning analytics and performance metrics

### 💰 Finance Management
- **Expense Tracking**: Categorize and monitor spending patterns
- **Budget Planning**: Set budgets with smart alerts and recommendations
- **Income Management**: Track multiple income sources and projections
- **Financial Goals**: Set savings targets with progress visualization

### 🤖 AI Integration
- **Smart Recommendations**: AI-powered task prioritization and scheduling
- **Productivity Insights**: Analyze patterns and suggest improvements
- **Automated Workflows**: Intelligent automation of repetitive tasks
- **Personal Assistant**: AI chat interface for quick actions and queries

### 🔒 Security & Privacy
- **End-to-End Encryption**: Secure data transmission and storage
- **Privacy-First Design**: Your data belongs to you, not us
- **Two-Factor Authentication**: Enhanced account security
- **Audit Logs**: Complete activity tracking for transparency

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

## 🤝 Contributing

We love your input! We want to make contributing to The Hub as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

### Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

### 📋 Contribution Guidelines

Please read our [Contributing Guide](docs/contributing.md) for detailed information on:

- **Development Setup**: Getting your development environment running
- **Code Standards**: Go and TypeScript/Vue.js coding conventions
- **Testing**: Writing and running tests
- **Documentation**: Keeping docs up to date
- **Pull Request Process**: How to submit your changes

### 🐛 Reporting Bugs

We use GitHub issues to track public bugs. Report a bug by opening a new issue with:

- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Screenshots if applicable

### 💡 Suggesting Features

Feature requests are welcome! Please provide:

- Clear description of the proposed feature
- Use case and problem it solves
- Proposed implementation approach
- Mockups or examples if UI-related

### 📖 Documentation

- [Architecture Overview](docs/architecture.md)
- [API Documentation](docs/api.md)
- [Backend Documentation](the-hub-backend/README.md)
- [Frontend Documentation](the-hub-frontend/README.md)
- [Deployment Guide](docs/deployment.md)
- [Contributing Guide](docs/contributing.md)

### 🏷️ Types of Contributions

- **🐛 Bug Fixes**: Fix existing issues
- **✨ Features**: Add new functionality
- **📚 Documentation**: Improve documentation
- **🎨 UI/UX**: Improve user interface and experience
- **⚡ Performance**: Improve application performance
- **🔒 Security**: Security enhancements
- **🧪 Testing**: Add or improve tests
- **🔧 Tools**: Development tools and scripts

## 🌍 Community

- **GitHub Issues**: [Report bugs and request features](https://github.com/your-org/the-hub/issues)
- **GitHub Discussions**: [Join the conversation](https://github.com/your-org/the-hub/discussions)
- **Discord**: [Join our community server](https://discord.gg/the-hub)
- **Twitter**: [Follow for updates](https://twitter.com/thehub_app)

### 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/your-org/the-hub?style=social)
![GitHub forks](https://img.shields.io/github/forks/your-org/the-hub?style=social)
![GitHub contributors](https://img.shields.io/github/contributors/your-org/the-hub)
![GitHub last commit](https://img.shields.io/github/last-commit/your-org/the-hub)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### 📋 License Summary

- ✅ **Commercial Use**: You can use this project for commercial purposes
- ✅ **Modification**: You can modify the source code
- ✅ **Distribution**: You can distribute the project
- ✅ **Private Use**: You can use privately
- ❌ **Liability**: No liability for damages
- ❌ **Warranty**: No warranty provided

---

## 🙏 Acknowledgments

- **Vue.js Team** for the amazing framework
- **Go Community** for the powerful language and ecosystem
- **Open Source Contributors** for their valuable contributions
- **Our Users** for their feedback and support

---

<div align="center">

**Made with ❤️ by the The Hub team**

[⭐ Star us on GitHub](https://github.com/your-org/the-hub) • [🐛 Report a bug](https://github.com/your-org/the-hub/issues) • [💡 Request a feature](https://github.com/your-org/the-hub/issues)

</div>
