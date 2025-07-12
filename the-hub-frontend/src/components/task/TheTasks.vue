<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import FormTasks from './FormTasks.vue'

const taskStore = useTaskStore()

const filter = ref<'all' | 'complete' | 'pending'>('all')
const searchQuery = ref('')
const showForm = ref(true)

const filteredTasks = computed(() => {
  let result = taskStore.tasks

  if (filter.value !== 'all') {
    result = result.filter(task => task.status === filter.value)
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(task =>
      task.title.toLowerCase().includes(query) ||
      task.description.toLowerCase().includes(query)
    )
  }

  return result
})


const editingTaskId = ref<number | null>(null)
const editFormData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
})

const startEdit = (task) => {
  editingTaskId.value = task.task_id
  Object.assign(editFormData, { ...task })
}

const cancelEdit = () => {
  editingTaskId.value = null
}

const saveEdit = async (id: number) => {
  const updated = { ...editFormData, task_id: id }
  if (updated.due_date) {
    const date = new Date(updated.due_date)
    updated.due_date = date.toISOString()
  }

  await taskStore.editTask(updated)
  editingTaskId.value = null
}


const completeTask = async (task) => {
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
  if (taskStore.tasks.length === 0) {
    taskStore.fetchTasks()
  }
})
</script>

<template>
  <div class="px-6">
    <div class="text-center">
      <h2 class="text-xl font-bold mb-4 ">Tasks</h2>
      <p>Organize your tasks and achieve your goals</p>
    </div>

    <!-- Filters -->
    <!-- Filters + Search -->
    <div class="flex flex-wrap gap-2 items-center mb-4">
      <input v-model="searchQuery" placeholder="Search tasks..."
        class="flex-grow border px-3 py-2 rounded w-full sm:w-auto" />

      <div class="flex gap-2">
        <button @click="filter = 'all'" :class="{ 'font-bold': filter === 'all' }">All</button>
        <button @click="filter = 'pending'" :class="{ 'font-bold': filter === 'pending' }">Pending</button>
        <button @click="filter = 'complete'" :class="{ 'font-bold': filter === 'complete' }">Complete</button>
      </div>
    </div>

    <!-- Add form -->
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

    <div v-else>
      <svg width="30px" height="30px" class="hover:text-red-700" @click="showForm = true" viewBox="0 0 48 48"
        fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect width="48px" height="48px" fill="white" fill-opacity="0.01" />
        <path d="M14 14L34 34" stroke="#000000" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" />
        <path d="M14 34L34 14" stroke="#000000" stroke-width="4" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
      <FormTasks />
    </div>

    <p v-if="taskStore.loading">Loading...</p>
    <ul v-else>
      <p v-if="taskStore.tasks.length === 0">No tasks added</p>
      <li v-for="task in filteredTasks" :key="task.task_id" class="bg-white shadow rounded-lg p-4 mb-4 border-l-4"
        :class="task.status === 'complete' ? 'border-green-500' : 'border-yellow-500'">

        <div @dblclick="startEdit(task)" class="flex justify-between items-start">
          <template v-if="editingTaskId !== task.task_id">
            <div>
              <h3 class="text-lg font-semibold">{{ task.title }}</h3>
              <p class="text-sm text-gray-600 mb-1">{{ task.description }}</p>
              <p class="text-sm text-gray-500">
                {{ task.due_date ? new Date(task.due_date).toLocaleString() : '' }}
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
              <!-- Trashcan svg -->
              <svg fill="#000000" height="200px" width="200px" version="1.1" id="Capa_1" class="w-8
                md:w-12
              h-8 md:h-12 hover:text-red-500" xmlns="http://www.w3.org/2000/svg"
                xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 60.167 60.167" xml:space="preserve">
                <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
                <g id="SVGRepo_iconCarrier">
                  <path
                    d="M54.5,11.667H39.88V3.91c0-2.156-1.754-3.91-3.91-3.91H24.196c-2.156,0-3.91,1.754-3.91,3.91v7.756H5.667 c-0.552,0-1,0.448-1,1s0.448,1,1,1h2.042v40.5c0,3.309,2.691,6,6,6h32.75c3.309,0,6-2.691,6-6v-40.5H54.5c0.552,0,1-0.448,1-1 S55.052,11.667,54.5,11.667z M22.286,3.91c0-1.053,0.857-1.91,1.91-1.91H35.97c1.053,0,1.91,0.857,1.91,1.91v7.756H22.286V3.91z M50.458,54.167c0,2.206-1.794,4-4,4h-32.75c-2.206,0-4-1.794-4-4v-40.5h40.75V54.167z M38.255,46.153V22.847c0-0.552,0.448-1,1-1 s1,0.448,1,1v23.306c0,0.552-0.448,1-1,1S38.255,46.706,38.255,46.153z M29.083,46.153V22.847c0-0.552,0.448-1,1-1s1,0.448,1,1 v23.306c0,0.552-0.448,1-1,1S29.083,46.706,29.083,46.153z M19.911,46.153V22.847c0-0.552,0.448-1,1-1s1,0.448,1,1v23.306 c0,0.552-0.448,1-1,1S19.911,46.706,19.911,46.153z">
                  </path>
                </g>
              </svg>

            </button>
          </template>

          <template v-else>
            <div class="flex flex-col w-full space-y-2">
              <input v-model="editFormData.title" class="border p-1 rounded" />
              <input v-model="editFormData.description" class="border p-1 rounded" />
              <input type="datetime-local" v-model="editFormData.due_date" class="border p-1 rounded" />
              <input type="number" min="1" max="5" v-model="editFormData.priority" class="border p-1 rounded" />
              <div class="flex gap-2">
                <button @click="saveEdit(task.task_id)" class="bg-green-500 text-white px-3 py-1 rounded">Save</button>
                <button @click="cancelEdit" class="bg-gray-400 text-white px-3 py-1 rounded">Cancel</button>
              </div>
            </div>
          </template>
        </div>
      </li>
    </ul>
  </div>
</template>
