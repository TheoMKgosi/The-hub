<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, RouterView } from 'vue-router'
import Nav from "./components/Nav.vue";
import Toast from '@/components/ui/Toast.vue'

const isOpen = ref(false)
</script>

<template>
  <Nav />

  <RouterView />
  <Toast />
</template>
