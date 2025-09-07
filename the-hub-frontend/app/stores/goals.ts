interface Task {
  task_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  order: number
}

interface Goal {
  goal_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  category?: string
  color: string
  progress: number
  total_tasks: number
  completed_tasks: number
  tasks: Task[]
}

export interface GoalsResponse {
  goals: Goal[]
}

export const useGoalStore = defineStore('goal', () => {
  const goals = ref<Goal[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()
  const { validateObject, schemas } = useValidation()

  async function fetchGoals() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedGoals = await $api<GoalsResponse>("/goals")

    if (fetchedGoals) {
      goals.value = fetchedGoals.goals

      // Initialize tasks array for each goal
      goals.value.forEach(goal => {
        if (!goal.tasks) {
          goal.tasks = []
        }
      })
    }

    loading.value = false
  }

  async function createGoal(payload: {
    title: string
    description: string
    due_date?: string
    priority?: number
    category?: string
    color?: string
  }) {
    const validation = validateObject(payload, schemas.goal.create)

    if (!validation.isValid) {
      const errorMessage = Object.values(validation.errors)[0]
      addToast(errorMessage, "error")
      return
    }

    // Create optimistic goal
    const optimisticGoal: Goal = {
      goal_id: `temp-${Date.now()}`,
      title: payload.title,
      description: payload.description,
      due_date: payload.due_date,
      priority: payload.priority,
      status: 'active',
      category: payload.category,
      color: payload.color || '#3B82F6',
      progress: 0,
      total_tasks: 0,
      completed_tasks: 0,
      tasks: []
    }

    // Optimistically add to local state
    goals.value.push(optimisticGoal)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Goal>('/goals', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic goal with real data
      const optimisticIndex = goals.value.findIndex(g => g.goal_id === optimisticGoal.goal_id)
      if (optimisticIndex !== -1) {
        goals.value[optimisticIndex] = { ...data, tasks: [] }
      }

      addToast("Goal added successfully", "success")
    } catch (err) {
      // Remove optimistic goal on error
      goals.value = goals.value.filter(g => g.goal_id !== optimisticGoal.goal_id)
      addToast(err?.message || "Goal not added", "error")
    }
  }

  async function updateGoal(payload: Goal) {
    // Store original goal for potential rollback
    const originalGoalIndex = goals.value.findIndex(g => g.goal_id === payload.goal_id)
    const originalGoal = originalGoalIndex !== -1 ? { ...goals.value[originalGoalIndex] } : null

    // Optimistically update the goal
    if (originalGoalIndex !== -1) {
      goals.value[originalGoalIndex] = { ...goals.value[originalGoalIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Goal>(`/goals/${payload.goal_id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          title: payload.title,
          description: payload.description,
          due_date: payload.due_date,
          priority: payload.priority,
          status: payload.status,
          category: payload.category,
          color: payload.color
        })
      })

      // Update with server response to ensure consistency
      if (originalGoalIndex !== -1 && data) {
        goals.value[originalGoalIndex] = data
      }

      addToast("Goal updated successfully", "success")
    } catch (err) {
      // Revert optimistic update on error
      if (originalGoal && originalGoalIndex !== -1) {
        goals.value[originalGoalIndex] = originalGoal
      }
      addToast(err?.message || "Goal update failed", "error")
    }
  }

  async function deleteGoal(id: string) {
    // Store the goal for potential rollback
    const goalToDelete = goals.value.find(g => g.goal_id === id)
    if (!goalToDelete) {
      addToast("Goal not found", "error")
      return
    }

    // Optimistically remove from local state
    goals.value = goals.value.filter((g) => g.goal_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`/goals/${id}`, {
        method: 'DELETE'
      })

      addToast("Goal deleted successfully", "success")
    } catch (err) {
      // Restore the goal on error
      goals.value.push(goalToDelete)
      addToast("Goal deletion failed", "error")
    }
  }

  async function fetchGoalTasks(goalId: string) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<{ tasks: Task[] }>(`/goals/${goalId}/tasks`)

      // Update the goal with its tasks
      const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
      if (goalIndex !== -1) {
        goals.value[goalIndex].tasks = data.tasks
      }

      return data.tasks
    } catch (err) {
      addToast("Failed to fetch goal tasks", "error")
      return []
    }
  }

  async function addTaskToGoal(goalId: string, taskData: {
    title: string
    description?: string
    priority?: number
    due_date?: string
  }) {
    // Find the goal
    const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
    if (goalIndex === -1) {
      addToast("Goal not found", "error")
      throw new Error("Goal not found")
    }

    // Create optimistic task
    const optimisticTask: Task = {
      task_id: `temp-goal-task-${Date.now()}`,
      title: taskData.title,
      description: taskData.description || '',
      due_date: taskData.due_date,
      priority: taskData.priority,
      status: 'pending',
      order: goals.value[goalIndex].tasks?.length || 0
    }

    // Optimistically add to goal's tasks
    if (!goals.value[goalIndex].tasks) {
      goals.value[goalIndex].tasks = []
    }
    goals.value[goalIndex].tasks.push(optimisticTask)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks`, {
        method: 'POST',
        body: JSON.stringify(taskData)
      })

      // Replace optimistic task with real data
      const optimisticIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === optimisticTask.task_id)
      if (optimisticIndex !== -1) {
        goals.value[goalIndex].tasks[optimisticIndex] = data
      }

      addToast("Task added to goal successfully", "success")
      return data
    } catch (err) {
      // Remove optimistic task on error
      goals.value[goalIndex].tasks = goals.value[goalIndex].tasks.filter(t => t.task_id !== optimisticTask.task_id)
      addToast("Failed to add task to goal", "error")
      throw err
    }
  }

  async function updateGoalTask(goalId: string, taskId: string, taskData: {
    title?: string
    description?: string
    priority?: number
    status?: string
    due_date?: string
  }) {
    // Find the goal and task
    const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
    if (goalIndex === -1 || !goals.value[goalIndex].tasks) {
      addToast("Goal or task not found", "error")
      throw new Error("Goal or task not found")
    }

    const taskIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === taskId)
    if (taskIndex === -1) {
      addToast("Task not found", "error")
      throw new Error("Task not found")
    }

    // Store original task for potential rollback
    const originalTask = { ...goals.value[goalIndex].tasks[taskIndex] }

    // Optimistically update the task
    goals.value[goalIndex].tasks[taskIndex] = { ...goals.value[goalIndex].tasks[taskIndex], ...taskData }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks/${taskId}`, {
        method: 'PATCH',
        body: JSON.stringify(taskData)
      })

      // Update with server response to ensure consistency
      goals.value[goalIndex].tasks[taskIndex] = data

      addToast("Task updated successfully", "success")
      return data
    } catch (err) {
      // Revert optimistic update on error
      goals.value[goalIndex].tasks[taskIndex] = originalTask
      addToast("Failed to update task", "error")
      throw err
    }
  }

  async function deleteGoalTask(goalId: string, taskId: string) {
    // Find the goal and task
    const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
    if (goalIndex === -1 || !goals.value[goalIndex].tasks) {
      addToast("Goal not found", "error")
      throw new Error("Goal not found")
    }

    const taskIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === taskId)
    if (taskIndex === -1) {
      addToast("Task not found", "error")
      throw new Error("Task not found")
    }

    // Store the task for potential rollback
    const taskToDelete = goals.value[goalIndex].tasks[taskIndex]

    // Optimistically remove from local state
    goals.value[goalIndex].tasks = goals.value[goalIndex].tasks.filter(t => t.task_id !== taskId)

    try {
      const { $api } = useNuxtApp()
      await $api(`/goals/${goalId}/tasks/${taskId}`, {
        method: 'DELETE'
      })

      addToast("Task deleted successfully", "success")
    } catch (err) {
      // Restore the task on error
      goals.value[goalIndex].tasks.push(taskToDelete)
      addToast("Failed to delete task", "error")
      throw err
    }
  }

  async function completeGoalTask(goalId: string, taskId: string) {
    // Find the goal and task
    const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
    if (goalIndex === -1 || !goals.value[goalIndex].tasks) {
      addToast("Goal not found", "error")
      throw new Error("Goal not found")
    }

    const taskIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === taskId)
    if (taskIndex === -1) {
      addToast("Task not found", "error")
      throw new Error("Task not found")
    }

    // Store original task for potential rollback
    const originalTask = { ...goals.value[goalIndex].tasks[taskIndex] }

    // Optimistically toggle the status
    const newStatus = goals.value[goalIndex].tasks[taskIndex].status === 'completed' ? 'pending' : 'completed'
    goals.value[goalIndex].tasks[taskIndex].status = newStatus

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks/${taskId}/complete`, {
        method: 'PATCH'
      })

      // Update with server response to ensure consistency
      goals.value[goalIndex].tasks[taskIndex] = data

      const statusMessage = data.status === 'completed' ? 'Task completed' : 'Task marked as pending'
      addToast(statusMessage, "success")
      return data
    } catch (err) {
      // Revert optimistic update on error
      goals.value[goalIndex].tasks[taskIndex] = originalTask
      addToast("Failed to update task status", "error")
      throw err
    }
  }

  function reset() {
    goals.value = []
  }

  return {
    goals,
    loading,
    fetchError,
    fetchGoals,
    createGoal,
    updateGoal,
    deleteGoal,
    fetchGoalTasks,
    addTaskToGoal,
    updateGoalTask,
    deleteGoalTask,
    completeGoalTask,
    reset,
  }
})
