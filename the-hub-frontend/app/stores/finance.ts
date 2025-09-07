interface Category {
  budget_category_id: number
  name: string
}

interface Budget {
  budget_id: number
  category_id: number
  amount: number
  start_date: Date
  end_date: Date
}

export interface CategoryResponse {
  categories: Category[]
}

export interface BudgetResponse {
  budgets: Budget[]
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
      budget_category_id: -Date.now(), // Use negative ID to avoid conflicts
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
      addToast("Category not added", "error")
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

  async function deleteCategory(id: number) {
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
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchBudget() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedBudget = await $api<BudgetResponse>('budgets')

    if (fetchedBudget) budgets.value = fetchedBudget.budgets

    loading.value = false
  }


  async function submitForm(payload: Budget) {
    // Create optimistic budget
    const optimisticBudget: Budget = {
      budget_id: -Date.now(), // Use negative ID to avoid conflicts
      category_id: payload.category_id,
      amount: payload.amount,
      start_date: payload.start_date,
      end_date: payload.end_date
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
        budgets.value[optimisticIndex] = data
      }

      incomeStore.fetchIncomes()
      addToast("Budget added succesfully", "success")

    } catch (err) {
      // Remove optimistic budget on error
      budgets.value = budgets.value.filter(b => b.budget_id !== optimisticBudget.budget_id)
      addToast("Budget not added", "error")
    }
  }

  async function editBudget(payload: Budget) {
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

      addToast("Edited budget succesfully", "success")
      incomeStore.fetchIncomes()
    } catch (err) {
      // Revert optimistic update on error
      if (originalBudget && originalBudgetIndex !== -1) {
        budgets.value[originalBudgetIndex] = originalBudget
      }
      addToast("Editing budget failed", "error")
    }
  }

  async function deleteBudget(budgetID: number, incomeID: number) {
    // Store the budget for potential rollback
    const budgetToDelete = budgets.value.find(b => b.budget_id === budgetID)
    if (!budgetToDelete) {
      addToast("Budget not found", "error")
      return
    }

    // Optimistically remove from local state
    budgets.value = budgets.value.filter((b) => b.budget_id !== budgetID)

    try {
      const { $api } = useNuxtApp()
      await $api(`budgets/${budgetID}/${incomeID}`, {
        method: 'DELETE'
      })

      incomeStore.fetchIncomes()
      addToast("Budget deleted succesfully", "success")
    } catch(err) {
      // Restore the budget on error
      budgets.value.push(budgetToDelete)
      addToast("Budget did not delete ", "error")
    }
  }

  function reset() {
    budgets.value = []
  }

  return {
    budgets,
    loading,
    fetchError,
    fetchBudget,
    editBudget,
    deleteBudget,
    submitForm,
    reset,
  }
})
