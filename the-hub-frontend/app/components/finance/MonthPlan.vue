<script setup lang="ts">
import { parseDate } from '@internationalized/date'
import type { SelectMenuItem } from '@nuxt/ui'
const budgetStore = useBudgetStore()
const categoryStore = useCategoryStore()
const incomeStore = useIncomeStore()
const financialGoalStore = useFinancialGoalStore()
const transactionStore = useTransactionStore()
const toast = useToast()

const showCreateBudget = ref(false)
const selectedMonth = ref(new Date().toISOString().slice(0, 7))

const createForm = reactive({
  category_id: '',
  amount: 0,
  start_date: shallowRef(parseDate(new Date().toLocaleDateString('sv-SE'))),
  end_date: shallowRef(parseDate(new Date().toLocaleDateString('sv-SE'))),
})

const editingBudgetId = ref<string | null>(null)
const editForm = reactive({
  budget_id: '',
  category_id: '',
  amount: 0,
  start_date: '',
  end_date: '',
})

const totalIncome = computed(() => {
  const [year, month] = selectedMonth.value.split('-').map(Number)
  return incomeStore.incomes
    .filter(i => {
      const d = new Date(i.received_at)
      return d.getFullYear() === year && d.getMonth() + 1 === month
    })
    .reduce((sum, i) => sum + i.amount, 0)
})

const categories = computed(() => {
  const items: SelectMenuItem[] = []
  categoryStore.categories.map(category => (items.push({ label: category.name, id: category.budget_category_id })))
  return items
})

const totalBudgeted = computed(() =>
  budgetStore.budgets
    .filter(b => isBudgetInMonth(b, selectedMonth.value))
    .reduce((sum, b) => sum + b.amount, 0)
)

const totalSpent = computed(() =>
  budgetStore.analytics?.reduce((sum, a) => sum + a.spent_amount, 0)
)

const remaining = computed(() => totalIncome.value - totalBudgeted.value)

const monthBudgets = computed(() => {
  return budgetStore.budgets
    .filter(b => isBudgetInMonth(b, selectedMonth.value))
    .map(b => ({
      ...b,
      analytics: budgetStore.analytics?.find(a => a.budget_id === b.budget_id),
      linkedGoals: financialGoalStore.goals.filter(
        g => g.category_id === b.category_id && g.status === 'active'
      ),
    }))
})

function isBudgetInMonth(budget: any, month: string) {
  const [year, m] = month.split('-').map(Number)
  const start = new Date(budget.start_date)
  const end = new Date(budget.end_date)
  const monthStart = new Date(year, m - 1, 1)
  const monthEnd = new Date(year, m, 0)
  return start <= monthEnd && end >= monthStart
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}

function getProgressPercent(analytics: any) {
  if (!analytics) return 0
  return Math.min(analytics.utilization_rate, 100)
}

function getStatusClass(status: string) {
  switch (status) {
    case 'over_budget': return 'bg-red-500'
    case 'warning': return 'bg-yellow-500'
    case 'caution': return 'bg-orange-500'
    default: return 'bg-green-500'
  }
}

function getStatusText(status: string) {
  switch (status) {
    case 'over_budget': return 'Over Budget'
    case 'warning': return 'Warning'
    case 'caution': return 'Caution'
    default: return 'On Track'
  }
}

function getStatusTextClass(status: string) {
  switch (status) {
    case 'over_budget': return 'text-red-600 dark:text-red-400'
    case 'warning': return 'text-yellow-600 dark:text-yellow-400'
    case 'caution': return 'text-orange-600 dark:text-orange-400'
    default: return 'text-green-600 dark:text-green-400'
  }
}

function getPriorityColor(priority: number | null) {
  const colors = ['', 'text-red-500', 'text-orange-500', 'text-yellow-500', 'text-blue-500', 'text-green-500']
  return priority ? colors[priority] || '' : ''
}

function getGoalProgress(goal: any) {
  if (!goal.target_amount) return 0
  return Math.min((goal.current_amount / goal.target_amount) * 100, 100)
}

