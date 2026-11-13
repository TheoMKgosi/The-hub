import { defineStore } from 'pinia'

export interface ActivityEntry {
  activity_id: string
  user_id: string
  title: string
  notes: string
  category: string
  tags: string[]
  activity_date: string
  duration_min: number
  created_at: string
  updated_at: string
}

interface ActivityStats {
  days: number
  total_entries: number
  total_minutes: number
  total_hours: number
  categories: Array<{
    category: string
    minutes: number
    entries: number
  }>
  daily_stats: Array<{
    date: string
    minutes: number
    entries: number
  }>
  weekly_stats: Array<{
    week_start: string
    minutes: number
    entries: number
  }>
  top_tags: Array<{
    tag: string
    minutes: number
    count: number
  }>
  current_streak: number
}

export interface ActivityEntryData {
  title: string
  notes?: string
  category?: string
  tags?: string[]
  activity_date?: string
  duration_min?: number
}

interface ActivityFilters {
  date_from?: string
  date_to?: string
  category?: string
  tag?: string
  search?: string
  limit?: number
  offset?: number
}

export const useActivityLogStore = defineStore('activity-log', () => {
  const entries = ref<ActivityEntry[]>([])
  const total = ref(0)
  const stats = ref<ActivityStats | null>(null)
  const loading = ref(false)
  const toast = useToast()

  const fetchEntries = async (filters?: ActivityFilters) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const queryParams = new URLSearchParams()
      if (filters?.date_from) queryParams.append('date_from', filters.date_from)
      if (filters?.date_to) queryParams.append('date_to', filters.date_to)
      if (filters?.category) queryParams.append('category', filters.category)
      if (filters?.tag) queryParams.append('tag', filters.tag)
      if (filters?.search) queryParams.append('search', filters.search)
      if (filters?.limit) queryParams.append('limit', filters.limit.toString())
      if (filters?.offset) queryParams.append('offset', filters.offset.toString())

      const response = await $api<{ activities: ActivityEntry[]; total: number }>(`/activities?${queryParams}`)
      entries.value = response.activities
      total.value = response.total
    } catch (err) {
      toast.add({ title: 'Error', description: 'Failed to fetch activities', color: 'error' })
      console.error('Error fetching activities:', err)
    } finally {
      loading.value = false
    }
  }

  const createEntry = async (data: ActivityEntryData) => {
    try {
      const { $api } = useNuxtApp()
      const newEntry = await $api<ActivityEntry>('/activities', {
        method: 'POST',
        body: JSON.stringify(data)
      })
      entries.value.unshift(newEntry)
      total.value++
      toast.add({ title: 'Activity', description: 'Activity logged', color: 'success' })
      return newEntry
    } catch (err) {
      toast.add({ title: 'Error', description: 'Failed to log activity', color: 'error' })
      console.error('Error creating activity:', err)
    }
  }

  const updateEntry = async (id: string, data: Partial<ActivityEntryData>) => {
    try {
      const { $api } = useNuxtApp()
      const updated = await $api<ActivityEntry>(`/activities/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data)
      })
      const idx = entries.value.findIndex(e => e.activity_id === id)
      if (idx !== -1) {
        entries.value[idx] = updated
      }
      toast.add({ title: 'Activity', description: 'Activity updated', color: 'success' })
      return updated
    } catch (err) {
      toast.add({ title: 'Error', description: 'Failed to update activity', color: 'error' })
      console.error('Error updating activity:', err)
    }
  }

  const deleteEntry = async (id: string) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/activities/${id}`, { method: 'DELETE' })
      entries.value = entries.value.filter(e => e.activity_id !== id)
      total.value = Math.max(0, total.value - 1)
      toast.add({ title: 'Activity', description: 'Activity deleted', color: 'success' })
    } catch (err) {
      toast.add({ title: 'Error', description: 'Failed to delete activity', color: 'error' })
      console.error('Error deleting activity:', err)
    }
  }

  const fetchStats = async (days: number = 30) => {
    try {
      const { $api } = useNuxtApp()
      const response = await $api<ActivityStats>(`/activities/stats?days=${days}`)
      stats.value = response
    } catch (err) {
      toast.add({ title: 'Error', description: 'Failed to fetch activity stats', color: 'error' })
      console.error('Error fetching activity stats:', err)
    }
  }

  // Recent categories/tags from loaded entries, used for quick-add suggestions
  const categorySuggestions = computed(() => {
    const seen = new Set<string>()
    for (const entry of entries.value) {
      const category = entry.category?.trim()
      if (category && !seen.has(category)) seen.add(category)
    }
    return [...seen]
  })

  const tagSuggestions = computed(() => {
    const seen = new Set<string>()
    for (const entry of entries.value) {
      for (const tag of entry.tags || []) {
        if (tag && !seen.has(tag)) seen.add(tag)
      }
    }
    return [...seen]
  })

  return {
    entries,
    total,
    stats,
    loading,
    fetchEntries,
    createEntry,
    updateEntry,
    deleteEntry,
    fetchStats,
    categorySuggestions,
    tagSuggestions
  }
})
