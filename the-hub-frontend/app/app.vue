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
  <NuxtLayout>
      <NuxtRouteAnnouncer />
      <NuxtPage />
      <UiToast />
  </NuxtLayout>
</template>
