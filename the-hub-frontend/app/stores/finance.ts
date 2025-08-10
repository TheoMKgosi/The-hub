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

  async function fetchCategory() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedCategories = await $api<CategoryResponse>('categories')

    if (fetchedCategories) categories.value = fetchedCategories.categories

    loading.value = false
  }


  // TODO: Change to object
  async function submitForm(payload: Category) {
    try {
      const { $api } = useNuxtApp()
      await $api('categories', {
        method: 'POST',
        body: payload
      })
      fetchCategory()
      addToast("Category added succesfully", "success")
    } catch (err) {
      addToast("Category not added", "error")
    }
  }

  async function editCategory(payload: Category) {
    try {
      const { $api } = useNuxtApp()
      await $api(`categories/${payload.category_id}`, {
        method: 'PATCH',
        body: payload
      })
      fetchCategory()
      addToast("Edited category succesfully", "success")
    } catch (err) {
      addToast("Editing category failed", "error")
    }
  }

  async function deleteCategory(id: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`categories/${id}`, {
        method: 'DELETE'
      })
      categories.value = categories.value.filter((c) => c.category_id !== id)
      addToast("Category deleted succesfully", "success")
    } catch (err) {
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
    try {
      const { $api } = useNuxtApp()
      await $api('budgets', {
        method: 'POST',
        body: payload
      })
      incomeStore.fetchIncomes()
      addToast("Budget added succesfully", "success")

    } catch (err) {
      addToast("Budget not added", "error")
    }
  }

  async function editBudget(payload: Budget) {
    try {
      const { $api } = useNuxtApp()
      await $api(`budgets/${payload.budget_id}`, {
        method: 'PATCH',
        body: payload
      })

      addToast("Edited budget succesfully", "success")
      incomeStore.fetchIncomes()
    } catch (err) {
      addToast("Editing budget failed", "error")
    }
  }

  async function deleteBudget(budgetID: number, incomeID: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`budgets/${budgetID}/${incomeID}`, {
        method: 'DELETE'
      })
      incomeStore.fetchIncomes()
      addToast("Budget deleted succesfully", "success")
    } catch(err) {
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
