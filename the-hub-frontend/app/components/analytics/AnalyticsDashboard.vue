<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import { useGoalStore } from '@/stores/goals'

interface AnalyticsData {
  user_id: string
  period: string
  start_date: string
  end_date: string
  total_tasks: number
  total_completed: number
  average_completion_rate: number
  average_productivity_score: number
  total_time_spent_minutes: number
  average_time_per_task_minutes: number
  completion_rate_trend: number
  productivity_trend: number
  best_day?: string
  best_day_score: number
  most_productive_hour?: number
  most_productive_day_of_week?: number
  goals_completed: number
  goal_success_rate: number
  recommendations: string[]
}

interface ChartData {
  labels: string[]
  data: number[]
  metric: string
}

const taskStore = useTaskStore()
const goalStore = useGoalStore()

const analytics = ref<AnalyticsData | null>(null)
const chartData = ref<ChartData | null>(null)
const loading = ref(false)
const selectedPeriod = ref('weekly')
const selectedDays = ref(30)
const selectedMetric = ref('completion_rate')

const metrics = [
  { value: 'completion_rate', label: 'Completion Rate', color: '#10B981' },
  { value: 'productivity_score', label: 'Productivity Score', color: '#3B82F6' },
  { value: 'time_spent', label: 'Time Spent (min)', color: '#F59E0B' },
  { value: 'tasks_completed', label: 'Tasks Completed', color: '#EF4444' },
]

const loadAnalytics = async () => {
  loading.value = true
  try {
    // Load analytics data
    const analyticsResponse = await $fetch<AnalyticsData>('/api/analytics/tasks', {
      method: 'GET',
      query: {
        period: selectedPeriod.value,
        days: selectedDays.value
      }
    })
    analytics.value = analyticsResponse

    // Load chart data
    const chartResponse = await $fetch<ChartData>('/api/analytics/tasks/chart', {
      method: 'GET',
      query: {
        metric: selectedMetric.value,
        days: selectedDays.value
      }
    })
    chartData.value = chartResponse
  } catch (error) {
    console.error('Failed to load analytics:', error)
  } finally {
    loading.value = false
  }
}

