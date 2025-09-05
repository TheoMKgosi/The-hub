import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

interface StudySession {
  id: string
  user_id: string
  topic_id?: string
  task_id?: string
  duration_min: number
  started_at: string
  ended_at: string
}

interface StudySessionStats {
  total_minutes: number
  total_hours: number
  days: number
  daily_stats: Array<{
    date: string
    minutes: number
  }>
  topic_stats: Array<{
    topic_id?: string
    topic_title: string
    minutes: number
    sessions: number
  }>
  average_daily: number
}

export const useStudySessionStore = defineStore('study-sessions', () => {
  const sessions = ref<StudySession[]>([])
  const stats = ref<StudySessionStats | null>(null)
  const loading = ref(false)
  const { addToast } = useToast()

  const createSession = async (data: {
    topic_id?: string
    task_id?: string
    duration_min: number
    notes?: string
  }) => {
    try {
      const { $api } = useNuxtApp()
      const newSession = await $api<StudySession>('/study-sessions', {
        method: 'POST',
        body: JSON.stringify(data)
      })

      if (newSession) {
        sessions.value.unshift(newSession)
        addToast('Study session logged successfully!', 'success')
        return newSession
      }
    } catch (err) {
      addToast('Failed to log study session', 'error')
      console.error('Error creating study session:', err)
    }
  }

  const fetchSessions = async (filters?: {
    topic_id?: string
    task_id?: string
    date_from?: string
    date_to?: string
    limit?: number
  }) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const queryParams = new URLSearchParams()

      if (filters?.topic_id) queryParams.append('topic_id', filters.topic_id)
      if (filters?.task_id) queryParams.append('task_id', filters.task_id)
      if (filters?.date_from) queryParams.append('date_from', filters.date_from)
      if (filters?.date_to) queryParams.append('date_to', filters.date_to)
      if (filters?.limit) queryParams.append('limit', filters.limit.toString())

      const response = await $api<{ study_sessions: StudySession[] }>(`/study-sessions?${queryParams}`)
      sessions.value = response.study_sessions
    } catch (err) {
      addToast('Failed to fetch study sessions', 'error')
      console.error('Error fetching study sessions:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchStats = async (days: number = 30) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<StudySessionStats>(`/study-sessions/stats?days=${days}`)
      stats.value = response
    } catch (err) {
      addToast('Failed to fetch study statistics', 'error')
      console.error('Error fetching study stats:', err)
    } finally {
      loading.value = false
    }
  }

  const getTotalStudyTime = (days: number = 30) => {
    if (!stats.value) return 0
    return stats.value.total_minutes
  }

  const getAverageDailyTime = (days: number = 30) => {
    if (!stats.value) return 0
    return stats.value.average_daily
  }

  const getTopicBreakdown = () => {
    if (!stats.value) return []
    return stats.value.topic_stats
  }

  const getDailyStats = () => {
    if (!stats.value) return []
    return stats.value.daily_stats
  }

  return {
    sessions,
    stats,
    loading,
    createSession,
    fetchSessions,
    fetchStats,
    getTotalStudyTime,
    getAverageDailyTime,
    getTopicBreakdown,
    getDailyStats
  }
})