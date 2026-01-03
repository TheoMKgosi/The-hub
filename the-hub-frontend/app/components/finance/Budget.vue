<script setup>
const formData = reactive({
  category_id: '',
  amount: 0,
  start_date: '',
  end_date: '',
})

const editFormData = reactive({
  budget_id: '',
  category_id: '',
  amount: 0,
  start_date: '',
  end_date: '',
})

const categoryStore = useCategoryStore()
const budgetStore = useBudgetStore()
const transactionStore = useTransactionStore()
const { addToast } = useToast()

const isEditing = ref(false)
const editingBudgetId = ref('')

// Computed properties for budget analytics
const getBudgetAnalytics = (budget) => {
  const analytics = budgetStore.analytics.find(a => a.budget_id === budget.budget_id)
  if (!analytics) return null

  return {
    ...analytics,
    utilizationPercentage: Math.round(analytics.utilization_rate * 100),
    statusColor: getStatusColor(analytics.status),
    statusText: getStatusText(analytics.status),
    projectedOverspend: calculateProjectedOverspend(analytics),
    daysRemaining: analytics.days_remaining,
    dailySpendingRate: analytics.daily_spending_rate
  }
}

const getStatusColor = (status) => {
  switch (status?.toLowerCase()) {
    case 'on_track': return 'green'
    case 'warning': return 'yellow'
    case 'over_budget': return 'red'
    default: return 'gray'
  }
}

const getStatusText = (status) => {
  switch (status?.toLowerCase()) {
    case 'on_track': return 'On Track'
    case 'warning': return 'Warning'
    case 'over_budget': return 'Over Budget'
    default: return 'Unknown'
  }
}

const calculateProjectedOverspend = (analytics) => {
  if (!analytics || analytics.days_remaining <= 0) return 0
  const projectedTotal = analytics.spent_amount + (analytics.daily_spending_rate * analytics.days_remaining)
  return Math.max(0, projectedTotal - analytics.budget_amount)
}

