interface Category {
  budget_category_id: string
  name: string
}

interface Budget {
  budget_id: string
  category_id: string
  Category: Category
  amount: number
  start_date: string
  end_date: string
  income_id?: string
}

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

export interface CategoryResponse {
  categories: Category[]
}

export interface BudgetResponse {
  budgets: Budget[]
}

export interface CreateBudgetRequest {
  category_id: string
  amount: number
  start_date: string
  end_date: string
  income_id?: string
}


export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()
  const { validateObject, schemas } = useValidation()

  async function fetchCategory() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedCategories = await $api<CategoryResponse>('categories')

    if (fetchedCategories) categories.value = fetchedCategories.categories

    loading.value = false
  }


  async function submitForm(payload: {name: string;}) {
    const validation = validateObject(payload, schemas.category.create)

    if (!validation.isValid) {
      const errorMessage = Object.values(validation.errors)[0]
      addToast(errorMessage, "error")
      return
    }

    // Create optimistic category
    const optimisticCategory: Category = {
      budget_category_id: `temp-${Date.now()}`, // Use temp ID to avoid conflicts
      name: payload.name
    }

    // Optimistically add to local state
    categories.value.push(optimisticCategory)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Category>('categories', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic category with real data
      const optimisticIndex = categories.value.findIndex(c => c.budget_category_id === optimisticCategory.budget_category_id)
      if (optimisticIndex !== -1) {
        categories.value[optimisticIndex] = data
      }

      addToast("Category added succesfully", "success")
    } catch (err) {
      // Remove optimistic category on error
      categories.value = categories.value.filter(c => c.budget_category_id !== optimisticCategory.budget_category_id)
      addToast("Failed to add category. Please try again.", "error")
    }
  }

  async function editCategory(payload: Category) {
    // Store original category for potential rollback
    const originalCategoryIndex = categories.value.findIndex(c => c.budget_category_id === payload.budget_category_id)
    const originalCategory = originalCategoryIndex !== -1 ? { ...categories.value[originalCategoryIndex] } : null

    // Optimistically update the category
    if (originalCategoryIndex !== -1) {
      categories.value[originalCategoryIndex] = { ...categories.value[originalCategoryIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Category>(`categories/${payload.budget_category_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalCategoryIndex !== -1 && data) {
        categories.value[originalCategoryIndex] = data
      }

      addToast("Edited category succesfully", "success")
    } catch (err) {
      // Revert optimistic update on error
      if (originalCategory && originalCategoryIndex !== -1) {
        categories.value[originalCategoryIndex] = originalCategory
      }
      addToast("Editing category failed", "error")
    }
  }

  async function deleteCategory(id: string) {
    // Store the category for potential rollback
    const categoryToDelete = categories.value.find(c => c.budget_category_id === id)
    if (!categoryToDelete) {
      addToast("Category not found", "error")
      return
    }

    // Optimistically remove from local state
    categories.value = categories.value.filter((c) => c.budget_category_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`categories/${id}`, {
        method: 'DELETE'
      })

      addToast("Category deleted succesfully", "success")
    } catch (err) {
      // Restore the category on error
      categories.value.push(categoryToDelete)
      addToast("Category did not delete", "error")
    }
  }

  function reset() {
    categories.value = []
  }

  return {
    categories,
    loading,
    fetchError,
    fetchCategory,
    editCategory,
    deleteCategory,
    submitForm,
    reset
  }
})


export const useBudgetStore = defineStore('budget', () => {
  const incomeStore = useIncomeStore()
  const budgets = ref<Budget[]>([])
  const loading = ref(false)
  const creating = ref(false)
  const updating = ref(false)
  const deleting = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  // Analytics and alerts state
  const analytics = ref<BudgetAnalytics[]>([])
  const alerts = ref<BudgetAlert[]>([])
  const suggestions = ref<BudgetSuggestion[]>([])
  const analyticsLoading = ref(false)
  const alertsLoading = ref(false)
  const suggestionsLoading = ref(false)

  async function fetchBudget() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedBudget = await $api<BudgetResponse>('budgets')

    if (fetchedBudget) budgets.value = fetchedBudget.budgets

    loading.value = false
  }


  async function submitForm(payload: CreateBudgetRequest) {
    creating.value = true

    // Create optimistic budget
    const optimisticBudget: Budget = {
      budget_id: `temp-${Date.now()}`, // Use temp ID to avoid conflicts
      category_id: payload.category_id,
      Category: categories.value.find(c => c.budget_category_id === payload.category_id) || { budget_category_id: payload.category_id, name: 'Loading...' },
      amount: payload.amount,
      start_date: payload.start_date,
      end_date: payload.end_date,
      income_id: payload.income_id
    }

    // Optimistically add to local state
    budgets.value.push(optimisticBudget)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Budget>('budgets', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic budget with real data
      const optimisticIndex = budgets.value.findIndex(b => b.budget_id === optimisticBudget.budget_id)
      if (optimisticIndex !== -1) {
        budgets.value[optimisticIndex] = {
          ...data,
          Category: categories.value.find(c => c.budget_category_id === data.category_id) || data.Category
        }
      }

      incomeStore.fetchIncomes()
      addToast("Budget added successfully", "success")

    } catch (err) {
      // Remove optimistic budget on error
      budgets.value = budgets.value.filter(b => b.budget_id !== optimisticBudget.budget_id)
      addToast("Failed to add budget. Please try again.", "error")
    } finally {
      creating.value = false
    }
  }

  async function editBudget(payload: Budget) {
    updating.value = true

    // Store original budget for potential rollback
    const originalBudgetIndex = budgets.value.findIndex(b => b.budget_id === payload.budget_id)
    const originalBudget = originalBudgetIndex !== -1 ? { ...budgets.value[originalBudgetIndex] } : null

    // Optimistically update the budget
    if (originalBudgetIndex !== -1) {
      budgets.value[originalBudgetIndex] = { ...budgets.value[originalBudgetIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Budget>(`budgets/${payload.budget_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalBudgetIndex !== -1 && data) {
        budgets.value[originalBudgetIndex] = data
      }

      addToast("Budget updated successfully", "success")
      incomeStore.fetchIncomes()
    } catch (err) {
      // Revert optimistic update on error
      if (originalBudget && originalBudgetIndex !== -1) {
        budgets.value[originalBudgetIndex] = originalBudget
      }
      addToast("Failed to update budget. Please try again.", "error")
    } finally {
      updating.value = false
    }
  }

  async function deleteBudget(budgetID: string) {
    deleting.value = true

    // Store the budget for potential rollback
    const budgetToDelete = budgets.value.find(b => b.budget_id === budgetID)
    if (!budgetToDelete) {
      addToast("Budget not found", "error")
      deleting.value = false
      return
    }

    // Optimistically remove from local state
    budgets.value = budgets.value.filter((b) => b.budget_id !== budgetID)

    try {
      const { $api } = useNuxtApp()
      await $api(`budgets/${budgetID}`, {
        method: 'DELETE'
      })

      incomeStore.fetchIncomes()
      addToast("Budget deleted successfully", "success")
    } catch(err) {
      // Restore the budget on error
      budgets.value.push(budgetToDelete)
      addToast("Failed to delete budget. Please try again.", "error")
    } finally {
      deleting.value = false
    }
  }

  async function fetchBudgetAnalytics(period: string = 'current') {
    analyticsLoading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<{ analytics: BudgetAnalytics[] }>(`budgets/analytics?period=${period}`)
      if (response) {
        analytics.value = response.analytics
      }
    } catch (error) {
      addToast('Failed to load budget analytics', 'error')
    } finally {
      analyticsLoading.value = false
    }
  }

  async function fetchBudgetAlerts() {
    alertsLoading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<{ alerts: BudgetAlert[] }>('budgets/alerts')
      if (response) {
        alerts.value = response.alerts
      }
    } catch (error) {
      addToast('Failed to load budget alerts', 'error')
    } finally {
      alertsLoading.value = false
    }
  }

  async function fetchBudgetSuggestions() {
    suggestionsLoading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<{
        suggestions: {
          categories: BudgetSuggestion[]
          analysis_period: string
          total_suggestions: number
        }
      }>('budgets/suggestions')
      if (response) {
        suggestions.value = response.suggestions.categories
        addToast(`Loaded ${response.suggestions.total_suggestions} budget suggestions`, 'success')
      }
    } catch (error) {
      addToast('Failed to load budget suggestions', 'error')
    } finally {
      suggestionsLoading.value = false
    }
  }

  function reset() {
    budgets.value = []
    analytics.value = []
    alerts.value = []
    suggestions.value = []
  }

  return {
    budgets,
    loading,
    creating,
    updating,
    deleting,
    fetchError,
    fetchBudget,
    editBudget,
    deleteBudget,
    submitForm,
    reset,
    // Analytics and alerts
    analytics,
    alerts,
    suggestions,
    analyticsLoading,
    alertsLoading,
    suggestionsLoading,
    fetchBudgetAnalytics,
    fetchBudgetAlerts,
    fetchBudgetSuggestions,
  }
})
