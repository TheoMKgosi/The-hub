# Frontend Documentation

This document provides an overview of The Hub's frontend architecture, components, and pages.

## Technology Stack

- **Framework:** Nuxt.js 3 (Vue.js 3)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **State Management:** Pinia stores and composables
- **API Integration:** Custom API composables with error handling
- **Authentication:** JWT-based with middleware and refresh tokens
- **PWA Support:** Service workers, offline functionality
- **Testing:** Vitest with Vue Test Utils
- **Build Tool:** Vite (via Nuxt.js)
- **Package Manager:** npm/yarn/bun
- **UI Components:** Custom Vue components

## Project Structure

```
the-hub-frontend/
├── app/                    # Nuxt app directory
│   ├── assets/css/         # Global styles
│   ├── components/         # Reusable Vue components
│   │   ├── finance/        # Finance-related components
│   │   ├── learning/       # Learning components
│   │   ├── task/           # Task management components
│   │   ├── ui/             # UI components (Button, NavLink, etc.)
│   │   ├── ConfirmDialog.vue
│   │   ├── Nav.vue
│   │   ├── Tabs.vue
│   │   └── TheCalendar.vue
│   ├── composables/       # Vue composables
│   │   ├── useApi.ts       # API integration
│   │   ├── useDarkMode.ts  # Dark mode management
│   │   ├── useDebounce.ts  # Debounce utility
│   │   └── useToast.ts     # Toast notifications
│   ├── layouts/           # Nuxt layouts
│   │   └── default.vue     # Default layout
│   ├── middleware/        # Route middleware
│   │   └── authenticated.global.ts
│   ├── pages/             # File-based routing
│   │   ├── learning/      # Learning pages
│   │   ├── dashboard.vue  # Main dashboard
│   │   ├── finance.vue    # Finance page
│   │   ├── index.vue      # Home page
│   │   ├── login.vue      # Login page
│   │   ├── register.vue   # Registration page
│   │   ├── settings.vue   # User settings
│   │   └── time.vue       # Time management
│   ├── plugins/           # Nuxt plugins
│   │   └── api.ts         # API plugin
│   └── stores/            # Pinia stores
│       ├── auth.ts        # Authentication store
│       ├── cards.ts       # Flashcard store
│       ├── decks.ts       # Deck store
│       ├── finance.ts     # Finance store
│       ├── goals.ts       # Goals store
│       ├── income.ts      # Income store
│       ├── schedule.ts    # Schedule store
│       ├── tags.ts        # Tags store
│       ├── task-learning.ts
│       ├── tasks.ts       # Tasks store
│       └── topics.ts      # Topics store
├── public/                # Static assets
└── composables/           # Additional composables
```

## Core Components

### UI Components

#### Button Component
```vue
<template>
  <button
    :class="buttonClasses"
    @click="$emit('click')"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false
})
</script>
```

**Usage:**
```vue
<Button variant="primary" size="lg" @click="handleSubmit">
  Submit
</Button>
```

#### NavLink Component
Navigation link component with active state management.

#### Toast Component
Notification toast component for user feedback.

#### Tabs Component
Tab navigation component for organizing content.

### Feature Components

#### Finance Components
- **Budget.vue:** Budget management interface
- **BudgetAlerts.vue:** Budget alert notifications
- **BudgetAnalytics.vue:** Budget analytics and insights
- **BudgetDashboard.vue:** Comprehensive budget dashboard
- **BudgetExport.vue:** Budget data export functionality
- **Category.vue:** Category management for transactions
- **Income.vue:** Income source management
- **Transaction.vue:** Transaction management interface

#### Learning Components
- **Explore.vue:** Learning content exploration
- **Flashcards.vue:** Flashcard review interface
- **Topics.vue:** Learning topic management

#### Task Components
- **Dashboard.vue:** Task dashboard view
- **FormTask.vue:** Task creation and editing form
- **RecentlyDeletedTasks.vue:** Recently deleted tasks recovery
- **StatsDashboard.vue:** Task statistics and analytics
- **TaskFilters.vue:** Advanced task filtering options
- **TheTasks.vue:** Task list and management
- **TimeTracker.vue:** Time tracking for tasks

#### Analytics Components
- **AnalyticsDashboard.vue:** Comprehensive analytics dashboard

#### Collaboration Components
- **TaskCollaboration.vue:** Task collaboration features

#### Calendar Components
- **CalendarIntegrations.vue:** External calendar integrations
- **CalendarZonesManager.vue:** Calendar zone management
- **TheCalendar.vue:** Calendar interface

#### Goal Components
- **TheGoals.vue:** Goal management interface
- **GoalTasks.vue:** Goal-task relationship management
- **FormGoal.vue:** Goal creation and editing form

#### Utility Components
- **CommandPalette.vue:** Command palette for quick actions
- **ConfirmDialog.vue:** Confirmation dialog component
- **EditTaskForm.vue:** Advanced task editing form
- **ErrorBoundary.vue:** Error boundary for error handling
- **AddTaskToGoal.vue:** Add task to goal functionality

## Composables

### useApi
Handles all API communication with the backend.

