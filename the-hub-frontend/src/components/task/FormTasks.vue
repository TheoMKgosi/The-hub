<script setup lang="ts">
import { reactive, ref } from "vue";
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
})

const taskForm = ref(null)

const submitForm = async () => {
  const dataToSend = { ...formData }
  if (dataToSend.due_date) {
    const date = new Date(dataToSend.due_date)
    dataToSend.due_date = date.toISOString()
  }

  taskStore.submitForm(dataToSend)
  Object.assign(formData, {
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
})
}

</script>

<template>
  <div>
    <form @submit.prevent="submitForm" ref="taskForm" class="m-4 max-w-md mx-auto bg-white shadow-md rounded-lg p-6 space-y-4">
      <h2 class="text-2xl font-semibold text-center mb-4">Create a Task</h2>

      <div class="flex flex-col">
        <label class="mb-1 font-medium">Title</label>
        <input type="text" v-model="formData.title" name="title"
          class="p-2 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none" placeholder="Task title" />
      </div>

      <div class="flex flex-col">
        <label class="mb-1 font-medium">Description</label>
        <input type="text" v-model="formData.description" name="description"
          class="p-2 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none"
          placeholder="Optional description" />
      </div>

      <div class="flex flex-col">
        <label class="mb-1 font-medium">Due Date</label>
        <input type="datetime-local" v-model="formData.due_date" name="due_date"
          class="p-2 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none" />
      </div>

      <div class="flex flex-col">
        <label class="mb-1 font-medium">Priority</label>
        <input type="number" max="5" min="1" v-model="formData.priority" name="priority"
          class="p-2 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none" />
      </div>

      <button type="submit"
        class="w-full bg-blue-600 text-white font-semibold py-2 rounded-md hover:bg-blue-700 transition">
        Create Task
      </button>
    </form>

  </div>
</template>
