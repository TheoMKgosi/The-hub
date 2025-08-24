<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useGoalStore } from '@/stores/goals'

interface Task {
  task_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  order: number
}

interface Props {
  goalId: string
}

const props = defineProps<Props>()

const goalStore = useGoalStore()
const tasks = ref<Task[]>([])
const loading = ref(false)
const editingTask = ref<Task | null>(null)
const showDeleteConfirm = ref<string | null>(null)

const fetchTasks = async () => {
  loading.value = true
  try {
    const goalTasks = await goalStore.fetchGoalTasks(props.goalId)
    tasks.value = goalTasks
  } catch (err) {
    console.error('Failed to fetch goal tasks:', err)
  } finally {
    loading.value = false
  }
}

const getPriorityColor = (priority?: number) => {
  if (!priority) return 'text-text-light/60 dark:text-text-dark/60'
  switch (priority) {
    case 1:
    case 2:
      return 'text-green-600 dark:text-green-400'
    case 3:
      return 'text-yellow-600 dark:text-yellow-400'
    case 4:
    case 5:
      return 'text-red-600 dark:text-red-400'
    default:
      return 'text-text-light/60 dark:text-text-dark/60'
  }
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'completed':
    case 'complete':
      return 'text-green-600 dark:text-green-400'
    case 'in_progress':
      return 'text-blue-600 dark:text-blue-400'
    case 'pending':
      return 'text-yellow-600 dark:text-yellow-400'
    default:
      return 'text-text-light/60 dark:text-text-dark/60'
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return null
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const editTask = (task: Task) => {
  editingTask.value = task
}

const closeEditForm = () => {
  editingTask.value = null
}

const onTaskUpdated = () => {
  fetchTasks() // Refresh the tasks list
}

const toggleComplete = async (task: Task) => {
  try {
    await goalStore.completeGoalTask(props.goalId, task.task_id)
    await fetchTasks() // Refresh the tasks list
  } catch (err) {
    console.error('Failed to toggle task completion:', err)
  }
}

const confirmDelete = (taskId: string) => {
  showDeleteConfirm.value = taskId
}

const cancelDelete = () => {
  showDeleteConfirm.value = null
}

const deleteTask = async (taskId: string) => {
  try {
    await goalStore.deleteGoalTask(props.goalId, taskId)
    await fetchTasks() // Refresh the tasks list
    showDeleteConfirm.value = null
  } catch (err) {
    console.error('Failed to delete task:', err)
  }
}

onMounted(() => {
  fetchTasks()
})
</script>

<template>
  <div class="space-y-3">
    <div v-if="loading" class="text-sm text-text-light/60 dark:text-text-dark/60">
      Loading tasks...
    </div>

    <div v-else-if="tasks.length === 0" class="text-sm text-text-light/60 dark:text-text-dark/60">
      No tasks added yet
    </div>

    <div v-else class="space-y-3">
      <h4 class="text-sm font-medium text-text-light dark:text-text-dark flex items-center gap-2">
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
        </svg>
        Tasks ({{ tasks.length }})
      </h4>

      <div
        v-for="task in tasks"
        :key="task.task_id"
        class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20 hover:border-surface-light/30 dark:hover:border-surface-dark/30 transition-colors"
      >
        <!-- Edit Form -->
        <EditTaskForm
          v-if="editingTask && editingTask.task_id === task.task_id"
          :goal-id="goalId"
          :task="task"
          :goal-title="''"
          @close="closeEditForm"
          @updated="onTaskUpdated"
        />

        <!-- Delete Confirmation -->
        <div
          v-else-if="showDeleteConfirm === task.task_id"
          class="bg-red-50/50 dark:bg-red-900/20 rounded-lg p-4 border border-red-200/50 dark:border-red-800/50"
        >
          <p class="text-sm text-red-700 dark:text-red-300 mb-3">
            Are you sure you want to delete "{{ task.title }}"?
          </p>
          <div class="flex gap-2">
            <UiButton
              @click="deleteTask(task.task_id)"
              variant="destructive"
              size="sm"
            >
              Delete
            </UiButton>
            <UiButton
              @click="cancelDelete"
              variant="default"
              size="sm"
            >
              Cancel
            </UiButton>
          </div>
        </div>

        <!-- Task Display -->
        <div v-else class="flex items-start gap-3">
          <!-- Complete Toggle -->
          <button
            @click="toggleComplete(task)"
            class="flex-shrink-0 w-5 h-5 border-2 border-text-light/40 dark:border-text-dark/40 rounded-full flex items-center justify-center hover:border-primary transition-colors mt-0.5"
            :class="{
              'bg-primary border-primary': task.status === 'completed',
              'bg-surface-light/20 dark:bg-surface-dark/20': task.status !== 'completed'
            }"
          >
            <svg
              v-if="task.status === 'completed'"
              class="w-3 h-3 text-white"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>

          <!-- Task Content -->
          <div class="flex-1 min-w-0">
            <h5
              class="text-sm font-medium leading-tight mb-1"
              :class="{
                'text-text-light dark:text-text-dark': task.status !== 'completed',
                'text-text-light/60 dark:text-text-dark/60 line-through': task.status === 'completed'
              }"
            >
              {{ task.title }}
            </h5>

            <p v-if="task.description" class="text-xs text-text-light/80 dark:text-text-dark/80 mb-2 line-clamp-2">
              {{ task.description }}
            </p>

            <div class="flex items-center gap-3 text-xs">
              <span v-if="task.priority" :class="getPriorityColor(task.priority)" class="flex items-center gap-1">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
                {{ task.priority }}
              </span>
              <span :class="getStatusColor(task.status)" class="flex items-center gap-1 capitalize">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                {{ task.status.replace('_', ' ') }}
              </span>
            </div>

            <div v-if="task.due_date" class="text-xs text-text-light/60 dark:text-text-dark/60 mt-1 flex items-center gap-1">
              <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" />
              </svg>
              Due: {{ formatDate(task.due_date) }}
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center gap-1">
            <UiButton
              @click="editTask(task)"
              variant="default"
              size="sm"
              class="p-1.5 hover:bg-surface-light/30 dark:hover:bg-surface-dark/30"
              title="Edit task"
            >
              <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
                <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
              </svg>
            </UiButton>

            <UiButton
              @click="confirmDelete(task.task_id)"
              variant="destructive"
              size="sm"
              class="p-1.5 hover:bg-red-100 dark:hover:bg-red-900/30"
              title="Delete task"
            >
              <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" clip-rule="evenodd" />
                <path fill-rule="evenodd" d="M10 5a2 2 0 00-2 2v8a2 2 0 002 2h4a2 2 0 002-2V7a2 2 0 00-2-2h-4zm3 2a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1V8a1 1 0 00-1-1h-2z" clip-rule="evenodd" />
                <path fill-rule="evenodd" d="M4 7a1 1 0 011-1h10a1 1 0 110 2H5a1 1 0 01-1-1z" clip-rule="evenodd" />
              </svg>
            </UiButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>