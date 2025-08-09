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
  <Nav />
  <div class="flex-1 ml-0 md:ml-64 bg-[#FAF9F6]">
    <NuxtRouteAnnouncer />
    <NuxtPage />
    <UiToast />
  </div>
</template>
