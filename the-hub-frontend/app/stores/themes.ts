import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Theme, ThemeSuggestion, ThemeWeekday, ApplyTheme, ThemeAssignment } from '~/types/theme'
import type { Task } from '~/types/task'

export interface ThemesResponse {
  themes: Theme[]
}

export interface ThemeSuggestionsResponse {
  suggestions: ThemeSuggestion[]
}

export const useThemeStore = defineStore('theme', () => {
  const themes = ref<Theme[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const enabled = ref(false)
  const suggestions = ref<ThemeSuggestion[]>([])
  const suggestLoading = ref(false)
  const toast = useToast()

  async function fetchThemes() {
    const { $api } = useNuxtApp()
    loading.value = true
    try {
      const response = await $api<ThemesResponse>('/themes')
      if (response) themes.value = response.themes
    } catch (error) {
      fetchError.value = error as Error
      toast.add({ title: 'Theme', description: 'Failed to fetch themes', color: 'error' })
    } finally {
      loading.value = false
    }
  }

  async function fetchEnabled() {
    const auth = useAuthStore()
    if (!auth.isLoggedIn || !auth.user?.user_id) return

    const { $api } = useNuxtApp()
    try {
      const response = await $api<{ settings: Record<string, unknown> }>(`/users/${auth.user.user_id}/settings`)
      if (response?.settings) {
        enabled.value = response.settings.theme_days_enabled === true
      }
    } catch (error) {
      console.warn('Failed to load theme days setting:', error)
    }
  }

  async function setEnabled(value: boolean) {
    const auth = useAuthStore()
    if (!auth.isLoggedIn || !auth.user?.user_id) {
      toast.add({ title: 'Theme', description: 'You must be logged in', color: 'error' })
      return
    }

    const { $api } = useNuxtApp()
    try {
      await $api(`/users/${auth.user.user_id}/settings`, {
        method: 'PATCH',
        body: JSON.stringify({ theme_days_enabled: value })
      })
      enabled.value = value
      toast.add({
        title: 'Theme Days',
        description: value ? 'Theme Days enabled' : 'Theme Days disabled',
        color: 'success'
      })
    } catch (error) {
      toast.add({ title: 'Theme', description: 'Failed to update setting', color: 'error' })
    }
  }

  async function createTheme(payload: { name: string; description: string; color: string; days: ThemeWeekday[] }) {
    const { $api } = useNuxtApp()
    try {
      const theme = await $api<Theme>('/themes', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      themes.value.push(theme)
      toast.add({ title: 'Theme', description: 'Theme created successfully', color: 'success' })
      return theme
    } catch (error: any) {
      const message = error?.data?.error || 'Failed to create theme'
      toast.add({ title: 'Theme', description: message, color: 'error' })
      throw error
    }
  }

  async function updateTheme(themeId: string, payload: { name: string; description: string; color: string; days: ThemeWeekday[] }) {
    const { $api } = useNuxtApp()
    const originalThemeIndex = themes.value.findIndex(t => t.theme_id === themeId)
    const originalTheme = originalThemeIndex !== -1 ? { ...themes.value[originalThemeIndex] } : null

    try {
      const theme = await $api<Theme>(`/themes/${themeId}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })
      if (originalThemeIndex !== -1 && theme) {
        themes.value[originalThemeIndex] = theme
      }
      toast.add({ title: 'Theme', description: 'Theme updated successfully', color: 'success' })
      return theme
    } catch (error: any) {
      if (originalTheme && originalThemeIndex !== -1) {
        themes.value[originalThemeIndex] = originalTheme
      }
      const message = error?.data?.error || 'Failed to update theme'
      toast.add({ title: 'Theme', description: message, color: 'error' })
      throw error
    }
  }

  async function deleteTheme(themeId: string) {
    const { $api } = useNuxtApp()
    const themeToDelete = themes.value.find(t => t.theme_id === themeId)

    try {
      await $api(`/themes/${themeId}`, { method: 'DELETE' })
      themes.value = themes.value.filter(t => t.theme_id !== themeId)
      toast.add({ title: 'Theme', description: 'Theme deleted successfully', color: 'success' })
    } catch (error) {
      toast.add({ title: 'Theme', description: 'Failed to delete theme', color: 'error' })
      throw error
    }
  }

  async function getAISuggestions(): Promise<ThemeSuggestion[]> {
    const { $api } = useNuxtApp()
    suggestLoading.value = true
    try {
      const response = await $api<ThemeSuggestionsResponse>('/themes/ai/suggest', { method: 'POST' })
      suggestions.value = response?.suggestions ?? []
      return suggestions.value
    } catch (error) {
      toast.add({ title: 'Theme', description: 'Failed to get AI theme suggestions', color: 'error' })
      return []
    } finally {
      suggestLoading.value = false
    }
  }

  async function applyAISuggestions(payload: {
    themes: ApplyTheme[]
    assignments: ThemeAssignment[]
  }): Promise<boolean> {
    const { $api } = useNuxtApp()
    suggestLoading.value = true
    try {
      await $api('/themes/ai/apply', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      toast.add({ title: 'Theme', description: 'Themes applied successfully', color: 'success' })
      suggestions.value = []
      await fetchThemes()
      return true
    } catch (error: any) {
      const message = error?.data?.error || 'Failed to apply themes'
      toast.add({ title: 'Theme', description: message, color: 'error' })
      return false
    } finally {
      suggestLoading.value = false
    }
  }

  const themeForWeekday = (day: string): Theme | undefined => {
    return themes.value.find(theme => theme.days.includes(day as ThemeWeekday))
  }

  const activeThemeToday = computed<Theme | undefined>(() => {
    const today = new Date().toLocaleDateString('en-US', { weekday: 'long' }).toLowerCase()
    return themeForWeekday(today)
  })

  // When Theme Days mode is enabled, only tasks assigned to today's theme are
  // shown. Tasks with no theme (or a theme not active today) are hidden.
  const filterTasksByTheme = (tasks: Task[]): Task[] => {
    if (!enabled.value) return tasks
    if (!activeThemeToday.value) return []
    return tasks.filter(task => task.theme_id === activeThemeToday.value?.theme_id)
  }

  function reset() {
    themes.value = []
    enabled.value = false
    suggestions.value = []
  }

  return {
    themes,
    loading,
    fetchError,
    enabled,
    suggestions,
    suggestLoading,
    activeThemeToday,
    fetchThemes,
    fetchEnabled,
    setEnabled,
    createTheme,
    updateTheme,
    deleteTheme,
    getAISuggestions,
    applyAISuggestions,
    filterTasksByTheme,
    reset
  }
})
