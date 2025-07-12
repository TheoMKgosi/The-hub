<script setup lang="ts">
import { ref } from "vue"
import { RouterLink, useRoute } from "vue-router"
import { useAuthStore } from "@/stores/auth"

const authStore = useAuthStore()
const route = useRoute()

const logout = () => {
  authStore.logout()
}

const isActive = (href: string) => route.path === href

const labels = ref([
  { label: "Home", link: "/dashboard" },
  { label: "Plan", link: "/plan" },
  { label: "Time", link: "/time" },
])

const showMenu = ref(false)
</script>

<template>
  <div class="bg-white shadow-md">
    <div class="flex justify-between items-center p-4">
      <!-- Logo -->
      <img src="/logo.svg" alt="Logo" class="h-12 w-12" />

      <!-- Hamburger -->
      <button class="md:hidden" @click="showMenu = !showMenu">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            v-if="!showMenu"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16"
          />
          <path
            v-else
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>

      <!-- Desktop Nav -->
      <div class="hidden md:flex items-center gap-4 bg-primary text-white px-3 py-1 rounded-2xl">
        <RouterLink
          v-for="item in labels"
          :key="item.link"
          :to="item.link"
          :class="[
            'p-2 rounded-2xl hover:text-blue-600 transition-all',
            isActive(item.link) ? 'bg-gray-300 text-black' : ''
          ]"
        >
          {{ item.label }}
        </RouterLink>
      </div>

      <!-- Logout -->
      <button
        @click="logout"
        class="hidden md:block p-2 bg-primary text-white rounded-2xl hover:text-gray-300 hover:bg-red-700 hover:scale-105 transition-all"
      >
        Logout
      </button>
    </div>

    <!-- Mobile Nav -->
    <div v-if="showMenu" class="md:hidden flex flex-col items-center gap-2 pb-4">
      <RouterLink
        v-for="item in labels"
        :key="item.link"
        :to="item.link"
        @click="showMenu = false"
        :class="[
          'w-full text-center p-2 rounded-md hover:bg-gray-200 transition',
          isActive(item.link) ? 'bg-gray-300 font-bold' : ''
        ]"
      >
        {{ item.label }}
      </RouterLink>
      <button
        @click="logout"
        class="w-full p-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
      >
        Logout
      </button>
    </div>
  </div>
</template>

