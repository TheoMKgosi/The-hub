import { defineStore } from 'pinia'
import { ref } from 'vue'

interface CalendarIntegration {
  id: string
  provider: 'google' | 'outlook' | 'apple'
  provider_user_id: string
  calendar_id: string
  is_active: boolean
  sync_enabled: boolean
  last_sync_at?: string
}

interface CalendarResponse {
  integrations: CalendarIntegration[]
}

export const useCalendarStore = defineStore('calendar', () => {
  const integrations = ref<CalendarIntegration[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchIntegrations() {
    const { $api } = useNuxtApp()
    loading.value = true
    try {
      const response = await $api<CalendarResponse>('/calendar/integrations')
      integrations.value = response.integrations || []
      fetchError.value = null
    } catch (err) {
      fetchError.value = err as Error
    } finally {
      loading.value = false
    }
  }

  async function initiateGoogleAuth() {
    const { $api } = useNuxtApp()
    try {
      const { auth_url } = await $api('/calendar/google/auth', {
        method: 'POST'
      })
      // Redirect to Google OAuth
      window.location.href = auth_url
    } catch (err) {
      fetchError.value = err as Error
    }
  }

  async function syncCalendar(integrationId: string) {
    const { $api } = useNuxtApp()
    try {
      await $api(`/calendar/integrations/${integrationId}/sync`, {
        method: 'POST'
      })
      // Refresh integrations after sync
      await fetchIntegrations()
    } catch (err) {
      fetchError.value = err as Error
    }
  }

  async function deleteIntegration(integrationId: string) {
    const { $api } = useNuxtApp()
    try {
      await $api(`/calendar/integrations/${integrationId}`, {
        method: 'DELETE'
      })
      // Remove from local state
      integrations.value = integrations.value.filter(i => i.id !== integrationId)
      fetchError.value = null
    } catch (err) {
      fetchError.value = err as Error
    }
  }

  function reset() {
    integrations.value = []
  }

  return {
    integrations,
    loading,
    fetchError,
    fetchIntegrations,
    initiateGoogleAuth,
    syncCalendar,
    deleteIntegration,
    reset,
  }
})