import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

interface Transaction {
  transaction_id: string
  description: string
  amount: number
  type: 'income' | 'expense'
  date: string
  category_id?: string
  category?: {
    budget_category_id: string
    name: string
  }
}

export interface TransactionResponse {
  transactions: Transaction[]
}

export const useTransactionStore = defineStore('transaction', () => {
  const transactions = ref<Transaction[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchTransactions() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedTransactions = await $api<TransactionResponse>('transactions')

    if (fetchedTransactions) transactions.value = fetchedTransactions.transactions

    loading.value = false
  }

  async function submitForm(payload: Omit<Transaction, 'transaction_id'>) {
    try {
      const { $api } = useNuxtApp()
      await $api('transactions', {
        method: 'POST',
        body: payload
      })
      fetchTransactions()
      addToast("Transaction added successfully", "success")
    } catch (err) {
      addToast("Transaction not added", "error")
    }
  }

  async function editTransaction(payload: Transaction) {
    try {
      const { $api } = useNuxtApp()
      await $api(`transactions/${payload.transaction_id}`, {
        method: 'PATCH',
        body: payload
      })

      fetchTransactions()
      addToast("Transaction edited successfully", "success")
    } catch (err) {
      addToast("Editing transaction failed", "error")
    }
  }

  async function deleteTransaction(id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`transactions/${id}`, {
        method: 'DELETE'
      })
      transactions.value = transactions.value.filter((t) => t.transaction_id !== id)
      addToast("Transaction deleted successfully", "success")
    } catch(err) {
      addToast('Transaction did not delete', 'error')
    }
  }

  function reset() {
    transactions.value = []
  }

  return {
    transactions,
    loading,
    fetchError,
    fetchTransactions,
    editTransaction,
    deleteTransaction,
    submitForm,
    reset,
  }
})