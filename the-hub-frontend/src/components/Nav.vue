<script setup lang="ts">
import { ref } from "vue";
import { RouterLink } from "vue-router";
import { useAuthStore } from "@/stores/auth"

const authStore = useAuthStore()


const logout = () => {
  authStore.logout()
}

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
  },
])
</script>

<template>
  <div class="flex justify-between p-3 items-center">
    <div class="flex justify-between items-center">
      <img src="/logo.svg" alt="Logo" class="h-12 w-12">
      <RouterLink v-for="item in labels" :to="item.link" class="p-4 hover:text-blue-600 hover:scale-110">{{ item.label
        }}</RouterLink>
    </div>
    <button @click="logout" class="p-4 hover:text-blue-600 hover:scale-110 cursor-pointer">Logout</button>
  </div>
</template>
