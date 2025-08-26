<script setup lang="ts">
const statsStore = useStatsStore()

// Period selection
const selectedPeriod = ref<'today' | 'week' | 'month' | 'year'>('today')
const customDateRange = ref({
  start: '',
  end: ''
})
const useCustomRange = ref(false)

// Trends period
const trendsDays = ref(30)

// Fetch stats when period changes
watch(selectedPeriod, (newPeriod) => {
  if (!useCustomRange.value) {
    statsStore.fetchStats(newPeriod)
  }
})

watch(useCustomRange, (useCustom) => {
  if (!useCustom) {
    statsStore.fetchStats(selectedPeriod.value)
  }
})

// Initial load
onMounted(async () => {
  await Promise.all([
    statsStore.fetchStats('today'),
    statsStore.fetchTrends(30)
  ])
})

// Handle custom date range
const applyCustomRange = () => {
  if (customDateRange.value.start && customDateRange.value.end) {
    statsStore.fetchStats('custom', customDateRange.value.start, customDateRange.value.end)
  }
}

// Refresh data
const refreshData = async () => {
  await Promise.all([
    statsStore.fetchStats(useCustomRange.value ? 'custom' : selectedPeriod.value,
                         customDateRange.value.start, customDateRange.value.end),
    statsStore.fetchTrends(trendsDays.value)
  ])
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header with controls -->
    <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Task Statistics</h2>
          <p class="text-text-light/70 dark:text-text-dark/70 mt-1">Track your productivity and task completion patterns</p>
        </div>

        <div class="flex flex-col sm:flex-row gap-3">
          <!-- Period selector -->
          <div class="flex gap-2">
            <label class="flex items-center gap-2 text-sm text-text-light dark:text-text-dark">
              <input v-model="useCustomRange" type="checkbox" class="rounded" />
              Custom Range
            </label>
          </div>

          <template v-if="!useCustomRange">
            <select v-model="selectedPeriod"
                    class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
              <option value="today">Today</option>
              <option value="week">This Week</option>
              <option value="month">This Month</option>
              <option value="year">This Year</option>
            </select>
          </template>

          <template v-else>
            <input v-model="customDateRange.start" type="date"
                   class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <input v-model="customDateRange.end" type="date"
                   class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <UiButton @click="applyCustomRange" variant="primary" size="sm">
              Apply
            </UiButton>
          </template>

          <UiButton @click="refreshData" variant="default" size="sm" :disabled="statsStore.loading">
            <span v-if="statsStore.loading" class="animate-spin">⟳</span>
            <span v-else>↻</span>
          </UiButton>
        </div>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="statsStore.loading" class="text-center py-12">
      <div class="animate-spin text-4xl">⟳</div>
      <p class="text-text-light/70 dark:text-text-dark/70 mt-4">Loading statistics...</p>
    </div>

    <!-- Error state -->
    <div v-else-if="statsStore.error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
      <p class="text-red-600 dark:text-red-400">{{ statsStore.error }}</p>
      <UiButton @click="refreshData" variant="danger" size="sm" class="mt-3">
        Try Again
      </UiButton>
    </div>

    <!-- Statistics content -->
    <template v-else-if="statsStore.currentStats">
      <!-- Key Metrics Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Completion Rate</p>
              <p class="text-3xl font-bold text-primary">{{ statsStore.completionRate.toFixed(1) }}%</p>
            </div>
            <div class="text-success text-2xl">✓</div>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Total Tasks</p>
              <p class="text-3xl font-bold text-text-light dark:text-text-dark">{{ statsStore.totalTasks }}</p>
            </div>
            <div class="text-primary text-2xl">📋</div>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Completed</p>
              <p class="text-3xl font-bold text-success">{{ statsStore.completedTasks }}</p>
            </div>
            <div class="text-success text-2xl">✅</div>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Overdue</p>
              <p class="text-3xl font-bold text-warning">{{ statsStore.overdueTasks }}</p>
            </div>
            <div class="text-warning text-2xl">⚠️</div>
          </div>
        </div>
      </div>

      <!-- Time-sensitive tasks -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Due Today</p>
              <p class="text-2xl font-bold text-warning">{{ statsStore.tasksDueToday }}</p>
            </div>
            <div class="text-warning text-xl">📅</div>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Due Tomorrow</p>
              <p class="text-2xl font-bold text-primary">{{ statsStore.tasksDueTomorrow }}</p>
            </div>
            <div class="text-primary text-xl">📅</div>
          </div>
        </div>

        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-text-light/70 dark:text-text-dark/70 text-sm font-medium">Due This Week</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">{{ statsStore.tasksDueThisWeek }}</p>
            </div>
            <div class="text-text-light dark:text-text-dark text-xl">📊</div>
          </div>
        </div>
      </div>

      <!-- Charts Section -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Priority Distribution Chart -->
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Priority Distribution</h3>
          <div class="space-y-3">
            <div v-for="item in statsStore.priorityData" :key="item.priority"
                 class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
              <div class="flex items-center gap-3">
                <div class="w-3 h-3 rounded-full"
                     :class="item.priority.includes('1') || item.priority.includes('2') ? 'bg-success' :
                             item.priority.includes('3') ? 'bg-primary' : 'bg-warning'"></div>
                <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ item.priority }}</span>
              </div>
              <div class="text-right">
                <div class="text-lg font-bold text-text-light dark:text-text-dark">{{ item.count }}</div>
                <div class="text-xs text-text-light/70 dark:text-text-dark/70">
                  {{ item.count > 0 ? Math.round((item.completed / item.count) * 100) : 0 }}% completed
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Goal Distribution Chart -->
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Goal Association</h3>
          <div class="space-y-3">
            <div v-for="item in statsStore.goalData" :key="item.name"
                 class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
              <div class="flex items-center gap-3">
                <div class="w-3 h-3 rounded-full"
                     :class="item.name.includes('With') ? 'bg-success' : 'bg-primary'"></div>
                <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ item.name }}</span>
              </div>
              <div class="text-lg font-bold text-text-light dark:text-text-dark">{{ item.value }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Trends Section -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Completion Trends</h3>
          <select v-model="trendsDays"
                  @change="statsStore.fetchTrends(trendsDays)"
                  class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
            <option :value="7">Last 7 days</option>
            <option :value="30">Last 30 days</option>
            <option :value="90">Last 90 days</option>
          </select>
        </div>

        <div v-if="statsStore.trendsLoading" class="text-center py-8">
          <div class="animate-spin text-2xl">⟳</div>
          <p class="text-text-light/70 dark:text-text-dark/70 mt-2">Loading trends...</p>
        </div>

        <div v-else-if="statsStore.trendsCompletionData.length > 0" class="space-y-2">
          <div v-for="trend in statsStore.trendsCompletionData.slice(-10)" :key="trend.date"
               class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
            <div class="flex items-center gap-3">
              <span class="text-sm font-medium text-text-light dark:text-text-dark">
                {{ new Date(trend.date).toLocaleDateString() }}
              </span>
            </div>
            <div class="flex items-center gap-4">
              <div class="text-right">
                <div class="text-sm font-bold text-text-light dark:text-text-dark">
                  {{ trend.completedTasks }}/{{ trend.totalTasks }}
                </div>
                <div class="text-xs text-text-light/70 dark:text-text-dark/70">
                  {{ trend.completionRate.toFixed(1) }}% completed
                </div>
              </div>
              <div class="w-20 bg-surface-light dark:bg-surface-dark rounded-full h-2">
                <div class="bg-primary h-2 rounded-full transition-all duration-300"
                     :style="{ width: trend.completionRate + '%' }"></div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-8 text-text-light/70 dark:text-text-dark/70">
          No trend data available
        </div>
      </div>
    </template>

    <!-- No data state -->
    <div v-else class="text-center py-12">
      <div class="text-6xl mb-4">📊</div>
      <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No Statistics Available</h3>
      <p class="text-text-light/70 dark:text-text-dark/70 mb-4">Create some tasks to start seeing your productivity statistics.</p>
      <UiButton variant="primary" class="inline-flex items-center gap-2">
        <span>+</span>
        Create Task
      </UiButton>
    </div>
  </div>
</template>