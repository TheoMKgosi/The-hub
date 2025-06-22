<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending'
})

const submitForm = async () => {
  const dataToSend = { ...formData }
  if (dataToSend.due_date) {
    const date = new Date(dataToSend.due_date)
    dataToSend.due_date = date.toISOString()
  }

  taskStore.submitForm(dataToSend)
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
  taskStore.fetchTasks()
})
</script>

<template>
  <form @submit.prevent="submitForm" class="m-4 max-w-md mx-auto bg-white shadow-md rounded-lg p-6 space-y-4">
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


  <div class="px-6 ">
    <h2 class="text-xl font-bold mb-4">Tasks</h2>
    <p v-if="taskStore.loading">Loading...</p>
    <ul v-else>
      <p v-if="taskStore.tasks.length === 0">No tasks added</p>
      <li v-for="task in taskStore.tasks" :key="task.task_id" class="bg-white shadow rounded-lg p-4 mb-4 border-l-4"
        :class="task.status === 'complete' ? 'border-green-500' : 'border-yellow-500'">
        <div class="flex justify-between items-start">
          <div>
            <h3 class="text-lg font-semibold">{{ task.title }}</h3>
            <p class="text-sm text-gray-600 mb-1">{{ task.description }}</p>
            <p class="text-sm text-gray-500">{{ task.due_date ? new Date(task.due_date).toLocaleString() : '' }}
            </p>
            <div class="flex items-center gap-2 mt-2">
              <input type="checkbox" @click.prevent="completeTask(task)" :checked="task.status === 'complete'"
                class="accent-green-600" />
              <span class="text-sm font-medium">{{ task.status }}</span>
            </div>
            <p class="text-sm text-gray-500 mt-1">Priority: {{ task.priority }}</p>
          </div>
          <button @click.prevent="deleteTask(task.task_id)" class="text-red-500 hover:text-red-700 transition"
            title="Delete task">
            ❌
          </button>
        </div>
      </li>

    </ul>
  </div>
</template>
