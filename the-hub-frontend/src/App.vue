<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="flex">
    <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />
    <div>
      <nav class="p-4">
        <RouterLink to="/" class="p-4 hover:text-blue-900">Home</RouterLink>
        <RouterLink to="/plan" class="p-4 hover:text-blue-900">Plan</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />

</template>

