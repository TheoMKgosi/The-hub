<template>
  <div>
    <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow">
      <div class="p-6 border-b border-surface-light dark:border-surface-dark">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-text-light dark:text-text-dark flex items-center">
            <span class="text-emerald-600 dark:text-emerald-400 mr-2">💰</span>
            Financial Management
          </h2>
          <NuxtLink to="/finance" class="text-text-light dark:text-text-dark/60 hover:text-primary transition-colors">
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
                  '0.00'}} / ${{ budget.amount.toFixed(2) }}
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
          <div v-if="budgetStore.budgets.length === 0" class="text-center py-4 text-text-light dark:text-text-dark/60">
            <p class="text-sm">No budgets created yet</p>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="mt-6 pt-4 border-t border-surface-light dark:border-surface-dark">
          <div class="flex gap-2">
            <NuxtLink to="/finance" class="flex-1">
              <BaseButton variant="primary" size="sm" text="Manage Finances" class="w-full">
              </BaseButton>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>
