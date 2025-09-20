import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task } from './tasks'

interface Schedule {
  id: string
  title: string
  start: Date
  end: Date
  task_id?: string
  task?: Task
  recurrence_rule_id?: string
  recurrence_rule?: any
  created_by_ai?: boolean
}

export interface ScheduleResponse {
  schedule: Schedule[]
}

interface Suggestion {
  id: string
  title: string
  start: string
  end: string
  task_id?: string
  created_by_ai?: boolean
}

export const useScheduleStore = defineStore('schedule', () => {
  const schedule = ref<Schedule[]>([])
  const suggestions = ref<Suggestion[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchSchedule() {
    const { $api } = useNuxtApp()
    loading.value = true
    try {
      const { schedule: fetchedSchedule } = await $api<ScheduleResponse>('/schedule')
      if (fetchedSchedule) schedule.value = fetchedSchedule.map(e => ({...e, start: new Date(e.start), end: new Date(e.end)}))
      fetchError.value = null
    } catch (err) {
      fetchError.value = err as Error
    } finally {
      loading.value = false
    }
  }

  async function submitForm(formData: any) {
    // Create optimistic schedule item
    const optimisticSchedule: Schedule = {
      id: `temp-${Date.now()}`,
      title: formData.title || 'New Schedule Item',
      start: new Date(formData.start),
      end: new Date(formData.end),
      task_id: formData.task_id,
      recurrence_rule_id: formData.recurrence_rule_id,
      created_by_ai: formData.created_by_ai || false
    }

    // Optimistically add to local state
    schedule.value.push(optimisticSchedule)

    const { $api } = useNuxtApp()
    try {
      const data = await $api<Schedule>('schedule', {
        method: 'POST',
        body: JSON.stringify(formData)
      })

      // Replace optimistic schedule with real data
      const optimisticIndex = schedule.value.findIndex(s => s.id === optimisticSchedule.id)
      if (optimisticIndex !== -1) {
        schedule.value[optimisticIndex] = {
          ...data,
          start: new Date(data.start),
          end: new Date(data.end)
        }
      }

      fetchError.value = null
    } catch (err) {
      // Remove optimistic schedule on error
      schedule.value = schedule.value.filter(s => s.id !== optimisticSchedule.id)
      fetchError.value = err as Error
    }
  }

  function validateScheduleInput(start: Date, end: Date): string | null {
    if (start >= end) {
      return "Start time must be before end time"
    }

    // Check for reasonable duration (max 24 hours)
    const durationMs = end.getTime() - start.getTime()
    if (durationMs > 24 * 60 * 60 * 1000) {
      return "Event duration cannot exceed 24 hours"
    }

    // Check for past dates (allow events up to 1 hour in the past for flexibility)
    const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000)
    if (start < oneHourAgo) {
      return "Cannot schedule events in the past"
    }

    return null
  }

  function hasTimeConflict(start: Date, end: Date, excludeId?: string): boolean {
    return schedule.value.some(s => {
      if (excludeId && s.id === excludeId) return false
      // Check if time ranges overlap
      return start < s.end && end > s.start
    })
  }

  async function updateSchedule(id: string, updateData: Partial<Schedule>): Promise<boolean> {
    // Find the schedule item for potential rollback
    const scheduleToUpdate = schedule.value.find(s => s.id === id)
    if (!scheduleToUpdate) {
      fetchError.value = new Error("Schedule item not found")
      return false
    }

    // Validate time fields if provided
    const newStart = updateData.start || scheduleToUpdate.start
    const newEnd = updateData.end || scheduleToUpdate.end

    if (newStart && newEnd) {
      const validationError = validateScheduleInput(newStart, newEnd)
      if (validationError) {
        fetchError.value = new Error(validationError)
        return false
      }

      // Check for conflicts with other events
      if (hasTimeConflict(newStart, newEnd, id)) {
        fetchError.value = new Error("This time slot conflicts with an existing scheduled event")
        return false
      }
    }

    // Store original data for rollback
    const originalData = { ...scheduleToUpdate }

    // Optimistically update local state
    const index = schedule.value.findIndex(s => s.id === id)
    schedule.value[index] = { ...scheduleToUpdate, ...updateData }

    const { $api } = useNuxtApp()
    try {
      const data = await $api<Schedule>(`schedule/${id}`, {
        method: 'PUT',
        body: JSON.stringify(updateData)
      })

      // Update with real data from server
      schedule.value[index] = {
        ...data,
        start: new Date(data.start),
        end: new Date(data.end)
      }

      fetchError.value = null
      return true
    } catch (err) {
      // Restore original data on error
      schedule.value[index] = originalData
      fetchError.value = err as Error
      return false
    }
  }

  async function deleteSchedule(id: string) {
    // Store the schedule item for potential rollback
    const scheduleToDelete = schedule.value.find(s => s.id === id)
    if (!scheduleToDelete) {
      fetchError.value = new Error("Schedule item not found")
      return
    }

    // Optimistically remove from local state
    schedule.value = schedule.value.filter((s) => s.id !== id)

    const { $api } = useNuxtApp()
    try {
      await $api(`schedule/${id}`, {
        method: 'DELETE'
      })
      fetchError.value = null
    } catch (err) {
      // Restore the schedule item on error
      schedule.value.push(scheduleToDelete)
      fetchError.value = err as Error
    }
  }

  async function getSuggestions() {
    const { $api } = useNuxtApp()
    try {
      const { suggestions: fetchedSuggestions } = await $api('/schedule/suggestions')
      fetchError.value = null

      // Store suggestions locally for persistence
      if (fetchedSuggestions && fetchedSuggestions.length > 0) {
        suggestions.value = fetchedSuggestions
        saveSuggestionsToStorage()
      }

      return fetchedSuggestions || []
    } catch (err) {
      fetchError.value = err as Error
      // Return cached suggestions if API fails
      return loadSuggestionsFromStorage()
    }
  }

  function saveSuggestionsToStorage() {
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem('schedule-suggestions', JSON.stringify(suggestions.value))
      } catch (error) {
        console.warn('Failed to save suggestions to localStorage:', error)
      }
    }
  }

  function loadSuggestionsFromStorage(): Suggestion[] {
    if (typeof window !== 'undefined') {
      try {
        const stored = localStorage.getItem('schedule-suggestions')
        if (stored) {
          const parsed = JSON.parse(stored)
          suggestions.value = parsed
          return parsed
        }
      } catch (error) {
        console.warn('Failed to load suggestions from localStorage:', error)
      }
    }
    return []
  }

  function removeSuggestion(suggestionId: string) {
    suggestions.value = suggestions.value.filter(s => s.id !== suggestionId)
    saveSuggestionsToStorage()
  }

  function clearSuggestions() {
    suggestions.value = []
    if (typeof window !== 'undefined') {
      localStorage.removeItem('schedule-suggestions')
    }
  }

  function reset() {
    schedule.value = []
  }

  return {
    schedule,
    suggestions,
    loading,
    fetchError,
    fetchSchedule,
    submitForm,
    updateSchedule,
    deleteSchedule,
    getSuggestions,
    removeSuggestion,
    clearSuggestions,
    reset,
  }
})
