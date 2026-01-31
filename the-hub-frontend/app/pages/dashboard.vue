<script setup lang="ts">
const taskStore = useTaskStore()
const authStore = useAuthStore()
const budgetStore = useBudgetStore()
const transactionStore = useTransactionStore()
const incomeStore = useIncomeStore()
const name = authStore.user?.name

// Financial overview calculations
const totalIncome = computed(() => {
  return incomeStore.incomes.reduce((sum, income) => sum + income.amount, 0)
})

const totalBudgeted = computed(() => {
  return budgetStore.budgets.reduce((sum, budget) => sum + budget.amount, 0)
})

const totalSpent = computed(() => {
  return budgetStore.analytics.reduce((sum, analytic) => sum + analytic.spent_amount, 0)
})

const budgetUtilization = computed(() => {
  if (totalBudgeted.value === 0) return 0
  return Math.round((totalSpent.value / totalBudgeted.value) * 100)
})

const budgetsOnTrack = computed(() => {
  return budgetStore.analytics.filter(analytic => analytic.status === 'on_track').length
})

const budgetsAtRisk = computed(() => {
  return budgetStore.analytics.filter(analytic => analytic.status === 'warning' || analytic.status === 'over_budget').length
})

const netPosition = computed(() => {
  return totalIncome.value - totalSpent.value
})

// Fetch data on mount
onMounted(async () => {
  await Promise.all([
    taskStore.fetchTasks(),
    // budgetStore.fetchBudget(),
    // budgetStore.fetchBudgetAnalytics(),
    // transactionStore.fetchTransactions(),
    // incomeStore.fetchIncomes()
  ])
})
</script>

