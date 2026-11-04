import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Task } from '../../../app/types/task'
import { useValidation } from '../../../app/composables/useValidation'

// The tasks store relies on Vue auto-imports (ref/computed), Pinia's
// defineStore, and Nuxt auto-imports (useNuxtApp, useToast, useValidation,
// useOffline). Stub them all before the store module loads (defineStore is
// referenced at module top level, so the store must be imported dynamically
// below); the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()
const storeDataLocallyMock = vi.fn()
const getLocalDataMock = vi.fn()
const syncPendingOperationsMock = vi.fn()
// Ref so tests can flip the online state; the store captures this instance.
const isOnlineRef = ref(true)

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('computed', computed)
vi.stubGlobal('defineStore', defineStore)
vi.stubGlobal('useValidation', useValidation)
vi.stubGlobal('useOffline', () => ({
  isOnline: isOnlineRef,
  storeDataLocally: storeDataLocallyMock,
  getLocalData: getLocalDataMock,
  syncPendingOperations: syncPendingOperationsMock
}))

const { useTaskStore } = await import('../../../app/stores/tasks')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  storeDataLocallyMock.mockReset()
  syncPendingOperationsMock.mockReset()
  getLocalDataMock.mockReset()
  // Empty cache by default; tests that need cached data override this.
  getLocalDataMock.mockResolvedValue([])
  isOnlineRef.value = true
})

const task = (id: string, overrides: Partial<Task> = {}): Task => ({
  task_id: id,
  title: `task-${id}`,
  description: '',
  status: 'pending',
  order: 0,
  time_spent_minutes: 0,
  is_recurring: false,
  ...overrides
} as Task)

