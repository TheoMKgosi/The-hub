import { ref, computed, watch, readonly, onUnmounted, getCurrentInstance } from 'vue'
import type { Task } from '~/types/task'

// Shared module-level state so every caller (layout, dashboard widgets,
// plan page) works with the same selection and pending navigation flags.
const selectedTaskIndex = ref(0)
const pendingNewTask = ref(false)
const pendingEditTaskId = ref<string | null>(null)
const lastDeletedTaskId = ref<string | null>(null)

// Components (plan page, dashboard widgets) can register the task list
// their UI is displaying, so navigation follows what is actually visible.
// providerVersion is bumped on registration so computeds re-evaluate
// (plain module variables are not reactive).
let tasksProvider: (() => Task[]) | null = null
const providerVersion = ref(0)

/**
 * Task management keyboard shortcut logic, shared between the global
 * shortcut registry, the dashboard widgets and the plan page.
 *
 * Must be called from a component setup context (it registers watches and
 * store access), and relies on Nuxt auto-imports for useTaskStore/useToast.
 */
export function useTaskShortcuts() {
  const taskStore = useTaskStore()
  const toast = useToast()
  const route = useRoute()

  const currentTasks = computed<Task[]>(() => {
    providerVersion.value
    return tasksProvider ? tasksProvider() : taskStore.tasks
  })

  const selectedTaskId = computed(() => currentTasks.value[selectedTaskIndex.value]?.task_id ?? null)

  // Keep the selection inside bounds when the visible list shrinks
  watch(currentTasks, (tasks) => {
    if (selectedTaskIndex.value >= tasks.length) {
      selectedTaskIndex.value = Math.max(0, tasks.length - 1)
    }
  })

  /**
   * Scope task navigation to a specific (e.g. filtered) list. The most
   * recent registration wins; the provider is unregistered automatically
   * when the registering component unmounts.
   */
  const registerTasks = (provider: () => Task[]) => {
    tasksProvider = provider
    providerVersion.value++
    selectedTaskIndex.value = 0
    const unregister = () => {
      if (tasksProvider === provider) {
        tasksProvider = null
        providerVersion.value++
      }
    }
    if (getCurrentInstance()) onUnmounted(unregister)
    return unregister
  }

  const getSelectedTask = () => currentTasks.value[selectedTaskIndex.value] ?? null

  const moveSelection = (delta: number) => {
    const tasks = currentTasks.value
    if (tasks.length === 0) {
      // Lazy-load tasks the first time a task shortcut is used
      if (taskStore.tasks.length === 0 && !taskStore.loading) taskStore.fetchTasks()
      return
    }
    const next = selectedTaskIndex.value + delta
    selectedTaskIndex.value = Math.min(Math.max(next, 0), tasks.length - 1)
  }

  const toggleComplete = (task?: Task) => {
    const target = task ?? getSelectedTask()
    if (!target) {
      toast.add({ title: 'Task', description: 'No task selected', color: 'info' })
      return
    }
    const newStatus = target.status === 'completed' ? 'pending' : 'completed'
    taskStore.editTask({ task_id: target.task_id, status: newStatus })
  }

  /** Open the edit modal for a task. The plan page owns the modal, so when
   * already on /plan we dispatch an event; otherwise navigate there first
   * and hand over the pending task id. */
  const openEdit = (taskId: string) => {
    if (route.path === '/plan') {
      window.dispatchEvent(new CustomEvent('shortcut:edit-task', { detail: { task_id: taskId } }))
    } else {
      pendingEditTaskId.value = taskId
      navigateTo('/plan')
    }
  }

  const editSelected = () => {
    const target = getSelectedTask()
    if (!target) {
      toast.add({ title: 'Task', description: 'No task selected', color: 'info' })
      return
    }
    openEdit(target.task_id)
  }

  const consumePendingEditTask = () => {
    const id = pendingEditTaskId.value
    pendingEditTaskId.value = null
    return id
  }

  const deleteTask = (task?: Task) => {
    const target = task ?? getSelectedTask()
    if (!target) {
      toast.add({ title: 'Task', description: 'No task selected', color: 'info' })
      return
    }
    lastDeletedTaskId.value = target.task_id
    // Silent: the toast below carries the undo action
    taskStore.deleteTask(target.task_id, { silent: true })
    toast.add({
      title: 'Task deleted',
      description: `"${target.title}" was deleted`,
      actions: [{ label: 'Undo', onClick: () => taskStore.undoDeleteTask(target.task_id) }]
    })
  }

  const undoDelete = async () => {
    let id = lastDeletedTaskId.value
    if (!id) {
      const recent = await taskStore.getRecentlyDeletedTasks()
      id = recent[0]?.task_id ?? null
    }
    if (!id) {
      toast.add({ title: 'Task', description: 'No recently deleted task to restore', color: 'info' })
      return
    }
    lastDeletedTaskId.value = null
    await taskStore.undoDeleteTask(id)
  }

  const newTask = () => {
    if (route.path === '/plan') {
      window.dispatchEvent(new CustomEvent('command:new-task'))
    } else {
      pendingNewTask.value = true
      navigateTo('/plan')
    }
  }

  const consumePendingNewTask = () => {
    if (!pendingNewTask.value) return false
    pendingNewTask.value = false
    return true
  }

  return {
    selectedTaskIndex: readonly(selectedTaskIndex),
    selectedTaskId,
    currentTasks,
    registerTasks,
    moveSelection,
    toggleComplete,
    editSelected,
    openEdit,
    consumePendingEditTask,
    deleteTask,
    undoDelete,
    newTask,
    consumePendingNewTask,
  }
}
