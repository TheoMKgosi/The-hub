<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth"

const authStore = useAuthStore()

const route = useRoute()

const logout = () => {
  authStore.logout()
}

const isActive = (href) => route.path === href

const labels = ref([
  {
    label: "Home",
    link: "/dashboard"
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
  <div class="flex justify-between p-3 items-center ">
    <img src="/logo.svg" alt="Logo" class="h-12 w-12">
    <div class="flex justify-between items-center bg-primary text-white rounded-2xl">
      <RouterLink v-for="item in labels" :to="item.link" :class="isActive(item.link) ?
      'bg-gray-300 text-black': ''" class="p-4 hover:text-blue-600 rounded-2xl
        active: hover:scale-110">{{ item.label
      }}</RouterLink>
    </div>
    <button @click="logout" class="p-4 hover:text-gray-300 hover:bg-red-700 hover:scale-110 cursor-pointer
      bg-primary text-white border-l-black rounded-2xl">Logout</button>
  </div>
</template>
