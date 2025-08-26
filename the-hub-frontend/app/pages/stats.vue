<template>
  <div class="p-4">
    <Tab :tabs="['Tasks','Learning', 'Finance']"/>
    <!-- Task Statistics Section -->
    <div class="mb-8">
      <TaskStatsDashboard />
    </div>
  </div>
</template>
