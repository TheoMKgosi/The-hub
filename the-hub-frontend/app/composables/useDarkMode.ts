type ThemeMode = 'light' | 'dark' | 'system'

export function useDarkMode() {
  const themeMode = ref<ThemeMode>('system')
  const isDark = ref(false)

  // Initialize theme on component mount
  onMounted(() => {
    // Get saved theme from localStorage or default to system
    const savedTheme = localStorage.getItem('theme') as ThemeMode | null
    if (savedTheme && ['light', 'dark', 'system'].includes(savedTheme)) {
      themeMode.value = savedTheme
    }

    // Apply the theme
    applyTheme(themeMode.value)
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
    toggleTheme
  }
}