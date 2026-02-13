<script setup lang="ts">
const taskStore = useTaskStore()
onMounted(() => {
  if (taskStore.tasks.length === 0) taskStore.fetchTasks()
})
</script>
<template>
  <div class="p-2">
    <TaskTemplate :taskList="taskStore.tasks" />
  </div>
</template>
