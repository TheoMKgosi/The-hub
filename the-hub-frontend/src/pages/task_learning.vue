<script setup lang="ts">
import { useTaskLearningStore } from '@/stores/task-learning';
import { useTopicStore } from '@/stores/topics';
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const topicID = parseInt(route.params.topic_id)

const topicStore = useTopicStore()
const taskLearningStore = useTaskLearningStore()


onMounted(() => {
  taskLearningStore.fetchTasks(topicID)
})

const addTask = () => {
  taskLearningStore.createTask()
}
</script>
<template>
  <div class="max-w-3xl mx-auto p-6 space-y-6">

    <!-- Header -->
    <header class="space-y-2">
      <h1 class="text-3xl font-bold">Master SQL Joins</h1>
      <div class="flex flex-wrap gap-2 text-sm text-gray-600">
        <span class="bg-gray-100 px-2 py-1 rounded">🎯 goal</span>
        <span class="bg-gray-100 px-2 py-1 rounded">📅 Due 2/15/2024</span>
        <span class="bg-gray-100 px-2 py-1 rounded">⏱️ 20h estimated</span>
      </div>
      <button @click="addTask" class="mt-2 px-4 py-2 bg-black text-white rounded hover:bg-gray-800">
        + Add Task
      </button>
    </header>

    <!-- Progress Overview -->
    <section class="bg-white p-4 shadow rounded space-y-3">
      <h2 class="text-xl font-semibold">Progress Overview</h2>
      <p class="text-sm text-gray-500">1 of 2 tasks completed</p>
      <div class="w-full bg-gray-200 h-2 rounded">
        <div class="bg-black h-2 rounded" style="width: 50%;"></div>
      </div>
      <div class="flex gap-2 mt-2">
        <span class="bg-gray-200 text-sm px-2 py-1 rounded">Databases</span>
        <span class="bg-gray-200 text-sm px-2 py-1 rounded">Backend</span>
      </div>
      <p class="bg-gray-100 text-sm p-2 rounded">
        <strong>Motivation:</strong> Need this for my new job
      </p>
    </section>

    <!-- Tasks & Lessons -->
    <section class="space-y-4">
      <h2 class="text-xl font-semibold">Tasks & Lessons</h2>
      <p class="text-sm text-gray-500">Break down your learning into manageable tasks</p>

      <!-- Task 1 - Completed -->
      <div class="bg-white p-4 rounded shadow space-y-2">
        <div class="flex items-start gap-2">
          <input type="checkbox" checked disabled class="mt-1" />
          <div class="flex flex-col flex-grow">
            <span class="line-through text-gray-500">Watch tutorial on INNER JOIN</span>
            <div class="flex gap-2 text-xs mt-1">
              <span class="bg-red-500 text-white px-2 py-0.5 rounded">high</span>
              <span class="text-gray-500">30min</span>
            </div>
          </div>
        </div>
        <div class="ml-6 text-sm">
          <a href="#" class="text-blue-600 underline">🔗 Resource Link</a>
          <div class="bg-blue-50 p-2 mt-2 rounded text-gray-700">
            💬 Great explanation, now I understand the concept<br />
            Difficulty: 3/5 ⭐ 4/5
          </div>
        </div>
      </div>

      <!-- Task 2 - Incomplete -->
      <div class="bg-white p-4 rounded shadow flex items-start gap-2">
        <input type="checkbox" class="mt-1" />
        <div class="flex flex-col flex-grow">
          <span>Practice LEFT JOIN exercises</span>
          <div class="flex gap-2 text-xs mt-1">
            <span class="bg-black text-white px-2 py-0.5 rounded">medium</span>
            <span class="text-gray-500">45min</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

