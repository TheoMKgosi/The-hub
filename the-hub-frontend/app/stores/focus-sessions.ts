import type { Task } from '~/types/task'
interface FocusSession {
  focus_session_id: string
  user_id: string
  task_id?: string
  goal_id?: string
  duration_min: number
  started_at: string
  ended_at?: string
  status: string
  session_type: string
  cycle_group?: string
  cycle_order: number
  notes: string
}

interface FocusSessionStats {
  total_minutes: number
  total_hours: number
  total_sessions: number
  days: number
  daily_stats: Array<{
    date: string
    minutes: number
    sessions_count: number
  }>
  task_breakdown: Array<{
    task_id?: string
    task_title: string
    minutes: number
    sessions_count: number
  }>
  weekly_stats: Array<{
    week_start: string
    minutes: number
    sessions_count: number
  }>
  average_daily: number
}

interface StartSessionData {
  task_id?: string
  goal_id?: string
  session_type?: string
  planned_duration?: number
  notes?: string
  cycle_group?: string
  cycle_order?: number
}

export const useFocusSessionStore = defineStore('focus-sessions', () => {
  const activeSession = ref<FocusSession | null>(null)
  const sessions = ref<FocusSession[]>([])
  const stats = ref<FocusSessionStats | null>(null)
  const loading = ref(false)
  const toast = useToast()

  const startSession = async (data: StartSessionData) => {
    const optimisticSession: FocusSession = {
      focus_session_id: `temp-${Date.now()}`,
      user_id: '',
      task_id: data.task_id,
      goal_id: data.goal_id,
      duration_min: 0,
      started_at: new Date().toISOString(),
      status: 'active',
      session_type: data.session_type || 'focus',
      cycle_group: data.cycle_group,
      cycle_order: data.cycle_order || 0,
      notes: data.notes || ''
    }

    activeSession.value = optimisticSession

    try {
      const { $api } = useNuxtApp()
      const newSession = await $api<FocusSession>('/focus-sessions', {
        method: 'POST',
        body: JSON.stringify(data)
      })
      activeSession.value = newSession
      const existingIndex = sessions.value.findIndex(s => s.focus_session_id === optimisticSession.focus_session_id)
      if (existingIndex !== -1) {
        sessions.value[existingIndex] = newSession
      } else {
        sessions.value.unshift(newSession)
      }
      toast.add({title: "Task", description: "Focus session started", color: "success"})
      return newSession
    } catch (err) {
      activeSession.value = null
      toast.add({title: "Error", description: "Failed to start focus session", color: "error"})
      console.error('Error starting session:', err)
    }
  }

  const stopSession = async (sessionId: string, status: 'completed' | 'interrupted' = 'completed') => {
    try {
      const { $api } = useNuxtApp()
      const result = await $api<FocusSession>(`/focus-sessions/${sessionId}`, {
        method: 'PATCH',
        body: JSON.stringify({ status })
      })
      if (activeSession.value?.focus_session_id === sessionId) {
        activeSession.value = null
      }
      const idx = sessions.value.findIndex(s => s.focus_session_id === sessionId)
      if (idx !== -1) {
        sessions.value[idx] = result
      }
      const label = status === 'completed' ? 'completed' : 'interrupted'
      toast.add({title: "Task", description: `Focus session ${label}`, color: "success"})
      return result
    } catch (err) {
      toast.add({title: "Error", description: "Failed to stop session", color: "error"})
      console.error('Error stopping session:', err)
    }
  }

  const fetchSessions = async (filters?: {
    task_id?: string
    status?: string
    session_type?: string
    date_from?: string
    date_to?: string
    limit?: number
  }) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const queryParams = new URLSearchParams()
      if (filters?.task_id) queryParams.append('task_id', filters.task_id)
      if (filters?.status) queryParams.append('status', filters.status)
      if (filters?.session_type) queryParams.append('session_type', filters.session_type)
      if (filters?.date_from) queryParams.append('date_from', filters.date_from)
      if (filters?.date_to) queryParams.append('date_to', filters.date_to)
      if (filters?.limit) queryParams.append('limit', filters.limit.toString())

      const response = await $api<{ focus_sessions: FocusSession[] }>(`/focus-sessions?${queryParams}`)
      sessions.value = response.focus_sessions
    } catch (err) {
      toast.add({title: "Error", description: "Failed to fetch focus sessions", color: "error"})
      console.error('Error fetching sessions:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchStats = async (days: number = 30, sessionType: string = 'focus') => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<FocusSessionStats>(`/focus-sessions/stats?days=${days}&session_type=${sessionType}`)
      stats.value = response
    } catch (err) {
      toast.add({title: "Error", description: "Failed to fetch focus stats", color: "error"})
      console.error('Error fetching stats:', err)
    } finally {
      loading.value = false
    }
  }

  const getTodayMinutes = computed(() => {
    if (!stats.value) return 0
    const today = new Date().toISOString().split('T')[0]
    const todayStats = stats.value.daily_stats?.find(d => d.date.startsWith(today!))
    return todayStats?.minutes || 0
  })

  const getWeeklyMinutes = computed(() => {
    if (!stats.value) return 0
    return stats.value.weekly_stats.reduce((sum, w) => sum + w.minutes, 0)
  })

  return {
    activeSession,
    sessions,
    stats,
    loading,
    startSession,
    stopSession,
    fetchSessions,
    fetchStats,
    getTodayMinutes,
    getWeeklyMinutes
  }
})
