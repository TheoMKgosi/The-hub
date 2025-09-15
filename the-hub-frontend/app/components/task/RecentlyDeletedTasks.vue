<template>
  <div class="recently-deleted-tasks">
    <div class="header">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
        Recently Deleted Tasks
      </h3>
      <button
        @click="showDeleted = !showDeleted"
        class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
      >
        {{ showDeleted ? 'Hide' : 'Show' }} ({{ deletedTasks.length }})
      </button>
    </div>

    <Transition name="slide">
      <div v-if="showDeleted && deletedTasks.length > 0" class="deleted-tasks-list mt-4 space-y-2">
        <div
          v-for="task in deletedTasks"
          :key="task.task_id"
          class="deleted-task-item bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <h4 class="font-medium text-gray-900 dark:text-white line-through">
                {{ task.title }}
              </h4>
              <p v-if="task.description" class="text-sm text-gray-600 dark:text-gray-400 mt-1 line-through">
                {{ task.description }}
              </p>
              <div class="flex items-center gap-2 mt-2 text-xs text-gray-500 dark:text-gray-400">
                <span v-if="task.priority" class="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">
                  Priority {{ task.priority }}
                </span>
                <span v-if="task.due_date" class="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">
                  Due: {{ formatDate(task.due_date) }}
                </span>
                <span class="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">
                  Deleted: {{ formatDeletedTime(task.deleted_at) }}
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2 ml-4">
              <button
                @click="undoDelete(task.task_id)"
                :disabled="undoLoading === task.task_id"
                class="px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-md transition-colors"
              >
                <span v-if="undoLoading === task.task_id">Restoring...</span>
                <span v-else>Undo</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="showDeleted && deletedTasks.length === 0" class="empty-state mt-4 text-center py-8">
        <p class="text-gray-500 dark:text-gray-400">No recently deleted tasks</p>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTaskStore } from '~/stores/tasks'
import { useToast } from '~/composables/useToast'

interface DeletedTask {
  task_id: string
  title: string
  description?: string
  priority?: number
  due_date?: string
  deleted_at: string
}

const taskStore = useTaskStore()
const { addToast } = useToast()

const showDeleted = ref(false)
const deletedTasks = ref<DeletedTask[]>([])
const undoLoading = ref<string | null>(null)

const fetchDeletedTasks = async () => {
  try {
    const tasks = await taskStore.getRecentlyDeletedTasks()
    deletedTasks.value = tasks
  } catch (error) {
    console.error('Failed to fetch deleted tasks:', error)
  }
}

const undoDelete = async (taskId: string) => {
  undoLoading.value = taskId
  try {
    await taskStore.undoDeleteTask(taskId)
    // Remove from deleted tasks list
    deletedTasks.value = deletedTasks.value.filter(task => task.task_id !== taskId)
    // Refresh main tasks list
    await taskStore.fetchTasks()
  } catch (error) {
    console.error('Failed to undo delete:', error)
  } finally {
    undoLoading.value = null
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}

const formatDeletedTime = (deletedAt: string) => {
  if (!deletedAt) return ''
  const deletedDate = new Date(deletedAt)
  const now = new Date()
  const diffMs = now.getTime() - deletedDate.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return 'Today'
  } else if (diffDays === 1) {
    return 'Yesterday'
  } else if (diffDays < 7) {
    return `${diffDays} days ago`
  } else {
    return deletedDate.toLocaleDateString()
  }
}

onMounted(() => {
  fetchDeletedTasks()
})

// Refresh deleted tasks when component becomes visible
watch(showDeleted, (newValue) => {
  if (newValue) {
    fetchDeletedTasks()
  }
})
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>