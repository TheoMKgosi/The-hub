<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="flex bg-gray-300">
    <img alt="Vue logo" class="p-3" src="@/assets/logo.svg" width="100" height="100"/>
    <div>
      <nav class="p-4">
        <RouterLink to="/" class="p-4 hover:text-blue-900">Home</RouterLink>
        <RouterLink to="/plan" class="p-4 hover:text-blue-900">Plan</RouterLink>
        <RouterLink to="/time" class="p-4 hover:text-blue-900">Time</RouterLink>

      </nav>
    </div>
  </header>

  <RouterView />

</template>
