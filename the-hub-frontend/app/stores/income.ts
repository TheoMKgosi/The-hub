
import { defineStore } from 'pinia'
import { ref } from 'vue'
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
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedIncome = await $api<IncomeResponse>('incomes')

    if (fetchedIncome) incomes.value = fetchedIncome.incomes

    loading.value = false
  }


  // TODO: Fix the payload for parameter to be object
  async function submitForm(payload: Income) {
    try {
      const { $api } = useNuxtApp()
      await $api('incomes', {
        method: 'POST',
        body: payload
      })
      fetchIncomes()
      addToast("Income added succesfully", "success")
    } catch (err) {
      addToast("Income not added", "error")

    }
  }

  async function editIncome(payload: Income) {
    try {
      const { $api } = useNuxtApp()
      await $api(`incomes/${payload.income_id}`, {
        method: 'PATCH',
        body: payload
      })

      // TODO: Filter locally than make a request
      fetchIncomes()
      addToast("Edited income succesfully", "success")
    } catch (err) {
      addToast("Editing income failed", "error")
    }
  }

  async function deleteIncome(id: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`incomes/${id}`, {
        method: 'DELETE'
      })
      incomes.value = incomes.value.filter((c) => c.income_id !== id)
      addToast("Income deleted succesfully", "success")

    } catch(err) {
      addToast('Income did not delete', 'error')
    }
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
