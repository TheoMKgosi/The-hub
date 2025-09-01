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

  async function createGoal(payload: { title: string; description: string }) {
    try {
      const validation = validateObject(payload, schemas.goal.create)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      const data = await $api<Goal>('/goals', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Initialize tasks array for the new goal
      data.tasks = []
      goals.value.push(data)
      addToast("Goal added successfully", "success")
    } catch (err) {
      addToast(err?.message || "Goal not added", "error")
    }
  }

  async function updateGoal(payload: Goal) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Goal>(`/goals/${payload.goal_id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          title: payload.title,
          description: payload.description
        })
      })

      const index = goals.value.findIndex(g => g.goal_id === payload.goal_id)
      if (index !== -1) {
        goals.value[index] = data
      }
      addToast("Goal updated successfully", "success")
    } catch (err) {
      addToast(err?.message || "Goal update failed", "error")
    }
  }

  async function deleteGoal(id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`/goals/${id}`, {
        method: 'DELETE'
      })

      goals.value = goals.value.filter((g) => g.goal_id !== id)
      addToast("Goal deleted successfully", "success")
    } catch (err) {
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
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks`, {
        method: 'POST',
        body: JSON.stringify(taskData)
      })

      // Update the goal's tasks
      const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
      if (goalIndex !== -1) {
        if (!goals.value[goalIndex].tasks) {
          goals.value[goalIndex].tasks = []
        }
        goals.value[goalIndex].tasks.push(data)
      }

      addToast("Task added to goal successfully", "success")
      return data
    } catch (err) {
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
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks/${taskId}`, {
        method: 'PATCH',
        body: JSON.stringify(taskData)
      })

      // Update the task in the goal's tasks
      const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
      if (goalIndex !== -1 && goals.value[goalIndex].tasks) {
        const taskIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === taskId)
        if (taskIndex !== -1) {
          goals.value[goalIndex].tasks[taskIndex] = data
        }
      }

      addToast("Task updated successfully", "success")
      return data
    } catch (err) {
      addToast("Failed to update task", "error")
      throw err
    }
  }

  async function deleteGoalTask(goalId: string, taskId: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`/goals/${goalId}/tasks/${taskId}`, {
        method: 'DELETE'
      })

      // Remove the task from the goal's tasks
      const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
      if (goalIndex !== -1 && goals.value[goalIndex].tasks) {
        goals.value[goalIndex].tasks = goals.value[goalIndex].tasks.filter(t => t.task_id !== taskId)
      }

      addToast("Task deleted successfully", "success")
    } catch (err) {
      addToast("Failed to delete task", "error")
      throw err
    }
  }

  async function completeGoalTask(goalId: string, taskId: string) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/goals/${goalId}/tasks/${taskId}/complete`, {
        method: 'PATCH'
      })

      // Update the task in the goal's tasks
      const goalIndex = goals.value.findIndex(g => g.goal_id === goalId)
      if (goalIndex !== -1 && goals.value[goalIndex].tasks) {
        const taskIndex = goals.value[goalIndex].tasks.findIndex(t => t.task_id === taskId)
        if (taskIndex !== -1) {
          goals.value[goalIndex].tasks[taskIndex] = data
        }
      }

      const statusMessage = data.status === 'completed' ? 'Task completed' : 'Task marked as pending'
      addToast(statusMessage, "success")
      return data
    } catch (err) {
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