const submitForm = async () => {
  // Validation
  if (!formData.category_id) {
    addToast("Please select a category", "error")
    return
  }

  if (!formData.amount || formData.amount <= 0) {
    addToast("Please enter a valid amount greater than 0", "error")
    return
  }

  if (!formData.start_date) {
    addToast("Please select a start date", "error")
    return
  }

  if (!formData.end_date) {
    addToast("Please select an end date", "error")
    return
  }

  const startDate = new Date(formData.start_date)
  const endDate = new Date(formData.end_date)

  if (endDate <= startDate) {
    addToast("End date must be after start date", "error")
    return
  }

  const dataToSend = { ...formData }
  budgetStore.submitForm(dataToSend)
  Object.assign(formData, {
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
}

const startEdit = (budget) => {
  isEditing.value = true
  editingBudgetId.value = budget.budget_id
  Object.assign(editFormData, {
    budget_id: budget.budget_id,
    category_id: budget.category_id,
    amount: budget.amount,
    start_date: budget.start_date,
    end_date: budget.end_date,
  })
}

const cancelEdit = () => {
  isEditing.value = false
  editingBudgetId.value = ''
  Object.assign(editFormData, {
    budget_id: '',
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
}

const submitEdit = async () => {
  // Validation
  if (!editFormData.category_id) {
    addToast("Please select a category", "error")
    return
  }

  if (!editFormData.amount || editFormData.amount <= 0) {
    addToast("Please enter a valid amount greater than 0", "error")
    return
  }

  if (!editFormData.start_date) {
    addToast("Please select a start date", "error")
    return
  }

  if (!editFormData.end_date) {
    addToast("Please select an end date", "error")
    return
  }

  const startDate = new Date(editFormData.start_date)
  const endDate = new Date(editFormData.end_date)

  if (endDate <= startDate) {
    addToast("End date must be after start date", "error")
    return
  }

  try {
    await budgetStore.editBudget(editFormData)
    cancelEdit()
  } catch (error) {
    addToast("Failed to update budget", "error")
  }
}

const deleteBudget = async (budgetId) => {
  if (confirm('Are you sure you want to delete this budget?')) {
    try {
      await budgetStore.deleteBudget(budgetId)
    } catch (error) {
      addToast("Failed to delete budget", "error")
    }
  }
}

onMounted(async () => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
  if (budgetStore.budgets.length === 0) {
    budgetStore.fetchBudget()
  }
  if (transactionStore.transactions.length === 0) {
    transactionStore.fetchTransactions()
  }

  // Fetch budget analytics
  await budgetStore.fetchBudgetAnalytics()
  await budgetStore.fetchBudgetAlerts()
})

const formatDate = (date) => new Date(date).toLocaleDateString()

const handleCategorySelect = (category) => {
  formData.category_id = category.budget_category_id
}

const handleCategoryCreate = async (categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(cat =>
      cat.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) {
      formData.category_id = newCategory.budget_category_id
    }
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

const handleEditCategorySelect = (category) => {
  editFormData.category_id = category.budget_category_id
}

const handleEditCategoryCreate = async (categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(cat =>
      cat.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) {
      editFormData.category_id = newCategory.budget_category_id
    }
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 text-text-light dark:text-text-dark">Budget Management</h2>

    <!-- Budget Creation Form -->
    <form @submit.prevent="submitForm"
      class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Create New Budget</h3>
      <div class="space-y-4">
         <div>
           <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Category</label>
           <ComboBox
             :model-value="formData.category_id"
             :categories="categoryStore.categories"
             placeholder="Select or create category..."
             @select="handleCategorySelect"
             @create="handleCategoryCreate"
           />
         </div>

        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
          <input v-model.number="formData.amount" type="number" placeholder="Enter amount" step="0.01" min="0.01"
            required
            class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Start Date</label>
            <input v-model="formData.start_date" type="date" required
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">End Date</label>
            <input v-model="formData.end_date" type="date" required
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
        </div>

        <div class="flex justify-end">
          <UiButton type="submit" variant="primary" size="md" :disabled="budgetStore.creating">
            {{ budgetStore.creating ? 'Creating...' : 'Create Budget' }}
          </UiButton>
        </div>
      </div>
    </form>

    <!-- Edit Budget Form -->
    <div v-if="isEditing"
      class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Edit Budget</h3>
      <form @submit.prevent="submitEdit" class="space-y-4">
         <div>
           <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Category</label>
           <ComboBox
             :model-value="editFormData.category_id"
             :categories="categoryStore.categories"
             placeholder="Select or create category..."
             @select="handleEditCategorySelect"
             @create="handleEditCategoryCreate"
           />
         </div>

        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
          <input v-model.number="editFormData.amount" type="number" placeholder="Enter amount" step="0.01" min="0.01"
            required
            class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Start Date</label>
            <input v-model="editFormData.start_date" type="date" required
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">End Date</label>
            <input v-model="editFormData.end_date" type="date" required
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <UiButton type="button" variant="default" size="md" @click="cancelEdit" :disabled="budgetStore.updating">
            Cancel
          </UiButton>
          <UiButton type="submit" variant="primary" size="md" :disabled="budgetStore.updating">
            {{ budgetStore.updating ? 'Updating...' : 'Update Budget' }}
          </UiButton>
        </div>
      </form>
    </div>

     <!-- Budget Analytics Overview -->
     <div v-if="budgetStore.alerts.length > 0" class="mb-6">
       <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Budget Alerts</h3>
       <div class="space-y-2">
         <div v-for="alert in budgetStore.alerts" :key="alert.budget_id"
           :class="[
             'p-3 rounded-lg border-l-4',
             alert.alert_type === 'over_budget' ? 'bg-red-50 dark:bg-red-900/20 border-red-500' :
             alert.alert_type === 'danger' ? 'bg-orange-50 dark:bg-orange-900/20 border-orange-500' :
             'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-500'
           ]">
           <div class="flex items-center justify-between">
             <div>
               <p class="font-medium text-text-light dark:text-text-dark">{{ alert.category_name }}</p>
               <p class="text-sm text-text-light dark:text-text-dark/80">{{ alert.message }}</p>
             </div>
             <span :class="[
               'px-2 py-1 text-xs font-medium rounded-full',
               alert.alert_type === 'over_budget' ? 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100' :
               alert.alert_type === 'danger' ? 'bg-orange-100 text-orange-800 dark:bg-orange-800 dark:text-orange-100' :
               'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100'
             ]">
               {{ alert.alert_type.replace('_', ' ').toUpperCase() }}
             </span>
           </div>
         </div>
       </div>
     </div>

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
             <div class="text-right">
               <div class="text-lg font-bold text-primary">${{ budget.amount.toFixed(2) }}</div>
               <div v-if="getBudgetAnalytics(budget)" class="text-sm text-text-light dark:text-text-dark/60">
                 Spent: ${{ getBudgetAnalytics(budget).spent_amount.toFixed(2) }}
               </div>
             </div>
           </div>

           <!-- Progress Bar -->
           <div v-if="getBudgetAnalytics(budget)" class="mb-3">
             <div class="flex justify-between text-sm mb-1">
               <span class="text-text-light dark:text-text-dark/80">Progress</span>
               <span :class="[
                 'font-medium',
                 getBudgetAnalytics(budget).statusColor === 'red' ? 'text-red-600 dark:text-red-400' :
                 getBudgetAnalytics(budget).statusColor === 'yellow' ? 'text-yellow-600 dark:text-yellow-400' :
                 'text-green-600 dark:text-green-400'
               ]">
                 {{ getBudgetAnalytics(budget).utilizationPercentage }}%
               </span>
             </div>
             <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
               <div
                 :class="[
                   'h-2 rounded-full transition-all duration-300',
                   getBudgetAnalytics(budget).statusColor === 'red' ? 'bg-red-500' :
                   getBudgetAnalytics(budget).statusColor === 'yellow' ? 'bg-yellow-500' :
                   'bg-green-500'
                 ]"
                 :style="{ width: Math.min(getBudgetAnalytics(budget).utilizationPercentage, 100) + '%' }"
               ></div>
             </div>
           </div>

           <!-- Status Badge -->
           <div v-if="getBudgetAnalytics(budget)" class="mb-2">
             <span :class="[
               'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
               getBudgetAnalytics(budget).statusColor === 'red' ? 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100' :
               getBudgetAnalytics(budget).statusColor === 'yellow' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100' :
               'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100'
             ]">
               {{ getBudgetAnalytics(budget).statusText }}
             </span>
           </div>

           <!-- Analytics Info -->
           <div v-if="getBudgetAnalytics(budget)" class="space-y-1 mb-3 text-sm">
             <div class="flex justify-between">
               <span class="text-text-light dark:text-text-dark/60">Remaining:</span>
               <span class="font-medium text-text-light dark:text-text-dark">
                 ${{ getBudgetAnalytics(budget).remaining_amount.toFixed(2) }}
               </span>
             </div>
             <div class="flex justify-between">
               <span class="text-text-light dark:text-text-dark/60">Daily Rate:</span>
               <span class="font-medium text-text-light dark:text-text-dark">
                 ${{ getBudgetAnalytics(budget).dailySpendingRate.toFixed(2) }}
               </span>
             </div>
             <div class="flex justify-between">
               <span class="text-text-light dark:text-text-dark/60">Days Left:</span>
               <span class="font-medium text-text-light dark:text-text-dark">
                 {{ getBudgetAnalytics(budget).daysRemaining }}
               </span>
             </div>
             <div v-if="getBudgetAnalytics(budget).projectedOverspend > 0" class="flex justify-between">
               <span class="text-red-600 dark:text-red-400">Projected Overspend:</span>
               <span class="font-medium text-red-600 dark:text-red-400">
                 ${{ getBudgetAnalytics(budget).projectedOverspend.toFixed(2) }}
               </span>
             </div>
           </div>

           <p class="text-sm text-text-light dark:text-text-dark/60 mb-3">
             {{ formatDate(budget.start_date) }} → {{ formatDate(budget.end_date) }}
           </p>

           <div class="flex gap-2">
             <UiButton variant="default" size="sm" class="flex-1" @click="startEdit(budget)"
               :disabled="budgetStore.updating || budgetStore.deleting">Edit</UiButton>
             <UiButton variant="danger" size="sm" class="flex-1" @click="deleteBudget(budget.budget_id)"
               :disabled="budgetStore.updating || budgetStore.deleting">
               {{ budgetStore.deleting ? 'Deleting...' : 'Delete' }}
             </UiButton>
           </div>
         </div>
       </div>
     </div>
  </div>
</template>
