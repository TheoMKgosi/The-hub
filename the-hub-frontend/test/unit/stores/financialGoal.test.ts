import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref } from 'vue'

// The financial-goal store relies on Vue auto-imports (ref) and Nuxt
// auto-imports (defineStore, useNuxtApp, useToast). Stub them all before the
// store module loads (defineStore is referenced at module top level, so the
// store must be imported dynamically below); the $api mock stands in for the
// API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('defineStore', defineStore)

const { useFinancialGoalStore } = await import('../../../app/stores/financialGoal')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'log').mockImplementation(() => {})
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const goal = (id: string, overrides: Record<string, unknown> = {}) => ({
  goal_id: id,
  title: `goal-${id}`,
  description: '',
  target_amount: 1000,
  current_amount: 100,
  type: 'savings',
  status: 'active',
  priority: null,
  target_date: null,
  category_id: null,
  color: '#3B82F6',
  created_at: '2026-01-01T00:00:00.000Z',
  ...overrides
})

describe('useFinancialGoalStore', () => {
  describe('fetchGoals', () => {
    it('populates goals and clears loading', async () => {
      const store = useFinancialGoalStore()
      apiMock.mockResolvedValue({ goals: [goal('g1'), goal('g2')] })

      await store.fetchGoals()

      expect(apiMock).toHaveBeenCalledWith('financial-goals')
      expect(store.goals).toEqual([goal('g1'), goal('g2')])
      expect(store.loading).toBe(false)
    })

    it('shows a toast and clears loading on failure', async () => {
      const store = useFinancialGoalStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchGoals()

      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to load financial goals', color: 'error' })
    })
  })

  describe('submitForm', () => {
    it('unshifts an optimistic goal, then replaces it with the server response', async () => {
      const store = useFinancialGoalStore()
      const payload = { title: 'Emergency fund', target_amount: 5000, current_amount: 500, type: 'savings' as const }
      apiMock.mockResolvedValue(goal('g9', { title: 'Emergency fund' }))

      await store.submitForm(payload)

      expect(apiMock).toHaveBeenCalledWith('financial-goals', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.goals).toHaveLength(1)
      expect(store.goals[0]).toEqual(goal('g9', { title: 'Emergency fund' }))
      expect(store.creating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Goal created successfully', color: 'success' })
    })

    it('applies defaults for optional fields in the optimistic goal', async () => {
      const store = useFinancialGoalStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.submitForm({ title: 'Bare minimum' })
      expect(store.creating).toBe(true)
      expect(store.goals[0]).toMatchObject({
        title: 'Bare minimum',
        description: '',
        target_amount: null,
        current_amount: 0,
        type: 'savings',
        status: 'active',
        priority: null,
        color: '#3B82F6'
      })

      resolveApi(goal('g9'))
      await pending
      expect(store.goals[0].goal_id).toBe('g9')
    })

    it('rolls back the optimistic goal on failure', async () => {
      const store = useFinancialGoalStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ title: 'Emergency fund' })

      expect(store.goals).toHaveLength(0)
      expect(store.creating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to create goal', color: 'error' })
    })
  })

  describe('editGoal', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { title: 'old' }))
      apiMock.mockResolvedValue(goal('g1', { title: 'new' }))

      await store.editGoal({ goal_id: 'g1', title: 'new' })

      expect(apiMock).toHaveBeenCalledWith('financial-goals/g1', {
        method: 'PATCH',
        body: JSON.stringify({ goal_id: 'g1', title: 'new' })
      })
      expect(store.goals[0]).toEqual(goal('g1', { title: 'new' }))
      expect(store.updating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Goal updated successfully', color: 'success' })
    })

    it('reverts to the original goal on failure', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { title: 'original' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editGoal({ goal_id: 'g1', title: 'changed' })

      expect(store.goals[0]).toEqual(goal('g1', { title: 'original' }))
      expect(store.updating).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to update goal', color: 'error' })
    })

    it('is a no-op for a goal that is not in the list', async () => {
      const store = useFinancialGoalStore()
      apiMock.mockResolvedValue(goal('g99'))

      await store.editGoal({ goal_id: 'g99', title: 'remote' })

      expect(store.goals).toHaveLength(0)
    })
  })

  describe('deleteGoal', () => {
    it('removes the goal optimistically and succeeds', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1'), goal('g2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteGoal('g1')

      expect(apiMock).toHaveBeenCalledWith('financial-goals/g1', { method: 'DELETE' })
      expect(store.goals.map(g => g.goal_id)).toEqual(['g2'])
      expect(store.deleting).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Goal deleted successfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useFinancialGoalStore()

      await store.deleteGoal('g404')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Goal not found', color: 'error' })
      expect(store.deleting).toBe(false)
    })

    it('restores the goal when deletion fails', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteGoal('g1')

      expect(store.goals.map(g => g.goal_id)).toEqual(['g1'])
      expect(store.deleting).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to delete goal', color: 'error' })
    })
  })

  describe('addProgress', () => {
    it('applies progress optimistically, then replaces the goal with the server response', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1'))
      apiMock.mockResolvedValue(goal('g1', { current_amount: 175 }))

      await store.addProgress('g1', 75)

      expect(apiMock).toHaveBeenCalledWith('financial-goals/g1/progress', {
        method: 'PATCH',
        body: JSON.stringify({ amount: 75 })
      })
      expect(store.goals[0]).toEqual(goal('g1', { current_amount: 175 }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Finance', description: 'Added $75.00 to goal', color: 'success' })
    })

    it('marks the goal completed locally when the target is reached', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { target_amount: 1000, current_amount: 950 }))
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.addProgress('g1', 50)
      expect(store.goals[0].current_amount).toBe(1000)
      expect(store.goals[0].status).toBe('completed')

      resolveApi(goal('g1', { current_amount: 1000, status: 'completed' }))
      await pending
      expect(store.goals[0].status).toBe('completed')
    })

    it('rolls back the amount and status on failure', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { target_amount: 1000, current_amount: 950 }))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.addProgress('g1', 100)

      expect(store.goals[0].current_amount).toBe(950)
      expect(store.goals[0].status).toBe('active')
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to update progress', color: 'error' })
    })

    it('is a no-op for a goal that is not in the list', async () => {
      const store = useFinancialGoalStore()

      await store.addProgress('g404', 10)

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('toggleChecklist', () => {
    it('flips the checklist status optimistically and persists it', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { type: 'checklist' }))

      await store.toggleChecklist('g1')

      expect(apiMock).toHaveBeenCalledWith('financial-goals/g1', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'completed' })
      })
      expect(store.goals[0].status).toBe('completed')
      // NOTE: toggleChecklist intentionally shows no success toast
      expect(toastAddMock).not.toHaveBeenCalled()
    })

    it('flips a completed checklist back to active', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { type: 'checklist', status: 'completed' }))
      apiMock.mockResolvedValue(undefined)

      await store.toggleChecklist('g1')

      expect(apiMock).toHaveBeenCalledWith('financial-goals/g1', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'active' })
      })
      expect(store.goals[0].status).toBe('active')
    })

    it('is a no-op for a savings goal', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1'))

      await store.toggleChecklist('g1')

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.goals[0].status).toBe('active')
    })

    it('reverts the status on failure', async () => {
      const store = useFinancialGoalStore()
      store.goals.push(goal('g1', { type: 'checklist' }))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.toggleChecklist('g1')

      expect(store.goals[0].status).toBe('active')
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to update goal', color: 'error' })
    })
  })

  it('reset clears all goals', () => {
    const store = useFinancialGoalStore()
    store.goals.push(goal('g1'))

    store.reset()

    expect(store.goals).toHaveLength(0)
  })
})
