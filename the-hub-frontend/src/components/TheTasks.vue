<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: '',
  priority: 3,
  status: 'pending'
})

const submitForm = async () => {
  taskStore.submitForm(formData)
}

const completeTask = async (task) => {
  // taskStore.completeTask(task)
  if (task.status == 'pending') {
    task.status = 'complete'
    taskStore.completeTask(task)
  } else {
    task.status = 'pending'
    taskStore.completeTask(task)
  }
}

const deleteTask = async (id: Number) => {
  taskStore.deleteTask(id)
}

onMounted(() => {
  taskStore.fetchGoals()
})
</script>

<template>
  <form @submit.prevent="submitForm" class="m-4">
    <fieldset class="flex flex-col border p-4 max-w-md">
      <legend class="text-center">Task entry</legend>
      <label>Title</label>
      <input type="text" v-model="formData.title" name="title" class="border rounded-xs max-w-sm">
      <label>Description</label>
      <input type="text" v-model="formData.description" name="description" class="border rounded-xs max-w-sm">
      <label>Priority</label>
      <input type="number" max="5" v-model="formData.priority" name="priority" class="border rounded-xs max-w-sm">
      <button type="submit" class="p-4 mt-3 border max-w-3xs">Create task</button>
    </fieldset>
  </form>

  <div>
    <h2 class="text-xl font-bold mb-4">Tasks</h2>
    <p v-if="taskStore.loading">Loading...</p>
    <ul v-else>
      <p v-if="taskStore.tasks.length === 0">No tasks added</p>
      <li v-for="task in taskStore.tasks" :key="task.task_id" class="bg-green-200 p-4 mb-4 ml-4 max-w-xl">
        <div class="flex justify-between">
          <div>
            <h3 class="font-bold">{{ task.title }}</h3>
            <p>{{ task.description }}</p>
            <input type="checkbox" @click.prevent="completeTask(task)" :checked="task.status === 'complete'"/><span>{{ task.status }}</span>

            <!--
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" class="sr-only peer">
          <div
            class="w-11 h-6 bg-gray-300 rounded-full peer-checked:bg-green-500 transition-colors duration-300 ring-2 ring-transparent peer-focus:ring-green-300">
            <div
              class="absolute left-0.5 top-0.5 w-5 h-5 bg-white rounded-full shadow-md transform
                peer-checked:translate-x-5 transition-transform duration-300 ring-1 ring-gray-400 peer-focus:ring-green-500">
            </div>
          </div>
          <span class="ml-3 text-sm font-medium text-gray-700">Toggle me</span>
        </label>

-->

            <p>{{ task.priority }}</p>
          </div>
          <button @click.prevent="deleteTask(task.task_id)">❌</button>
        </div>
      </li>
    </ul>
  </div>
</template>
