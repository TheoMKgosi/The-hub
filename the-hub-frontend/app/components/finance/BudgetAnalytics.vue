<script setup>
const budgetStore = useBudgetStore()
const transactionStore = useTransactionStore()

// Analytics period selection
const selectedPeriod = ref('current')
const periods = [
  { value: 'current', label: 'Current Period' },
  { value: 'last_month', label: 'Last Month' },
  { value: 'last_3_months', label: 'Last 3 Months' }
]

// Computed analytics
const analyticsData = computed(() => {
  return budgetStore.analytics
})

const totalBudgeted = computed(() => {
  return analyticsData.value.reduce((sum, analytic) => sum + analytic.budget_amount, 0)
})

const totalSpent = computed(() => {
  return analyticsData.value.reduce((sum, analytic) => sum + analytic.spent_amount, 0)
})

const totalRemaining = computed(() => {
  return analyticsData.value.reduce((sum, analytic) => sum + analytic.remaining_amount, 0)
})

const overallUtilization = computed(() => {
  if (totalBudgeted.value === 0) return 0
  return Math.round((totalSpent.value / totalBudgeted.value) * 100)
})

const topSpendingCategories = computed(() => {
  return analyticsData.value
    .sort((a, b) => b.spent_amount - a.spent_amount)
    .slice(0, 5)
})

const budgetPerformanceData = computed(() => {
  return analyticsData.value.map(analytic => ({
    category: analytic.category_name,
    budgeted: analytic.budget_amount,
    spent: analytic.spent_amount,
    remaining: analytic.remaining_amount,
    utilization: analytic.utilization_rate,
    status: analytic.status
  }))
})

// Fetch analytics when period changes
watch(selectedPeriod, async (newPeriod) => {
  await budgetStore.fetchBudgetAnalytics(newPeriod)
})

onMounted(async () => {
  await budgetStore.fetchBudgetAnalytics(selectedPeriod.value)
  if (transactionStore.transactions.length === 0) {
    transactionStore.fetchTransactions()
  }
})
</script>

