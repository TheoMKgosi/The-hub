type ThemeMode = 'light' | 'dark' | 'system'

export function useDarkMode() {
  const themeMode = ref<ThemeMode>('system')
  const isDark = ref(false)

  // Initialize theme on component mount
  onMounted(async () => {
    // First try to get theme from localStorage for immediate response
    const savedTheme = localStorage.getItem('theme') as ThemeMode | null
    if (savedTheme && ['light', 'dark', 'system'].includes(savedTheme)) {
      themeMode.value = savedTheme
      applyTheme(themeMode.value)
    }

    // Then try to load from backend if user is logged in
    await loadThemeFromBackend()
  })

  // Watch for theme changes
  watch(themeMode, (newTheme) => {
    applyTheme(newTheme)
    // Save to localStorage
    if (newTheme === 'system') {
      localStorage.removeItem('theme')
    } else {
      localStorage.setItem('theme', newTheme)
    }
  })

  const applyTheme = (theme: ThemeMode) => {
    const root = document.documentElement

    if (theme === 'dark') {
      root.classList.add('dark')
      isDark.value = true
    } else if (theme === 'light') {
      root.classList.remove('dark')
      isDark.value = false
    } else if (theme === 'system') {
      // Check system preference
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.classList.toggle('dark', prefersDark)
      isDark.value = prefersDark
    }
  }

  const setTheme = (theme: ThemeMode) => {
    themeMode.value = theme
  }

  // Load theme from backend settings
  const loadThemeFromBackend = async () => {
    try {
      const auth = useAuthStore()
      if (!auth.isLoggedIn || !auth.user?.id) return

      const { $api } = useNuxtApp()
      const response = await $api(`/users/${auth.user.id}/settings`)

      if (response.settings?.theme?.mode) {
        const backendTheme = response.settings.theme.mode as ThemeMode
        if (['light', 'dark', 'system'].includes(backendTheme)) {
          themeMode.value = backendTheme
          applyTheme(backendTheme)
        }
      }
    } catch (error) {
      // Silently fail - localStorage theme will be used as fallback
      console.warn('Failed to load theme from backend:', error)
    }
  }

  const toggleTheme = () => {
    if (themeMode.value === 'light') {
      setTheme('dark')
    } else if (themeMode.value === 'dark') {
      setTheme('system')
    } else {
      setTheme('light')
    }
  }

  return {
    themeMode: readonly(themeMode),
    isDark: readonly(isDark),
    setTheme,
    toggleTheme,
    loadThemeFromBackend
  }
}