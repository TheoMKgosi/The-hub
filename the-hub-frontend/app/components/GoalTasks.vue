<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useGoalStore } from '@/stores/goals'

interface Props {
  goalId: string
}

const props = defineProps<Props>()

const goalStore = useGoalStore()
const tasks = ref([])
const loading = ref(false)

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

onMounted(() => {
  fetchTasks()
})
</script>

<template>
  <div class="mt-4">
    <div v-if="loading" class="text-sm text-text-light/60 dark:text-text-dark/60">
      Loading tasks...
    </div>

    <div v-else-if="tasks.length === 0" class="text-sm text-text-light/60 dark:text-text-dark/60">
      No tasks added yet
    </div>

    <div v-else class="space-y-2">
      <h4 class="text-sm font-medium text-text-light dark:text-text-dark mb-3">
        Tasks ({{ tasks.length }})
      </h4>

      <div
        v-for="task in tasks"
        :key="task.task_id"
        class="bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg p-3 border border-surface-light/10 dark:border-surface-dark/10"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <h5 class="text-sm font-medium text-text-light dark:text-text-dark truncate">
              {{ task.title }}
            </h5>

            <p v-if="task.description" class="text-xs text-text-light/80 dark:text-text-dark/80 mt-1 line-clamp-2">
              {{ task.description }}
            </p>

            <div class="flex items-center gap-3 mt-2 text-xs">
              <span v-if="task.priority" :class="getPriorityColor(task.priority)">
                Priority: {{ task.priority }}
              </span>
              <span :class="getStatusColor(task.status)" class="capitalize">
                {{ task.status.replace('_', ' ') }}
              </span>
            </div>

            <div v-if="task.due_date" class="text-xs text-text-light/60 dark:text-text-dark/60 mt-1">
              Due: {{ formatDate(task.due_date) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>