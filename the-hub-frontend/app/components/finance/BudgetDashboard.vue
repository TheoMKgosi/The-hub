<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useBudgetStore } from '@/stores/finance'
import { useCategoryStore } from '@/stores/finance'
import { addToast } from '@/composables/useToast'

interface BudgetAnalytics {
  budget_id: string
  category_name: string
  budget_amount: number
  spent_amount: number
  remaining_amount: number
  utilization_rate: number
  days_remaining: number
  daily_spending_rate: number
  status: string
}

interface BudgetAlert {
  budget_id: string
  category_name: string
  alert_type: string
  message: string
  current_usage: number
  budget_amount: number
  days_remaining: number
}

interface BudgetSuggestion {
  category_id: string
  category_name: string
  current_average: number
  suggested_budget: number
  transaction_count: number
  variability: number
  confidence: string
}

const budgetStore = useBudgetStore()
const categoryStore = useCategoryStore()

const analytics = ref<BudgetAnalytics[]>([])
const alerts = ref<BudgetAlert[]>([])
const suggestions = ref<BudgetSuggestion[]>([])
const loading = ref(false)
const loadingSuggestions = ref(false)
const selectedPeriod = ref('current')

// Computed properties for dashboard metrics
const totalBudget = computed(() => budgetStore.analytics.reduce((sum, item) => sum + item.budget_amount, 0))
const totalSpent = computed(() => budgetStore.analytics.reduce((sum, item) => sum + item.spent_amount, 0))
const totalRemaining = computed(() => budgetStore.analytics.reduce((sum, item) => sum + item.remaining_amount, 0))
const averageUtilization = computed(() => {
  if (budgetStore.analytics.length === 0) return 0
  return budgetStore.analytics.reduce((sum, item) => sum + item.utilization_rate, 0) / budgetStore.analytics.length
})

const criticalAlerts = computed(() => budgetStore.alerts.filter(alert => alert.alert_type === 'over_budget' || alert.alert_type === 'danger'))
const warningAlerts = computed(() => budgetStore.alerts.filter(alert => alert.alert_type === 'warning'))

const loadBudgetData = async () => {
  await Promise.all([
    budgetStore.fetchBudgetAnalytics(selectedPeriod.value),
    budgetStore.fetchBudgetAlerts()
  ])
}

const loadBudgetSuggestions = async () => {
  await budgetStore.fetchBudgetSuggestions()
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'on_track': return 'text-green-600'
    case 'warning': return 'text-yellow-600'
    case 'over_budget': return 'text-red-600'
    default: return 'text-gray-600'
  }
}

const getStatusBgColor = (status: string) => {
  switch (status) {
    case 'on_track': return 'bg-green-100 dark:bg-green-900/20'
    case 'warning': return 'bg-yellow-100 dark:bg-yellow-900/20'
    case 'over_budget': return 'bg-red-100 dark:bg-red-900/20'
    default: return 'bg-gray-100 dark:bg-gray-900/20'
  }
}

const getAlertIcon = (alertType: string) => {
  switch (alertType) {
    case 'over_budget': return '🚨'
    case 'danger': return '⚠️'
    case 'warning': return '🔔'
    default: return 'ℹ️'
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount)
}

const formatPercentage = (value: number) => {
  return `${value.toFixed(1)}%`
}

onMounted(() => {
  loadBudgetData()
})

