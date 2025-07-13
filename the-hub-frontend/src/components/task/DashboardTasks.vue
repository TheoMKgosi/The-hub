<script setup lang="ts">
import { onMounted, computed } from "vue";
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

onMounted(() => {
  if (taskStore.tasks.length === 0) {
    taskStore.fetchTasks()
  }
})

</script>

<template>
  <div>
    <div class="p-6 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-900 flex items-center">
          <span class="text-green-600 mr-2">✓</span>
          Task Management
        </h2>
        <button class="text-gray-400 hover:text-gray-600">
          <span class="text-lg">+</span>
        </button>
      </div>
    </div>
    <p v-if="taskStore.loading" class="p-4">Loading...</p>
    <ul v-else class="p-4">
      <p v-if="taskStore.tasks.length === 0">No tasks added</p>
      <li v-for="task in taskStore.tasks" :key="task.task_id" class="bg-white shadow rounded-lg p-4 mb-4 border-l-4"
        :class="task.status === 'complete' ? 'border-green-500' : 'border-yellow-500'">

        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
          <div class="flex items-center">
            <input type="checkbox" class="h-4 w-4 text-green-600 rounded" />
            <span class="ml-3 text-sm text-gray-900">{{ task.title }}</span>
          </div>
          <span class="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full">{{ task.priority }}</span>
        </div>
      </li>
    </ul>
  </div>
</template>