<template>
  <main>
    <h1 class="text-center text-text-light dark:text-text-dark mb-8 pt-8">The Hub welcomes you, {{ name || 'stranger' }}
    </h1>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mb-8">
      <UiBanner />
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

        <!-- Financial Overview Cards -->
        <!--
         <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow p-6">
           <div class="flex items-center">
             <div class="h-8 w-8 bg-green-100 dark:bg-green-900/20 rounded-lg flex items-center justify-center">
               <span class="text-green-600 dark:text-green-400 font-bold">$</span>
             </div>
             <div class="ml-4">
               <p class="text-sm font-medium text-text-light dark:text-text-dark">Total Income</p>
               <p class="text-2xl font-semibold text-green-600 dark:text-green-400">${{ totalIncome.toLocaleString() }}</p>
             </div>
           </div>
         </div>

         <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow p-6">
           <div class="flex items-center">
             <div class="h-8 w-8 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
               <span class="text-blue-600 dark:text-blue-400 font-bold">📊</span>
             </div>
             <div class="ml-4">
               <p class="text-sm font-medium text-text-light dark:text-text-dark">Budget Utilization</p>
               <p class="text-2xl font-semibold text-blue-600 dark:text-blue-400">{{ budgetUtilization }}%</p>
             </div>
           </div>
         </div>

         <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow p-6">
           <div class="flex items-center">
             <div :class="[
               'h-8 w-8 rounded-lg flex items-center justify-center',
               netPosition >= 0 ? 'bg-emerald-100 dark:bg-emerald-900/20' : 'bg-red-100 dark:bg-red-900/20'
             ]">
               <span :class="[
                 'font-bold',
                 netPosition >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'
               ]">
                 {{ netPosition >= 0 ? '↗' : '↘' }}
               </span>
             </div>
             <div class="ml-4">
               <p class="text-sm font-medium text-text-light dark:text-text-dark">Net Position</p>
               <p :class="[
                 'text-2xl font-semibold',
                 netPosition >= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'
               ]">
                 ${{ Math.abs(netPosition).toLocaleString() }}
               </p>
             </div>
           </div>
         </div>
        <!--
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
              <span class="text-blue-600 font-bold">⏰</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Hours Today</p>
              <p class="text-2xl font-semibold text-gray-900">6.5h</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-purple-100 rounded-lg flex items-center justify-center">
              <span class="text-purple-600 font-bold">📚</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Courses Active</p>
              <p class="text-2xl font-semibold text-gray-900">3</p>
            </div>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="h-8 w-8 bg-emerald-100 rounded-lg flex items-center justify-center">
              <span class="text-emerald-600 font-bold">$</span>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Budget Left</p>
              <p class="text-2xl font-semibold text-gray-900">$2,340</p>
            </div>
          </div>
        </div>
        -->
      </div>

      <!-- Management Sections Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow">
          <TaskDashboard />
        </div>
        <!-- Time Management -->
        <!--
        <div class="bg-white rounded-lg shadow">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 flex items-center">
                <span class="text-blue-600 mr-2">⏰</span>
                Time Management
              </h2>
              <button class="text-gray-400 hover:text-gray-600">
                <span class="text-lg">📅</span>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center justify-between p-3 bg-blue-50 rounded-lg border-l-4 border-blue-500">
                <div>
                  <p class="text-sm font-medium text-gray-900">Team Standup</p>
                  <p class="text-xs text-gray-600">9:00 AM - 9:30 AM</p>
                </div>
                <span class="text-xs text-blue-600 font-medium">Now</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p class="text-sm font-medium text-gray-900">Design Review</p>
                  <p class="text-xs text-gray-600">11:00 AM - 12:00 PM</p>
                </div>
                <span class="text-xs text-gray-500">Upcoming</span>
              </div>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p class="text-sm font-medium text-gray-900">Client Presentation</p>
                  <p class="text-xs text-gray-600">2:00 PM - 3:00 PM</p>
                </div>
                <span class="text-xs text-gray-500">Upcoming</span>
              </div>
            </div>
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-900 mb-3">Time Tracking</h3>
              <div class="space-y-2">
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Development</span>
                  <span class="text-sm font-medium text-gray-900">4.5h</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Meetings</span>
                  <span class="text-sm font-medium text-gray-900">1.5h</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-sm text-gray-600">Planning</span>
                  <span class="text-sm font-medium text-gray-900">0.5h</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        -->

        <!-- Learning Management -->
        <!--
        <div class="bg-white rounded-lg shadow">
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 flex items-center">
                <span class="text-purple-600 mr-2">📚</span>
                Learning Management
              </h2>
              <button class="text-gray-400 hover:text-gray-600">
                <span class="text-lg">🎯</span>
              </button>
            </div>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="p-4 bg-purple-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">Advanced React Patterns</h3>
                  <span class="text-xs text-purple-600 font-medium">85%</span>
                </div>
                <div class="w-full bg-purple-200 rounded-full h-2 mb-2">
                  <div class="bg-purple-600 h-2 rounded-full" style="width: 85%"></div>
                </div>
                <p class="text-xs text-gray-600">12/14 lessons completed</p>
              </div>
              <div class="p-4 bg-blue-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">TypeScript Fundamentals</h3>
                  <span class="text-xs text-blue-600 font-medium">60%</span>
                </div>
                <div class="w-full bg-blue-200 rounded-full h-2 mb-2">
                  <div class="bg-blue-600 h-2 rounded-full" style="width: 60%"></div>
                </div>
                <p class="text-xs text-gray-600">9/15 lessons completed</p>
              </div>
              <div class="p-4 bg-green-50 rounded-lg">
                <div class="flex justify-between items-start mb-2">
                  <h3 class="text-sm font-medium text-gray-900">System Design Basics</h3>
                  <span class="text-xs text-green-600 font-medium">30%</span>
                </div>
                <div class="w-full bg-green-200 rounded-full h-2 mb-2">
                  <div class="bg-green-600 h-2 rounded-full" style="width: 30%"></div>
                </div>
                <p class="text-xs text-gray-600">3/10 lessons completed</p>
              </div>
            </div>
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-900 mb-3">Study Goals</h3>
              <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <span class="text-sm text-gray-900">Daily study time</span>
                <span class="text-sm font-medium text-green-600">2h / 2h ✓</span>
              </div>
            </div>
          </div>
        </div>
        -->

        <!-- Financial Management -->
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow">
          <div class="p-6 border-b border-surface-light dark:border-surface-dark">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-text-light dark:text-text-dark flex items-center">
                <span class="text-emerald-600 dark:text-emerald-400 mr-2">💰</span>
                Financial Management
              </h2>
              <NuxtLink to="/finance"
                class="text-text-light dark:text-text-dark/60 hover:text-primary transition-colors">
                <span class="text-lg">📈</span>
              </NuxtLink>
            </div>
          </div>
          <div class="p-6">
            <!-- Budget Status Summary -->
            <div class="grid grid-cols-3 gap-4 mb-6">
              <div class="text-center p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
                <p class="text-xs text-green-700 dark:text-green-300 font-medium">On Track</p>
                <p class="text-lg font-semibold text-green-600 dark:text-green-400">{{ budgetsOnTrack }}</p>
              </div>
              <div class="text-center p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
                <p class="text-xs text-yellow-700 dark:text-yellow-300 font-medium">At Risk</p>
                <p class="text-lg font-semibold text-yellow-600 dark:text-yellow-400">{{ budgetsAtRisk }}</p>
              </div>
              <div class="text-center p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                <p class="text-xs text-blue-700 dark:text-blue-300 font-medium">Total</p>
                <p class="text-lg font-semibold text-blue-600 dark:text-blue-400">{{ budgetStore.budgets.length }}</p>
              </div>
            </div>

            <!-- Recent Transactions -->
            <div class="space-y-3 mb-6">
              <h3 class="text-sm font-medium text-text-light dark:text-text-dark">Recent Transactions</h3>
              <div v-for="transaction in transactionStore.transactions.slice(0, 3)" :key="transaction.transaction_id"
                class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
                <div class="flex items-center">
                  <div :class="[
                    'w-8 h-8 rounded-full flex items-center justify-center',
                    transaction.type === 'income' ? 'bg-green-100 dark:bg-green-900/20' : 'bg-red-100 dark:bg-red-900/20'
                  ]">
                    <span :class="[
                      'text-sm',
                      transaction.type === 'income' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
                    ]">
                      {{ transaction.type === 'income' ? '+' : '-' }}
                    </span>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-text-light dark:text-text-dark">{{ transaction.description }}</p>
                    <p class="text-xs text-text-light dark:text-text-dark/60">{{ new
                      Date(transaction.date).toLocaleDateString() }}</p>
                  </div>
                </div>
                <span :class="[
                  'text-sm font-medium',
                  transaction.type === 'income' ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
                ]">
                  {{ transaction.type === 'income' ? '+' : '-' }}${{ transaction.amount.toFixed(2) }}
                </span>
              </div>
              <div v-if="transactionStore.transactions.length === 0"
                class="text-center py-4 text-text-light dark:text-text-dark/60">
                <p class="text-sm">No transactions yet</p>
              </div>
            </div>

            <!-- Budget Overview -->
            <div class="space-y-3">
              <h3 class="text-sm font-medium text-text-light dark:text-text-dark mb-3">Budget Overview</h3>
              <div v-for="budget in budgetStore.budgets.slice(0, 3)" :key="budget.budget_id" class="space-y-2">
                <div class="flex justify-between items-center">
                  <span class="text-sm text-text-light dark:text-text-dark/80">{{ budget.Category.name }}</span>
                  <span class="text-sm text-text-light dark:text-text-dark">
                    ${{budgetStore.analytics.find(a => a.budget_id === budget.budget_id)?.spent_amount.toFixed(2) ||
                    '0.00' }} / ${{ budget.amount.toFixed(2) }}
                  </span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div :class="[
                    'h-2 rounded-full transition-all duration-300',
                    (budgetStore.analytics.find(a => a.budget_id === budget.budget_id)?.utilization_rate || 0) >= 100 ? 'bg-red-500' :
                      (budgetStore.analytics.find(a => a.budget_id === budget.budget_id)?.utilization_rate || 0) >= 80 ? 'bg-yellow-500' :
                        'bg-green-500'
                  ]"
                    :style="{ width: Math.min(budgetStore.analytics.find(a => a.budget_id === budget.budget_id)?.utilization_rate || 0, 100) + '%' }">
                  </div>
                </div>
              </div>
              <div v-if="budgetStore.budgets.length === 0"
                class="text-center py-4 text-text-light dark:text-text-dark/60">
                <p class="text-sm">No budgets created yet</p>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="mt-6 pt-4 border-t border-surface-light dark:border-surface-dark">
              <div class="flex gap-2">
                <NuxtLink to="/finance" class="flex-1">
                  <UiButton variant="primary" size="sm" class="w-full">
                    Manage Finances
                  </UiButton>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </main>
</template>
