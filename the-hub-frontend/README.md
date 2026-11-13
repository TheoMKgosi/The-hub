# The Hub Frontend

A modern Vue.js frontend for The Hub productivity application, built with Nuxt.js 3, TypeScript, and Tailwind CSS.

## 🚀 Features

- **Modern Vue.js**: Built with Vue 3 and Nuxt.js 3 framework
- **TypeScript**: Full TypeScript support for better development experience
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **State Management**: Pinia stores for centralized state management
- **Component Architecture**: Reusable Vue components organized by feature
- **Authentication**: JWT-based authentication with route protection
- **Real-time Updates**: Live data synchronization with backend
- **Progressive Web App**: PWA-ready with offline capabilities
- **Accessibility**: WCAG compliant components and interfaces

## 🏗 Technology Stack

- **Framework**: Nuxt.js 3 (Vue.js 3)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Pinia
- **HTTP Client**: Custom composables with fetch API
- **Routing**: Nuxt file-based routing
- **Build Tool**: Vite (via Nuxt)
- **Testing**: Vitest + Vue Test Utils
- **Linting**: ESLint + Prettier

## 📁 Project Structure

```
the-hub-frontend/
├── app/                    # Nuxt app directory
│   ├── assets/css/         # Global styles and Tailwind config
│   ├── components/         # Reusable Vue components
│   │   ├── ui/             # UI components (Button, NavLink, etc.)
│   │   ├── finance/        # Finance-related components
│   │   ├── learning/       # Learning components
│   │   ├── task/           # Task management components
│   │   └── shared/         # Shared components
│   ├── composables/        # Vue composables for logic reuse
│   │   ├── useApi.ts       # API communication
│   │   ├── useAuth.ts      # Authentication logic
│   │   ├── useToast.ts     # Notification system
│   │   └── useValidation.ts # Form validation
│   ├── layouts/            # Page layouts
│   │   └── default.vue     # Main application layout
│   ├── middleware/         # Route middleware
│   │   └── authenticated.global.ts # Auth protection
│   ├── pages/              # File-based routing
│   │   ├── index.vue       # Home page
│   │   ├── login.vue       # Authentication
│   │   ├── dashboard.vue   # Main dashboard
│   │   ├── finance.vue     # Finance management
│   │   ├── learning/       # Learning pages
│   │   └── settings.vue    # User settings
│   ├── plugins/            # Nuxt plugins
│   │   └── api.ts          # API plugin setup
│   ├── stores/             # Pinia stores
│   │   ├── auth.ts         # Authentication store
│   │   ├── tasks.ts        # Task management store
│   │   ├── goals.ts        # Goal tracking store
│   │   ├── finance.ts      # Finance store
│   │   └── learning.ts     # Learning store
│   ├── types/              # TypeScript type definitions
│   └── utils/              # Utility functions
├── public/                 # Static assets
├── test/                   # Test files
│   ├── components/         # Component tests
│   ├── composables/        # Composable tests
│   └── setup.ts            # Test configuration
├── nuxt.config.ts          # Nuxt configuration
├── tailwind.config.js      # Tailwind CSS configuration
├── tsconfig.json           # TypeScript configuration
├── vitest.config.ts        # Vitest configuration
└── package.json            # Dependencies and scripts
```

## 📋 Prerequisites

- Node.js 16+ or Bun
- npm, yarn, or bun package manager
- Backend API server running (see backend README)

## 🛠 Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository_url>
cd the-hub-frontend
```

### 2. Install Dependencies

#### Using npm
```bash
npm install
```

#### Using yarn
```bash
yarn install
```

#### Using bun (recommended)
```bash
bun install
```

### 3. Environment Configuration
Copy the example environment file:
```bash
cp .env.example .env
```

Configure your environment variables:
```bash
# API Configuration
NUXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1

# Authentication
NUXT_PUBLIC_JWT_SECRET=your_jwt_secret

# Application
NUXT_PUBLIC_APP_NAME="The Hub"
NUXT_PUBLIC_APP_VERSION="1.0.0"

# Development
NUXT_PUBLIC_DEV=true
```

### 4. Development Server

#### Using npm
```bash
npm run dev
```

#### Using yarn
```bash
yarn dev
```

#### Using bun
```bash
bun run dev
```

The application will be available at `http://localhost:3000`

## 🧪 Testing

### Run All Tests
```bash
# npm
npm run test

# yarn
yarn test

# bun
bun run test
```

### Run Tests in Watch Mode
```bash
# npm
npm run test:watch

# yarn
yarn test:watch

# bun
bun run test:watch
```

### Run Tests with UI
```bash
# npm
npm run test:ui

# yarn
yarn test:ui

# bun
bun run test:ui
```

### Run Specific Test
```bash
# npm
npm run test Button.test.ts

# yarn
yarn test Button.test.ts

# bun
bun run test Button.test.ts
```

## 📦 Build & Deployment

### Development Build
```bash
# npm
npm run build

# yarn
yarn build

# bun
bun run build
```

### Production Build
```bash
# npm
NODE_ENV=production npm run build

# yarn
NODE_ENV=production yarn build

# bun
NODE_ENV=production bun run build
```

### Preview Production Build
```bash
# npm
npm run preview

# yarn
yarn preview

# bun
bun run preview
```

