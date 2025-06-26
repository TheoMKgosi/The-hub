<script setup lang="ts">
import { ref } from "vue";
import { RouterLink } from "vue-router";

const labels = ref([
  {
    label: "Home",
    link: "/"
  },
  {
    label: "Plan",
    link: "/plan"
  },
  {
    label: "Time",
    link: "/time"
  }
])
</script>

<template>
  <div class="flex">
    <RouterLink v-for="item in labels" :to="item.link">{{ item.label }}</RouterLink>
  </div>
</template>
