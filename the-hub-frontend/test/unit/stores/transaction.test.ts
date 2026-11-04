import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useTransactionStore } from '../../../app/stores/transaction'

// The transaction store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them all; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
})

const transaction = (id: string, overrides: Record<string, unknown> = {}) => ({
  transaction_id: id,
  description: `txn-${id}`,
  amount: 25,
  type: 'expense',
  date: '2026-01-15',
  ...overrides
})

describe('useTransactionStore', () => {
  describe('fetchTransactions', () => {
    it('populates transactions and clears loading', async () => {
      const store = useTransactionStore()
      apiMock.mockResolvedValue({ transactions: [transaction('t1'), transaction('t2')] })

      await store.fetchTransactions()

      expect(apiMock).toHaveBeenCalledWith('transactions')
      expect(store.transactions).toEqual([transaction('t1'), transaction('t2')])
      expect(store.loading).toBe(false)
    })

    it('keeps existing transactions when the API returns nothing', async () => {
      const store = useTransactionStore()
      store.transactions.push(transaction('t1'))
      apiMock.mockResolvedValue(null)

      await store.fetchTransactions()

      expect(store.transactions).toEqual([transaction('t1')])
      expect(store.loading).toBe(false)
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useTransactionStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchTransactions()
      expect(store.loading).toBe(true)

      resolveApi({ transactions: [] })
      await pending
      expect(store.loading).toBe(false)
    })

    it('leaves loading stuck when the API rejects', async () => {
      const store = useTransactionStore()
      apiMock.mockRejectedValue(new Error('network down'))

      // NOTE: documents current behavior — fetchTransactions has no
      // try/finally, so loading stays true after a failed fetch and the
      // error propagates to the caller.
      await expect(store.fetchTransactions()).rejects.toThrow('network down')
      expect(store.loading).toBe(true)
    })
  })

  describe('submitForm', () => {
    it('adds the transaction optimistically, then replaces it with the server response', async () => {
      const store = useTransactionStore()
      const payload = { description: 'Coffee', amount: 5, type: 'expense' as const, date: '2026-01-15' }
      apiMock.mockResolvedValue(transaction('t9', { description: 'Coffee', amount: 5 }))

      await store.submitForm(payload)

      expect(apiMock).toHaveBeenCalledWith('transactions', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.transactions).toHaveLength(1)
      expect(store.transactions[0]).toEqual(transaction('t9', { description: 'Coffee', amount: 5 }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Transaction added successfully', color: 'success' })
    })

    it('rolls back the optimistic transaction on failure', async () => {
      const store = useTransactionStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ description: 'Coffee', amount: 5, type: 'expense', date: '2026-01-15' })

      expect(store.transactions).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Transaction not added', color: 'error' })
    })
  })

  describe('editTransaction', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useTransactionStore()
      store.transactions.push(transaction('t1', { amount: 10 }))
      apiMock.mockResolvedValue(transaction('t1', { amount: 99 }))

      await store.editTransaction(transaction('t1', { amount: 99 }) as any)

      expect(apiMock).toHaveBeenCalledWith('transactions/t1', {
        method: 'PATCH',
        body: JSON.stringify(transaction('t1', { amount: 99 }))
      })
      expect(store.transactions[0]).toEqual(transaction('t1', { amount: 99 }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Transaction edited successfully', color: 'success' })
    })

    it('reverts to the original transaction on failure', async () => {
      const store = useTransactionStore()
      store.transactions.push(transaction('t1', { amount: 10, description: 'original' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editTransaction(transaction('t1', { amount: 99 }) as any)

      expect(store.transactions[0]).toEqual(transaction('t1', { amount: 10, description: 'original' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Editing transaction failed', color: 'error' })
    })

    it('is a no-op for a transaction that is not in the list', async () => {
      const store = useTransactionStore()
      apiMock.mockResolvedValue(transaction('t99'))

      await store.editTransaction(transaction('t99') as any)

      expect(store.transactions).toHaveLength(0)
    })
  })

  describe('deleteTransaction', () => {
    it('removes the transaction optimistically and succeeds', async () => {
      const store = useTransactionStore()
      store.transactions.push(transaction('t1'), transaction('t2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTransaction('t1')

      expect(apiMock).toHaveBeenCalledWith('transactions/t1', { method: 'DELETE' })
      expect(store.transactions.map(t => t.transaction_id)).toEqual(['t2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Transaction deleted successfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useTransactionStore()

      await store.deleteTransaction('t404')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Transaction not found', color: 'error' })
    })

    it('restores the transaction when deletion fails', async () => {
      const store = useTransactionStore()
      store.transactions.push(transaction('t1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteTransaction('t1')

      expect(store.transactions.map(t => t.transaction_id)).toEqual(['t1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Transaction did not delete', color: 'error' })
    })
  })

  it('reset clears all transactions', () => {
    const store = useTransactionStore()
    store.transactions.push(transaction('t1'))

    store.reset()

    expect(store.transactions).toHaveLength(0)
  })
})
