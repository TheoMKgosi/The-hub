<script setup lang="ts">
useHead({ title: 'Activity Log' })
</script>

<template>
  <div class="p-4 space-y-8">
    <div>
      <h1 class="text-2xl font-bold">Activity Log</h1>
      <p class="text-dimmed">Record what you have been doing so you can analyse it later.</p>
    </div>

    <div class="rounded-lg border border-surface-light dark:border-surface-dark p-4">
      <h2 class="mb-3 text-lg font-semibold">Log an activity</h2>
      <ActivityQuickAdd />
    </div>

    <div>
      <h2 class="mb-3 text-lg font-semibold">History</h2>
      <ActivityTimeline />
    </div>

    <ActivityStats />
  </div>
</template>
