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
    const { $api } = useNuxtApp()
    try {
      const data = await $api<Schedule>('schedule', {
        method: 'POST',
        body: JSON.stringify(formData)
      })
      if (data) {
        schedule.value.push({
          ...data,
          start: new Date(data.start),
          end: new Date(data.end)
        })
      }
      fetchError.value = null
    } catch (err) {
      fetchError.value = err as Error
    }
  }

  async function deleteSchedule(id: string) {
    const { $api } = useNuxtApp()
    try {
      await $api(`schedule/${id}`, {
        method: 'DELETE'
      })
      schedule.value = schedule.value.filter((t) => t.id !== id)
      fetchError.value = null
    } catch (err) {
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
