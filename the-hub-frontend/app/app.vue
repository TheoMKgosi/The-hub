<script setup lang="ts">
useHead({
  htmlAttrs: {
    lang: 'en'
  },
  link: [
    {
      rel: 'icon',
      type: 'image/png',
      href: '/favicon.png'
    }
  ],
  title: 'The Hub'
})

</script>
<template>
  <ErrorBoundary>
    <NuxtLayout>
        <NuxtRouteAnnouncer />
        <NuxtPage />
        <UiToast />
    </NuxtLayout>
  </ErrorBoundary>
</template>
