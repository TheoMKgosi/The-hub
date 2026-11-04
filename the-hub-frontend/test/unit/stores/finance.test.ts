import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref } from 'vue'

// finance.ts contains three stores (category, budget, receipt) and relies on
// Vue auto-imports (ref), Nuxt auto-imports (defineStore, useNuxtApp,
// useToast, useIncomeStore) and the useValidation composable. Stub them all
// before the store module loads (defineStore is referenced at module top
// level, so the store must be imported dynamically below); the $api mock
// stands in for the API plugin. useValidation gets the real implementation.
const apiMock = vi.fn()
const toastAddMock = vi.fn()
const incomeFetchMock = vi.fn()

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('defineStore', defineStore)
const { useValidation } = await import('../../../app/composables/useValidation')
vi.stubGlobal('useValidation', useValidation)
// The budget store instantiates the income store during setup; a light fake
// keeps $api call assertions clean.
vi.stubGlobal('useIncomeStore', () => ({ fetchIncomes: incomeFetchMock }))

const { useCategoryStore, useBudgetStore, useReceiptStore } = await import('../../../app/stores/finance')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  incomeFetchMock.mockClear()
  vi.spyOn(console, 'log').mockImplementation(() => {})
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const category = (id: string, name = `cat-${id}`) => ({ budget_category_id: id, name })

const budget = (id: string, overrides: Record<string, unknown> = {}) => ({
  budget_id: id,
  category_id: 'c1',
  Category: category('c1'),
  amount: 500,
  start_date: '2026-01-01',
  end_date: '2026-02-01',
  ...overrides
})

const receipt = (id: string, overrides: Record<string, unknown> = {}) => ({
  receipt_id: id,
  title: `receipt-${id}`,
  image_path: '/uploads/receipt.png',
  created_at: '2026-01-01T00:00:00.000Z',
  ...overrides
})

const analyticsRow = (id: string, spent = 50) => ({
  budget_id: id,
  category_name: 'Food',
  budget_amount: 100,
  spent_amount: spent,
  remaining_amount: 100 - spent,
  utilization_rate: spent,
  days_remaining: 10,
  daily_spending_rate: 5,
  status: 'ok'
})

