<script setup lang="ts">
import dayjs from 'dayjs';
import type { Task } from '~/types/task'

interface Props {
  task_id: string,
  status: string,
  title: string,
  description?: string,
  due_date?: Date,
  priority?: number,
  time_estimate_minutes?: number
}

// const props = defineProps({
//   task_id: { type: String, required: true },
//   status: { type: String, default: 'pending' },
//   title: { type: String, default: 'Untitled Task' },
//   description: { type: String, default: '' },
//   due_date: { type: Date, default: null },
//   priority: { type: Number, default: 3 },
//   time_estimate_minutes: { type: Number, default: 0 }
// })

const { task_id, status = 'pending', title = 'Untitled Task', description = '', due_date = null, priority = 3, time_estimate_minutes = 0 } = defineProps<Props>()

const emit = defineEmits<{
  (e: 'completeTask', id: string): void;
  (e: 'deleteTask', id: string): void;
  (e: 'moveTaskUp', id: string): void;
  (e: 'moveTaskDown', id: string): void;
  (e: 'edit', id: string, updates: any): void;
}>()

const draft = reactive({
  title: title,
  description: description,
  priority: priority,
  due_date: due_date
})

const isMenuOpen = ref(false)
const isEditing = ref(false)
const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}
const startEdit = () => {
  isEditing.value = true
  isMenuOpen.value = false
}

const completeBtnClick = () => {
  emit('completeTask', task_id)
}

const deleteBtnClick = () => {
  emit('deleteTask', task_id)
}

const moveUpBtnClick = () => {
  emit('moveTaskUp', task_id)
}

const moveDownBtnClick = () => {
  emit('moveTaskDown', task_id)
}

const cancelEdit = () => {
  // Reset draft to original prop values
  draft.title = title
  draft.description = description
  draft.priority = priority
  draft.due_date = due_date
  isEditing.value = false
}

const saveEdit = () => {
  emit('edit', task_id, { ...draft })
  isEditing.value = false
}
</script>

<template>
  <div
    class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
    :class="[status === 'complete' ? 'border-success' : 'border-warning',]">
    <div class="flex justify-between">
      <div>
        <div v-if="isEditing" class="space-y-3">
          <input v-model="draft.title" class="w-full p-2 border rounded dark:bg-gray-800 dark:text-white"
            placeholder="Task Title" />
          <textarea v-model="draft.description"
            class="w-full p-2 border rounded dark:bg-gray-800 dark:text-white text-sm"
            placeholder="Description"></textarea>
          <div class="flex gap-2">
            <button @click="saveEdit" class="px-3 py-1 bg-success text-white rounded text-sm font-bold">Save</button>
            <button @click="cancelEdit" class="px-3 py-1 bg-gray-500 text-white rounded text-sm">Cancel</button>
          </div>
        </div>
        <div class="justify-between items-start" v-else>
          <div class="flex items-center gap-2 mb-2">
            <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
              {{ title }}
            </h3>
          </div>
          <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">
            {{ description }}
          </p>
          <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
            {{ due_date ? dayjs(due_date).fromNow() : "" }}
          </p>
          <div class="flex items-center gap-2 mt-2">
            <input type="checkbox" @click="completeBtnClick" :checked="status === 'complete'"
              class="accent-success w-4 h-4" />
            <span class="text-sm font-medium text-text-light dark:text-text-dark capitalize">{{ status }}</span>
          </div>
          <div>
            <p class="text-sm text-text-light dark:text-text-dark/60 mt-1">
              Priority: {{ priority }}
            </p>
          </div>
        </div>

        <div class="flex flex-col justify-between">
          <!-- Three-dot menu button -->
          <div class="ml-auto">
            <button @click="toggleMenu"
              class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors duration-200"
              title="Task options">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z">
                </path>
              </svg>
            </button>

            <!-- Dropdown menu -->
            <div v-if="isMenuOpen"
              class="absolute right-4 mt-2 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-2xl border border-surface-light/20 dark:border-surface-dark/20 z-10">
              <div class="py-1">
                <button @click="startEdit"
                  class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z">
                    </path>
                  </svg>
                  Edit
                </button>
                <button @click="moveUpBtnClick"
                  class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7">
                    </path>
                  </svg>
                  Move Up
                </button>
                <button @click="moveDownBtnClick"
                  class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7">
                    </path>
                  </svg>
                  Move Down
                </button>
                <div class="border-t border-surface-light/20 dark:border-surface-dark/20"></div>
                <button @click="deleteBtnClick"
                  class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-200">
                  <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                    </path>
                  </svg>
                  Delete
                </button>
              </div>
            </div>
          </div>
          <div v-if="time_estimate_minutes" class="flex items-center gap-1 mt-1">
            <span class="hidden sm:inline">⏱️</span>
            <span class="text-sm text-text-light dark:text-text-dark/60">
              Est: {{ Math.floor(time_estimate_minutes / 60) }}h {{ time_estimate_minutes % 60 }}m
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
