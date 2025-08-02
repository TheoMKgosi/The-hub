<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, RouterView } from 'vue-router'
import Nav from "./components/Nav.vue";
import Toast from '@/components/ui/Toast.vue'

const isOpen = ref(false)
</script>

<template>
  <Nav />
  <div class="mt-16 md:ml-64 md:mt-0">
    <RouterView />
  </div>
  <Toast />
</template>
