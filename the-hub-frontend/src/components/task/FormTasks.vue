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
const showForm = ref(true)

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
  <div v-if="showForm" @click="showForm = false">
    <div class="bg-white shadow rounded-lg p-4 mb-4 border-l-4">
      <p class="text-center">Add Task</p>
      <svg fill="#000000" height="50px" width="50px" class="mx-auto" version="1.1" id="Layer_1"
        xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 210.414 210.414"
        xml:space="preserve">
        <g>
          <g>
            <g>
              <path
                d="M105.207,0C47.196,0,0,47.196,0,105.207c0,58.011,47.196,105.207,105.207,105.207     c58.011,0,105.207-47.196,105.207-105.207C210.414,47.196,163.218,0,105.207,0z M105.207,202.621     c-53.715,0-97.414-43.699-97.414-97.414c0-53.715,43.699-97.414,97.414-97.414c53.715,0,97.414,43.699,97.414,97.414     C202.621,158.922,158.922,202.621,105.207,202.621z" />
              <path
                d="M155.862,101.31h-46.759V54.552c0-2.152-1.745-3.897-3.897-3.897s-3.897,1.745-3.897,3.897v46.759H54.552     c-2.152,0-3.897,1.745-3.897,3.897c0,2.152,1.745,3.897,3.897,3.897h46.759v46.759c0,2.152,1.745,3.897,3.897,3.897     s3.897-1.745,3.897-3.897v-46.759h46.759c2.152,0,3.897-1.745,3.897-3.897C159.759,103.055,158.014,101.31,155.862,101.31z" />
            </g>
          </g>
        </g>
      </svg>
    </div>
  </div>

  <!-- Modal Overlay -->
  <div v-else class="fixed inset-0 bg-black/50  flex items-center justify-center p-4 z-50" @click="showForm = true">

    <!-- Modal Content -->
    <div class="bg-white rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl" @click.stop>

      <!-- Modal Header -->
      <div class="flex items-center justify-between p-4 border-b">
        <h2 class="text-xl font-semibold">Create a Task</h2>
        <button @click="showForm = true" class="text-gray-500 hover:text-gray-700 text-2xl leading-none">
          ×
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-4">
        <form @submit.prevent="submitForm" ref="taskForm" class="space-y-4">

          <div class="flex flex-col">
            <label class="mb-1 font-medium text-sm">Title</label>
            <input type="text" v-model="formData.title" name="title"
              class="p-3 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none text-base"
              placeholder="Task title" required />
          </div>

          <div class="flex flex-col">
            <label class="mb-1 font-medium text-sm">Description</label>
            <textarea v-model="formData.description" name="description" rows="3"
              class="p-3 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none text-base resize-none"
              placeholder="Optional description"></textarea>
          </div>

          <div class="flex flex-col">
            <label class="mb-1 font-medium text-sm">Due Date</label>
            <input type="datetime-local" v-model="formData.due_date" name="due_date"
              class="p-3 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none text-base" />
          </div>

          <div class="flex flex-col">
            <label class="mb-1 font-medium text-sm">Priority</label>
            <select v-model="formData.priority" name="priority"
              class="p-3 border rounded-md focus:ring-2 focus:ring-blue-400 focus:outline-none text-base">
              <option value="1">1 - Low</option>
              <option value="2">2</option>
              <option value="3">3 - Medium</option>
              <option value="4">4</option>
              <option value="5">5 - High</option>
            </select>
          </div>

          <!-- Modal Footer -->
          <div class="flex flex-col-reverse sm:flex-row gap-2 pt-4">
            <button type="button" @click="showForm = true"
              class="w-full sm:w-auto px-4 py-3 border border-gray-300 text-gray-700 font-medium rounded-md hover:bg-gray-50 transition">
              Cancel
            </button>
            <button type="submit"
              class="w-full sm:w-auto px-4 py-3 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 transition">
              Create Task
            </button>
          </div>

        </form>
      </div>
    </div>
  </div>
</template>
