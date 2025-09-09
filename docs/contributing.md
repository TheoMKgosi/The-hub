# Contributing Guide

Welcome to The Hub! We're excited that you're interested in contributing to our productivity platform. This guide will help you get started with development, understand our processes, and ensure your contributions are high-quality and consistent with our standards.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Code Review Guidelines](#code-review-guidelines)
- [Documentation](#documentation)
- [Issue Reporting](#issue-reporting)
- [Community Guidelines](#community-guidelines)

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.19+** (for backend development)
- **Node.js 18+** or **Bun** (for frontend development)
- **PostgreSQL** (production database) or **SQLite** (development)
- **Git** (version control)
- **Docker** (optional, for containerized development)

### Repository Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/the-hub.git
   cd the-hub
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/original-org/the-hub.git
   ```

4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 🛠 Development Environment

### Backend Setup

```bash
cd the-hub-backend

# Install dependencies
go mod download
go mod tidy

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
# DB_HOST, DB_PASSWORD, JWT_SECRET, etc.

# Run the application
go run main.go
```

### Frontend Setup

```bash
cd the-hub-frontend

# Install dependencies (choose one)
npm install
# or
yarn install
# or
bun install

# Copy environment file
cp .env.example .env

# Run development server
npm run dev
# or
yarn dev
# or
bun run dev
```

### Database Setup

#### PostgreSQL (Recommended for development)
```bash
# Install PostgreSQL
# macOS
brew install postgresql
brew services start postgresql

# Ubuntu
sudo apt install postgresql postgresql-contrib

# Create database
createdb the_hub

# Set up user and permissions
createuser thehub_user
psql -c "ALTER USER thehub_user PASSWORD 'your_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE the_hub TO thehub_user;"
```

#### SQLite (Quick setup)
SQLite is automatically created when you first run the backend application.

### Docker Setup (Alternative)

```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🔄 Development Workflow

### Branch Naming Convention

Use descriptive branch names following this pattern:

```
feature/add-user-authentication
bugfix/fix-login-validation
hotfix/critical-security-patch
refactor/cleanup-api-routes
docs/update-contributing-guide
```

### Commit Message Format

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(auth): add JWT token refresh functionality

fix(api): resolve user creation validation error

docs(readme): update installation instructions

refactor(tasks): simplify task status logic
```

### Development Process

1. **Choose an issue** from GitHub Issues or create one
2. **Create a feature branch** from `main`
3. **Write code** following our standards
4. **Write tests** for your changes
5. **Run tests** and ensure they pass
6. **Update documentation** if needed
7. **Commit changes** with descriptive messages
8. **Push branch** to your fork
9. **Create Pull Request** to main repository

## 📝 Coding Standards

### Backend (Go) Standards

#### Code Style
- Follow standard Go formatting: `go fmt`
- Use `goimports` for import organization
- Maximum line length: 120 characters
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

#### Project Structure
```
internal/
├── handlers/     # HTTP request handlers
├── models/       # Database models
├── config/       # Configuration
├── routes/       # Route definitions
├── util/         # Utilities
└── middleware/   # Custom middleware
```

#### Error Handling
```go
// Good: Explicit error handling
user, err := getUserByID(id)
if err != nil {
    log.Printf("Failed to get user %d: %v", id, err)
    return nil, fmt.Errorf("user not found: %w", err)
}

// Avoid: Ignoring errors
user, _ := getUserByID(id) // Don't do this
```

#### Database Operations
```go
// Good: Use transactions for multiple operations
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### Frontend (Vue.js/TypeScript) Standards

#### Component Structure
```vue
<template>
  <div class="component-name">
    <!-- Template content -->
  </div>
</template>

<script setup lang="ts">
interface Props {
  propName: string
}

const props = withDefaults(defineProps<Props>(), {
  propName: 'default'
})

const emit = defineEmits<{
  eventName: [payload: string]
}>()

// Reactive data
const data = ref('value')

// Computed properties
const computedValue = computed(() => {
  return data.value.toUpperCase()
})

// Methods
const handleEvent = () => {
  emit('eventName', 'payload')
}
</script>

<style scoped>
/* Component styles */
</style>
```

#### Composables
```typescript
// useExample.ts
export const useExample = () => {
  const data = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchData = async () => {
    loading.value = true
    try {
      const response = await $fetch('/api/endpoint')
      data.value = response
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  return {
    data: readonly(data),
    loading: readonly(loading),
    error: readonly(error),
    fetchData
  }
}
```

#### State Management (Pinia)
```typescript
// stores/example.ts
export const useExampleStore = defineStore('example', () => {
  const items = ref([])
  const loading = ref(false)

  const fetchItems = async () => {
    loading.value = true
    try {
      const response = await $fetch('/api/items')
      items.value = response
    } finally {
      loading.value = false
    }
  }

  return {
    items: readonly(items),
    loading: readonly(loading),
    fetchItems
  }
})
```

## 🧪 Testing Guidelines

### Backend Testing

#### Unit Tests
```go
func TestCreateUser(t *testing.T) {
    // Setup
    db := setupTestDB()
    service := NewUserService(db)

    // Test data
    user := &User{Name: "Test User", Email: "test@example.com"}

    // Execute
    err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "Test User", user.Name)
}
```

#### Integration Tests
```go
func TestUserAPI(t *testing.T) {
    router := setupTestRouter()

    // Create request
    req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(`{
        "name": "Test User",
        "email": "test@example.com"
    }`))
    req.Header.Set("Content-Type", "application/json")

    // Execute
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, 201, w.Code)

    var response User
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Test User", response.Name)
}
```

### Frontend Testing

#### Component Tests
```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Button from '../components/Button.vue'

describe('Button', () => {
  it('renders correctly', () => {
    const wrapper = mount(Button, {
      props: { variant: 'primary' }
    })
    expect(wrapper.text()).toContain('Button')
  })

  it('emits click event', async () => {
    const wrapper = mount(Button)
    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })
})
```

#### Composable Tests
```typescript
import { describe, it, expect } from 'vitest'
import { useCounter } from '../composables/useCounter'

describe('useCounter', () => {
  it('increments counter', () => {
    const { count, increment } = useCounter()
    increment()
    expect(count.value).toBe(1)
  })
})
```

### Test Coverage

- **Backend**: Aim for 80%+ code coverage
- **Frontend**: Aim for 70%+ code coverage
- **Critical paths**: 90%+ coverage for authentication, payments, etc.

### Running Tests

```bash
# Backend
go test ./... -v                    # Run all tests
go test -cover ./...               # With coverage
go test -race ./...                # Race condition detection

# Frontend
npm run test                       # Run all tests
npm run test:watch                 # Watch mode
npm run test:coverage              # With coverage
```

## 🔄 Pull Request Process

### Before Submitting

1. **Update your branch** with the latest changes:
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-branch
   git rebase main
   ```

2. **Run all tests** and ensure they pass:
   ```bash
   # Backend
   go test ./... -v

   # Frontend
   npm run test
   ```

3. **Run linting** and fix any issues:
   ```bash
   # Backend
   go fmt ./...
   go vet ./...

   # Frontend
   npm run lint
   npm run typecheck
   ```

4. **Update documentation** if needed

### Creating a Pull Request

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create PR** on GitHub with:
   - Clear title describing the change
   - Detailed description of what was changed and why
   - Reference to any related issues
   - Screenshots for UI changes
   - Test results

3. **PR Template**:
   ```markdown
   ## Description
   Brief description of the changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   - [ ] Unit tests added/updated
   - [ ] Integration tests added/updated
   - [ ] Manual testing completed

   ## Screenshots (if applicable)
   Add screenshots of UI changes

   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Tests pass
   - [ ] Documentation updated
   - [ ] No breaking changes
   ```

## 👀 Code Review Guidelines

### For Reviewers

**Check for:**
- Code quality and adherence to standards
- Test coverage and quality
- Security implications
- Performance considerations
- Documentation updates
- Breaking changes

**Provide feedback:**
- Be constructive and specific
- Explain reasoning for suggestions
- Suggest improvements, don't just point out problems
- Acknowledge good practices

### For Contributors

**When receiving feedback:**
- Don't take criticism personally
- Ask questions if something is unclear
- Address all feedback before requesting re-review
- Update PR description with changes made

### Review Checklist

- [ ] **Functionality**: Does the code work as expected?
- [ ] **Code Quality**: Is the code clean, readable, and well-structured?
- [ ] **Testing**: Are there adequate tests? Do they pass?
- [ ] **Documentation**: Is documentation updated?
- [ ] **Security**: Are there any security concerns?
- [ ] **Performance**: Are there any performance issues?
- [ ] **Compatibility**: Does it work with existing code?

## 📚 Documentation

### When to Update Documentation

- Adding new features or APIs
- Changing existing functionality
- Fixing bugs that affect user experience
- Updating dependencies with breaking changes
- Adding new configuration options

### Documentation Standards

- Use clear, concise language
- Include code examples where helpful
- Keep screenshots up to date
- Test all instructions on a fresh setup
- Use consistent formatting

### API Documentation

```go
// Good: Document exported functions
// CreateUser creates a new user account
// It validates the input data and hashes the password
func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

## 🐛 Issue Reporting

### Bug Reports

**Good bug report includes:**
- Clear title describing the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, browser, versions)
- Screenshots or error messages
- Code snippets if applicable

**Bug Report Template:**
```markdown
## Bug Description
Brief description of the bug

## Steps to Reproduce
1. Go to '...'
2. Click on '...'
3. See error

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., Windows 10]
- Browser: [e.g., Chrome 91]
- Version: [e.g., v1.0.0]

## Additional Context
Any other information
```

### Feature Requests

**Good feature request includes:**
- Clear description of the proposed feature
- Use case or problem it solves
- Proposed implementation (if any)
- Mockups or examples (if UI-related)
- Priority level (nice-to-have vs must-have)

## 🤝 Community Guidelines

### Code of Conduct

- **Be respectful** and inclusive
- **Be constructive** in feedback
- **Help others** when possible
- **Follow project standards**
- **Communicate clearly**

### Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Documentation**: Check existing docs first
- **Community**: Join our community channels

### Recognition

Contributors are recognized through:
- GitHub contributor statistics
- Mention in release notes
- Project maintainer status for active contributors
- Community acknowledgments

## 🎯 Development Best Practices

### General Principles

1. **DRY (Don't Repeat Yourself)**: Avoid code duplication
2. **KISS (Keep It Simple, Stupid)**: Simple solutions are better
3. **SOLID Principles**: Write maintainable, extensible code
4. **YAGNI (You Aren't Gonna Need It)**: Don't over-engineer
5. **TDD (Test-Driven Development)**: Write tests first when possible

### Performance Considerations

- Optimize database queries
- Use appropriate data structures
- Implement caching where beneficial
- Monitor memory usage
- Profile performance bottlenecks

### Security Best Practices

- Validate all user inputs
- Use parameterized queries
- Implement proper authentication
- Keep dependencies updated
- Follow OWASP guidelines

## 📞 Support

If you need help contributing:

1. **Check existing documentation**
2. **Search GitHub Issues** for similar questions
3. **Create a GitHub Discussion** for questions
4. **Join our community** channels
5. **Contact maintainers** for specific guidance

## 🙏 Acknowledgments

Thank you for contributing to The Hub! Your efforts help make this project better for everyone.

---

*This contributing guide is living documentation. Please suggest improvements via pull request.*