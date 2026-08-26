<script setup lang="ts">
const themeStore = useThemeStore()

const toggle = async () => {
  const next = !themeStore.enabled
  await themeStore.setEnabled(next)
  if (next && themeStore.themes.length === 0) {
    await themeStore.fetchThemes()
  }
}
</script>

<template>
  <UButton
    :label="themeStore.enabled ? 'Theme Days: On' : 'Theme Days: Off'"
    :icon="themeStore.enabled ? 'i-lucide-sun' : 'i-lucide-moon'"
    :active="themeStore.enabled"
    variant="outline"
    @click="toggle"
  />
</template>
