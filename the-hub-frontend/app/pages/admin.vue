<script setup lang="ts">
definePageMeta({
  middleware: ['authenticated', 'admin']
})

const adminStore = useAdminStore()
const auth = useAuthStore()
const { stats, loading } = storeToRefs(adminStore)
const router = useRouter()

onMounted(() => {
  adminStore.getSystemStats()
})

const statCards = [
  {
    title: 'Total Users',
    value: stats.value?.total_users || 0,
    icon: '👥',
    color: 'blue'
  },
  {
    title: 'Active Users',
    value: stats.value?.active_users || 0,
    icon: '🟢',
    color: 'green'
  },
  {
    title: 'Total Tasks',
    value: stats.value?.total_tasks || 0,
    icon: '📋',
    color: 'purple'
  },
  {
    title: 'Total Goals',
    value: stats.value?.total_goals || 0,
    icon: '🎯',
    color: 'orange'
  },
  {
    title: 'Total Decks',
    value: stats.value?.total_decks || 0,
    icon: '📚',
    color: 'indigo'
  },
  {
    title: 'Study Sessions',
    value: stats.value?.total_study_sessions || 0,
    icon: '📖',
    color: 'pink'
  }
]

const getColorClasses = (color: string) => {
  const colors = {
    blue: 'bg-blue-500/10 text-blue-600 dark:bg-blue-500/20 dark:text-blue-400',
    green: 'bg-green-500/10 text-green-600 dark:bg-green-500/20 dark:text-green-400',
    purple: 'bg-purple-500/10 text-purple-600 dark:bg-purple-500/20 dark:text-purple-400',
    orange: 'bg-orange-500/10 text-orange-600 dark:bg-orange-500/20 dark:text-orange-400',
    indigo: 'bg-indigo-500/10 text-indigo-600 dark:bg-indigo-500/20 dark:text-indigo-400',
    pink: 'bg-pink-500/10 text-pink-600 dark:bg-pink-500/20 dark:text-pink-400'
  }
  return colors[color as keyof typeof colors] || colors.blue
}

const handleLogout = async () => {
  await auth.logout()
  await router.push('/admin/login')
}

const refreshStats = () => {
  adminStore.getSystemStats()
}
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <!-- Admin Header -->
    <div class="bg-surface-light dark:bg-surface-dark border-b border-surface-light/50 dark:border-surface-dark/50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-3">
              <img src="/logo.svg" alt="The Hub Admin" class="h-8 w-8" />
              <div>
                <h1 class="text-xl font-bold text-text-light dark:text-text-dark">
                  Admin Dashboard
                </h1>
                <p class="text-sm text-text-light/60 dark:text-text-dark/60">
                  System overview and management
                </p>
              </div>
            </div>
          </div>
          
          <div class="flex items-center space-x-4">
            <UiButton
              @click="refreshStats"
              variant="secondary"
              size="sm"
              :loading="loading"
            >
              Refresh
            </UiButton>
            
            <div class="flex items-center space-x-3">
              <div class="text-right">
                <p class="text-sm font-medium text-text-light dark:text-text-dark">
                  {{ auth.user?.name }}
                </p>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">
                  Administrator
                </p>
              </div>
              <div class="w-8 h-8 bg-red-500 text-white rounded-full flex items-center justify-center text-sm font-bold">
                A
              </div>
            </div>
            
            <UiButton
              @click="handleLogout"
              variant="danger"
              size="sm"
            >
              Logout
            </UiButton>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-light dark:border-primary-dark"></div>
      </div>

      <!-- Stats Grid -->
      <div v-else-if="stats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div
          v-for="stat in statCards"
          :key="stat.title"
          class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50 hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-text-light/60 dark:text-text-dark/60">
                {{ stat.title }}
              </p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark mt-2">
                {{ stat.value.toLocaleString() }}
              </p>
            </div>
            <div :class="['p-3 rounded-full', getColorClasses(stat.color)]">
              <span class="text-2xl">{{ stat.icon }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-text-light dark:text-text-dark">
              User Management
            </h2>
            <svg class="h-6 w-6 text-text-light/40 dark:text-text-dark/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-4">
            Manage user accounts, roles, and permissions
          </p>
          <UiButton
            @click="$router.push('/admin/users')"
            variant="primary"
            size="sm"
            class="w-full"
          >
            Manage Users
          </UiButton>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-text-light dark:text-text-dark">
              System Health
            </h2>
            <svg class="h-6 w-6 text-text-light/40 dark:text-text-dark/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-4">
            Monitor system performance and status
          </p>
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-sm text-text-light/70 dark:text-text-dark/70">
              All systems operational
            </span>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-text-light dark:text-text-dark">
              Quick Actions
            </h2>
            <svg class="h-6 w-6 text-text-light/40 dark:text-text-dark/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div class="space-y-2">
            <UiButton
              @click="$router.push('/dashboard')"
              variant="secondary"
              size="sm"
              class="w-full"
            >
              View User Dashboard
            </UiButton>
            <UiButton
              @click="$router.push('/admin/users')"
              variant="secondary"
              size="sm"
              class="w-full"
            >
              Create New Admin
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50">
        <h2 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">
          Recent Admin Activity
        </h2>
        <div class="text-center py-8 text-text-light/50 dark:text-text-dark/50">
          <svg class="h-12 w-12 mx-auto mb-4 text-text-light/30 dark:text-text-dark/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-sm">
            No recent activity to display
          </p>
        </div>
      </div>
    </div>
  </div>
</template>