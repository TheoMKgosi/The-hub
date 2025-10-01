<script setup lang="ts">
definePageMeta({
  layout: false
})

const router = useRouter()

const navigateToLogin = () => {
  router.push('/login')
}

const navigateToRegister = () => {
  router.push('/register')
}

const features = [
  {
    icon: '🎯',
    title: 'Task & Goal Management',
    description: 'Create, categorize, and prioritize tasks with advanced filtering and search. Set long-term goals with progress monitoring and achievement tracking.'
  },
  {
    icon: '⏰',
    title: 'Intelligent Time Management',
    description: 'AI-powered scheduling with smart time blocking and calendar optimization. Sync with Google Calendar and track productivity patterns.'
  },
  {
    icon: '📚',
    title: 'Advanced Learning Management',
    description: 'Scientific flashcard system using spaced repetition. Structure learning materials into topics, decks, and study paths with detailed progress analytics.'
  },
  {
    icon: '💰',
    title: 'Comprehensive Finance Management',
    description: 'Track expenses, create budgets, and monitor spending patterns. Set savings targets with progress visualization and detailed transaction analytics.'
  },
  {
    icon: '🤖',
    title: 'AI Integration & Automation',
    description: 'Smart recommendations, natural language processing, and automated workflows. AI-powered task prioritization and personalized productivity insights.'
  },
  {
    icon: '📊',
    title: 'Analytics & Reporting',
    description: 'Comprehensive dashboard with key metrics and insights. Generate detailed reports on tasks, time, learning, and finances with trend analysis.'
  }
]
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-background-light via-surface-light to-background-light dark:from-background-dark dark:via-surface-dark dark:to-background-dark">
    <!-- Hero Section -->
    <section class="relative px-6 pt-14 pb-16 sm:px-6 sm:pt-20 sm:pb-24 lg:px-8">
      <div class="mx-auto max-w-7xl">
        <div class="text-center">
          <!-- Logo -->
          <div class="mx-auto mb-8 flex justify-center">
            <img src="/logo.svg" alt="The Hub Logo" class="h-16 w-16 sm:h-20 sm:w-20" />
          </div>

          <!-- Main Heading -->
          <h1 class="text-4xl font-bold tracking-tight text-text-light dark:text-text-dark sm:text-6xl lg:text-7xl">
            The Hub
          </h1>

          <!-- Tagline -->
          <p class="mx-auto mt-6 max-w-2xl text-lg leading-8 text-text-light/80 dark:text-text-dark/80 sm:text-xl">
            Your Intelligent Productivity Ecosystem
          </p>

          <!-- Description -->
          <p class="mx-auto mt-6 max-w-3xl text-base leading-7 text-text-light/70 dark:text-text-dark/70 sm:text-lg">
            Streamline your daily life through intelligent task management, learning systems, financial tracking, and time optimization.
            Built with modern web technologies for a seamless experience across devices.
          </p>

          <!-- CTA Buttons -->
          <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
            <UiButton
              @click="navigateToRegister"
              variant="primary"
              size="lg"
              class="w-full sm:w-auto px-8 py-3 text-lg"
            >
              Get Started
            </UiButton>
            <UiButton
              @click="navigateToLogin"
              variant="secondary"
              size="lg"
              class="w-full sm:w-auto px-8 py-3 text-lg"
            >
              Sign In
            </UiButton>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 sm:py-24 bg-surface-light/50 dark:bg-surface-dark/50">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="mx-auto max-w-2xl text-center">
          <h2 class="text-3xl font-bold tracking-tight text-text-light dark:text-text-dark sm:text-4xl">
            Everything you need to stay productive
          </h2>
          <p class="mt-6 text-lg leading-8 text-text-light/70 dark:text-text-dark/70">
            Powerful features designed to adapt to your workflow and help you achieve your goals more efficiently.
          </p>
        </div>

        <div class="mx-auto mt-16 max-w-7xl">
          <div class="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="feature in features"
              :key="feature.title"
              class="relative bg-background-light dark:bg-background-dark rounded-2xl p-8 shadow-lg border border-surface-light dark:border-surface-dark hover:shadow-xl transition-shadow duration-300"
            >
              <div class="text-4xl mb-4">{{ feature.icon }}</div>
              <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-3">
                {{ feature.title }}
              </h3>
              <p class="text-text-light/70 dark:text-text-dark/70 leading-relaxed">
                {{ feature.description }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-surface-light dark:bg-surface-dark border-t border-surface-light dark:border-surface-dark">
      <div class="mx-auto max-w-7xl px-6 py-12 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Logo and Description -->
          <div class="col-span-1 md:col-span-2">
            <div class="flex items-center space-x-3 mb-4">
              <img src="/logo.svg" alt="The Hub Logo" class="h-8 w-8" />
              <span class="text-xl font-bold text-text-light dark:text-text-dark">The Hub</span>
            </div>
            <p class="text-text-light/70 dark:text-text-dark/70 max-w-md">
              A comprehensive personal productivity platform designed to streamline your daily life through intelligent task management and goal achievement.
            </p>
          </div>

          <!-- Links -->
          <div>
            <h3 class="text-sm font-semibold text-text-light dark:text-text-dark uppercase tracking-wider mb-4">
              Resources
            </h3>
            <ul class="space-y-2">
              <li>
                <a href="https://github.com/your-org/the-hub" target="_blank" class="text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark transition-colors">
                  GitHub
                </a>
              </li>
              <li>
                <a href="/docs/api.md" class="text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark transition-colors">
                  API Documentation
                </a>
              </li>
              <li>
                <a href="/docs/contributing.md" class="text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark transition-colors">
                  Contributing
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div class="mt-8 pt-8 border-t border-surface-light dark:border-surface-dark">
          <p class="text-center text-text-light/50 dark:text-text-dark/50 text-sm">
            © 2024 The Hub. Made with ❤️ by the The Hub team.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>