describe('useCategoryStore', () => {
  describe('fetchCategory', () => {
    it('populates categories and clears loading', async () => {
      const store = useCategoryStore()
      apiMock.mockResolvedValue({ categories: [category('c1'), category('c2')] })

      await store.fetchCategory()

      expect(apiMock).toHaveBeenCalledWith('categories')
      expect(store.categories).toEqual([category('c1'), category('c2')])
      expect(store.loading).toBe(false)
    })

    it('keeps existing categories when the API returns nothing', async () => {
      const store = useCategoryStore()
      store.categories.push(category('c1'))
      apiMock.mockResolvedValue(null)

      await store.fetchCategory()

      expect(store.categories).toEqual([category('c1')])
      expect(store.loading).toBe(false)
    })

    it('leaves loading stuck when the API rejects', async () => {
      const store = useCategoryStore()
      apiMock.mockRejectedValue(new Error('network down'))

      // NOTE: documents current behavior — fetchCategory has no try/finally,
      // so loading stays true after a failed fetch (unlike other stores).
      await expect(store.fetchCategory()).rejects.toThrow('network down')
      expect(store.loading).toBe(true)
    })
  })

  describe('submitForm', () => {
    it('rejects an invalid payload with a toast and no API call', async () => {
      const store = useCategoryStore()

      await store.submitForm({ name: '' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.categories).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'name is required', color: 'error' })
    })

    it('adds the category optimistically, then replaces it with the server response', async () => {
      const store = useCategoryStore()
      apiMock.mockResolvedValue(category('c9', 'work'))

      await store.submitForm({ name: 'work' })

      expect(apiMock).toHaveBeenCalledWith('categories', {
        method: 'POST',
        body: JSON.stringify({ name: 'work' })
      })
      expect(store.categories).toEqual([category('c9', 'work')])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Category added successfully', color: 'success' })
    })

    it('rolls back the optimistic category on failure', async () => {
      const store = useCategoryStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ name: 'work' })

      expect(store.categories).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to add category. Please try again.', color: 'error' })
    })
  })

  describe('editCategory', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useCategoryStore()
      store.categories.push(category('c1', 'old-name'))
      apiMock.mockResolvedValue(category('c1', 'new-name'))

      await store.editCategory(category('c1', 'new-name'))

      expect(apiMock).toHaveBeenCalledWith('categories/c1', {
        method: 'PATCH',
        body: JSON.stringify(category('c1', 'new-name'))
      })
      expect(store.categories[0]).toEqual(category('c1', 'new-name'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Edited category succesfully', color: 'success' })
    })

    it('reverts to the original category on failure', async () => {
      const store = useCategoryStore()
      store.categories.push(category('c1', 'original'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editCategory(category('c1', 'changed'))

      expect(store.categories[0]).toEqual(category('c1', 'original'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing category failed', color: 'error' })
    })

    it('is a no-op for a category that is not in the list', async () => {
      const store = useCategoryStore()
      apiMock.mockResolvedValue(category('c99', 'remote'))

      await store.editCategory(category('c99', 'remote'))

      expect(store.categories).toHaveLength(0)
    })
  })

  describe('deleteCategory', () => {
    it('removes the category optimistically via soft-delete and succeeds', async () => {
      const store = useCategoryStore()
      store.categories.push(category('c1'), category('c2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteCategory('c1')

      expect(apiMock).toHaveBeenCalledWith('categories/c1/soft-delete', { method: 'PATCH' })
      expect(store.categories.map(c => c.budget_category_id)).toEqual(['c2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Category deleted successfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useCategoryStore()

      await store.deleteCategory('c404')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Category not found', color: 'error' })
    })

    it('restores the category when deletion fails', async () => {
      const store = useCategoryStore()
      store.categories.push(category('c1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteCategory('c1')

      expect(store.categories.map(c => c.budget_category_id)).toEqual(['c1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Category deletion failed', color: 'error' })
    })
  })

  it('reset clears all categories', () => {
    const store = useCategoryStore()
    store.categories.push(category('c1'))

    store.reset()

    expect(store.categories).toHaveLength(0)
  })
})

describe('useBudgetStore', () => {
  describe('fetchBudget', () => {
    it('populates budgets and clears loading', async () => {
      const store = useBudgetStore()
      apiMock.mockResolvedValue({ budgets: [budget('b1'), budget('b2')] })

      await store.fetchBudget()

      expect(apiMock).toHaveBeenCalledWith('budgets')
      expect(store.budgets).toHaveLength(2)
      expect(store.loading).toBe(false)
    })
  })

  describe('submitForm', () => {
    it('adds the budget optimistically with a resolved Category, then swaps in the server budget', async () => {
      const categoryStore = useCategoryStore()
      categoryStore.categories.push(category('c1', 'Food'))
      const store = useBudgetStore()
      apiMock.mockResolvedValue(budget('b9'))

      await store.submitForm({ category_id: 'c1', amount: 500, start_date: '2026-01-01', end_date: '2026-02-01' })

      expect(apiMock).toHaveBeenCalledWith('budgets', {
        method: 'POST',
        body: JSON.stringify({ category_id: 'c1', amount: 500, start_date: '2026-01-01', end_date: '2026-02-01' })
      })
      expect(store.budgets).toHaveLength(1)
      expect(store.budgets[0].budget_id).toBe('b9')
      // The Category is taken from the category store, not the server payload
      expect(store.budgets[0].Category.name).toBe('Food')
      expect(incomeFetchMock).toHaveBeenCalledTimes(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Budget added successfully', color: 'success' })
      expect(store.creating).toBe(false)
    })

    it('uses a Loading... placeholder for an unknown category while the request is in flight', async () => {
      const store = useBudgetStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.submitForm({ category_id: 'unknown', amount: 100, start_date: '2026-01-01', end_date: '2026-02-01' })
      expect(store.budgets[0].Category.name).toBe('Loading...')

      resolveApi(budget('b9'))
      await pending
      expect(store.budgets[0].budget_id).toBe('b9')
    })

    it('rolls back the optimistic budget on failure and does not refresh incomes', async () => {
      const store = useBudgetStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ category_id: 'c1', amount: 500, start_date: '2026-01-01', end_date: '2026-02-01' })

      expect(store.budgets).toHaveLength(0)
      expect(incomeFetchMock).not.toHaveBeenCalled()
      expect(store.creating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to add budget', color: 'error' })
    })
  })

  describe('editBudget', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useBudgetStore()
      store.budgets.push(budget('b1', { amount: 100 }))
      apiMock.mockResolvedValue(budget('b1', { amount: 900 }))

      await store.editBudget(budget('b1', { amount: 900 }))

      expect(apiMock).toHaveBeenCalledWith('budgets/b1', {
        method: 'PATCH',
        body: JSON.stringify(budget('b1', { amount: 900 }))
      })
      expect(store.budgets[0].amount).toBe(900)
      expect(incomeFetchMock).toHaveBeenCalledTimes(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Budget updated successfully', color: 'success' })
      expect(store.updating).toBe(false)
    })

    it('reverts to the original budget on failure', async () => {
      const store = useBudgetStore()
      store.budgets.push(budget('b1', { amount: 100 }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editBudget(budget('b1', { amount: 900 }))

      expect(store.budgets[0].amount).toBe(100)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to update budget. Please try again.', color: 'error' })
      expect(store.updating).toBe(false)
    })
  })

  describe('deleteBudget', () => {
    it('removes the budget optimistically and refreshes incomes', async () => {
      const store = useBudgetStore()
      store.budgets.push(budget('b1'), budget('b2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteBudget('b1')

      expect(apiMock).toHaveBeenCalledWith('budgets/b1', { method: 'DELETE' })
      expect(store.budgets.map(b => b.budget_id)).toEqual(['b2'])
      expect(incomeFetchMock).toHaveBeenCalledTimes(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Budget deleted successfully', color: 'success' })
      expect(store.deleting).toBe(false)
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useBudgetStore()

      await store.deleteBudget('b404')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Budget not found', color: 'error' })
      expect(store.deleting).toBe(false)
    })

    it('restores the budget when deletion fails', async () => {
      const store = useBudgetStore()
      store.budgets.push(budget('b1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteBudget('b1')

      expect(store.budgets.map(b => b.budget_id)).toEqual(['b1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to delete budget', color: 'error' })
      expect(store.deleting).toBe(false)
    })
  })

  describe('analytics, alerts and suggestions', () => {
    it('fetchBudgetAnalytics stores the analytics payload', async () => {
      const store = useBudgetStore()
      apiMock.mockResolvedValue({ analytics: [analyticsRow('b1')] })

      await store.fetchBudgetAnalytics('current')

      expect(apiMock).toHaveBeenCalledWith('budgets/analytics?period=current')
      expect(store.analytics).toEqual([analyticsRow('b1')])
      expect(store.analyticsLoading).toBe(false)
    })

    it('fetchBudgetAnalytics shows a toast on failure and clears loading', async () => {
      const store = useBudgetStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchBudgetAnalytics()

      expect(store.analyticsLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load budget analytics', color: 'error' })
    })

    it('fetchBudgetAlerts stores the alerts payload', async () => {
      const store = useBudgetStore()
      const alert = { budget_id: 'b1', category_name: 'Food', alert_type: 'warning', message: 'close to limit', current_usage: 90, budget_amount: 100, days_remaining: 3 }
      apiMock.mockResolvedValue({ alerts: [alert] })

      await store.fetchBudgetAlerts()

      expect(apiMock).toHaveBeenCalledWith('budgets/alerts')
      expect(store.alerts).toEqual([alert])
      expect(store.alertsLoading).toBe(false)
    })

    it('fetchBudgetAlerts shows a toast on failure', async () => {
      const store = useBudgetStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchBudgetAlerts()

      expect(store.alertsLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load budget alerts', color: 'error' })
    })

    it('fetchBudgetSuggestions stores categories and reports the count', async () => {
      const store = useBudgetStore()
      const suggestion = { category_id: 'c1', category_name: 'Food', current_average: 120, suggested_budget: 150, transaction_count: 8, variability: 0.2, confidence: 'high' }
      apiMock.mockResolvedValue({ suggestions: { categories: [suggestion], analysis_period: 'last_90_days', total_suggestions: 1 } })

      await store.fetchBudgetSuggestions()

      expect(apiMock).toHaveBeenCalledWith('budgets/suggestions')
      expect(store.suggestions).toEqual([suggestion])
      expect(store.suggestionsLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Loaded 1 budget suggestions', color: 'success' })
    })

    it('fetchBudgetSuggestions shows a toast on failure', async () => {
      const store = useBudgetStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchBudgetSuggestions()

      expect(store.suggestionsLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load budget suggestions', color: 'error' })
    })
  })

  it('reset clears budgets, analytics, alerts and suggestions', () => {
    const store = useBudgetStore()
    store.budgets.push(budget('b1'))
    store.analytics.push(analyticsRow('b1') as any)
    store.alerts.push({ budget_id: 'b1' } as any)
    store.suggestions.push({ category_id: 'c1' } as any)

    store.reset()

    expect(store.budgets).toHaveLength(0)
    expect(store.analytics).toHaveLength(0)
    expect(store.alerts).toHaveLength(0)
    expect(store.suggestions).toHaveLength(0)
  })
})

describe('useReceiptStore', () => {
  describe('fetchReceipts', () => {
    it('populates receipts and clears loading', async () => {
      const store = useReceiptStore()
      apiMock.mockResolvedValue({ receipts: [receipt('r1'), receipt('r2')] })

      await store.fetchReceipts()

      expect(apiMock).toHaveBeenCalledWith('receipts')
      expect(store.receipts).toEqual([receipt('r1'), receipt('r2')])
      expect(store.loading).toBe(false)
    })

    it('shows a toast and clears loading on failure', async () => {
      const store = useReceiptStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchReceipts()

      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load receipts', color: 'error' })
    })
  })

  describe('submitForm', () => {
    it('unshifts the server receipt on success (no optimistic insert)', async () => {
      const store = useReceiptStore()
      const payload = { title: 'Lunch', image_data: 'base64-data', amount: 12.5, date: '2026-01-02', category_id: 'c1' }
      apiMock.mockResolvedValue(receipt('r9', { title: 'Lunch' }))

      await store.submitForm(payload)

      expect(apiMock).toHaveBeenCalledWith('receipts', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.receipts).toEqual([receipt('r9', { title: 'Lunch' })])
      expect(store.creating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Receipt added successfully', color: 'success' })
    })

    it('adds nothing on failure and clears creating', async () => {
      const store = useReceiptStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ title: 'Lunch', image_data: 'base64-data' })

      expect(store.receipts).toHaveLength(0)
      expect(store.creating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to add receipt. Please try again.', color: 'error' })
    })
  })

  describe('updateReceipt', () => {
    it('applies the update optimistically, then reconciles with the server response', async () => {
      const store = useReceiptStore()
      store.receipts.push(receipt('r1', { title: 'old' }))
      apiMock.mockResolvedValue(receipt('r1', { title: 'new' }))

      await store.updateReceipt({ receipt_id: 'r1', title: 'new' })

      expect(apiMock).toHaveBeenCalledWith('receipts/r1', {
        method: 'PATCH',
        body: JSON.stringify({ receipt_id: 'r1', title: 'new' })
      })
      expect(store.receipts[0]).toEqual(receipt('r1', { title: 'new' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Receipt updated successfully', color: 'success' })
      expect(store.updating).toBe(false)
    })

    it('reverts to the original receipt on failure', async () => {
      const store = useReceiptStore()
      store.receipts.push(receipt('r1', { title: 'original' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.updateReceipt({ receipt_id: 'r1', title: 'changed' })

      expect(store.receipts[0]).toEqual(receipt('r1', { title: 'original' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to update receipt. Please try again.', color: 'error' })
      expect(store.updating).toBe(false)
    })
  })

  describe('deleteReceipt', () => {
    it('removes the receipt optimistically and succeeds', async () => {
      const store = useReceiptStore()
      store.receipts.push(receipt('r1'), receipt('r2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteReceipt('r1')

      expect(apiMock).toHaveBeenCalledWith('receipts/r1', { method: 'DELETE' })
      expect(store.receipts.map(r => r.receipt_id)).toEqual(['r2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Receipt deleted successfully', color: 'success' })
      expect(store.deleting).toBe(false)
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useReceiptStore()

      await store.deleteReceipt('r404')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Receipt not found', color: 'error' })
      expect(store.deleting).toBe(false)
    })

    it('restores the receipt when deletion fails', async () => {
      const store = useReceiptStore()
      store.receipts.push(receipt('r1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteReceipt('r1')

      expect(store.receipts.map(r => r.receipt_id)).toEqual(['r1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to delete receipt. Please try again.', color: 'error' })
      expect(store.deleting).toBe(false)
    })
  })

  it('reset clears all receipts', () => {
    const store = useReceiptStore()
    store.receipts.push(receipt('r1'))

    store.reset()

    expect(store.receipts).toHaveLength(0)
  })
})