async function submitCreateBudget() {
  if (!createForm.category_id || createForm.amount <= 0 || !createForm.start_date || !createForm.end_date) {
    toast.add({ title: "Error", description: "Please fill all fields", color: "error" })
    return
  }

  console.log(createForm.start_date.toString())

  await budgetStore.submitForm({ ...createForm, start_date: createForm.start_date.toString(), end_date: createForm.end_date.toString() })
  await budgetStore.fetchBudgetAnalytics()

  Object.assign(createForm, {
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
  showCreateBudget.value = false
}

function startEdit(budget: any) {
  editingBudgetId.value = budget.budget_id
  Object.assign(editForm, {
    budget_id: budget.budget_id,
    category_id: budget.category_id,
    amount: budget.amount,
    start_date: budget.start_date.slice(0, 10),
    end_date: budget.end_date.slice(0, 10),
  })
}

function cancelEdit() {
  editingBudgetId.value = null
  Object.assign(editForm, {
    budget_id: '',
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
}

async function submitEdit() {
  if (!editForm.amount || editForm.amount <= 0) {
    toast.add({ title: "Error", description: "Please enter a valid amount", color: "error" })
    return
  }

  await budgetStore.editBudget({
    budget_id: editForm.budget_id,
    category_id: editForm.category_id,
    Category: { budget_category_id: editForm.category_id, name: '' },
    amount: editForm.amount,
    start_date: editForm.start_date,
    end_date: editForm.end_date,
  } as any)
  await budgetStore.fetchBudgetAnalytics()

  cancelEdit()
}

async function deleteBudgetItem(budgetId: string) {
  if (confirm('Delete this budget?')) {
    await budgetStore.deleteBudget(budgetId)
    await budgetStore.fetchBudgetAnalytics()
  }
}

async function handleGoalContribute(goal: any) {
  const monthlyTarget = goal.target_amount && goal.target_date
    ? Math.ceil((goal.target_amount - goal.current_amount) / Math.max(1, getMonthsUntil(goal.target_date)))
    : 50

  const amount = prompt(`Contribute to "${goal.title}"`, monthlyTarget.toString())
  if (!amount) return

  const numAmount = parseFloat(amount)
  if (isNaN(numAmount) || numAmount <= 0) {
    toast.add({ title: "Error", description: "Please enter a valid amount", color: "error" })
    return
  }

  const remaining = goal.target_amount ? goal.target_amount - goal.current_amount : numAmount
  if (numAmount > remaining) {
    toast.add({ title: "Error", description: `Amount exceeds remaining target of $${remaining.toFixed(2)}`, color: "error" })
    return
  }

  await financialGoalStore.addProgress(goal.goal_id, numAmount)
}

function getMonthsUntil(dateStr: string) {
  const target = new Date(dateStr)
  const now = new Date()
  return Math.max(1, (target.getFullYear() - now.getFullYear()) * 12 + target.getMonth() - now.getMonth())
}

function changeMonth(delta: number) {
  const [year, month] = selectedMonth.value.split('-').map(Number)
  const d = new Date(year, month - 1 + delta, 1)
  selectedMonth.value = d.toISOString().slice(0, 7)
}

function formatMonth(month: string) {
  const [year, m] = month.split('-').map(Number)
  return new Date(year, m - 1).toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}

const handleCategorySelect = (category: any) => {
  createForm.category_id = category.budget_category_id
}

const handleCategoryCreate = async (categoryName: string) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    const newCategory = categoryStore.categories.find(
      c => c.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) {
      createForm.category_id = newCategory.budget_category_id
    }
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

onMounted(async () => {
  if (incomeStore.incomes.length === 0) incomeStore.fetchIncomes()
  if (categoryStore.categories.length === 0) categoryStore.fetchCategory()
  if (budgetStore.budgets.length === 0) budgetStore.fetchBudget()
  if (financialGoalStore.goals.length === 0) financialGoalStore.fetchGoals()
  await budgetStore.fetchBudgetAnalytics()
})
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">Monthly Budget Plan</h2>
      <UButton :label="showCreateBudget ? 'Cancel' : '+ Add Budget'" @click="showCreateBudget = !showCreateBudget" />
    </div>

    <!-- Month Selector -->
    <div class="flex items-center justify-center gap-4">
      <UButton icon="i-lucide-chevron-left" variant="soft" color="neutral" @click="changeMonth(-1)" />
      <span class="text-xl font-semibold  min-w-48 text-center">
        {{ selectedMonth }}
      </span>
      <UButton icon="i-lucide-chevron-right" variant="soft" color="neutral" @click="changeMonth(1)" />
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <div
        class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
        <p class="text-sm text-text-light/60 dark:text-text-dark/60">Income</p>
        <p class="text-xl font-bold text-green-600 dark:text-green-400">${{ totalIncome?.toFixed(2) }}</p>
      </div>
      <div
        class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
        <p class="text-sm text-text-light/60 dark:text-text-dark/60">Budgeted</p>
        <p class="text-xl font-bold text-primary">${{ totalBudgeted?.toFixed(2) }}</p>
      </div>
      <div
        class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
        <p class="text-sm text-text-light/60 dark:text-text-dark/60">Spent</p>
        <p class="text-xl font-bold text-text-light dark:text-text-dark">${{ totalSpent?.toFixed(2) }}</p>
      </div>
      <div :class="[
        'p-4 rounded-lg border',
        remaining >= 0
          ? 'bg-surface-light dark:bg-surface-dark border-surface-light dark:border-surface-dark'
          : 'bg-red-50 dark:bg-red-900/20 border-red-300 dark:border-red-700',
      ]">
        <p class="text-sm text-text-light/60 dark:text-text-dark/60">Remaining</p>
        <p
          :class="['text-xl font-bold', remaining >= 0 ? 'text-text-light dark:text-text-dark' : 'text-red-600 dark:text-red-400']">
          ${{ remaining?.toFixed(2) }}
        </p>
      </div>
    </div>

    <!-- Create Budget Form -->
    <div v-if="showCreateBudget"
      class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
      <h3 class="text-lg font-semibold mb-3 text-text-light dark:text-text-dark">Create Budget</h3>
      <UForm :state="createForm" @submit="submitCreateBudget" class="space-y-3">
        <UFormField label="Category">
          <USelectMenu v-model="createForm.category_id" value-key="id" :items="categories" placeholder="Select category"
            create-item @create="handleCategoryCreate" />
        </UFormField>
        <UFormField label="Amount">
          <UInput v-model="createForm.amount" type="number" />
        </UFormField>
        <UFormField label="Start Date">
          <UInputDate v-model="createForm.start_date" />
        </UFormField>
        <UFormField label="End Date">
          <UInputDate v-model="createForm.end_date" />
        </UFormField>
        <UButton label="Create Budget" type="submit" />
      </UForm>
    </div>

    <!-- Budget Cards -->
    <div v-if="monthBudgets?.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
      <p class="text-lg mb-2">No budgets planned for {{ formatMonth(selectedMonth) }}</p>
      <p class="text-sm">Click "+ Add Budget" to start planning</p>
    </div>

    <div class="space-y-4">
      <div v-for="budget in monthBudgets" :key="budget.budget_id"
        class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
        <!-- Budget Header -->
        <div class="flex items-center justify-between mb-2">
          <div class="flex-1">
            <h3 class="font-semibold text-text-light dark:text-text-dark">{{ budget.Category.name }}</h3>
            <p class="text-xs text-text-light/60 dark:text-text-dark/60">
              {{ formatDate(budget.start_date) }} - {{ formatDate(budget.end_date) }}
            </p>
          </div>
          <div class="text-right ml-4">
            <template v-if="editingBudgetId === budget.budget_id">
              <div class="flex items-center gap-2">
                <input v-model.number="editForm.amount" type="number" step="0.01" min="0.01"
                  class="w-24 px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded focus:outline-none focus:ring-1 focus:ring-primary" />
                <button @click="submitEdit" class="text-green-600 hover:text-green-800 p-1"
                  :disabled="budgetStore.updating">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
                <button @click="cancelEdit" class="text-gray-500 hover:text-gray-700 p-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </template>
            <template v-else>
              <div class="flex items-center gap-2">
                <div class="text-right">
                  <p class="text-lg font-bold text-primary">${{ budget.amount.toFixed(2) }}</p>
                  <p v-if="budget.analytics" class="text-xs text-text-light/60 dark:text-text-dark/60">
                    Spent: ${{ budget.analytics.spent_amount.toFixed(2) }}
                  </p>
                </div>
                <div class="flex gap-1">
                  <UTooltip text="Edit Budget">
                    <UButton icon="i-lucide-square-pen" @click="startEdit(budget)" variant="outline" size="xs" />
                  </UTooltip>
                  <UTooltip text="Delete Budget">
                    <UButton icon="i-lucide-trash-2" @click="deleteBudgetItem(budget.budget_id)" variant="outline"
                      color="error" size="xs" />
                  </UTooltip>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- Progress Bar -->
        <div v-if="budget.analytics" class="mb-3">
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
            <div :class="['h-2.5 rounded-full transition-all duration-300', getStatusClass(budget.analytics.status)]"
              :style="{ width: getProgressPercent(budget.analytics) + '%' }"></div>
          </div>
          <div class="flex justify-between text-xs mt-1">
            <span :class="getStatusTextClass(budget.analytics.status)">
              {{ getStatusText(budget.analytics.status) }} - {{ getProgressPercent(budget.analytics).toFixed(1) }}%
            </span>
            <span class="text-text-light/60 dark:text-text-dark/60">
              {{ budget.analytics.days_remaining }}d left
            </span>
          </div>
        </div>

        <!-- Linked Goals -->
        <div v-if="budget.linkedGoals.length > 0"
          class="mt-3 pt-3 border-t border-surface-light dark:border-surface-dark">
          <p class="text-sm font-medium text-text-light dark:text-text-dark mb-2">Linked Goals</p>
          <div v-for="goal in budget.linkedGoals" :key="goal.goal_id"
            class="flex items-center justify-between py-1 text-sm">
            <div class="flex items-center gap-2">
              <span :class="['w-2 h-2 rounded-full', getPriorityColor(goal.priority)]"
                :style="{ backgroundColor: goal.priority ? undefined : goal.color }"></span>
              <span class="text-text-light dark:text-text-dark">{{ goal.title }}</span>
              <span v-if="goal.target_amount" class="text-xs text-text-light/60 dark:text-text-dark/60">
                ${{ goal.current_amount.toFixed(0) }}/${{ goal.target_amount.toFixed(0) }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <div v-if="goal.target_amount" class="w-16 bg-gray-200 dark:bg-gray-700 rounded-full h-1.5">
                <div class="h-1.5 rounded-full bg-primary" :style="{ width: getGoalProgress(goal) + '%' }"></div>
              </div>
              <UButton v-if="goal.type === 'savings'" label="Contribute" @click="handleGoalContribute(goal)" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Unallocated Goals -->
    <div v-if="financialGoalStore?.goals.filter(g => !g.category_id && g.status === 'active').length > 0"
      class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
      <h3 class="text-lg font-semibold mb-3 text-text-light dark:text-text-dark">Unallocated Goals</h3>
      <p class="text-sm text-text-light/60 dark:text-text-dark/60 mb-2">
        These goals are not linked to any budget category. Edit them to link to a category for better tracking.
      </p>
      <div v-for="goal in financialGoalStore.goals.filter(g => !g.category_id && g.status === 'active')"
        :key="goal.goal_id"
        class="flex items-center justify-between py-2 border-b border-surface-light/20 last:border-0">
        <span class="text-text-light dark:text-text-dark">{{ goal.title }}</span>
        <span class="text-xs px-2 py-1 rounded-full bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100">
          {{ goal.type }}
        </span>
      </div>
    </div>
  </div>
</template>
