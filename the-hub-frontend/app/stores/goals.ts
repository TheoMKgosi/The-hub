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
      addToast("Goal not added", "error")
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
      addToast("Goal update failed", "error")
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
    reset,
  }
})
