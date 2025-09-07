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

export const useScheduleStore = defineStore('schedule', () => {
  const schedule = ref<Schedule[]>([])
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

  async function getAISuggestions() {
    const { $api } = useNuxtApp()
    try {
      const { suggestions } = await $api('/ai/suggestions')
      fetchError.value = null
      return suggestions || []
    } catch (err) {
      fetchError.value = err as Error
      return []
    }
  }

  function reset() {
    schedule.value = []
  }

  return {
    schedule,
    loading,
    fetchError,
    fetchSchedule,
    submitForm,
    deleteSchedule,
    getAISuggestions,
    reset,
  }
})
