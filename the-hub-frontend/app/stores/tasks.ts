interface Task {
  task_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  order?: number
  goal_id?: string
  parent_task_id?: string
  subtasks?: Task[]
  time_estimate_minutes?: number
  time_spent_minutes: number
  is_recurring: boolean
  template_id?: string
}

interface TimeEntry {
  time_entry_id: string
  task_id: string
  description: string
  start_time: string
  end_time?: string
  duration_minutes: number
  is_running: boolean
}

interface TaskTemplate {
  template_id: string
  name: string
  description: string
  category: string
  title_template: string
  description_template: string
  priority?: number
  time_estimate_minutes?: number
  tags: string
  is_public: boolean
  usage_count: number
}

interface RecurrenceRule {
  recurrence_rule_id: string
  name: string
  description: string
  frequency: string
  interval: number
  by_day?: string
  by_month_day?: number
  by_month?: number
  start_date?: string
  end_date?: string
  count?: number
  title_template: string
  description_template: string
  priority?: number
  time_estimate_minutes?: number
  due_date_offset_days?: number
}

export interface TaskResponse {
  tasks: Task[]
}

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()
  const { validateObject, schemas } = useValidation()

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

    const { tasks: fetchedTasks } = await $api<TaskResponse>(url)

    if (fetchedTasks) tasks.value = fetchedTasks
    loading.value = false
  }


  async function submitForm(payload: { title: string; description: string; due_date?: string; priority?: number; status?: string; natural_language_input?: string; use_natural_language?: boolean; parent_task_id?: string }) {
    try {
      // Validate payload
      let validationSchema = schemas.task.create

      if (payload.use_natural_language && payload.natural_language_input) {
        validationSchema = schemas.task.naturalLanguage
      }

      const validation = validateObject(payload, validationSchema)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      const data = await $api<Task>('tasks', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      tasks.value.push(data)
      addToast("Task added successfully", "success")

    } catch (err) {
      addToast(err?.message || "Task not added", "error")
    }
  }

  async function editTask(payload: Task) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api(`tasks/${payload.task_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      fetchTasks()
      addToast("Edited task succesfully", "success")

    } catch (err) {
      addToast(err?.message || "Editing task failed", "error")
    }
  }

  async function reorderTask(payload: {task_id: string, order: number}[]) {
    const { $api } = useNuxtApp()
    await $api("/tasks/reorder", {
      method: 'PUT',
      body: JSON.stringify({task_orders: payload} )
    })
  }

  async function completeTask(payload: Task) {
    const { $api } = useNuxtApp()
    await $api(`tasks/${payload.task_id}`, {
      method: 'PATCH',
      body: JSON.stringify({ status: payload.status })
    })
  }

  async function deleteTask(id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`tasks/${id}`, {
        method: 'DELETE'
      })
      tasks.value = tasks.value.filter((t) => t.task_id !== id)
      addToast("Task deleted succesfully", "success")

    } catch (err) {
      addToast(err?.message || "Task did not delete", "error")
    }
  }

  async function createSubtask(parentTaskId: string, payload: { title: string; description?: string; priority?: number; due_date?: string }) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>('tasks', {
        method: 'POST',
        body: JSON.stringify({
          ...payload,
          parent_task_id: parentTaskId
        })
      })

      // Add subtask to parent task
      const parentTask = tasks.value.find(t => t.task_id === parentTaskId)
      if (parentTask) {
        if (!parentTask.subtasks) parentTask.subtasks = []
        parentTask.subtasks.push(data)
      }

      addToast("Subtask added successfully", "success")
      return data
    } catch (err) {
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
    try {
      const { $api } = useNuxtApp()
      const task = await $api<Task>(`task-templates/${templateId}/create-task`, {
        method: 'POST',
        body: JSON.stringify(taskData)
      })

      tasks.value.push(task)
      addToast("Task created from template successfully", "success")
      return task
    } catch (err) {
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
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ tasks: Task[] }>(`recurrence-rules/${ruleId}/generate-tasks?count=${count}`, {
        method: 'POST'
      })

      // Add new tasks to the store
      tasks.value.push(...data.tasks)
      addToast(`Generated ${data.tasks.length} recurring tasks`, "success")
      return data.tasks
    } catch (err) {
      addToast("Failed to generate recurring tasks", "error")
      throw err
    }
  }

  function reset() {
    tasks.value = []
  }


  return {
    tasks,
    completedTasks,
    loading,
    fetchError,
    fetchTasks,
    editTask,
    reorderTask,
    completeTask,
    deleteTask,
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
    reset,
  }
})
