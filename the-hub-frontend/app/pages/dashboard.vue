<script setup lang="ts">
const taskStore = useTaskStore()
const authStore = useAuthStore()
const name = authStore.user?.name

onMounted(async () => {
  await Promise.all([
    taskStore.fetchTasks(),
  ])
})
</script>

<template>
  <main>
    <h1 class="text-center text-text-light dark:text-text-dark mb-8 pt-8">The Hub welcomes you, {{ name || 'stranger' }}
    </h1>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mb-8">
      <Banner />
    </div>
    <!-- Main Dashboard -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Overview Stats -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-success/10 dark:bg-success/20 rounded-lg flex items-center justify-center">
              <span class="text-success font-bold">✓</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-text-light dark:text-text-dark">Tasks completed</p>
              <p class="text-2xl font-semibold text-text-light dark:text-text-dark">{{ taskStore.completedTasks.length
              }}/ {{
                  taskStore.tasks.length }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Management Sections Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow">
          <TaskDashboard />
        </div>
      </div>
    </main>
  </main>
</template>
