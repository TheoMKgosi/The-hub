<script setup lang="ts">
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import FormBudget from '@/components/finance/FormBudget.vue'

const incomeStore = useIncomeStore()
const categoryStore = useCategoryStore()
const budgetStore = useBudgetStore()

const activeIncomeId = ref<number | null>(null)
const showDialog = ref(false)
const showIncomeModal = ref(true)
const searchQuery = ref('')

const budgetID = ref(0)
const incomeID = ref(0)

const formData = reactive({
  source: '',
  amount: 0,
  received_at: null
})

const filteredIncome = computed(() => {
  let result = incomeStore.incomes

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(income =>
      income.source.toLowerCase().includes(query)
    )
  }

  return result
})

const deleteItem = (id: number, incomeID: number) => {
  budgetStore.deleteBudget(id, incomeID)
}

const submitForm = async () => {
  const dataToSend = { ...formData }
  incomeStore.submitForm(dataToSend)
  Object.assign(formData, {
    source: '',
    amount: 0,
    received_at: null
  })
  showIncomeModal.value = true
}

onMounted(() => {
  if (incomeStore.incomes.length === 0) {
    incomeStore.fetchIncomes()
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})

const formatDate = (date) => new Date(date).toLocaleDateString()
const openForm = (id: number) => {
  activeIncomeId.value = id
}

const closeForm = () => {
  activeIncomeId.value = null
}

const remainingAmount = (amount, budgets) => {
  let remaining = amount
  for (const budget of budgets) {
    remaining -= budget.amount
  }
  return remaining
}
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Income Management</h2>

    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search income sources..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <!-- Floating Action Button -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showIncomeModal" @click="showIncomeModal = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
          <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
            <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Income Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="!showIncomeModal" @click="showIncomeModal = true" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50">
          <div class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark" @click.stop>

            <!-- Modal Header -->
            <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
              <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Add New Income</h2>
              <UiButton @click="showIncomeModal = true" variant="default" size="sm" class="p-2">
                ×
              </UiButton>
            </div>

            <!-- Modal Body -->
            <div class="p-6">
              <form @submit.prevent="submitForm" class="space-y-4">

                <div>
                  <label for="source" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Income Source</label>
                  <input type="text" id="source" v-model="formData.source" placeholder="e.g., Salary, Freelance, Business"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <div>
                  <label for="amount" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
                  <input type="number" id="amount" v-model="formData.amount" placeholder="0.00" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <div>
                  <label for="received" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Received Date</label>
                  <input type="date" id="received" v-model="formData.received_at"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <!-- Modal Footer -->
                <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                  <UiButton type="button" @click="showIncomeModal = true" variant="default" size="md" class="w-full sm:w-auto">
                    Cancel
                  </UiButton>
                  <UiButton type="submit" variant="primary" size="md" class="w-full sm:w-auto">
                    Create Income
                  </UiButton>
                </div>

              </form>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <p class="text-sm text-text-light dark:text-text-dark/60 text-center">Double-click a budget to delete it</p>

    <!-- Income Cards -->
    <div class="space-y-4">
      <div v-if="incomeStore.incomes.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No income sources added yet</p>
        <p class="text-sm">Create your first income source above to get started</p>
      </div>

      <div v-for="income in filteredIncome" :key="income.income_id"
        class="p-6 rounded-lg shadow-md bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
        <!-- Income Header -->
        <div class="flex justify-between items-start mb-4">
          <div class="flex-1">
            <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-1">{{ income.source }}</h3>
            <p class="text-sm text-text-light dark:text-text-dark/60">{{ formatDate(income.received_at) }}</p>
          </div>
          <p class="text-xl font-bold text-success">${{ income.amount.toFixed(2) }}</p>
        </div>

        <!-- Budgets -->
        <div class="space-y-3">
          <h4 class="font-medium text-text-light dark:text-text-dark">Allocated Budgets</h4>
          <div v-if="income.budgets.length === 0" class="text-sm text-text-light dark:text-text-dark/60 italic">
            No budgets created yet
          </div>
          <div v-else v-for="budget in income.budgets" :key="budget.budget_id"
            class="flex justify-between items-center p-3 rounded-md bg-surface-light/50 dark:bg-surface-dark/50 border border-surface-light dark:border-surface-dark hover:bg-red-50 dark:hover:bg-red-900/20 hover:border-red-200 dark:hover:border-red-800 transition-colors cursor-pointer"
            @dblclick="showDialog = true; budgetID = budget.budget_id; incomeID = income.income_id;">
            <div>
              <p class="font-medium text-text-light dark:text-text-dark">{{ budget.Category.name }}</p>
              <p class="text-xs text-text-light dark:text-text-dark/60">{{ formatDate(budget.start_date) }} - {{ formatDate(budget.end_date) }}</p>
            </div>
            <p class="font-semibold text-text-light dark:text-text-dark">P{{ budget.amount.toFixed(2) }}</p>
            <ConfirmDialog v-model:show="showDialog" :message="`Delete budget for ${budget.Category.name}?`"
              @confirm="deleteItem(budgetID, incomeID)" />
          </div>
        </div>

        <hr class="my-4 border-surface-light dark:border-surface-dark" />

        <!-- Remaining -->
        <div class="flex justify-between items-center font-medium">
          <p class="text-text-light dark:text-text-dark">Remaining:</p>
          <p class="text-lg" :class="remainingAmount(income.amount, income.budgets) >= 0 ? 'text-success' : 'text-red-500 dark:text-red-400'">
            ${{ remainingAmount(income.amount, income.budgets).toFixed(2) }}
          </p>
        </div>

        <!-- Budget Form Toggle -->
        <div v-if="activeIncomeId !== income.income_id" class="mt-4">
          <UiButton @click="openForm(income.income_id)" variant="default" size="sm" class="w-full">
            Create Budget for This Income
          </UiButton>
        </div>

        <!-- Use the new FormBudget component -->
        <FormBudget v-if="activeIncomeId === income.income_id" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.25s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.fade-scale-enter-to {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-from {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
