<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: ''
})

const submitForm = async () => {
  taskStore.submitForm(formData)
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

  <!-- TODO: Have tailwind css that strikes through based on status -->
  <div>
    <h2 class="text-xl font-bold mb-4">Tasks</h2>
    <p v-if="taskStore.loading">Loading...</p>
    <ul v-else>
      <p v-if="taskStore.tasks == 0">No tasks added</p>
      <li v-for="task in taskStore.tasks" :key="task.task_id" class="bg-green-200 p-4 mb-4 ml-4 max-w-xl">
        <h3>{{ task.title }}</h3>
        <p>{{ task.description }}</p>
        <p>{{ task.priority }}</p>
        <p>{{ task.status }}</p>
      </li>
    </ul>
  </div>
</template>
