<script setup lang="ts">


const authStore = useAuthStore()
const route = useRoute()

const logout = () => {
  authStore.logout()
}

const isActive = (href: string) => route.path === href

const labels = ref([
  { label: "Home", link: "/dashboard" },
  { label: "Plan", link: "/plan" },
  // { label: "Time", link: "/time" },
  { label: "Learning", link: "/learning" },
  { label: "Finance", link: "/finance" },
  { label: "Settings", link: "/settings" },
])

const showMenu = ref(false)
</script>

<template>
  <!-- Desktop Sidebar -->
  <div class="hidden md:flex fixed left-0 top-0 h-full w-64 bg-white/20 backdrop-blur-md border-r border-white/10 shadow-lg z-50">
    <div class="flex flex-col w-full p-6">
      <!-- Logo -->
      <div class="mb-8 bg-white/10 backdrop-blur-sm rounded-full p-3 border border-white/20 w-fit mx-auto">
        <img src="/logo.svg" alt="Logo" class="h-10 w-10" />
      </div>

      <!-- Navigation Links -->
      <nav class="flex-1 space-y-3">
        <NuxtLink
          v-for="item in labels"
          :key="item.link"
          :to="item.link"
          :class="[
            'flex items-center px-4 py-3 rounded-xl text-gray-700 font-medium transition-all duration-300 hover:shadow-md w-full',
            isActive(item.link)
              ? 'bg-gray-800/60 backdrop-blur-sm text-white shadow-lg border border-gray-700/40'
              : 'hover:text-gray-900 hover:bg-white/30'
          ]"
        >
          {{ item.label }}
        </NuxtLink>
      </nav>

      <!-- Logout Button -->
      <button
        @click="logout"
        class="mt-6 w-full px-4 py-3 bg-gradient-to-r from-red-500/20 to-red-600/20 backdrop-blur-sm text-red-700 rounded-xl border border-red-200/30 hover:from-red-500/30 hover:to-red-600/30 hover:text-red-800 hover:shadow-md hover:scale-105 transition-all duration-300 font-medium"
      >
        Logout
      </button>
    </div>
  </div>

  <!-- Mobile Top Navigation -->
  <div class="md:hidden fixed top-0 left-0 right-0 bg-white/20 backdrop-blur-md border-b border-white/10 shadow-lg z-50">
    <div class="flex justify-between items-center p-4">
      <!-- Logo -->
      <div class="bg-white/10 backdrop-blur-sm rounded-full p-2 border border-white/20">
        <img src="/logo.svg" alt="Logo" class="h-8 w-8" />
      </div>

      <!-- Hamburger Menu -->
      <button
        class="bg-white/10 backdrop-blur-sm rounded-lg p-2 border border-white/20 hover:bg-white/20 transition-all duration-300"
        @click="showMenu = !showMenu"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path v-if="!showMenu" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Mobile Menu Dropdown -->
    <div
      v-if="showMenu"
      class="bg-white/10 backdrop-blur-md border-t border-white/10 mx-4 mb-4 rounded-lg shadow-lg"
    >
      <div class="flex flex-col p-4 gap-2">
        <NuxtLink
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
        </NuxtLink>

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