const formatTime = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}h ${mins}m`
}

const formatPercentage = (value: number) => {
  return `${(value * 100).toFixed(1)}%`
}

const getTrendIcon = (trend: number) => {
  if (trend > 0) return '↗️'
  if (trend < 0) return '↘️'
  return '➡️'
}

const getTrendColor = (trend: number) => {
  if (trend > 0) return 'text-green-600'
  if (trend < 0) return 'text-red-600'
  return 'text-gray-600'
}

onMounted(() => {
  loadAnalytics()
})

// Watch for changes in filters
watch([selectedPeriod, selectedDays], () => {
  loadAnalytics()
})

watch(selectedMetric, () => {
  loadAnalytics()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-text-light dark:text-text-dark">Analytics Dashboard</h1>
        <p class="text-text-light/70 dark:text-text-dark/70">Track your productivity and task completion patterns</p>
      </div>

      <!-- Filters -->
      <div class="flex gap-4">
        <select v-model="selectedPeriod" class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md">
          <option value="daily">Daily</option>
          <option value="weekly">Weekly</option>
          <option value="monthly">Monthly</option>
        </select>

        <select v-model="selectedDays" class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md">
          <option :value="7">7 days</option>
          <option :value="30">30 days</option>
          <option :value="90">90 days</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>

    <!-- Analytics Content -->
    <div v-else-if="analytics" class="space-y-6">

      <!-- Key Metrics Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Completion Rate</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatPercentage(analytics.average_completion_rate) }}
              </p>
            </div>
            <div class="text-2xl">✅</div>
          </div>
          <div class="mt-2 flex items-center text-sm">
            <span :class="getTrendColor(analytics.completion_rate_trend)">
              {{ getTrendIcon(analytics.completion_rate_trend) }}
              {{ Math.abs(analytics.completion_rate_trend * 100).toFixed(1) }}%
            </span>
            <span class="text-text-light/50 dark:text-text-dark/50 ml-1">vs previous period</span>
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Productivity Score</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ analytics.average_productivity_score.toFixed(1) }}
              </p>
            </div>
            <div class="text-2xl">📈</div>
          </div>
          <div class="mt-2 flex items-center text-sm">
            <span :class="getTrendColor(analytics.productivity_trend)">
              {{ getTrendIcon(analytics.productivity_trend) }}
              {{ Math.abs(analytics.productivity_trend).toFixed(1) }}
            </span>
            <span class="text-text-light/50 dark:text-text-dark/50 ml-1">vs previous period</span>
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Time Spent</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatTime(analytics.total_time_spent_minutes) }}
              </p>
            </div>
            <div class="text-2xl">⏱️</div>
          </div>
          <div class="mt-2 text-sm text-text-light/70 dark:text-text-dark/70">
            {{ formatTime(analytics.average_time_per_task_minutes) }} per task
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Tasks Completed</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ analytics.total_completed }}
              </p>
            </div>
            <div class="text-2xl">🎯</div>
          </div>
          <div class="mt-2 text-sm text-text-light/70 dark:text-text-dark/70">
            of {{ analytics.total_tasks }} total tasks
          </div>
        </div>
      </div>

      <!-- Chart Section -->
      <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-6 border border-surface-light/20 dark:border-surface-dark/20">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-text-light dark:text-text-dark">Productivity Trends</h2>
          <select v-model="selectedMetric" class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md">
            <option v-for="metric in metrics" :key="metric.value" :value="metric.value">
              {{ metric.label }}
            </option>
          </select>
        </div>

        <div v-if="chartData" class="h-64">
          <!-- Simple chart visualization (you can replace with a proper charting library) -->
          <div class="flex items-end justify-between h-full space-x-1">
            <div
              v-for="(value, index) in chartData.data"
              :key="index"
              class="flex-1 bg-primary/20 rounded-t flex flex-col items-center justify-end pb-2"
              :style="{ height: `${Math.max((value / Math.max(...chartData.data)) * 100, 5)}%` }"
            >
              <div class="text-xs text-text-light/70 dark:text-text-dark/70 transform -rotate-45 whitespace-nowrap">
                {{ value.toFixed(1) }}
              </div>
            </div>
          </div>
          <div class="flex justify-between mt-2 text-xs text-text-light/70 dark:text-text-dark/70">
            <span v-for="(label, index) in chartData.labels" :key="index" class="flex-1 text-center">
              {{ new Date(label).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Recommendations -->
      <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-6 border border-surface-light/20 dark:border-surface-dark/20">
        <h2 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Recommendations</h2>
        <div class="space-y-3">
          <div
            v-for="(recommendation, index) in analytics.recommendations"
            :key="index"
            class="flex items-start gap-3 p-3 bg-primary/5 dark:bg-primary/10 rounded-lg border border-primary/20"
          >
            <div class="text-primary mt-0.5">💡</div>
            <p class="text-text-light dark:text-text-dark">{{ recommendation }}</p>
          </div>
        </div>
      </div>

      <!-- Best Performance -->
      <div v-if="analytics.best_day" class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-6 border border-surface-light/20 dark:border-surface-dark/20">
        <h2 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Best Performance</h2>
        <div class="flex items-center gap-4">
          <div class="text-4xl">🏆</div>
          <div>
            <p class="text-text-light dark:text-text-dark font-medium">
              {{ new Date(analytics.best_day).toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}
            </p>
            <p class="text-text-light/70 dark:text-text-dark/70">
              Productivity Score: {{ analytics.best_day_score.toFixed(1) }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- No Data State -->
    <div v-else class="text-center py-12">
      <div class="text-6xl mb-4">📊</div>
      <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No Analytics Data Yet</h3>
      <p class="text-text-light/70 dark:text-text-dark/70">Complete some tasks to start seeing your productivity insights!</p>
    </div>
  </div>
</template>