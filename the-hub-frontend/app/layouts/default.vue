<script setup lang="ts">
// Initialize dark mode composable
const { isDark } = useDarkMode()

// Watch for system theme changes
onMounted(() => {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleChange = () => {
    // The composable will handle system theme changes automatically
    // This is just to ensure consistency
    if (localStorage.getItem('theme') === 'system') {
      document.documentElement.classList.toggle('dark', mediaQuery.matches)
    }
  }

  mediaQuery.addEventListener('change', handleChange)

  // Cleanup
  onUnmounted(() => {
    mediaQuery.removeEventListener('change', handleChange)
  })
})
</script>
<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <UiNav />
    <div class="flex-1 ml-0 mt-20 md:ml-64 md:mt-0 min-h-[calc(100vh-5rem)] md:min-h-screen">
      <slot />
    </div>
  </div>
</template>
