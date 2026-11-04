import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref } from 'vue'
import { useValidation } from '../../../app/composables/useValidation'

// The goals store relies on Vue auto-imports (ref), Nuxt auto-imports
// (defineStore, useNuxtApp, useToast) and the useValidation composable.
// Stub them all before the store module loads (defineStore is referenced at
// module top level, so the store must be imported dynamically below); the
// $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('defineStore', defineStore)
vi.stubGlobal('useValidation', useValidation)

const { useGoalStore } = await import('../../../app/stores/goals')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
})

const goal = (id: string, overrides: Record<string, unknown> = {}) => ({
  goal_id: id,
  title: `Goal ${id}`,
  description: '',
  due_date: null,
  priority: 1,
  status: 'active',
  color: '#3B82F6',
  progress: 0,
  total_tasks: 0,
  completed_tasks: 0,
  tasks: [],
  ...overrides
}) as any

const task = (id: string, overrides: Record<string, unknown> = {}) => ({
  task_id: id,
  title: `Task ${id}`,
  description: '',
  status: 'pending',
  order: 0,
  ...overrides
}) as any

describe('useGoalStore', () => {
  describe('fetchGoals', () => {
    it('populates goals and clears loading', async () => {
      const store = useGoalStore()
      apiMock.mockResolvedValue({ goals: [goal('g1'), goal('g2')] })

      await store.fetchGoals()

      expect(apiMock).toHaveBeenCalledWith('/goals')
      expect(store.goals).toEqual([goal('g1'), goal('g2')])
      expect(store.loading).toBe(false)
    })

    it('initializes a tasks array for goals without tasks', async () => {
      const store = useGoalStore()
      apiMock.mockResolvedValue({ goals: [goal('g1', { tasks: undefined })] })

      await store.fetchGoals()

      expect(Array.isArray(store.goals[0].tasks)).toBe(true)
      expect(store.goals[0].tasks).toEqual([])
    })

    it('propagates the error and leaves loading stuck on true', async () => {
      // NOTE (actual behavior): fetchGoals has no try/catch, so the error
      // escapes the action and loading is never reset.
      const store = useGoalStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await expect(store.fetchGoals()).rejects.toThrow('network down')

      expect(store.loading).toBe(true)
      expect(store.fetchError).toBeNull()
    })
  })

  describe('createGoal', () => {
    it('adds the goal optimistically, then replaces it with the server response', async () => {
      const store = useGoalStore()
      const serverGoal = goal('g1', { total_tasks: 2 })
      apiMock.mockResolvedValue(serverGoal)

      const payload = {
        title: 'Learn Go',
        description: 'A goal',
        due_date: null as string | null,
        priority: 3,
        color: '#FF0000'
      }
      await store.createGoal(payload)

      expect(apiMock).toHaveBeenCalledWith('/goals', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.goals).toHaveLength(1)
      expect(store.goals[0]).toEqual({ ...serverGoal, tasks: [] })
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Goal added successfully',
        color: 'success'
      })
    })

    it('defaults the color to blue and forces active status', async () => {
      const store = useGoalStore()
      apiMock.mockResolvedValue(goal('g1'))

      await store.createGoal({ title: 'Learn Go', description: '' })

      expect(store.goals[0].status).toBe('active')
      expect(store.goals[0].color).toBe('#3B82F6')
    })

    it('rejects an invalid payload without calling the API', async () => {
      const store = useGoalStore()

      await store.createGoal({ title: '', description: '' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.goals).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'title is required',
        color: 'error'
      })
    })

    it('rolls back the optimistic goal on failure', async () => {
      const store = useGoalStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.createGoal({ title: 'Learn Go', description: '' })

      expect(store.goals).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'network down',
        color: 'error'
      })
    })
  })

  describe('updateGoal', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { title: 'old' }))
      const payload = goal('g1', { title: 'new', progress: 40 })
      const serverGoal = goal('g1', { title: 'new', progress: 50 })
      apiMock.mockResolvedValue(serverGoal)

      await store.updateGoal(payload)

      expect(apiMock).toHaveBeenCalledWith('/goals/g1', {
        method: 'PATCH',
        body: JSON.stringify({
          title: 'new',
          description: '',
          due_date: null,
          priority: 1,
          status: 'active',
          color: '#3B82F6'
        })
      })
      // The server response replaces the goal entirely
      expect(store.goals[0]).toEqual(serverGoal)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Goal updated successfully',
        color: 'success'
      })
    })

    it('reverts to the original goal on failure', async () => {
      const store = useGoalStore()
      const original = goal('g1', { title: 'original' })
      store.goals.push(original)
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.updateGoal(goal('g1', { title: 'changed' }))

      expect(store.goals[0]).toEqual(original)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'conflict',
        color: 'error'
      })
    })

    it('still calls the API for a goal that is not in the list', async () => {
      const store = useGoalStore()
      apiMock.mockResolvedValue(goal('g99', { title: 'remote' }))

      await store.updateGoal(goal('g99', { title: 'remote' }))

      expect(apiMock).toHaveBeenCalledWith('/goals/g99', { method: 'PATCH', body: expect.any(String) })
      expect(store.goals).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Goal updated successfully',
        color: 'success'
      })
    })
  })

  describe('deleteGoal', () => {
    it('removes the goal optimistically and succeeds', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1'), goal('g2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteGoal('g1')

      expect(apiMock).toHaveBeenCalledWith('/goals/g1', { method: 'DELETE' })
      expect(store.goals.map(g => g.goal_id)).toEqual(['g2'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Goal deleted successfully',
        color: 'success'
      })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useGoalStore()

      await store.deleteGoal('g42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Goal not found',
        color: 'error'
      })
    })

    it('restores the goal when deletion fails', async () => {
      // NOTE (actual behavior): the restored goal is appended at the end of
      // the list rather than back at its original index.
      const store = useGoalStore()
      store.goals.push(goal('g1'), goal('g2'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteGoal('g1')

      expect(store.goals.map(g => g.goal_id)).toEqual(['g2', 'g1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Goal deletion failed',
        color: 'error'
      })
    })
  })

  describe('fetchGoalTasks', () => {
    it('stores fetched tasks on the goal and returns them', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1'))
      const tasks = [task('t1'), task('t2')]
      apiMock.mockResolvedValue({ tasks })

      const result = await store.fetchGoalTasks('g1')

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/tasks')
      expect(store.goals[0].tasks).toEqual(tasks)
      expect(result).toEqual(tasks)
    })

    it('shows a toast and returns an empty list on failure', async () => {
      const store = useGoalStore()
      apiMock.mockRejectedValue(new Error('server error'))

      const result = await store.fetchGoalTasks('g1')

      expect(result).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch goal tasks',
        color: 'error'
      })
    })

    it('returns the tasks without storing them when the goal is unknown', async () => {
      const store = useGoalStore()
      const tasks = [task('t1')]
      apiMock.mockResolvedValue({ tasks })

      const result = await store.fetchGoalTasks('g-missing')

      expect(result).toEqual(tasks)
      expect(store.goals).toHaveLength(0)
    })
  })

  describe('addTaskToGoal', () => {
    it('adds the task optimistically, then replaces it with the server response', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1'))
      const serverTask = task('t1', { status: 'pending', order: 1 })
      apiMock.mockResolvedValue(serverTask)

      const result = await store.addTaskToGoal('g1', { title: 'Read docs', priority: 2 })

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/tasks', {
        method: 'POST',
        body: JSON.stringify({ title: 'Read docs', priority: 2 })
      })
      expect(store.goals[0].tasks).toEqual([serverTask])
      expect(result).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Task added to goal successfully',
        color: 'success'
      })
    })

    it('throws and shows a toast when the goal is unknown', async () => {
      const store = useGoalStore()

      await expect(store.addTaskToGoal('g-missing', { title: 'Read docs' }))
        .rejects.toThrow('Goal not found')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Goal not found',
        color: 'error'
      })
    })

    it('removes the optimistic task and rethrows on failure', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t0')] }))
      apiMock.mockRejectedValue(new Error('server error'))

      await expect(store.addTaskToGoal('g1', { title: 'Read docs' }))
        .rejects.toThrow('server error')

      expect(store.goals[0].tasks.map(t => t.task_id)).toEqual(['t0'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to add task to goal',
        color: 'error'
      })
    })
  })

  describe('updateGoalTask', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t1')] }))
      const serverTask = task('t1', { title: 'Renamed', status: 'pending' })
      apiMock.mockResolvedValue(serverTask)

      const result = await store.updateGoalTask('g1', 't1', { title: 'Renamed' })

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/tasks/t1', {
        method: 'PATCH',
        body: JSON.stringify({ title: 'Renamed' })
      })
      expect(store.goals[0].tasks[0]).toEqual(serverTask)
      expect(result).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Task updated successfully',
        color: 'success'
      })
    })

    it('throws without calling the API when the task is unknown', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [] }))

      await expect(store.updateGoalTask('g1', 't-missing', { title: 'x' }))
        .rejects.toThrow('Task not found')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Task not found',
        color: 'error'
      })
    })

    it('reverts the task and rethrows on failure', async () => {
      const store = useGoalStore()
      const original = task('t1', { status: 'pending' })
      store.goals.push(goal('g1', { tasks: [original] }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await expect(store.updateGoalTask('g1', 't1', { status: 'completed' }))
        .rejects.toThrow('conflict')

      expect(store.goals[0].tasks[0].status).toBe('pending')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to update task',
        color: 'error'
      })
    })
  })

  describe('deleteGoalTask', () => {
    it('removes the task optimistically and succeeds', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t1'), task('t2')] }))
      apiMock.mockResolvedValue(undefined)

      await store.deleteGoalTask('g1', 't1')

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/tasks/t1', { method: 'DELETE' })
      expect(store.goals[0].tasks.map(t => t.task_id)).toEqual(['t2'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Task deleted successfully',
        color: 'success'
      })
    })

    it('restores the task and rethrows on failure', async () => {
      const store = useGoalStore()
      const original = task('t1')
      store.goals.push(goal('g1', { tasks: [original] }))
      apiMock.mockRejectedValue(new Error('server error'))

      await expect(store.deleteGoalTask('g1', 't1')).rejects.toThrow('server error')

      expect(store.goals[0].tasks).toEqual([original])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to delete task',
        color: 'error'
      })
    })
  })

  describe('completeGoalTask', () => {
    it('toggles a pending task to completed and reports success', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t1', { status: 'pending' })] }))
      const serverTask = task('t1', { status: 'completed' })
      apiMock.mockResolvedValue(serverTask)

      const result = await store.completeGoalTask('g1', 't1')

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/tasks/t1/complete', { method: 'PATCH' })
      expect(store.goals[0].tasks[0]).toEqual(serverTask)
      expect(result).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Task completed',
        color: 'success'
      })
    })

    it('toggles a completed task back to pending', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t1', { status: 'completed' })] }))
      apiMock.mockResolvedValue(task('t1', { status: 'pending' }))

      await store.completeGoalTask('g1', 't1')

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Goal',
        description: 'Task marked as pending',
        color: 'success'
      })
    })

    it('reverts the optimistic status and rethrows on failure', async () => {
      const store = useGoalStore()
      store.goals.push(goal('g1', { tasks: [task('t1', { status: 'pending' })] }))
      apiMock.mockRejectedValue(new Error('server error'))

      await expect(store.completeGoalTask('g1', 't1')).rejects.toThrow('server error')

      expect(store.goals[0].tasks[0].status).toBe('pending')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to update task status',
        color: 'error'
      })
    })
  })

  describe('fetchGoalTaskRecommendations', () => {
    it('returns the AI recommendations response', async () => {
      const store = useGoalStore()
      const response = {
        recommendations: [{ title: 'Task A', description: '', priority: 2, estimated_hours: 1, reasoning: 'why' }],
        goal: goal('g1')
      }
      apiMock.mockResolvedValue(response)

      const result = await store.fetchGoalTaskRecommendations('g1')

      expect(apiMock).toHaveBeenCalledWith('/goals/g1/ai/recommendations')
      expect(result).toEqual(response)
      expect(toastAddMock).not.toHaveBeenCalled()
    })

    it('shows a toast and rethrows on failure', async () => {
      const store = useGoalStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await expect(store.fetchGoalTaskRecommendations('g1')).rejects.toThrow('server error')

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch AI recommendations',
        color: 'error'
      })
    })
  })

  it('reset clears all goals', () => {
    const store = useGoalStore()
    store.goals.push(goal('g1'))

    store.reset()

    expect(store.goals).toHaveLength(0)
  })
})
