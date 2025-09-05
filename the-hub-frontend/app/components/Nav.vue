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
  { label: "Time", link: "/time" },
  { label: "Learning", link: "/learning" },
  { label: "Finance", link: "/finance" },
  // { label: "Stats", link: "/stats"},
  { label: "Settings", link: "/settings" },
])

const activeIndex = ref(labels.value.findIndex(item => item.link === route.path) || 0)

watch(() => route.path, (newPath) => {
  const idx = labels.value.findIndex(item => item.link === newPath)
  if (idx !== -1) activeIndex.value = idx
})

const showMenu = ref(false)
</script>

<template>
  <!-- Desktop Sidebar -->
  <div
    class="hidden md:flex fixed left-0 top-0 h-full w-64 bg-surface-light dark:bg-surface-dark backdrop-blur-md border-r border-surface-light dark:border-surface-dark shadow-lg z-50">
    <div class="flex flex-col w-full p-6">
      <!-- Logo -->
      <div
        class="mb-8 backdrop-blur-sm  p-3 w-fit mx-auto">
        <img src="/logo.svg" alt="Logo" class="h-10 w-10" />
      </div>

      <!-- Navigation Links -->
      <nav class="relative flex-1 space-y-3">
        <!-- Active indicator bar -->
        <div class="absolute left-0 w-1 bg-primary rounded h-12 transition-all duration-300"
          :style="{ top: `${activeIndex * 60}px` }"></div>

        <UiNavLink v-for="(item, index) in labels" :key="item.link" :to="item.link"
          :active="isActive(item.link)" variant="nav" @click="activeIndex = index">
          {{ item.label }}
        </UiNavLink>
      </nav>

      <!-- Logout Button -->
      <UiButton @click="logout" variant="danger" size="md" class="mt-2 w-full">
        Logout
      </UiButton>
    </div>
  </div>

  <!-- Mobile Top Navigation -->
  <div
    class="md:hidden fixed top-0 left-0 right-0 bg-surface-light dark:bg-surface-dark backdrop-blur-md border-b border-surface-light dark:border-surface-dark shadow-lg z-50">
    <div class="flex justify-between items-center p-4">
      <!-- Logo -->
      <div class="backdrop-blur-sm rounded-full p-2">
        <img src="/logo.svg" alt="Logo" class="h-8 w-8" />
      </div>

      <!-- Hamburger Menu -->
      <button
        class="bg-surface-light/10 dark:bg-background-dark backdrop-blur-sm rounded-lg p-2 border border-white/20 hover:bg-white/20 transition-all duration-300"
        @click="showMenu = !showMenu">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-text-light dark:text-text-dark" fill="none" viewBox="0 0 24 24"
          stroke="currentColor">
          <path v-if="!showMenu" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16" />
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Mobile Menu Dropdown -->
    <div v-if="showMenu" class="bg-surface-light/10 dark:bg-surface-dark/10 backdrop-blur-md border-t border-surface-light/10 dark:border-surface-dark mx-4 mb-4 rounded-lg shadow-lg">
      <div class="flex flex-col p-4 gap-2">
        <UiNavLink v-for="item in labels" :key="item.link" :to="item.link"
          :active="isActive(item.link)" variant="nav" class="w-full text-center" @click="showMenu = false">
          {{ item.label }}
        </UiNavLink>

        <UiButton @click="logout" variant="danger" size="md" class="w-full mt-2">
          Logout
        </UiButton>
      </div>
    </div>
  </div>
</template>
