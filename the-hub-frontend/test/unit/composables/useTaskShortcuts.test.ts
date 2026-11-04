import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { reactive } from 'vue'
import type { Task } from '../../../app/types/task'

// useTaskShortcuts relies on Nuxt auto-imports (useTaskStore/useToast/
// useRoute/navigateTo). Stub them with a lightweight store mock so the
// composable works in the plain node unit environment. The module holds
// shared selection state, so it is re-imported per test via resetModules.
// reactive() is used so refs are unwrapped like in a real Pinia store.
const toastMock = { add: vi.fn() }
const storeMock = reactive({
  tasks: [] as Task[],
  loading: false,
  fetchTasks: vi.fn(),
  editTask: vi.fn(),
  deleteTask: vi.fn(),
  undoDeleteTask: vi.fn(),
  getRecentlyDeletedTasks: vi.fn(),
})
const route = { path: '/plan' }
const navigateTo = vi.fn()
const dispatchEvent = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useTaskStore', () => storeMock)
  vi.stubGlobal('useToast', () => toastMock)
  vi.stubGlobal('useRoute', () => route)
  vi.stubGlobal('navigateTo', navigateTo)

  // The composable dispatches CustomEvents on window; provide fakes for
  // the node environment when they are missing.
  if (typeof CustomEvent === 'undefined') {
    vi.stubGlobal('CustomEvent', class {
      type: string
      detail: unknown
      constructor(type: string, opts?: { detail?: unknown }) {
        this.type = type
        this.detail = opts?.detail ?? {}
      }
    })
  }
  vi.stubGlobal('window', { dispatchEvent, addEventListener: vi.fn(), removeEventListener: vi.fn() })
})

const makeTask = (id: string, overrides: Partial<Task> = {}): Task => ({
  task_id: id,
  title: `Task ${id}`,
  description: '',
  status: 'pending',
  start_time: new Date(),
  time_spent_minutes: 0,
  is_recurring: false,
  ...overrides
})

let useTaskShortcuts: typeof import('../../../app/composables/useTaskShortcuts')['useTaskShortcuts']

beforeEach(async () => {
  vi.clearAllMocks()
  route.path = '/plan'
  storeMock.tasks = [makeTask('t1'), makeTask('t2'), makeTask('t3')]
  storeMock.loading = false
  // Fresh module state (selection index, pending flags, registered provider)
  vi.resetModules()
  ;({ useTaskShortcuts } = await import('../../../app/composables/useTaskShortcuts'))
})