```typescript
const { data, error, loading, execute } = useApi('/tasks')

// Make API call
await execute()

// With parameters
await execute({
  method: 'POST',
  body: { title: 'New Task' }
})
```

### useDarkMode
Manages dark mode theme switching.

```typescript
const { isDark, toggleDarkMode } = useDarkMode()
```

### useToast
Provides toast notification functionality.

```typescript
const toast = useToast()

toast.success('Task created successfully')
toast.error('Failed to create task')
```

### useDebounce
Debounces function calls for performance optimization.

```typescript
const debouncedSearch = useDebounce((query: string) => {
  // Search logic here
}, 300)
```

## Pages

### Authentication Pages

#### Login Page (`/login`)
User authentication form with email/password fields.

**Features:**
- Form validation
- Error handling
- Redirect after successful login

#### Register Page (`/register`)
User registration form.

**Features:**
- Form validation
- Password confirmation
- Automatic login after registration

### Main Application Pages

#### Dashboard (`/dashboard`)
Main dashboard showing overview of tasks, goals, and recent activity.

**Features:**
- Task summary
- Goal progress
- Quick actions
- Recent activity feed

#### Tasks (`/tasks`)
Task management interface.

**Features:**
- Task list with filtering and sorting
- Task creation and editing
- Status updates
- Due date management

#### Goals (`/goals`)
Goal tracking and management.

**Features:**
- Goal creation and editing
- Progress tracking
- Category organization
- Achievement milestones

#### Finance (`/finance`)
Financial management dashboard.

**Features:**
- Transaction tracking
- Budget management
- Income monitoring
- Financial reports

#### Learning (`/learning`)
Learning management system with advanced features.

**Features:**
- Flashcard deck management with spaced repetition
- Learning analytics and progress tracking
- Topic organization and exploration
- Learning paths and structured courses
- Review sessions with performance metrics
- Browse and discover learning content

**Sub-pages:**
- `/learning/browse/[deck_id]`: Browse specific deck
- `/learning/cards/[deck_id]`: Manage cards in deck
- `/learning/review/[deck_id]`: Review flashcards
- `/learning/[id]`: Learning topic details
- `/learning/analytics`: Learning analytics dashboard
- `/learning/paths`: Learning path management

#### Time (`/time`)
Time management and scheduling.

**Features:**
- Calendar view
- Time blocking
- Schedule management
- Time tracking

#### Settings (`/settings`)
User settings and preferences.

**Features:**
- Profile management
- Theme preferences
- Notification settings
- Account settings

#### Statistics (`/stats`)
Comprehensive statistics and analytics dashboard.

**Features:**
- Task completion statistics
- Goal progress analytics
- Learning progress metrics
- Financial analytics
- Time tracking reports

#### Plan (`/plan`)
Planning and goal-setting interface.

**Features:**
- Long-term planning
- Goal setting and tracking
- Progress visualization
- Milestone management

#### Password Management
- **Forgot Password (`/forgot-password`)**: Password recovery
- **Reset Password (`/reset-password`)**: Password reset functionality

## State Management

The application uses Pinia stores for state management:

### Auth Store
Manages user authentication state and JWT tokens.

```typescript
const auth = useAuthStore()

// Check if user is authenticated
const isAuthenticated = auth.isAuthenticated

// Get current user
const user = auth.user

// Login/logout methods
await auth.login(credentials)
auth.logout()
```

### Task Store
Manages task-related state and operations.

```typescript
const tasks = useTasksStore()

// Get tasks
const userTasks = tasks.tasks

// Create new task
await tasks.createTask(taskData)

// Update task
await tasks.updateTask(id, updates)
```

## Styling

The application uses Tailwind CSS for styling with a custom design system:

### Color Scheme
- **Primary:** Blue tones for main actions
- **Secondary:** Gray tones for secondary elements
- **Success:** Green for positive actions
- **Warning:** Yellow for warnings
- **Error:** Red for errors

### Dark Mode
Full dark mode support with CSS custom properties and Tailwind's dark mode utilities.

### Responsive Design
Mobile-first responsive design using Tailwind's responsive utilities.

## Development Guidelines

### Component Structure
- Use Vue 3 Composition API with `<script setup>`
- Type all props and emits with TypeScript interfaces
- Follow Vue.js naming conventions
- Use slots for flexible content

### Code Style
- Use TypeScript for all new code
- Follow Vue.js style guide
- Use ESLint and Prettier for code formatting
- Write comprehensive JSDoc comments

### Testing
- Unit tests for composables and utilities
- Component tests for UI components
- Integration tests for critical user flows

### Performance
- Use lazy loading for routes
- Implement proper loading states
- Optimize images and assets
- Use Vue's keep-alive for frequently accessed components

## Build and Deployment

### Development
```bash
bun run dev  # Start development server
```

### Build
```bash
bun run build  # Build for production
```

### Preview
```bash
bun run preview  # Preview production build
```

### Deployment
The application can be deployed to various platforms:
- Vercel
- Netlify
- Docker containers
- Static hosting with API backend

## Environment Variables

```bash
# API Configuration
NUXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Authentication
NUXT_PUBLIC_JWT_SECRET=your_jwt_secret

# Development
NUXT_PUBLIC_DEV=true
```