<template>
  <div class="p-4 max-w-6xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 text-text-light dark:text-text-dark">Budget Analytics & Reports</h2>

    <!-- Period Selector -->
    <div class="mb-6">
      <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Analysis Period</label>
      <select v-model="selectedPeriod"
        class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
        <option v-for="period in periods" :key="period.value" :value="period.value">
          {{ period.label }}
        </option>
      </select>
    </div>

    <!-- Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <div class="flex items-center">
          <div class="h-8 w-8 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
            <span class="text-blue-600 dark:text-blue-400 font-bold">$</span>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-light dark:text-text-dark">Total Budgeted</p>
            <p class="text-2xl font-semibold text-blue-600 dark:text-blue-400">${{ totalBudgeted.toLocaleString() }}</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <div class="flex items-center">
          <div class="h-8 w-8 bg-red-100 dark:bg-red-900/20 rounded-lg flex items-center justify-center">
            <span class="text-red-600 dark:text-red-400 font-bold">💸</span>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-light dark:text-text-dark">Total Spent</p>
            <p class="text-2xl font-semibold text-red-600 dark:text-red-400">${{ totalSpent.toLocaleString() }}</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <div class="flex items-center">
          <div class="h-8 w-8 bg-green-100 dark:bg-green-900/20 rounded-lg flex items-center justify-center">
            <span class="text-green-600 dark:text-green-400 font-bold">💰</span>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-light dark:text-text-dark">Remaining</p>
            <p class="text-2xl font-semibold text-green-600 dark:text-green-400">${{ totalRemaining.toLocaleString() }}</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <div class="flex items-center">
          <div class="h-8 w-8 bg-purple-100 dark:bg-purple-900/20 rounded-lg flex items-center justify-center">
            <span class="text-purple-600 dark:text-purple-400 font-bold">📊</span>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-light dark:text-text-dark">Utilization</p>
            <p class="text-2xl font-semibold text-purple-600 dark:text-purple-400">{{ overallUtilization }}%</p>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Budget Performance Chart -->
      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Budget Performance</h3>
        <div class="space-y-4">
          <div v-for="budget in budgetPerformanceData" :key="budget.category"
            class="space-y-2">
            <div class="flex justify-between items-center">
              <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ budget.category }}</span>
              <span class="text-sm text-text-light dark:text-text-dark/60">
                ${{ budget.spent.toFixed(2) }} / ${{ budget.budgeted.toFixed(2) }}
              </span>
            </div>
             <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
               <div
                 :class="[
                   'h-3 rounded-full transition-all duration-300',
                   budget.status === 'over_budget' ? 'bg-red-500' :
                   budget.status === 'warning' ? 'bg-yellow-500' :
                   budget.status === 'caution' ? 'bg-orange-500' :
                   'bg-green-500'
                 ]"
                 :style="{ width: Math.min(budget.utilization, 100) + '%' }"
               ></div>
             </div>
             <div class="flex justify-between text-xs text-text-light dark:text-text-dark/60">
               <span>{{ budget.utilization.toFixed(1) }}% used</span>
               <span :class="[
                 budget.status === 'on_track' ? 'text-green-600 dark:text-green-400' :
                 budget.status === 'caution' ? 'text-orange-600 dark:text-orange-400' :
                 budget.status === 'warning' ? 'text-yellow-600 dark:text-yellow-400' :
                 'text-red-600 dark:text-red-400'
               ]">
                 {{ budget.status.replace('_', ' ').toUpperCase() }}
               </span>
             </div>
          </div>
        </div>
      </div>

      <!-- Top Spending Categories -->
      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
        <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Top Spending Categories</h3>
        <div class="space-y-4">
          <div v-for="(category, index) in topSpendingCategories" :key="category.category_name"
            class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
            <div class="flex items-center gap-3">
              <div :class="[
                'w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold',
                index === 0 ? 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-600 dark:text-yellow-400' :
                index === 1 ? 'bg-gray-100 dark:bg-gray-900/20 text-gray-600 dark:text-gray-400' :
                'bg-orange-100 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400'
              ]">
                {{ index + 1 }}
              </div>
              <div>
                <p class="font-medium text-text-light dark:text-text-dark">{{ category.category_name }}</p>
                <p class="text-sm text-text-light dark:text-text-dark/60">
                  {{ category.utilization_rate.toFixed(1) }}% of budget
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-semibold text-text-light dark:text-text-dark">${{ category.spent_amount.toFixed(2) }}</p>
              <p class="text-sm text-text-light dark:text-text-dark/60">
                of ${{ category.budget_amount.toFixed(2) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Detailed Analytics Table -->
    <div class="mt-8 bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
      <div class="p-6 border-b border-surface-light dark:border-surface-dark">
        <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Detailed Budget Analysis</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-surface-light/50 dark:bg-surface-dark/50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Category</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Budget</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Spent</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Remaining</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Utilization</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Daily Rate</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-text-light dark:text-text-dark/60 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-surface-light dark:divide-surface-dark">
            <tr v-for="analytic in analyticsData" :key="analytic.budget_id">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-text-light dark:text-text-dark">
                {{ analytic.category_name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light dark:text-text-dark">
                ${{ analytic.budget_amount.toFixed(2) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light dark:text-text-dark">
                ${{ analytic.spent_amount.toFixed(2) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light dark:text-text-dark">
                ${{ analytic.remaining_amount.toFixed(2) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light dark:text-text-dark">
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  analytic.utilization_rate >= 100 ? 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100' :
                  analytic.utilization_rate >= 80 ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100' :
                  'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100'
                ]">
                  {{ analytic.utilization_rate.toFixed(1) }}%
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light dark:text-text-dark">
                ${{ analytic.daily_spending_rate.toFixed(2) }}
              </td>
               <td class="px-6 py-4 whitespace-nowrap text-sm">
                 <span :class="[
                   'px-2 py-1 text-xs font-medium rounded-full',
                   analytic.status === 'on_track' ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100' :
                   analytic.status === 'caution' ? 'bg-orange-100 text-orange-800 dark:bg-orange-800 dark:text-orange-100' :
                   analytic.status === 'warning' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100' :
                   'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100'
                 ]">
                   {{ analytic.status.replace('_', ' ').toUpperCase() }}
                 </span>
               </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>