describe('useTaskStore', () => {
  describe('initial state', () => {
    it('starts with empty tasks and no loading state', () => {
      const store = useTaskStore()

      expect(store.tasks).toEqual([])
      expect(store.loading).toBe(false)
      expect(store.fetchError).toBeNull()
      expect(store.pendingOperations.size).toBe(0)
      expect(store.completedTasks).toEqual([])
      expect(store.unscheduledTasks).toEqual([])
      expect(store.aiTaskPreview).toEqual([])
    })
  })

  describe('fetchTasks', () => {
    it('populates tasks, stores them locally and clears loading', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [task('1'), task('2')] })

      await store.fetchTasks()

      expect(apiMock).toHaveBeenCalledWith('/tasks')
      expect(store.tasks.map(t => t.task_id)).toEqual(['1', '2'])
      expect(storeDataLocallyMock).toHaveBeenCalledWith('tasks', [task('1'), task('2')])
      expect(store.loading).toBe(false)
    })

    it('builds a filtered query string, skipping empty values', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [] })

      await store.fetchTasks({ status: 'pending', priority: 2, search: '', goal_id: undefined })

      expect(apiMock).toHaveBeenCalledWith('/tasks?status=pending&priority=2')
    })

    it('applies cached tasks instantly while the network request is in flight', async () => {
      const store = useTaskStore()
      const cached = [task('c1')]
      getLocalDataMock.mockResolvedValue(cached)
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchTasks()
      await new Promise(resolve => setTimeout(resolve, 0))
      expect(store.tasks).toEqual(cached)

      resolveApi({ tasks: [task('s1')] })
      await pending
      expect(store.tasks.map(t => t.task_id)).toEqual(['s1'])
    })

    it('toasts and uses cached data when offline', async () => {
      const store = useTaskStore()
      isOnlineRef.value = false
      getLocalDataMock.mockResolvedValue([task('c1')])

      await store.fetchTasks()

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks.map(t => t.task_id)).toEqual(['c1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: "You're offline. Showing cached data.",
        color: 'info'
      })
      expect(store.loading).toBe(false)
    })

    it('records the error and clears loading when the request fails', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchTasks()

      expect(store.fetchError).toBeInstanceOf(Error)
      expect(store.loading).toBe(false)
    })
  })

  describe('submitForm', () => {
    it('validates the payload before touching the API', async () => {
      const store = useTaskStore()

      await store.submitForm({ title: '', description: '' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'title is required',
        color: 'error'
      })
    })

    it('rejects a due date in the past', async () => {
      const store = useTaskStore()

      await store.submitForm({ title: 'late', description: '', due_date: '2020-01-01' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Due date cannot be in the past',
        color: 'error'
      })
    })

    // NOTE (bug): schemas.task.naturalLanguage validates the field
    // `natural_language`, but submitForm sends `natural_language_input`.
    // The payload can never satisfy the schema, so natural-language task
    // creation always fails validation on the client and never reaches the API.
    it('fails validation for natural-language input because the schema key does not match the payload key', async () => {
      const store = useTaskStore()

      await store.submitForm({
        title: 'pay rent',
        description: '',
        use_natural_language: true,
        natural_language_input: 'pay rent tomorrow morning please'
      })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'natural_language is required',
        color: 'error'
      })
    })

    it('adds the task optimistically, then replaces it with the server response', async () => {
      const store = useTaskStore()
      const serverTask = task('7', { title: 'work', status: 'pending' })
      apiMock.mockResolvedValue(serverTask)

      await store.submitForm({ title: 'work', description: 'desc', priority: 3 })

      expect(apiMock).toHaveBeenCalledWith('tasks', {
        method: 'POST',
        body: JSON.stringify({ title: 'work', description: 'desc', priority: 3 })
      })
      expect(store.tasks).toEqual([serverTask])
      expect(store.pendingOperations.size).toBe(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task added successfully',
        color: 'success'
      })
    })

    it('rolls back the optimistic task on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ title: 'work', description: '' })

      expect(store.tasks).toHaveLength(0)
      expect(store.pendingOperations.size).toBe(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'network down',
        color: 'error'
      })
    })

    // NOTE (bug): the offline branch calls the store's local
    // `addPendingOperation(operationId: string)` helper with the queue payload
    // object instead of the offline plugin's queue function (it was never
    // destructured from useOffline), so the operation is dropped into the
    // string Set and never queued for sync.
    it('keeps the optimistic task and queues locally when offline', async () => {
      const store = useTaskStore()
      isOnlineRef.value = false

      await store.submitForm({ title: 'work', description: '' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks).toHaveLength(1)
      expect(store.tasks[0].title).toBe('work')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task saved locally. Will sync when online.',
        color: 'info'
      })
      // Documents the misdirected call: the payload object ends up in the
      // string Set instead of being queued on the offline plugin.
      expect(store.pendingOperations.size).toBe(1)
    })
  })

  describe('editTask', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'))
      const serverTask = task('1', { title: 'renamed' })
      apiMock.mockResolvedValue(serverTask)

      await store.editTask({ task_id: '1', title: 'renamed' })

      expect(apiMock).toHaveBeenCalledWith('tasks/1', {
        method: 'PATCH',
        body: JSON.stringify({ task_id: '1', title: 'renamed' })
      })
      expect(store.tasks[0]).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Edited task succesfully',
        color: 'success'
      })
    })

    it('reverts to the original task on failure', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1', { title: 'original' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editTask({ task_id: '1', title: 'changed' })

      expect(store.tasks[0].title).toBe('original')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Editing task failed',
        color: 'error'
      })
    })

    it('updates locally without an API call when offline', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1', { title: 'original' }))
      isOnlineRef.value = false

      await store.editTask({ task_id: '1', title: 'renamed' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks[0].title).toBe('renamed')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task updated locally. Will sync when online.',
        color: 'info'
      })
    })

    it('leaves local state untouched when the task is unknown', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue(task('99'))

      await store.editTask({ task_id: '99', title: 'remote' })

      expect(apiMock).toHaveBeenCalledWith('tasks/99', expect.anything())
      expect(store.tasks).toHaveLength(0)
    })
  })

  describe('reorderTask', () => {
    it('applies the new order and sorts by it', async () => {
      const store = useTaskStore()
      store.tasks.push(task('a', { order: 0 }), task('b', { order: 1 }))
      apiMock.mockResolvedValue(undefined)

      await store.reorderTask([{ task_id: 'b', order: 0 }, { task_id: 'a', order: 1 }])

      expect(apiMock).toHaveBeenCalledWith('/tasks/reorder', {
        method: 'PUT',
        body: JSON.stringify({ task_orders: [{ task_id: 'b', order: 0 }, { task_id: 'a', order: 1 }] })
      })
      expect(store.tasks.map(t => t.task_id)).toEqual(['b', 'a'])
      expect(store.tasks[0].order).toBe(0)
      expect(store.tasks[1].order).toBe(1)
    })

    it('restores the original order when the request fails', async () => {
      const store = useTaskStore()
      store.tasks.push(task('a', { order: 0 }), task('b', { order: 1 }))
      apiMock.mockRejectedValue(new Error('boom'))

      await store.reorderTask([{ task_id: 'b', order: 0 }])

      expect(store.tasks.map(t => t.task_id)).toEqual(['a', 'b'])
      expect(store.tasks[0].order).toBe(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'boom',
        color: 'error'
      })
    })
  })

  describe('deleteTask', () => {
    it('toasts and does not call the API for an unknown id', async () => {
      const store = useTaskStore()

      await store.deleteTask('42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Task not found',
        color: 'error'
      })
    })

    it('removes the task optimistically and succeeds', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'), task('2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTask('1')

      expect(apiMock).toHaveBeenCalledWith('tasks/1', { method: 'DELETE' })
      expect(store.tasks.map(t => t.task_id)).toEqual(['2'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task deleted successfully',
        color: 'success'
      })
    })

    it('restores the task when deletion fails', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteTask('1')

      expect(store.tasks.map(t => t.task_id)).toEqual(['1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'server error',
        color: 'error'
      })
    })

    it('removes locally and queues when offline', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'))
      isOnlineRef.value = false

      await store.deleteTask('1')

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.tasks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task deleted locally. Will sync when online.',
        color: 'info'
      })
    })
  })

  describe('completeTask', () => {
    // NOTE (bug): completeTask only PATCHes the status; it never updates the
    // local task's status or shows a toast, so the store list can drift from
    // the server until the next fetch.
    it('patches the status without touching local state or toasting', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'))
      apiMock.mockResolvedValue(undefined)

      await store.completeTask(task('1', { status: 'completed' }))

      expect(apiMock).toHaveBeenCalledWith('tasks/1', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'completed' })
      })
      expect(store.tasks[0].status).toBe('pending')
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('undoDeleteTask', () => {
    it('restores the task and returns it', async () => {
      const store = useTaskStore()
      const restored = task('1')
      apiMock.mockResolvedValue(restored)

      const result = await store.undoDeleteTask('1')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/undo-delete', { method: 'PATCH' })
      expect(store.tasks).toEqual([restored])
      expect(result).toEqual(restored)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task restored successfully',
        color: 'success'
      })
    })

    it('returns null and toasts when offline', async () => {
      const store = useTaskStore()
      isOnlineRef.value = false

      const result = await store.undoDeleteTask('1')

      expect(apiMock).not.toHaveBeenCalled()
      expect(result).toBeNull()
      // NOTE (bug): the message contains a typo ("deleteion").
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Cannot undo task deleteion while offline',
        color: 'error'
      })
    })

    it('returns null on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('gone'))

      const result = await store.undoDeleteTask('1')

      expect(result).toBeNull()
      expect(store.tasks).toHaveLength(0)
    })
  })

  describe('getRecentlyDeletedTasks', () => {
    it('returns the tasks from the API', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [task('1')] })

      const result = await store.getRecentlyDeletedTasks()

      expect(apiMock).toHaveBeenCalledWith('tasks/recently-deleted')
      expect(result).toEqual([task('1')])
    })

    it('returns an empty list and toasts when offline', async () => {
      const store = useTaskStore()
      isOnlineRef.value = false

      const result = await store.getRecentlyDeletedTasks()

      expect(apiMock).not.toHaveBeenCalled()
      expect(result).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Cannot fetch deleted tasks while offline',
        color: 'error'
      })
    })
  })

  describe('createSubtask', () => {
    it('throws when the parent task is missing', async () => {
      const store = useTaskStore()

      await expect(store.createSubtask('nope', { title: 'sub' })).rejects.toThrow('Parent task not found')
      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Parent task not found',
        color: 'error'
      })
    })

    it('adds the subtask optimistically, then swaps in the server response', async () => {
      const store = useTaskStore()
      const serverSubtask = task('s1', { parent_task_id: '1' })
      apiMock.mockResolvedValue(serverSubtask)
      store.tasks.push(task('1'))

      const result = await store.createSubtask('1', { title: 'sub' })

      expect(apiMock).toHaveBeenCalledWith('tasks', {
        method: 'POST',
        body: JSON.stringify({ title: 'sub', parent_task_id: '1' })
      })
      expect(store.tasks[0].subtasks).toEqual([serverSubtask])
      expect(result).toEqual(serverSubtask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Subtask added successfully',
        color: 'success'
      })
    })

    it('removes the optimistic subtask and rethrows on failure', async () => {
      const store = useTaskStore()
      store.tasks.push(task('1'))
      apiMock.mockRejectedValue(new Error('boom'))

      await expect(store.createSubtask('1', { title: 'sub' })).rejects.toThrow('boom')

      expect(store.tasks[0].subtasks).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'boom',
        color: 'error'
      })
    })
  })

  describe('getTaskSubtasks', () => {
    it('stores subtasks on the task and returns them', async () => {
      const store = useTaskStore()
      const subtasks = [task('s1', { parent_task_id: '1' })]
      apiMock.mockResolvedValue({ subtasks })
      store.tasks.push(task('1'))

      const result = await store.getTaskSubtasks('1')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/subtasks')
      expect(store.tasks[0].subtasks).toEqual(subtasks)
      expect(result).toEqual(subtasks)
    })

    it('returns an empty list and toasts on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      const result = await store.getTaskSubtasks('1')

      expect(result).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch subtasks',
        color: 'error'
      })
    })
  })

  describe('dependencies', () => {
    it('creates a dependency and toasts success', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue(undefined)

      await store.createTaskDependency('1', '2')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/dependencies', {
        method: 'POST',
        body: JSON.stringify({ depends_on_id: '2' })
      })
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Dependency created successfully',
        color: 'success'
      })
    })

    it('rethrows and toasts on dependency creation failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('cycle'))

      await expect(store.createTaskDependency('1', '2')).rejects.toThrow('cycle')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'cycle',
        color: 'error'
      })
    })

    it('returns dependencies from the API', async () => {
      const store = useTaskStore()
      const data = { dependencies: [task('2')], dependents: [task('3')] }
      apiMock.mockResolvedValue(data)

      const result = await store.getTaskDependencies('1')

      expect(result).toEqual(data)
    })

    it('returns empty lists and toasts when fetching dependencies fails', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      const result = await store.getTaskDependencies('1')

      expect(result).toEqual({ dependencies: [], dependents: [] })
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch dependencies',
        color: 'error'
      })
    })

    it('deletes a dependency and toasts success', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue(undefined)

      await store.deleteTaskDependency('1', '2')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/dependencies/2', { method: 'DELETE' })
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Dependency deleted successfully',
        color: 'success'
      })
    })
  })

  describe('time tracking', () => {
    it('starts time tracking and returns the time entry', async () => {
      const store = useTaskStore()
      const entry = { time_entry_id: 'te1', task_id: '1', is_running: true }
      apiMock.mockResolvedValue(entry)

      const result = await store.startTimeTracking('1', 'working on it')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/time/start', {
        method: 'POST',
        body: JSON.stringify({ description: 'working on it' })
      })
      expect(result).toEqual(entry)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Time tracking started',
        color: 'success'
      })
    })

    it('rethrows and toasts when starting fails', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      await expect(store.startTimeTracking('1')).rejects.toThrow('boom')
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Failed to start time tracking',
        color: 'error'
      })
    })

    it('stops time tracking and returns the time entry', async () => {
      const store = useTaskStore()
      const entry = { time_entry_id: 'te1', task_id: '1', is_running: false }
      apiMock.mockResolvedValue(entry)

      const result = await store.stopTimeTracking('1')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/time/stop', { method: 'POST' })
      expect(result).toEqual(entry)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Time tracking stopped',
        color: 'success'
      })
    })

    it('returns time entries from the API', async () => {
      const store = useTaskStore()
      const entries = [{ time_entry_id: 'te1', task_id: '1' }]
      apiMock.mockResolvedValue({ time_entries: entries })

      const result = await store.getTaskTimeEntries('1')

      expect(apiMock).toHaveBeenCalledWith('tasks/1/time')
      expect(result).toEqual(entries)
    })
  })

  describe('task templates', () => {
    it('creates a template and returns it', async () => {
      const store = useTaskStore()
      const template = { template_id: 'tpl1', name: 'Weekly review' }
      apiMock.mockResolvedValue(template)

      const result = await store.createTaskTemplate({ name: 'Weekly review', title_template: 'Review' })

      expect(apiMock).toHaveBeenCalledWith('task-templates', {
        method: 'POST',
        body: JSON.stringify({ name: 'Weekly review', title_template: 'Review' })
      })
      expect(result).toEqual(template)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task template created successfully',
        color: 'success'
      })
    })

    it('fetches templates from the API', async () => {
      const store = useTaskStore()
      const templates = [{ template_id: 'tpl1' }]
      apiMock.mockResolvedValue({ templates })

      const result = await store.getTaskTemplates()

      expect(result).toEqual(templates)
    })

    it('creates a task from a template, replacing the optimistic task', async () => {
      const store = useTaskStore()
      const serverTask = task('9', { title: 'From template' })
      apiMock.mockResolvedValue(serverTask)

      const result = await store.createTaskFromTemplate('tpl1', { title: 'From template', priority: 2 })

      expect(apiMock).toHaveBeenCalledWith('task-templates/tpl1/create-task', {
        method: 'POST',
        body: JSON.stringify({ title: 'From template', priority: 2 })
      })
      expect(store.tasks).toEqual([serverTask])
      expect(result).toEqual(serverTask)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Task created from template successfully',
        color: 'success'
      })
    })

    it('falls back to a default optimistic title when no task data is given', async () => {
      const store = useTaskStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.createTaskFromTemplate('tpl1')
      await new Promise(resolve => setTimeout(resolve, 0))
      // Optimistic task uses a default title when taskData is empty
      expect(store.tasks[0].title).toBe('New Task')

      resolveApi(task('9'))
      await pending
      expect(apiMock).toHaveBeenCalledWith('task-templates/tpl1/create-task', {
        method: 'POST',
        body: JSON.stringify({})
      })
    })

    it('rolls back the optimistic task and rethrows on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      await expect(store.createTaskFromTemplate('tpl1', { title: 'x' })).rejects.toThrow('boom')

      expect(store.tasks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to create task from template',
        color: 'error'
      })
    })
  })

  describe('recurrence rules', () => {
    it('creates a recurrence rule and returns it', async () => {
      const store = useTaskStore()
      const rule = { recurrence_rule_id: 'r1', name: 'daily', frequency: 'daily' }
      apiMock.mockResolvedValue(rule)

      const result = await store.createRecurrenceRule({ name: 'daily', frequency: 'daily', title_template: 'Daily' })

      expect(apiMock).toHaveBeenCalledWith('recurrence-rules', {
        method: 'POST',
        body: JSON.stringify({ name: 'daily', frequency: 'daily', title_template: 'Daily' })
      })
      expect(result).toEqual(rule)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Recurrence rule created successfully',
        color: 'success'
      })
    })

    it('fetches recurrence rules from the API', async () => {
      const store = useTaskStore()
      const rules = [{ recurrence_rule_id: 'r1' }]
      apiMock.mockResolvedValue({ recurrence_rules: rules })

      const result = await store.getRecurrenceRules()

      expect(result).toEqual(rules)
    })

    it('generates recurring tasks, replacing the optimistic ones', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [task('r1', { is_recurring: true }), task('r2', { is_recurring: true })] })

      const result = await store.generateRecurringTasks('r1', 2)

      expect(apiMock).toHaveBeenCalledWith('recurrence-rules/r1/generate-tasks?count=2', { method: 'POST' })
      expect(store.tasks.map(t => t.task_id)).toEqual(['r1', 'r2'])
      expect(result).toHaveLength(2)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Generated 2 recurring tasks',
        color: 'success'
      })
    })

    it('removes all optimistic recurring tasks and rethrows on failure', async () => {
      const store = useTaskStore()
      store.tasks.push(task('keep'))
      apiMock.mockRejectedValue(new Error('boom'))

      await expect(store.generateRecurringTasks('r1', 3)).rejects.toThrow('boom')

      expect(store.tasks.map(t => t.task_id)).toEqual(['keep'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to generate recurring tasks',
        color: 'error'
      })
    })
  })

  describe('syncOfflineChanges', () => {
    it('syncs pending operations and refreshes when online', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [] })

      await store.syncOfflineChanges()

      expect(syncPendingOperationsMock).toHaveBeenCalled()
      expect(apiMock).toHaveBeenCalledWith('/tasks')
    })

    it('does nothing when offline', async () => {
      const store = useTaskStore()
      isOnlineRef.value = false

      await store.syncOfflineChanges()

      expect(syncPendingOperationsMock).not.toHaveBeenCalled()
      expect(apiMock).not.toHaveBeenCalled()
    })
  })

  describe('AI task suggestions', () => {
    it('stores the preview and returns the response', async () => {
      const store = useTaskStore()
      const preview = [{
        task_id: '1',
        original_title: 't',
        title: 'better t',
        description: 'd',
        priority: 3,
        estimated_hours: 2,
        subtasks: []
      }]
      const response = { preview, message: 'ok' }
      apiMock.mockResolvedValue(response)

      const result = await store.getAITaskPreview()

      expect(apiMock).toHaveBeenCalledWith('/tasks/ai-check', { method: 'POST' })
      expect(store.aiTaskPreview).toEqual(preview)
      expect(result).toEqual(response)
    })

    it('returns null and toasts on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      const result = await store.getAITaskPreview()

      expect(result).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to get AI recommendations',
        color: 'error'
      })
    })

    it('applies the selection, clears the preview and refreshes tasks', async () => {
      const store = useTaskStore()
      apiMock.mockResolvedValue({ tasks: [task('1')] })

      const result = await store.applyAITasks([{ task_id: '1', selected: true }])

      expect(apiMock).toHaveBeenCalledWith('/tasks/ai-check/apply', {
        method: 'POST',
        body: JSON.stringify({ applied_tasks: [{ task_id: '1', selected: true }] })
      })
      expect(result).toBe(true)
      expect(store.aiTaskPreview).toEqual([])
      // fetchTasks refresh happens after applying
      expect(apiMock).toHaveBeenCalledWith('/tasks')
      expect(store.tasks.map(t => t.task_id)).toEqual(['1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Tasks updated successfully',
        color: 'success'
      })
    })

    it('returns false and toasts on failure', async () => {
      const store = useTaskStore()
      apiMock.mockRejectedValue(new Error('boom'))

      const result = await store.applyAITasks([])

      expect(result).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to apply changes',
        color: 'error'
      })
    })
  })

  describe('computed getters', () => {
    it('completedTasks filters by completed status', () => {
      const store = useTaskStore()
      store.tasks.push(
        task('1', { status: 'completed' }),
        task('2'),
        task('3', { status: 'completed' })
      )

      expect(store.completedTasks.map(t => t.task_id)).toEqual(['1', '3'])
    })

    it('unscheduledTasks filters tasks without a start_time', () => {
      const store = useTaskStore()
      store.tasks.push(
        task('1', { start_time: new Date('2026-01-01') }),
        task('2', { start_time: null as any })
      )

      expect(store.unscheduledTasks.map(t => t.task_id)).toEqual(['2'])
    })
  })

  it('reset clears all tasks', () => {
    const store = useTaskStore()
    store.tasks.push(task('1'))

    store.reset()

    expect(store.tasks).toHaveLength(0)
  })
})