// Watch for period changes
watch(selectedPeriod, () => {
  loadBudgetData()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Budget Dashboard</h2>
        <p class="text-text-light/70 dark:text-text-dark/70">Monitor your spending and budget performance</p>
      </div>

      <!-- Period Selector -->
      <select
        v-model="selectedPeriod"
        class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md"
      >
        <option value="current">Current Period</option>
        <option value="last_month">Last Month</option>
        <option value="last_3_months">Last 3 Months</option>
      </select>
    </div>

    <!-- Loading State -->
    <div v-if="budgetStore.analyticsLoading || budgetStore.alertsLoading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>

    <div v-else>
      <!-- Key Metrics Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Total Budget</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatCurrency(totalBudget) }}
              </p>
            </div>
            <div class="text-2xl">💰</div>
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Total Spent</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatCurrency(totalSpent) }}
              </p>
            </div>
            <div class="text-2xl">💸</div>
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Remaining</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatCurrency(totalRemaining) }}
              </p>
            </div>
            <div class="text-2xl">📈</div>
          </div>
        </div>

        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">Avg Utilization</p>
              <p class="text-2xl font-bold text-text-light dark:text-text-dark">
                {{ formatPercentage(averageUtilization) }}
              </p>
            </div>
            <div class="text-2xl">📊</div>
          </div>
        </div>
      </div>

      <!-- Alerts Section -->
      <div v-if="budgetStore.alerts.length > 0" class="mb-6">
        <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Budget Alerts</h3>

        <!-- Critical Alerts -->
        <div v-if="criticalAlerts.length > 0" class="mb-4">
          <h4 class="text-md font-medium text-red-600 dark:text-red-400 mb-2">Critical Alerts</h4>
          <div class="space-y-2">
            <div
              v-for="alert in criticalAlerts"
              :key="alert.budget_id"
              class="flex items-start gap-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
            >
              <div class="text-lg">{{ getAlertIcon(alert.alert_type) }}</div>
              <div class="flex-1">
                <p class="text-text-light dark:text-text-dark font-medium">{{ alert.category_name }}</p>
                <p class="text-sm text-text-light/70 dark:text-text-dark/70">{{ alert.message }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Warning Alerts -->
        <div v-if="warningAlerts.length > 0">
          <h4 class="text-md font-medium text-yellow-600 dark:text-yellow-400 mb-2">Warnings</h4>
          <div class="space-y-2">
            <div
              v-for="alert in warningAlerts"
              :key="alert.budget_id"
              class="flex items-start gap-3 p-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg"
            >
              <div class="text-lg">{{ getAlertIcon(alert.alert_type) }}</div>
              <div class="flex-1">
                <p class="text-text-light dark:text-text-dark font-medium">{{ alert.category_name }}</p>
                <p class="text-sm text-text-light/70 dark:text-text-dark/70">{{ alert.message }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Budget Performance Chart -->
      <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-6 border border-surface-light/20 dark:border-surface-dark/20 mb-6">
        <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">Budget Performance</h3>

        <div v-if="budgetStore.analytics.length === 0" class="text-center py-8 text-text-light/60 dark:text-text-dark/60">
          No budget data available for the selected period
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="item in budgetStore.analytics"
            :key="item.budget_id"
            class="p-4 rounded-lg border border-surface-light/20 dark:border-surface-dark/20"
            :class="getStatusBgColor(item.status)"
          >
            <div class="flex items-center justify-between mb-2">
              <h4 class="font-semibold text-text-light dark:text-text-dark">{{ item.category_name }}</h4>
              <span :class="getStatusColor(item.status)" class="text-sm font-medium capitalize">
                {{ item.status.replace('_', ' ') }}
              </span>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-3">
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Budget</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ formatCurrency(item.budget_amount) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Spent</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ formatCurrency(item.spent_amount) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Remaining</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ formatCurrency(item.remaining_amount) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Days Left</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ item.days_remaining }}</p>
              </div>
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-surface-light/30 dark:bg-surface-dark/30 rounded-full h-2 mb-2">
              <div
                class="h-2 rounded-full transition-all duration-300"
                :class="item.utilization_rate >= 100 ? 'bg-red-500' : item.utilization_rate >= 80 ? 'bg-yellow-500' : 'bg-green-500'"
                :style="{ width: Math.min(item.utilization_rate, 100) + '%' }"
              ></div>
            </div>

            <div class="flex justify-between text-xs text-text-light/60 dark:text-text-dark/60">
              <span>{{ formatPercentage(item.utilization_rate) }} used</span>
              <span>${{ item.daily_spending_rate.toFixed(2) }}/day avg</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Budget Suggestions -->
      <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-6 border border-surface-light/20 dark:border-surface-dark/20">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Smart Suggestions</h3>
          <button
            @click="loadBudgetSuggestions"
            :disabled="budgetStore.suggestionsLoading"
            class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            {{ budgetStore.suggestionsLoading ? 'Loading...' : 'Get AI Suggestions' }}
          </button>
        </div>

        <div v-if="budgetStore.suggestions.length === 0" class="text-center py-4">
          <p class="text-text-light/60 dark:text-text-dark/60">
            Click "Get AI Suggestions" to receive personalized budget recommendations based on your spending patterns
          </p>
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="suggestion in budgetStore.suggestions"
            :key="suggestion.category_id"
            class="p-4 bg-primary/5 dark:bg-primary/10 rounded-lg border border-primary/20"
          >
            <div class="flex items-center justify-between mb-2">
              <h4 class="font-semibold text-text-light dark:text-text-dark">{{ suggestion.category_name }}</h4>
              <span class="text-xs px-2 py-1 rounded-full"
                    :class="suggestion.confidence === 'high' ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400' :
                           suggestion.confidence === 'medium' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400' :
                           'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'">
                {{ suggestion.confidence }} confidence
              </span>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-3">
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Current Avg</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ formatCurrency(suggestion.current_average) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Suggested Budget</p>
                <p class="font-semibold text-primary">{{ formatCurrency(suggestion.suggested_budget) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Transactions</p>
                <p class="font-semibold text-text-light dark:text-text-dark">{{ suggestion.transaction_count }}</p>
              </div>
              <div>
                <p class="text-xs text-text-light/60 dark:text-text-dark/60">Variability</p>
                <p class="font-semibold text-text-light dark:text-text-dark">${{ suggestion.variability.toFixed(2) }}</p>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <div class="flex-1 bg-surface-light/30 dark:bg-surface-dark/30 rounded-full h-2">
                <div
                  class="h-2 bg-primary rounded-full transition-all duration-300"
                  :style="{ width: Math.min((suggestion.current_average / suggestion.suggested_budget) * 100, 100) + '%' }"
                ></div>
              </div>
              <span class="text-xs text-text-light/60 dark:text-text-dark/60">
                {{ formatPercentage((suggestion.current_average / suggestion.suggested_budget) * 100) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>