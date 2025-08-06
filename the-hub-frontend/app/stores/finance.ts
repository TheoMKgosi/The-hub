import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'
import { useIncomeStore } from "./income.ts";



interface Category {
  category_id: number
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
    loading.value = true
    const { data, error } = await useMyFetch('categories').json<CategoryResponse>()

    if (data.value) categories.value = data.value.categories
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: Category) {
    const { data, error } = await useMyFetch('categories').post(formData).json()
    if (!data.value.category_id) {
      data.value.category_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Category not added", "error")
    } else {
      categories.value.push(data.value)
      addToast("Category added succesfully", "success")
    }
  }

  async function editCategory(category: Category) {
    const { error } = await useMyFetch(`categories/${category.category_id}`).patch(category).json()

    if (!error.value) {
      const index = categories.value.findIndex(c => c.category_id === category.category_id)
      if (index !== -1) {
        categories.value[index] = { ...categories.value[index], ...category }
        addToast("Edited category succesfully", "success")
      } else {
        addToast("Editing category failed", "error")
      }
    } else {
      addToast("Editing category failed", "error")
    }
  }

  async function deleteCategory(id: number) {
    await useMyFetch(`categories/${id}`).delete().json()
    categories.value = categories.value.filter((c) => c.category_id !== id)
    addToast("Category deleted succesfully", "success")
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
    loading.value = true
    const { data, error } = await useMyFetch('budgets').json<BudgetResponse>()

    if (data.value) budgets.value = data.value.budgets
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: Budget) {
    const { data, error } = await useMyFetch('budgets').post(formData).json()
    if (!data.value.budget_id) {
      data.value.budget_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Budget not added", "error")
    } else {
      budgets.value.push(data.value)
      addToast("Budget added succesfully", "success")
    }
    incomeStore.fetchIncomes()
  }

  async function editBudget(budget: Budget) {
    const { error } = await useMyFetch(`budgets/${budget.budget_id}`).patch(budget).json()

    if (!error.value) {
      const index = budgets.value.findIndex(c => c.budget_id === budget.budget_id)
      if (index !== -1) {
        budgets.value[index] = { ...budgets.value[index], ...budget }
        addToast("Edited budget succesfully", "success")
      } else {
        addToast("Editing budget failed", "error")
      }
    } else {
      addToast("Editing budget failed", "error")
    }
  }

  async function deleteBudget(id: number) {
    await useMyFetch(`budgets/${id}`).delete().json()
    budgets.value = budgets.value.filter((c) => c.budget_id !== id)
    addToast("Budget deleted succesfully", "success")
    incomeStore.fetchIncomes()
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
