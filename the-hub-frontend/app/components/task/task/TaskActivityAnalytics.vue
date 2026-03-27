<script setup lang="ts">
interface TaskActivityResponse {
  start_date: string
  end_date: string
  tasks_created: number
  tasks_completed: number
}

const { $api } = useNuxtApp()
const { addToast } = useToast()

const loading = ref(false)
const error = ref<string | null>(null)
const data = ref<TaskActivityResponse | null>(null)

const toDateInput = (d: Date) => {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const endDate = ref(toDateInput(new Date()))
const startDate = ref((() => {
  const d = new Date()
  d.setDate(d.getDate() - 7)
  return toDateInput(d)
})())

const load = async () => {
  if (!startDate.value || !endDate.value) return

  loading.value = true
  error.value = null

  try {
    const params = new URLSearchParams({
      start_date: startDate.value,
      end_date: endDate.value,
    })

    data.value = await $api<TaskActivityResponse>(`/stats/tasks/activity?${params.toString()}`)
  } catch (err: any) {
    const message = err?.message || 'Failed to load task analytics'
    error.value = message
    addToast(message, 'error')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
})
</script>

<template>
  <div class="space-y-4 p-4">
    <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-4">
      <div class="flex flex-col md:flex-row md:items-end md:justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Task Analytics</h3>
          <p class="text-sm text-text-light/70 dark:text-text-dark/70">Created vs done in a selected range</p>
        </div>

        <div class="flex flex-col sm:flex-row gap-2">
          <input v-model="startDate" type="date"
                 class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          <input v-model="endDate" type="date"
                 class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />

          <UiButton variant="primary" size="sm" :disabled="loading" @click="load">
            <span v-if="loading" class="animate-spin">⟳</span>
            <span v-else>Apply</span>
          </UiButton>
        </div>
      </div>
    </div>

    <div v-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
      <p class="text-red-600 dark:text-red-400">{{ error }}</p>
    </div>

    <div v-else-if="data" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-5">
        <p class="text-sm text-text-light/70 dark:text-text-dark/70">Tasks Created</p>
        <p class="text-3xl font-bold text-text-light dark:text-text-dark mt-1">{{ data.tasks_created }}</p>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark p-5">
        <p class="text-sm text-text-light/70 dark:text-text-dark/70">Tasks Done</p>
        <p class="text-3xl font-bold text-success mt-1">{{ data.tasks_completed }}</p>
      </div>
    </div>

    <div v-else class="text-center py-10 text-text-light/70 dark:text-text-dark/70">
      Select a date range to see analytics.
    </div>
  </div>
</template>