### Static Site Generation (SSG)
```bash
# npm
npm run generate

# yarn
yarn generate

# bun
bun run generate
```

## 🎨 Styling & Design

### Tailwind CSS
The application uses Tailwind CSS for styling with a custom design system:

- **Color Palette**: Custom color scheme with dark mode support
- **Typography**: Consistent font scales and spacing
- **Components**: Pre-built component classes
- **Responsive**: Mobile-first responsive design

### Dark Mode
Built-in dark mode support with:
- System preference detection
- Manual toggle option
- Persistent user preference
- Smooth transitions

## 🔧 Development Guidelines

### Component Structure
```vue
<template>
  <div class="component-name">
    <!-- Template content -->
  </div>
</template>

<script setup lang="ts">
// TypeScript interfaces
interface Props {
  propName: string
}

// Props with defaults
const props = withDefaults(defineProps<Props>(), {
  propName: 'default'
})

// Emits
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

### Composables
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

### State Management (Pinia)
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

  const addItem = (item) => {
    items.value.push(item)
  }

  return {
    items: readonly(items),
    loading: readonly(loading),
    fetchItems,
    addItem
  }
})
```

## 🚀 Deployment Options

### Vercel (Recommended)
1. Connect your repository to Vercel
2. Configure build settings:
   - **Build Command**: `npm run build`
   - **Output Directory**: `.output/public`
   - **Install Command**: `npm install`
3. Add environment variables
4. Deploy

### Netlify
1. Connect your repository to Netlify
2. Configure build settings:
   - **Build Command**: `npm run build`
   - **Publish Directory**: `.output/public`
3. Add environment variables
4. Deploy

### Server-Side Rendering
```bash
# Build for SSR
npm run build

# Start production server
npm run start
```

### Static Hosting
For static hosting platforms:
```bash
npm run generate
```
Deploy the generated `dist/` folder to your hosting provider.

## 🔧 Configuration

### Nuxt Configuration (`nuxt.config.ts`)
```typescript
export default defineNuxtConfig({
  // Modules
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],

  // Runtime config
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL
    }
  },

  // CSS
  css: ['~/assets/css/main.css'],

  // Build
  build: {
    transpile: ['@headlessui/vue']
  },

  // TypeScript
  typescript: {
    strict: true
  }
})
```

### Environment Variables
```bash
# API
NUXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1

# Authentication
NUXT_PUBLIC_JWT_SECRET=your_secret

# App
NUXT_PUBLIC_APP_NAME=The Hub
NUXT_PUBLIC_APP_VERSION=1.0.0

# Features
NUXT_PUBLIC_ENABLE_ANALYTICS=false
```

## 🧪 Testing Strategy

### Unit Tests
- Component logic and interactions
- Composables functionality
- Store actions and getters
- Utility functions

### Integration Tests
- Page navigation and routing
- Form submissions and validation
- API integration
- Cross-component interactions

### E2E Tests (Future)
- User workflows
- Critical user journeys
- Performance testing

## 📱 Progressive Web App (PWA)

### Features
- Offline capability
- Installable on mobile devices
- Push notifications (planned)
- Background sync (planned)

### Configuration
PWA features are configured in `nuxt.config.ts`:
```typescript
export default defineNuxtConfig({
  modules: [
    '@nuxtjs/pwa'
  ],

  pwa: {
    manifest: {
      name: 'The Hub',
      short_name: 'TheHub',
      description: 'Personal productivity platform'
    }
  }
})
```

## 🔍 Performance Optimization

### Code Splitting
- Automatic route-based code splitting
- Dynamic imports for heavy components
- Lazy loading of non-critical features

### Bundle Analysis
```bash
# Analyze bundle size
npm run build --analyze

# View bundle analyzer
npx nuxi build --analyze
```

### Image Optimization
- Automatic image optimization
- WebP format support
- Responsive images
- Lazy loading

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`npm run test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style
- Use TypeScript for all new code
- Follow Vue.js composition API patterns
- Use ESLint and Prettier for code formatting
- Write comprehensive JSDoc comments
- Follow component naming conventions

## 📝 Scripts

```json
{
  "scripts": {
    "dev": "nuxt dev",
    "build": "nuxt build",
    "start": "nuxt start",
    "generate": "nuxt generate",
    "preview": "nuxt preview",
    "test": "vitest",
    "test:watch": "vitest --watch",
    "test:ui": "vitest --ui",
    "lint": "eslint .",
    "lint:fix": "eslint . --fix",
    "typecheck": "nuxt typecheck"
  }
}
```

## 🐛 Troubleshooting

### Common Issues

**Build Errors**
- Clear node_modules: `rm -rf node_modules && npm install`
- Clear Nuxt cache: `rm -rf .nuxt`
- Check Node.js version compatibility

**TypeScript Errors**
- Run type checking: `npm run typecheck`
- Check TypeScript configuration
- Ensure proper type imports

**Styling Issues**
- Check Tailwind configuration
- Verify CSS imports in nuxt.config.ts
- Clear browser cache

### Development Tips
- Use Vue DevTools for debugging
- Enable source maps in development
- Use hot reload for faster development
- Check browser console for errors

## 📞 Support

For support and questions:
- Check the main project documentation
- Create an issue on GitHub
- Review Nuxt.js documentation
- Join our community discussions

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Happy coding! 🎉**
