<script setup lang="ts">
// Keyboard shortcuts and command palette
const { commandPaletteOpen, closeCommandPalette } = useKeyboardShortcuts()

useHead({
  htmlAttrs: {
    lang: 'en'
  },
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no' },
    { name: 'mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
    { name: 'apple-mobile-web-app-title', content: 'The Hub' },
    { name: 'theme-color', content: '#3B82F6' },
    { name: 'msapplication-TileColor', content: '#3B82F6' }
  ],
  link: [
    {
      rel: 'icon',
      type: 'image/png',
      href: '/favicon.png'
    },
    {
      rel: 'apple-touch-icon',
      href: '/icon-192x192.png'
    },
    {
      rel: 'manifest',
      href: '/manifest.json'
    }
  ],
  title: 'The Hub'
})

</script>
<template>
  <ErrorBoundary>
    <NuxtLoadingIndicator />
    <NuxtLayout>
      <NuxtRouteAnnouncer />
      <NuxtPage />
      <Toast />
      <!-- Command Palette -->
      <CommandPalette :is-open="commandPaletteOpen" @close="closeCommandPalette" />
    </NuxtLayout>
  </ErrorBoundary>
</template>
