<script setup>
const formData = reactive({
  category_id: '',
  amount: 0,
  start_date: '',
  end_date: '',
})

const categoryStore = useCategoryStore()
const budgetStore = useBudgetStore()

const submitForm = async () => {
  const dataToSend = { ...formData }
  budgetStore.submitForm(dataToSend)
  Object.assign(formData, {
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
}

onMounted(() => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
  if (budgetStore.budgets.length === 0) {
    budgetStore.fetchBudget()
  }
})

const formatDate = (date) => new Date(date).toLocaleDateString()

</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 text-text-light dark:text-text-dark">Budget Management</h2>

    <!-- Budget Creation Form -->
    <form @submit.prevent="submitForm" class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Create New Budget</h3>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Category</label>
          <select v-model="formData.category_id" class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
            <option disabled value="">Select category</option>
            <option v-for="cat in categoryStore.categories" :key="cat.budget_category_id" :value="cat.budget_category_id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
          <input v-model.number="formData.amount" type="number" placeholder="Enter amount" step="0.01"
            class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Start Date</label>
            <input v-model="formData.start_date" type="date"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">End Date</label>
            <input v-model="formData.end_date" type="date"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
        </div>

        <div class="flex justify-end">
          <UiButton type="submit" variant="primary" size="md">
            Create Budget
          </UiButton>
        </div>
      </div>
    </form>

    <!-- Budget List -->
    <div v-if="budgetStore.budgets.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
      No budgets created yet. Create your first budget above.
    </div>
    <div v-else>
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Your Budgets</h3>
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="budget in budgetStore.budgets" :key="budget.budget_id"
          class="bg-surface-light dark:bg-surface-dark p-4 rounded-lg shadow-md border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
          <div class="flex items-center justify-between mb-2">
            <h4 class="font-semibold text-text-light dark:text-text-dark">{{ budget.Category.name }}</h4>
            <span class="text-lg font-bold text-primary">${{ budget.amount.toFixed(2) }}</span>
          </div>
          <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">
            Budget period
          </p>
          <p class="text-sm text-text-light dark:text-text-dark/60">
            {{ formatDate(budget.start_date) }} → {{ formatDate(budget.end_date) }}
          </p>
          <div class="mt-3 flex gap-2">
            <UiButton variant="default" size="sm" class="flex-1">Edit</UiButton>
            <UiButton variant="danger" size="sm" class="flex-1">Delete</UiButton>
          </div>
        </div>
      </div>
    </div>
  </div>
      -->

      <div class="grid gap-4">
        <div v-for="budget in budgetStore.budgets" :key="budget.budget_id" class="border border-gray-200 dark:border-gray-700 p-4 rounded shadow bg-surface-light dark:bg-surface-dark">
          <h3 class="font-bold text-text-light dark:text-text-dark">{{ budget.Category.name }}</h3>
          <p class="text-text-light dark:text-text-dark">Amount: ${{ budget.amount.toFixed(2) }}</p>
          <p class="text-sm text-gray-600 dark:text-gray-400">
            {{ formatDate(budget.start_date) }} → {{ formatDate(budget.end_date) }}
          </p>
        </div>
      </div>
    </div>
  </div>

</template>
