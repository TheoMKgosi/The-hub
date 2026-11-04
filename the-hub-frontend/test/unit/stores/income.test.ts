import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useIncomeStore } from '../../../app/stores/income'

// The income store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them all; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
})

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
})

const income = (id: number, amount = 100, source = 'salary') => ({
  income_id: id,
  source,
  amount,
  received_at: new Date('2026-01-15')
})

describe('useIncomeStore', () => {
  describe('fetchIncomes', () => {
    it('populates incomes and clears loading', async () => {
      const store = useIncomeStore()
      apiMock.mockResolvedValue({ incomes: [income(1, 2500), income(2, 180)] })

      await store.fetchIncomes()

      expect(apiMock).toHaveBeenCalledWith('incomes')
      expect(store.incomes).toHaveLength(2)
      expect(store.loading).toBe(false)
    })
  })

  describe('submitForm', () => {
    it('adds the income optimistically, then replaces it with the server response', async () => {
      const store = useIncomeStore()
      const payload = income(0, 300, 'freelance')
      apiMock.mockResolvedValue(income(9, 300, 'freelance'))

      await store.submitForm(payload)

      expect(apiMock).toHaveBeenCalledWith('incomes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.incomes).toHaveLength(1)
      expect(store.incomes[0]).toEqual(income(9, 300, 'freelance'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Income added succesfully', color: 'success' })
    })

    it('rolls back the optimistic income on failure', async () => {
      const store = useIncomeStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm(income(0, 300))

      expect(store.incomes).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Income not added', color: 'error' })
    })
  })

  describe('editIncome', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useIncomeStore()
      store.incomes.push(income(1, 100))
      apiMock.mockResolvedValue(income(1, 500))

      await store.editIncome(income(1, 500, 'salary'))

      expect(apiMock).toHaveBeenCalledWith('incomes/1', {
        method: 'PATCH',
        body: JSON.stringify(income(1, 500, 'salary'))
      })
      expect(store.incomes[0].amount).toBe(500)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Edited income succesfully', color: 'success' })
    })

    it('reverts to the original income on failure', async () => {
      const store = useIncomeStore()
      store.incomes.push(income(1, 100, 'salary'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editIncome(income(1, 500, 'salary'))

      expect(store.incomes[0].amount).toBe(100)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing income failed', color: 'error' })
    })
  })

  describe('deleteIncome', () => {
    it('removes the income optimistically and succeeds', async () => {
      const store = useIncomeStore()
      store.incomes.push(income(1), income(2))
      apiMock.mockResolvedValue(undefined)

      await store.deleteIncome(1)

      expect(apiMock).toHaveBeenCalledWith('incomes/1', { method: 'DELETE' })
      expect(store.incomes.map(i => i.income_id)).toEqual([2])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Income deleted succesfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useIncomeStore()

      await store.deleteIncome(42)

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Income not found', color: 'error' })
    })

    it('restores the income when deletion fails', async () => {
      const store = useIncomeStore()
      store.incomes.push(income(1, 200, 'bonus'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteIncome(1)

      expect(store.incomes.map(i => i.income_id)).toEqual([1])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Income did not delete', color: 'error' })
    })
  })

  it('reset clears all incomes', () => {
    const store = useIncomeStore()
    store.incomes.push(income(1))

    store.reset()

    expect(store.incomes).toHaveLength(0)
  })
})
