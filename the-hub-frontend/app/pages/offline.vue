<template>
  <div class="min-h-screen flex items-center justify-center bg-surface-light dark:bg-surface-dark">
    <div class="text-center p-8">
      <div class="text-6xl mb-4">📶</div>
      <h1 class="text-2xl font-bold text-text-light dark:text-text-dark mb-2">
        You're Offline
      </h1>
      <p class="text-text-light/70 dark:text-text-dark/70 mb-6">
        It looks like you're not connected to the internet. Some features may not be available.
      </p>
      <button
        @click="retryConnection"
        class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors"
      >
        Try Again
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const retryConnection = () => {
  if (navigator.onLine) {
    window.location.reload();
  } else {
    // Show a message that they're still offline
    alert('You are still offline. Please check your internet connection.');
  }
};

// Listen for online/offline events
const handleOnline = () => {
  // Optionally auto-reload when back online
  // window.location.reload();
};

const handleOffline = () => {
  // Could show additional offline messaging
};

onMounted(() => {
  window.addEventListener('online', handleOnline);
  window.addEventListener('offline', handleOffline);
});

onUnmounted(() => {
  window.removeEventListener('online', handleOnline);
  window.removeEventListener('offline', handleOffline);
});
</script>