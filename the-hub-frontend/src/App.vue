<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, RouterView } from 'vue-router'
import Nav from "./components/Nav.vue";

const isOpen = ref(false)
</script>

<template>
  <Nav />

  <RouterView />
</template>
