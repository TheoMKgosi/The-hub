import { defineStore } from 'pinia'
import { ref } from 'vue'

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
  const toast = useToast()

  async function fetchIncomes() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedIncome = await $api<IncomeResponse>('incomes')

    if (fetchedIncome) incomes.value = fetchedIncome.incomes

    loading.value = false
  }


  async function submitForm(payload: Income) {
    // Create optimistic income
    const optimisticIncome: Income = {
      income_id: -Date.now(), // Use negative ID to avoid conflicts
      source: payload.source,
      amount: payload.amount,
      received_at: payload.received_at
    }

    // Optimistically add to local state
    incomes.value.push(optimisticIncome)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Income>('incomes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic income with real data
      const optimisticIndex = incomes.value.findIndex(i => i.income_id === optimisticIncome.income_id)
      if (optimisticIndex !== -1) {
        incomes.value[optimisticIndex] = data
      }

      toast.add({ title: "Income added succesfully", color: 'success' })
    } catch (err) {
      // Remove optimistic income on error
      incomes.value = incomes.value.filter(i => i.income_id !== optimisticIncome.income_id)
      toast.add({ title: "Income not added", color: 'error' })
    }
  }

  async function editIncome(payload: Income) {
    // Store original income for potential rollback
    const originalIncomeIndex = incomes.value.findIndex(i => i.income_id === payload.income_id)
    const originalIncome = originalIncomeIndex !== -1 ? { ...incomes.value[originalIncomeIndex] } : null

    // Optimistically update the income
    if (originalIncomeIndex !== -1) {
      incomes.value[originalIncomeIndex] = { ...incomes.value[originalIncomeIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Income>(`incomes/${payload.income_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalIncomeIndex !== -1 && data) {
        incomes.value[originalIncomeIndex] = data
      }

      toast.add({ title: "Edited income succesfully", color: 'success' })
    } catch (err) {
      // Revert optimistic update on error
      if (originalIncome && originalIncomeIndex !== -1) {
        incomes.value[originalIncomeIndex] = originalIncome
      }
      toast.add({ title: "Editing income failed", color: 'error' })
    }
  }

  async function deleteIncome(id: number) {
    // Store the income for potential rollback
    const incomeToDelete = incomes.value.find(i => i.income_id === id)
    if (!incomeToDelete) {
      toast.add({ title: "Income not found", color: 'error' })
      return
    }

    // Optimistically remove from local state
    incomes.value = incomes.value.filter((i) => i.income_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`incomes/${id}`, {
        method: 'DELETE'
      })

      toast.add({ title: "Income deleted succesfully", color: 'success' })
    } catch (err) {
      // Restore the income on error
      incomes.value.push(incomeToDelete)
      toast.add({ title: 'Income did not delete', color: 'error' })
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
