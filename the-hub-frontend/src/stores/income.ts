
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch'
import { useToast } from '@/composables/useToast'

interface Income {
  income_id: number
  source: string
  amount: number
  received_at: Date
}

export interface IncomeResponse {
  incomes: Income[]
}


export const useIncomeStore = defineStore('income', () => {
  const incomes = ref<Income[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchIncomes() {
    loading.value = true
    const { data, error } = await useMyFetch('incomes').json<IncomeResponse>()

    if (data.value) incomes.value = data.value.incomes
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: Income) {
    const { data, error } = await useMyFetch('incomes').post(formData).json()
    if (!data.value.income_id) {
      data.value.income_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Income not added", "error")
    } else {
      incomes.value.push(data.value)
      addToast("Income added succesfully", "success")
    }
  }

  async function editIncome(income: Income) {
    const { error } = await useMyFetch(`incomes/${income.income_id}`).patch(income).json()

    if (!error.value) {
      const index = incomes.value.findIndex(c => c.income_id === income.income_id)
      if (index !== -1) {
        incomes.value[index] = { ...incomes.value[index], ...income }
        addToast("Edited income succesfully", "success")
      } else {
        addToast("Editing income failed", "error")
      }
    } else {
      addToast("Editing income failed", "error")
    }
  }

  async function deleteIncome(id: number) {
    await useMyFetch(`incomes/${id}`).delete().json()
    incomes.value = incomes.value.filter((c) => c.income_id !== id)
    addToast("Income deleted succesfully", "success")
  }

  function reset() {
    incomes.value = []
  }

  return {
    incomes,
    loading,
    fetchError,
    fetchIncomes,
    editIncome,
    deleteIncome,
    submitForm,
    reset,
  }
})
