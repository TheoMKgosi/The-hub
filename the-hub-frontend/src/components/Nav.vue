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
  { label: "Learning", link: "/learning" },
  { label: "Settings", link: "/settings" },
])

const showMenu = ref(false)
</script>

<template>
  <div class="bg-white/20 backdrop-blur-md border-b border-white/10 shadow-lg sticky">
    <div class="flex justify-between items-center p-4">
      <!-- Logo -->
      <div class="bg-white/10 backdrop-blur-sm rounded-full p-2 border border-white/20">
        <img src="/logo.svg" alt="Logo" class="h-8 w-8" />
      </div>

      <!-- Hamburger -->
      <button
        class="md:hidden bg-white/10 backdrop-blur-sm rounded-lg p-2 border border-white/20 hover:bg-white/20 transition-all duration-300"
        @click="showMenu = !showMenu"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6 text-gray-700"
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
      <div class="hidden md:flex items-center gap-2 bg-white/10 backdrop-blur-md border border-white/20 px-4 py-2 rounded-2xl shadow-lg">
        <RouterLink
          v-for="item in labels"
          :key="item.link"
          :to="item.link"
          :class="[
            'px-4 py-2 rounded-xl text-gray-700 font-medium transition-all duration-300 hover:bg-white/30 hover:shadow-md',
            isActive(item.link)
              ? 'bg-gray-800/60 backdrop-blur-sm text-white shadow-lg border border-gray-700/40'
              : 'hover:text-gray-900'
          ]"
        >
          {{ item.label }}
        </RouterLink>
      </div>

      <!-- Logout -->
      <button
        @click="logout"
        class="hidden md:block px-4 py-2 bg-gradient-to-r from-red-500/20 to-red-600/20 backdrop-blur-sm text-red-700 rounded-xl border border-red-200/30 hover:from-red-500/30 hover:to-red-600/30 hover:text-red-800 hover:shadow-md hover:scale-105 transition-all duration-300 font-medium"
      >
        Logout
      </button>
    </div>

    <!-- Mobile Nav -->
    <div
      v-if="showMenu"
      class="md:hidden bg-white/10 backdrop-blur-md border-t border-white/10 mx-4 mb-4 rounded-lg shadow-lg"
    >
      <div class="flex flex-col p-4 gap-2">
        <RouterLink
          v-for="item in labels"
          :key="item.link"
          :to="item.link"
          @click="showMenu = false"
          :class="[
            'w-full text-center p-3 rounded-lg text-gray-700 font-medium transition-all duration-300 hover:bg-white/20 hover:shadow-md',
            isActive(item.link)
              ? 'bg-gray-800/60 backdrop-blur-sm text-white shadow-lg border border-gray-700/40'
              : ''
          ]"
        >
          {{ item.label }}
        </RouterLink>

        <button
          @click="logout"
          class="w-full p-3 bg-gradient-to-r from-red-500/20 to-red-600/20 backdrop-blur-sm text-red-700 rounded-lg border border-red-200/30 hover:from-red-500/30 hover:to-red-600/30 hover:text-red-800 hover:shadow-md transition-all duration-300 font-medium mt-2"
        >
          Logout
        </button>
      </div>
    </div>
  </div>
</template>
