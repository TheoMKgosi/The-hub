<script setup lang="ts">
import { onMounted } from "vue";
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
            <input type="checkbox" checked class="h-4 w-4 text-green-600 rounded" />
            <span class="ml-3 text-sm text-gray-900 line-through">{{ task.title }}</span>
          </div>
          <span class="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full">{{ task.priority }}</span>
        </div>
      </li>
    </ul>
  </div>
  <!-- Task Management -->

  <!--
  <div class="bg-white rounded-lg shadow">
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
    -->
  <!--
    <div class="p-6">
      <div class="space-y-4">
        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
          <div class="flex items-center">
            <input type="checkbox" class="h-4 w-4 text-green-600 rounded" />
            <span class="ml-3 text-sm text-gray-900">Review team feedback</span>
          </div>
          <span class="px-2 py-1 text-xs bg-yellow-100 text-yellow-800 rounded-full">Medium</span>
        </div>
        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
          <div class="flex items-center">
            <input type="checkbox" class="h-4 w-4 text-green-600 rounded" />
            <span class="ml-3 text-sm text-gray-900">Update documentation</span>
          </div>
          <span class="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full">Low</span>
        </div>
        <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
          <div class="flex items-center">
            <input type="checkbox" class="h-4 w-4 text-green-600 rounded" />
            <span class="ml-3 text-sm text-gray-900">Schedule client meeting</span>
          </div>
          <span class="px-2 py-1 text-xs bg-green-100 text-green-800 rounded-full">High</span>
        </div>
      </div>
      <div class="mt-6">
        <div class="flex justify-between text-sm text-gray-600 mb-2">
          <span>Progress</span>
          <span>75%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div class="bg-green-600 h-2 rounded-full" style="width: 75%"></div>
        </div>
      </div>
    </div>
  </div>
-->

</template>
