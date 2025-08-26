interface TaskStats {
  period: {
    start_date: string
    end_date: string
  }
  summary: {
    total_tasks: number
    completed_tasks: number
    pending_tasks: number
    overdue_tasks: number
    completion_rate: number
  }
  priority_distribution: {
    '1': number
    '2': number
    '3': number
    '4': number
    '5': number
  }
  priority_completion: {
    '1': number
    '2': number
    '3': number
    '4': number
    '5': number
  }
  goal_distribution: {
    with_goals: number
    without_goals: number
  }
  time_based: {
    due_today: number
    due_tomorrow: number
    due_this_week: number
  }
}

interface TaskStatsTrend {
  date: string
  period: {
    start_date: string
    end_date: string
  }
  summary: {
    total_tasks: number
    completed_tasks: number
    pending_tasks: number
    overdue_tasks: number
    completion_rate: number
  }
  priority_distribution: {
    '1': number
    '2': number
    '3': number
    '4': number
    '5': number
  }
  priority_completion: {
    '1': number
    '2': number
    '3': number
    '4': number
    '5': number
  }
  goal_distribution: {
    with_goals: number
    without_goals: number
  }
  time_based: {
    due_today: number
    due_tomorrow: number
    due_this_week: number
  }
}

export const useStatsStore = defineStore('stats', () => {
  const currentStats = ref<TaskStats | null>(null)
  const trends = ref<TaskStatsTrend[]>([])
  const loading = ref(false)
  const trendsLoading = ref(false)
  const error = ref<string | null>(null)
  const { addToast } = useToast()

  // Fetch current task statistics
  async function fetchStats(period: string = 'today', startDate?: string, endDate?: string) {
    const { $api } = useNuxtApp()
    loading.value = true
    error.value = null

    try {
      const params = new URLSearchParams({ period })
      if (startDate) params.append('start_date', startDate)
      if (endDate) params.append('end_date', endDate)

      const response = await $api<{ stats: TaskStats }>(`/stats/tasks?${params.toString()}`)
      if (response.stats) {
        currentStats.value = response.stats
      }
    } catch (err) {
      error.value = 'Failed to fetch task statistics'
      addToast('Failed to load statistics', 'error')
      console.error('Error fetching stats:', err)
    } finally {
      loading.value = false
    }
  }

  // Fetch task statistics trends over time
  async function fetchTrends(days: number = 30) {
    const { $api } = useNuxtApp()
    trendsLoading.value = true
    error.value = null

    try {
      const response = await $api<{ trends: TaskStatsTrend[] }>(`/stats/tasks/trends?days=${days}`)
      if (response.trends) {
        trends.value = response.trends
      }
    } catch (err) {
      error.value = 'Failed to fetch statistics trends'
      addToast('Failed to load trends', 'error')
      console.error('Error fetching trends:', err)
    } finally {
      trendsLoading.value = false
    }
  }

  // Computed properties for easy access to common metrics
  const completionRate = computed(() => {
    return currentStats.value?.summary.completion_rate || 0
  })

  const totalTasks = computed(() => {
    return currentStats.value?.summary.total_tasks || 0
  })

  const completedTasks = computed(() => {
    return currentStats.value?.summary.completed_tasks || 0
  })

  const pendingTasks = computed(() => {
    return currentStats.value?.summary.pending_tasks || 0
  })

  const overdueTasks = computed(() => {
    return currentStats.value?.summary.overdue_tasks || 0
  })

  const tasksDueToday = computed(() => {
    return currentStats.value?.time_based.due_today || 0
  })

  const tasksDueTomorrow = computed(() => {
    return currentStats.value?.time_based.due_tomorrow || 0
  })

  const tasksDueThisWeek = computed(() => {
    return currentStats.value?.time_based.due_this_week || 0
  })

  // Priority distribution as array for charts
  const priorityData = computed(() => {
    if (!currentStats.value) return []

    return [
      { priority: '1 - Low', count: currentStats.value.priority_distribution['1'], completed: currentStats.value.priority_completion['1'] },
      { priority: '2 - Low', count: currentStats.value.priority_distribution['2'], completed: currentStats.value.priority_completion['2'] },
      { priority: '3 - Medium', count: currentStats.value.priority_distribution['3'], completed: currentStats.value.priority_completion['3'] },
      { priority: '4 - High', count: currentStats.value.priority_distribution['4'], completed: currentStats.value.priority_completion['4'] },
      { priority: '5 - High', count: currentStats.value.priority_distribution['5'], completed: currentStats.value.priority_completion['5'] },
    ]
  })

  // Goal distribution for charts
  const goalData = computed(() => {
    if (!currentStats.value) return []

    return [
      { name: 'With Goals', value: currentStats.value.goal_distribution.with_goals },
      { name: 'Without Goals', value: currentStats.value.goal_distribution.without_goals },
    ]
  })

  // Trends data for time series charts
  const trendsCompletionData = computed(() => {
    return trends.value.map(trend => ({
      date: trend.date,
      completionRate: trend.summary.completion_rate,
      completedTasks: trend.summary.completed_tasks,
      totalTasks: trend.summary.total_tasks,
    }))
  })

  // Reset store
  function reset() {
    currentStats.value = null
    trends.value = []
    loading.value = false
    trendsLoading.value = false
    error.value = null
  }

  return {
    // State
    currentStats,
    trends,
    loading,
    trendsLoading,
    error,

    // Actions
    fetchStats,
    fetchTrends,
    reset,

    // Computed properties
    completionRate,
    totalTasks,
    completedTasks,
    pendingTasks,
    overdueTasks,
    tasksDueToday,
    tasksDueTomorrow,
    tasksDueThisWeek,
    priorityData,
    goalData,
    trendsCompletionData,
  }
})