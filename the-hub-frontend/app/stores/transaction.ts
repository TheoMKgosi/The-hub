import { defineStore } from 'pinia'
import { ref } from 'vue'

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
  const toast = useToast()

  async function fetchTransactions() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedTransactions = await $api<TransactionResponse>('transactions')

    if (fetchedTransactions) transactions.value = fetchedTransactions.transactions

    loading.value = false
  }

  async function submitForm(payload: Omit<Transaction, 'transaction_id'>) {
    // Create optimistic transaction
    const optimisticTransaction: Transaction = {
      transaction_id: `temp-${Date.now()}`,
      description: payload.description,
      amount: payload.amount,
      type: payload.type,
      date: payload.date,
      category_id: payload.category_id,
      category: payload.category
    }

    // Optimistically add to local state
    transactions.value.push(optimisticTransaction)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Transaction>('transactions', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic transaction with real data
      const optimisticIndex = transactions.value.findIndex(t => t.transaction_id === optimisticTransaction.transaction_id)
      if (optimisticIndex !== -1) {
        transactions.value[optimisticIndex] = data
      }

      toast.add({title: "Finance", description: "Transaction added successfully", color: "success"})
    } catch (err) {
      // Remove optimistic transaction on error
      transactions.value = transactions.value.filter(t => t.transaction_id !== optimisticTransaction.transaction_id)
      toast.add({title: "Error", description: "Transaction not added", color: "error"})
    }
  }

  async function editTransaction(payload: Transaction) {
    // Store original transaction for potential rollback
    const originalTransactionIndex = transactions.value.findIndex(t => t.transaction_id === payload.transaction_id)
    const originalTransaction = originalTransactionIndex !== -1 ? { ...transactions.value[originalTransactionIndex] } : null

    // Optimistically update the transaction
    if (originalTransactionIndex !== -1) {
      transactions.value[originalTransactionIndex] = { ...transactions.value[originalTransactionIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Transaction>(`transactions/${payload.transaction_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalTransactionIndex !== -1 && data) {
        transactions.value[originalTransactionIndex] = data
      }

      toast.add({title: "Finance", description: "Transaction edited successfully", color: "success"})
    } catch (err) {
      // Revert optimistic update on error
      if (originalTransaction && originalTransactionIndex !== -1) {
        transactions.value[originalTransactionIndex] = originalTransaction
      }
      toast.add({title: "Error", description: "Editing transaction failed", color: "error"})
    }
  }

  async function deleteTransaction(id: string) {
    // Store the transaction for potential rollback
    const transactionToDelete = transactions.value.find(t => t.transaction_id === id)
    if (!transactionToDelete) {
      toast.add({title: "Error", description: "Transaction not found", color: "error"})
      return
    }

    // Optimistically remove from local state
    transactions.value = transactions.value.filter((t) => t.transaction_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`transactions/${id}`, {
        method: 'DELETE'
      })

      toast.add({title: "Finance", description: "Transaction deleted successfully", color: "success"})
    } catch (err) {
      // Restore the transaction on error
      transactions.value.push(transactionToDelete)
      toast.add({title: "Error", description: "Transaction did not delete", color: "error"})
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
