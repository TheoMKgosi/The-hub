interface Goal {
  goal_id: number
  title: string
  description: string
  tasks: null
}

export interface GoalsResponse {
  goals: Goal[]
}

export const useGoalStore = defineStore('goal', () => {
  const goals = ref<Goal[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchGoals() {
    loading.value = true
    const { data, error } = await useFetch('http://localhost:8080/goals').json<GoalsResponse>()

    if (data.value) goals.value = data.value.goals
    fetchError.value = error.value

    loading.value = false
  }

  function reset() {
    goals.value = []
  }

  return {
    goals,
    loading,
    fetchError,
    fetchGoals,
    reset,
  }
})
