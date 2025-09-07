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
    // Store original integration for potential rollback
    const originalIntegrationIndex = integrations.value.findIndex(i => i.id === integrationId)
    const originalIntegration = originalIntegrationIndex !== -1 ? { ...integrations.value[originalIntegrationIndex] } : null

    // Optimistically update the last_sync_at
    if (originalIntegrationIndex !== -1) {
      integrations.value[originalIntegrationIndex] = {
        ...integrations.value[originalIntegrationIndex],
        last_sync_at: new Date().toISOString()
      }
    }

    const { $api } = useNuxtApp()
    try {
      const data = await $api<CalendarIntegration>(`/calendar/integrations/${integrationId}/sync`, {
        method: 'POST'
      })

      // Update with server response to ensure consistency
      if (originalIntegrationIndex !== -1 && data) {
        integrations.value[originalIntegrationIndex] = data
      }
    } catch (err) {
      // Revert optimistic update on error
      if (originalIntegration && originalIntegrationIndex !== -1) {
        integrations.value[originalIntegrationIndex] = originalIntegration
      }
      fetchError.value = err as Error
    }
  }

  async function deleteIntegration(integrationId: string) {
    // Store the integration for potential rollback
    const integrationToDelete = integrations.value.find(i => i.id === integrationId)
    if (!integrationToDelete) {
      fetchError.value = new Error("Integration not found")
      return
    }

    // Optimistically remove from local state
    integrations.value = integrations.value.filter(i => i.id !== integrationId)

    const { $api } = useNuxtApp()
    try {
      await $api(`/calendar/integrations/${integrationId}`, {
        method: 'DELETE'
      })
      fetchError.value = null
    } catch (err) {
      // Restore the integration on error
      integrations.value.push(integrationToDelete)
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