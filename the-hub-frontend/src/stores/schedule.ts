import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch'


interface Schedule {
  title: string
  start: Date
  end: Date
}

export interface ScheduleResponse {
  schedule: Schedule[]
}

export const useScheduleStore = defineStore('schedule', () => {
  const schedule = ref<Schedule[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchSchedule() {
    loading.value = true
    const { data, error } = await useMyFetch('schedule').json<ScheduleResponse>()

    if (data.value) schedule.value = data.value.schedule.map(e => ({...e, start: new Date(e.start), end: new Date(e.end)}))
    fetchError.value = error.value
    loading.value = false
  }

  async function submitForm(formData: Schedule) {
    const { data, error } = await useMyFetch('schedule').post(formData).json()
    schedule.value.push(data.value)
    fetchError.value = error.value
  }

  async function deleteSchedule(id: number) {
    await useMyFetch(`schedule/${id}`).delete().json()
    schedule.value = schedule.value.filter((t) => t.task_id !== id)
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
    reset,
  }
})
