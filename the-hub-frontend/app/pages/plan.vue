<script setup lang="ts">
const taskStore = useTaskStore()
onMounted(() => {
  taskStore.fetchTasks()
})
</script>
<template>
  <div id="plan" class="p-2">
    <TaskTemplate :taskList="taskStore.tasks">
    </TaskTemplate>
  </div>
</template>