describe('useTaskShortcuts', () => {
  describe('moveSelection', () => {
    it('moves the selection down/up and clamps at the bounds', () => {
      const { moveSelection, selectedTaskIndex, selectedTaskId } = useTaskShortcuts()

      moveSelection(1)
      expect(selectedTaskIndex.value).toBe(1)
      expect(selectedTaskId.value).toBe('t2')

      moveSelection(10)
      expect(selectedTaskIndex.value).toBe(2)
      expect(selectedTaskId.value).toBe('t3')

      moveSelection(-10)
      expect(selectedTaskIndex.value).toBe(0)
      expect(selectedTaskId.value).toBe('t1')
    })

    it('lazily fetches tasks the first time a shortcut is used with an empty store', () => {
      storeMock.tasks = []
      const { moveSelection } = useTaskShortcuts()

      moveSelection(1)
      expect(storeMock.fetchTasks).toHaveBeenCalledTimes(1)
    })

    it('operates on the registered provider list instead of the raw store list', () => {
      const { moveSelection, selectedTaskIndex, selectedTaskId, registerTasks } = useTaskShortcuts()

      registerTasks(() => [storeMock.tasks[0]])

      moveSelection(5)
      expect(selectedTaskIndex.value).toBe(0)
      expect(selectedTaskId.value).toBe('t1')

      moveSelection(-1)
      expect(selectedTaskIndex.value).toBe(0)
    })

    it('resets the selection when a new provider is registered', () => {
      const { moveSelection, selectedTaskIndex, registerTasks } = useTaskShortcuts()

      moveSelection(2)
      expect(selectedTaskIndex.value).toBe(2)

      registerTasks(() => storeMock.tasks)
      expect(selectedTaskIndex.value).toBe(0)
    })
  })

  describe('toggleComplete', () => {
    it('completes the selected task', () => {
      const { moveSelection, toggleComplete } = useTaskShortcuts()

      moveSelection(1) // t2
      toggleComplete()

      expect(storeMock.editTask).toHaveBeenCalledWith({ task_id: 't2', status: 'completed' })
    })

    it('reopens a completed task', () => {
      const { toggleComplete } = useTaskShortcuts()

      toggleComplete(makeTask('t1', { status: 'completed' }))

      expect(storeMock.editTask).toHaveBeenCalledWith({ task_id: 't1', status: 'pending' })
    })

    it('shows an info toast when nothing is selected', () => {
      storeMock.tasks = []
      const { toggleComplete } = useTaskShortcuts()

      toggleComplete()

      expect(toastMock.add).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Task', color: 'info' })
      )
      expect(storeMock.editTask).not.toHaveBeenCalled()
    })
  })

  describe('deleteTask', () => {
    it('deletes silently and shows a toast with an undo action', () => {
      const { moveSelection, deleteTask } = useTaskShortcuts()

      moveSelection(1) // t2
      deleteTask()

      expect(storeMock.deleteTask).toHaveBeenCalledWith('t2', { silent: true })
      expect(toastMock.add).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Task deleted',
          actions: [expect.objectContaining({ label: 'Undo' })]
        })
      )
    })

    it('undoes the last deleted task', async () => {
      const { deleteTask, undoDelete } = useTaskShortcuts()

      deleteTask(makeTask('t2'))

      storeMock.undoDeleteTask.mockResolvedValueOnce(null)
      await undoDelete()

      expect(storeMock.undoDeleteTask).toHaveBeenCalledWith('t2')
    })
  })

  describe('undoDelete', () => {
    it('falls back to the most recently deleted task', async () => {
      storeMock.getRecentlyDeletedTasks.mockResolvedValueOnce([makeTask('t9')])
      const { undoDelete } = useTaskShortcuts()

      await undoDelete()

      expect(storeMock.getRecentlyDeletedTasks).toHaveBeenCalledTimes(1)
      expect(storeMock.undoDeleteTask).toHaveBeenCalledWith('t9')
    })

    it('shows an info toast when there is nothing to restore', async () => {
      storeMock.getRecentlyDeletedTasks.mockResolvedValueOnce([])
      const { undoDelete } = useTaskShortcuts()

      await undoDelete()

      expect(storeMock.undoDeleteTask).not.toHaveBeenCalled()
      expect(toastMock.add).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Task', color: 'info' })
      )
    })
  })

  describe('newTask', () => {
    it('dispatches the command event directly on the plan page', () => {
      const { newTask, consumePendingNewTask } = useTaskShortcuts()

      newTask()

      expect(dispatchEvent).toHaveBeenCalledTimes(1)
      expect(dispatchEvent.mock.calls[0][0].type).toBe('command:new-task')
      expect(navigateTo).not.toHaveBeenCalled()
      expect(consumePendingNewTask()).toBe(false)
    })

    it('navigates to /plan and hands over a pending flag elsewhere', () => {
      route.path = '/dashboard'
      const { newTask, consumePendingNewTask } = useTaskShortcuts()

      newTask()

      expect(dispatchEvent).not.toHaveBeenCalled()
      expect(navigateTo).toHaveBeenCalledWith('/plan')
      expect(consumePendingNewTask()).toBe(true)
      expect(consumePendingNewTask()).toBe(false)
    })
  })

  describe('openEdit', () => {
    it('dispatches the edit event on the plan page', () => {
      const { openEdit } = useTaskShortcuts()

      openEdit('t2')

      expect(dispatchEvent).toHaveBeenCalledTimes(1)
      expect(dispatchEvent.mock.calls[0][0].type).toBe('shortcut:edit-task')
      expect(dispatchEvent.mock.calls[0][0].detail).toEqual({ task_id: 't2' })
      expect(navigateTo).not.toHaveBeenCalled()
    })

    it('navigates to /plan with a pending task id elsewhere', () => {
      route.path = '/dashboard'
      const { openEdit, consumePendingEditTask } = useTaskShortcuts()

      openEdit('t3')

      expect(navigateTo).toHaveBeenCalledWith('/plan')
      expect(consumePendingEditTask()).toBe('t3')
      expect(consumePendingEditTask()).toBeNull()
    })
  })
})
