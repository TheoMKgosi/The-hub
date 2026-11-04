import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref } from 'vue'

// The task-learning store relies on Vue auto-imports (ref) and Nuxt
// auto-imports (useNuxtApp, useToast). Stub them all before the store module
// loads; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('defineStore', defineStore)

const { useTaskLearningStore } = await import('../../../app/stores/task-learning')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
})

afterEach(() => {
  vi.restoreAllMocks()
})

const task = (id: string, overrides: Record<string, unknown> = {}) => ({
  task_learning_id: id,
  topic_id: 'topic-1',
  title: `task-${id}`,
  notes: '',
  status: 'not_started' as const,
  resources: [],
  ...overrides
})

describe('useTaskLearningStore', () => {
  it('starts with no tasks and no loading', () => {
    const store = useTaskLearningStore()

    expect(store.tasks).toEqual([])
    expect(store.loading).toBe(false)
  })

  describe('fetchTasks', () => {
    it('populates tasks for the topic and clears loading', async () => {
      const store = useTaskLearningStore()
      const serverTasks = [task('t1'), task('t2')]
      apiMock.mockResolvedValue({ task_learnings: serverTasks })

      await store.fetchTasks('topic-1')

      expect(apiMock).toHaveBeenCalledWith('/topics/topic-1/task-learnings')
      expect(store.tasks).toEqual(serverTasks)
      expect(store.loading).toBe(false)
    })

    it('shows an error toast on failure and clears loading', async () => {
      const store = useTaskLearningStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchTasks('topic-1')

      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch tasks.' })
      expect(store.loading).toBe(false)
    })
  })

  describe('createTask', () => {
    it('adds the task optimistically, then replaces it with the server response', async () => {
      const store = useTaskLearningStore()
      const serverTask = task('t1', { title: 'New task' })
      apiMock.mockResolvedValue(serverTask)

      await store.createTask('topic-1', { topic_id: 'topic-1', title: 'New task' })

      expect(apiMock).toHaveBeenCalledWith('/task-learnings', {
        method: 'POST',
        body: JSON.stringify({ topic_id: 'topic-1', title: 'New task' })
      })
      // Optimistic temp entry is swapped for the confirmed task
      expect(store.tasks).toHaveLength(1)
      expect(store.tasks[0]).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Task created!' })
    })

    it('rolls back the optimistic task on failure', async () => {
      const store = useTaskLearningStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.createTask('topic-1', { topic_id: 'topic-1', title: 'New task' })

      expect(store.tasks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to create task.' })
    })
  })

  describe('updateTask', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useTaskLearningStore()
      store.tasks.push(task('t1'))
      const serverTask = task('t1', { status: 'completed' })
      apiMock.mockResolvedValue(serverTask)

      await store.updateTask('topic-1', 't1', { status: 'completed' })

      expect(apiMock).toHaveBeenCalledWith('/task-learnings/t1', {
        method: 'PUT',
        body: JSON.stringify({ status: 'completed' })
      })
      expect(store.tasks[0]).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Task updated!' })
    })

    it('reverts to the original task on failure', async () => {
      const store = useTaskLearningStore()
      store.tasks.push(task('t1', { status: 'not_started' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.updateTask('topic-1', 't1', { status: 'completed' })

      expect(store.tasks[0].status).toBe('not_started')
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to update task.' })
    })

    it('does not mutate state for a task that is not in the list', async () => {
      const store = useTaskLearningStore()
      apiMock.mockResolvedValue(task('t9'))

      await store.updateTask('topic-1', 't9', { status: 'completed' })

      expect(store.tasks).toHaveLength(0)
    })
  })

  describe('deleteTask', () => {
    it('removes the task optimistically and succeeds', async () => {
      const store = useTaskLearningStore()
      store.tasks.push(task('t1'), task('t2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTask('topic-1', 't1')

      expect(apiMock).toHaveBeenCalledWith('/task-learnings/t1', { method: 'DELETE' })
      expect(store.tasks.map(t => t.task_learning_id)).toEqual(['t2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Task deleted.' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useTaskLearningStore()

      await store.deleteTask('topic-1', 'missing')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Task not found.' })
    })

    it('restores the task when deletion fails', async () => {
      const store = useTaskLearningStore()
      store.tasks.push(task('t1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteTask('topic-1', 't1')

      expect(store.tasks.map(t => t.task_learning_id)).toEqual(['t1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to delete task.' })
    })
  })

  describe('completeTask', () => {
    it('applies the status optimistically, then reconciles with the server response', async () => {
      const store = useTaskLearningStore()
      const optimistic = task('t1', { status: 'not_started' })
      store.tasks.push(optimistic)
      const serverTask = task('t1', { status: 'completed' })
      apiMock.mockResolvedValue(serverTask)

      await store.completeTask('topic-1', task('t1', { status: 'completed' }))

      expect(apiMock).toHaveBeenCalledWith('/task-learnings/t1', {
        method: 'PUT',
        body: JSON.stringify({ status: 'completed' })
      })
      expect(store.tasks[0]).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Task status updated!', color: 'success' })
    })

    it('reverts the status on failure', async () => {
      const store = useTaskLearningStore()
      store.tasks.push(task('t1', { status: 'not_started' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.completeTask('topic-1', task('t1', { status: 'completed' }))

      expect(store.tasks[0].status).toBe('not_started')
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to update task status', color: 'error' })
    })
  })
})
