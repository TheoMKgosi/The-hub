<script setup>

const sidebarOpen = ref(false)

const logout = () => {
  console.log('Logging out...') // You'd probably trigger an auth logout here
}
</script>

<template>
  <div class="flex h-screen overflow-hidden">
    <!-- Sidebar -->
    <aside
      class="bg-gray-900 text-white w-64 flex flex-col justify-between p-4 fixed inset-y-0 left-0 transform transition-transform duration-200 ease-in-out z-50
              md:translate-x-0 -translate-x-full md:relative md:block"
      :class="{ 'translate-x-0': sidebarOpen }"
    >
      <!-- Top (Logo) -->
      <div>
        <h1 class="text-2xl font-bold mb-8">The Hub</h1>
        <!-- Navigation -->
        <nav class="space-y-4">
          <NuxtLink to="/Dashboard" class="block hover:text-green-400">Dashboard</NuxtLink>
          <NuxtLink to="/plan" class="block hover:text-green-400">Plan</NuxtLink>
          <!--
          <NuxtLink to="/time" class="block hover:text-green-400">Time</NuxtLink>
          <NuxtLink to="/learning" class="block hover:text-green-400">Learning</NuxtLink>
          <NuxtLink to="/finance" class="block hover:text-green-400">Finance</NuxtLink>
          -->
        </nav>
      </div>

      <!-- Bottom (Logout) -->
      <button class="mt-8 text-red-400 hover:text-red-600" @click="logout">Logout</button>
    </aside>

    <!-- Mobile menu toggle -->
    <button
      class="absolute top-4 left-4 z-50 md:hidden text-white bg-gray-800 p-2 rounded"
      @click="sidebarOpen = !sidebarOpen"
    >
      ☰
    </button>

    <!-- Page content -->
    <div class="flex-1 ml-0 md:ml-64 p-4 overflow-auto w-full">
      <NuxtRouteAnnouncer />
      <NuxtPage />
    </div>
  </div>
</template>

