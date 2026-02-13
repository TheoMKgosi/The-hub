import type { Task, RecurrenceRule, TimeEntry, TaskTemplate, TaskUpdate } from "~/types/task";


export interface TaskResponse {
  tasks: Task[]
}

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const pendingOperations = ref<Set<string>>(new Set())
  const { addToast } = useToast()
  const { validateObject, schemas } = useValidation()
  const { isOnline, storeDataLocally, getLocalData, syncPendingOperations } = useOffline()

  // Helper functions for managing pending operations
  const addPendingOperation = (operationId: string) => {
    pendingOperations.value.add(operationId)
  }

  const removePendingOperation = (operationId: string) => {
    pendingOperations.value.delete(operationId)
  }

  const isOperationPending = (operationId: string) => {
    return pendingOperations.value.has(operationId)
  }

  const completedTasks = computed(() => {
    return tasks.value.filter(task => task.status === 'complete')
  })

  async function fetchTasks(filters?: {
    status?: string
    priority?: number
    goal_id?: string
    search?: string
    due_before?: string
    due_after?: string
    order_by?: string
    sort?: string
  }) {
    const { $api } = useNuxtApp()
    loading.value = true

    try {
      // Try to load from local storage first for instant UI
      const localTasks = await getLocalData('tasks')
      if (localTasks.length > 0) {
        tasks.value = localTasks
      }

      // Build query parameters
      const params = new URLSearchParams()
      if (filters) {
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== '') {
            params.append(key, value.toString())
          }
        })
      }

      const queryString = params.toString()
      const url = queryString ? `/tasks?${queryString}` : '/tasks'

      if (isOnline.value) {
        const { tasks: fetchedTasks } = await $api<TaskResponse>(url)

        if (fetchedTasks) {
          tasks.value = fetchedTasks
          // Store locally for offline access
          await storeDataLocally('tasks', fetchedTasks)
        }
      } else {
        addToast("You're offline. Showing cached data.", "info")
      }
    } catch (error) {
      fetchError.value = error as Error
      // If online request fails, use local data
      if (!isOnline.value) {
        const localTasks = await getLocalData('tasks')
        if (localTasks.length > 0) {
          tasks.value = localTasks
          addToast("Using cached data due to offline status", "info")
        }
      }
    } finally {
      loading.value = false
    }
  }


  async function submitForm(payload: { title: string; description: string; due_date?: string; priority?: number; status?: string; natural_language_input?: string; use_natural_language?: boolean; parent_task_id?: string }) {
    // Validate payload first
    let validationSchema = schemas.task.create

    if (payload.use_natural_language && payload.natural_language_input) {
      validationSchema = schemas.task.naturalLanguage
    }

    const validation = validateObject(payload, validationSchema)

    if (!validation.isValid) {
      const errorMessage = Object.values(validation.errors)[0]
      addToast(errorMessage, "error")
      return
    }

    // Create optimistic task with temporary ID
    const optimisticTask: Task = {
      task_id: `temp-${Date.now()}`,
      title: payload.title,
      description: payload.description || '',
      due_date: payload.due_date,
      priority: payload.priority,
      status: payload.status || 'pending',
      order: tasks.value.length,
      time_estimate_minutes: 0,
      time_spent_minutes: 0,
      is_recurring: false,
      parent_task_id: payload.parent_task_id
    }

    const operationId = `create-task-${optimisticTask.task_id}`
    addPendingOperation(operationId)

    // Optimistically add to local state
    tasks.value.push(optimisticTask)

    try {
      if (isOnline.value) {
        const { $api } = useNuxtApp()
        const data = await $api<Task>('tasks', {
          method: 'POST',
          body: JSON.stringify(payload)
        })

        // Replace optimistic task with real data
        const optimisticIndex = tasks.value.findIndex(t => t.task_id === optimisticTask.task_id)
        if (optimisticIndex !== -1) {
          tasks.value[optimisticIndex] = data
        }

        addToast("Task added successfully", "success")
      } else {
        // Queue operation for when back online
        await addPendingOperation({
          url: '/tasks',
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
          store: 'tasks',
          operation: 'create'
        })

        addToast("Task saved locally. Will sync when online.", "info")
      }

    } catch (err) {
      // Remove optimistic task on error
      tasks.value = tasks.value.filter(t => t.task_id !== optimisticTask.task_id)
      addToast(err?.message || "Task not added", "error")
    } finally {
      removePendingOperation(operationId)
    }
  }

  async function editTask(payload: TaskUpdate) {
    // Store original task for potential rollback
    const originalTaskIndex = tasks.value.findIndex(t => t.task_id === payload.task_id)
    const originalTask = originalTaskIndex !== -1 ? { ...tasks.value[originalTaskIndex] } : null

    // Optimistically update the task
    if (originalTaskIndex !== -1) {
      tasks.value[originalTaskIndex] = { ...tasks.value[originalTaskIndex], ...payload }
    }

    try {
      if (isOnline.value) {
        const { $api } = useNuxtApp()
        const data = await $api(`tasks/${payload.task_id}`, {
          method: 'PATCH',
          body: JSON.stringify(payload)
        })

        // Update with server response to ensure consistency
        if (originalTaskIndex !== -1 && data) {
          tasks.value[originalTaskIndex] = data
        }

        addToast("Edited task succesfully", "success")
      } else {
        // Queue operation for when back online
        await addPendingOperation({
          url: `/tasks/${payload.task_id}`,
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
          store: 'tasks',
          operation: 'update'
        })

        addToast("Task updated locally. Will sync when online.", "info")
      }

    } catch (err) {
      // Revert optimistic update on error
      if (originalTask && originalTaskIndex !== -1) {
        tasks.value[originalTaskIndex] = originalTask
      }
      addToast(err?.message || "Editing task failed", "error")
    }
  }

  async function reorderTask(payload: { task_id: string, order: number }[]) {
    const { $api } = useNuxtApp()
    await $api("/tasks/reorder", {
      method: 'PUT',
      body: JSON.stringify({ task_orders: payload })
    })
  }

  async function deleteTask(id: string) {
    // Store the task for potential rollback
    const taskToDelete = tasks.value.find(t => t.task_id === id)
    if (!taskToDelete) {
      addToast("Task not found", "error")
      return
    }

    // Optimistically remove from local state
    tasks.value = tasks.value.filter((t) => t.task_id !== id)

    try {
      if (isOnline.value) {
        const { $api } = useNuxtApp()
        await $api(`tasks/${id}`, {
          method: 'DELETE'
        })
        addToast("Task deleted successfully", "success")
      } else {
        // Queue operation for when back online
        await addPendingOperation({
          url: `/tasks/${id}`,
          method: 'DELETE',
          store: 'tasks',
          operation: 'delete'
        })

        addToast("Task deleted locally. Will sync when online.", "info")
      }

    } catch (err) {
      // Restore the task on error
      tasks.value.push(taskToDelete)
      addToast(err?.message || "Task did not delete", "error")
    }
  }

  async function completeTask(payload: Task) {
    const { $api } = useNuxtApp()
    await $api(`tasks/${payload.task_id}`, {
      method: 'PATCH',
      body: JSON.stringify({ status: payload.status })
    })
  }

  const unscheduledTasks = computed(() => {
    return tasks.value.filter((task: Task) => task.start_time === null)
  })

  async function undoDeleteTask(id: string) {
    try {
      if (isOnline.value) {
        const { $api } = useNuxtApp()
        const restoredTask = await $api<Task>(`tasks/${id}/undo-delete`, {
          method: 'PATCH'
        })

        // Add the restored task back to the local state
        tasks.value.push(restoredTask)
        addToast("Task restored successfully", "success")

        return restoredTask
      } else {
        addToast("Cannot undo task deletion while offline", "error")
        return null
      }
    } catch (err) {
      addToast(err?.message || "Failed to restore task", "error")
      return null
    }
  }

  async function getRecentlyDeletedTasks() {
    try {
      if (isOnline.value) {
        const { $api } = useNuxtApp()
        const data = await $api<{ tasks: Task[] }>('tasks/recently-deleted')
        return data.tasks
      } else {
        addToast("Cannot fetch deleted tasks while offline", "error")
        return []
      }
    } catch (err) {
      addToast("Failed to fetch recently deleted tasks", "error")
      return []
    }
  }

  async function createSubtask(parentTaskId: string, payload: { title: string; description?: string; priority?: number; due_date?: string }) {
    // Find parent task
    const parentTask = tasks.value.find(t => t.task_id === parentTaskId)
    if (!parentTask) {
      addToast("Parent task not found", "error")
      throw new Error("Parent task not found")
    }

    // Create optimistic subtask
    const optimisticSubtask: Task = {
      task_id: `temp-sub-${Date.now()}`,
      title: payload.title,
      description: payload.description || '',
      due_date: payload.due_date,
      priority: payload.priority,
      status: 'pending',
      order: 0,
      time_estimate_minutes: 0,
      time_spent_minutes: 0,
      is_recurring: false,
      parent_task_id: parentTaskId
    }

    // Optimistically add to parent task
    if (!parentTask.subtasks) parentTask.subtasks = []
    parentTask.subtasks.push(optimisticSubtask)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>('tasks', {
        method: 'POST',
        body: JSON.stringify({
          ...payload,
          parent_task_id: parentTaskId
        })
      })

      // Replace optimistic subtask with real data
      const optimisticIndex = parentTask.subtasks.findIndex(t => t.task_id === optimisticSubtask.task_id)
      if (optimisticIndex !== -1) {
        parentTask.subtasks[optimisticIndex] = data
      }

      addToast("Subtask added successfully", "success")
      return data
    } catch (err) {
      // Remove optimistic subtask on error
      parentTask.subtasks = parentTask.subtasks.filter(t => t.task_id !== optimisticSubtask.task_id)
      addToast(err?.message || "Subtask not added", "error")
      throw err
    }
  }

  async function getTaskSubtasks(taskId: string) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ subtasks: Task[] }>(`tasks/${taskId}/subtasks`)

      // Update task with subtasks
      const task = tasks.value.find(t => t.task_id === taskId)
      if (task) {
        task.subtasks = data.subtasks
      }

      return data.subtasks
    } catch (err) {
      addToast("Failed to fetch subtasks", "error")
      return []
    }
  }

  async function createTaskDependency(taskId: string, dependsOnId: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`tasks/${taskId}/dependencies`, {
        method: 'POST',
        body: JSON.stringify({ depends_on_id: dependsOnId })
      })

      addToast("Dependency created successfully", "success")
    } catch (err) {
      addToast(err?.message || "Failed to create dependency", "error")
      throw err
    }
  }

  async function getTaskDependencies(taskId: string) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ dependencies: Task[]; dependents: Task[] }>(`tasks/${taskId}/dependencies`)
      return data
    } catch (err) {
      addToast("Failed to fetch dependencies", "error")
      return { dependencies: [], dependents: [] }
    }
  }

  async function deleteTaskDependency(taskId: string, dependsOnId: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`tasks/${taskId}/dependencies/${dependsOnId}`, {
        method: 'DELETE'
      })

      addToast("Dependency deleted successfully", "success")
    } catch (err) {
      addToast("Failed to delete dependency", "error")
      throw err
    }
  }

  async function startTimeTracking(taskId: string, description: string = '') {
    try {
      const { $api } = useNuxtApp()
      const timeEntry = await $api<TimeEntry>(`tasks/${taskId}/time/start`, {
        method: 'POST',
        body: JSON.stringify({ description })
      })

      addToast("Time tracking started", "success")
      return timeEntry
    } catch (err) {
      addToast("Failed to start time tracking", "error")
      throw err
    }
  }

  async function stopTimeTracking(taskId: string) {
    try {
      const { $api } = useNuxtApp()
      const timeEntry = await $api<TimeEntry>(`tasks/${taskId}/time/stop`, {
        method: 'POST'
      })

      addToast("Time tracking stopped", "success")
      return timeEntry
    } catch (err) {
      addToast("Failed to stop time tracking", "error")
      throw err
    }
  }

  async function getTaskTimeEntries(taskId: string) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ time_entries: TimeEntry[] }>(`tasks/${taskId}/time`)
      return data.time_entries
    } catch (err) {
      addToast("Failed to fetch time entries", "error")
      return []
    }
  }

  async function createTaskTemplate(templateData: {
    name: string
    description?: string
    category?: string
    title_template: string
    description_template?: string
    priority?: number
    time_estimate_minutes?: number
    tags?: string
    is_public?: boolean
  }) {
    try {
      const { $api } = useNuxtApp()
      const template = await $api<TaskTemplate>('task-templates', {
        method: 'POST',
        body: JSON.stringify(templateData)
      })

      addToast("Task template created successfully", "success")
      return template
    } catch (err) {
      addToast("Failed to create task template", "error")
      throw err
    }
  }

  async function getTaskTemplates() {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ templates: TaskTemplate[] }>('task-templates')
      return data.templates
    } catch (err) {
      addToast("Failed to fetch task templates", "error")
      return []
    }
  }

  async function createTaskFromTemplate(templateId: string, taskData: Partial<{
    title: string
    description: string
    priority: number
    due_date: string
    time_estimate_minutes: number
    goal_id: string
    parent_task_id: string
  }> = {}) {
    // Create optimistic task with template data
    const optimisticTask: Task = {
      task_id: `temp-template-${Date.now()}`,
      title: taskData.title || 'New Task',
      description: taskData.description || '',
      due_date: taskData.due_date,
      priority: taskData.priority,
      status: 'pending',
      order: tasks.value.length,
      time_estimate_minutes: taskData.time_estimate_minutes || 0,
      time_spent_minutes: 0,
      is_recurring: false,
      goal_id: taskData.goal_id,
      parent_task_id: taskData.parent_task_id
    }

    // Optimistically add to local state
    tasks.value.push(optimisticTask)

    try {
      const { $api } = useNuxtApp()
      const task = await $api<Task>(`task-templates/${templateId}/create-task`, {
        method: 'POST',
        body: JSON.stringify(taskData)
      })

      // Replace optimistic task with real data
      const optimisticIndex = tasks.value.findIndex(t => t.task_id === optimisticTask.task_id)
      if (optimisticIndex !== -1) {
        tasks.value[optimisticIndex] = task
      }

      addToast("Task created from template successfully", "success")
      return task
    } catch (err) {
      // Remove optimistic task on error
      tasks.value = tasks.value.filter(t => t.task_id !== optimisticTask.task_id)
      addToast("Failed to create task from template", "error")
      throw err
    }
  }

  async function createRecurrenceRule(ruleData: {
    name: string
    description?: string
    frequency: string
    interval?: number
    by_day?: string
    by_month_day?: number
    by_month?: number
    start_date?: string
    end_date?: string
    count?: number
    title_template: string
    description_template?: string
    priority?: number
    time_estimate_minutes?: number
    due_date_offset_days?: number
  }) {
    try {
      const { $api } = useNuxtApp()
      const rule = await $api<RecurrenceRule>('recurrence-rules', {
        method: 'POST',
        body: JSON.stringify(ruleData)
      })

      addToast("Recurrence rule created successfully", "success")
      return rule
    } catch (err) {
      addToast("Failed to create recurrence rule", "error")
      throw err
    }
  }

  async function getRecurrenceRules() {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ recurrence_rules: RecurrenceRule[] }>('recurrence-rules')
      return data.recurrence_rules
    } catch (err) {
      addToast("Failed to fetch recurrence rules", "error")
      return []
    }
  }

  async function generateRecurringTasks(ruleId: string, count: number = 1) {
    // Create optimistic tasks
    const optimisticTasks: Task[] = []
    for (let i = 0; i < count; i++) {
      optimisticTasks.push({
        task_id: `temp-recurring-${Date.now()}-${i}`,
        title: 'Recurring Task',
        description: '',
        status: 'pending',
        order: tasks.value.length + i,
        time_estimate_minutes: 0,
        time_spent_minutes: 0,
        is_recurring: true
      })
    }

    // Optimistically add to local state
    tasks.value.push(...optimisticTasks)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ tasks: Task[] }>(`recurrence-rules/${ruleId}/generate-tasks?count=${count}`, {
        method: 'POST'
      })

      // Replace optimistic tasks with real data
      optimisticTasks.forEach((_, index) => {
        const optimisticIndex = tasks.value.findIndex(t => t.task_id === optimisticTasks[index].task_id)
        if (optimisticIndex !== -1 && data.tasks[index]) {
          tasks.value[optimisticIndex] = data.tasks[index]
        }
      })

      addToast(`Generated ${data.tasks.length} recurring tasks`, "success")
      return data.tasks
    } catch (err) {
      // Remove optimistic tasks on error
      tasks.value = tasks.value.filter(t => !t.task_id.startsWith('temp-recurring-'))
      addToast("Failed to generate recurring tasks", "error")
      throw err
    }
  }

  function reset() {
    tasks.value = []
  }


  // Sync function to be called when back online
  async function syncOfflineChanges() {
    if (isOnline.value) {
      await syncPendingOperations()
      // Refresh data after sync
      await fetchTasks()
    }
  }

  return {
    tasks,
    unscheduledTasks,
    completedTasks,
    loading,
    fetchError,
    pendingOperations,
    isOperationPending,
    fetchTasks,
    editTask,
    reorderTask,
    completeTask,
    deleteTask,
    undoDeleteTask,
    getRecentlyDeletedTasks,
    submitForm,
    createSubtask,
    getTaskSubtasks,
    createTaskDependency,
    getTaskDependencies,
    deleteTaskDependency,
    startTimeTracking,
    stopTimeTracking,
    getTaskTimeEntries,
    createTaskTemplate,
    getTaskTemplates,
    createTaskFromTemplate,
    createRecurrenceRule,
    getRecurrenceRules,
    generateRecurringTasks,
    syncOfflineChanges,
    reset,
  }
})
