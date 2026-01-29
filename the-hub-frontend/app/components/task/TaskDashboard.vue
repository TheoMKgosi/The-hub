<script setup lang="ts">
import dayjs from "dayjs"
import relativeTime from 'dayjs/plugin/relativeTime.js'

dayjs.extend(relativeTime)

const taskStore = useTaskStore()
const selectedTaskIndex = ref(0)

callOnce(async () => {
  if (taskStore.tasks.length === 0) await taskStore.fetchTasks()
})



// Filter tasks to only show those not linked with goals
const standaloneTasks = computed(() => {
  return taskStore.tasks.filter(task => !task.goal_id)
})

const completeTask = async (task) => {
  if (task.status == "pending") {
    task.status = "complete";
    taskStore.completeTask(task);
  } else {
    task.status = "pending";
    taskStore.completeTask(task);
  }
};

// Task-specific keyboard shortcuts
const handleTaskKeyboard = (event: KeyboardEvent) => {
  // Don't trigger if typing in input fields
  if (event.target instanceof HTMLInputElement ||
    event.target instanceof HTMLTextAreaElement) {
    return
  }

  const tasks = standaloneTasks.value.slice(0, 5)

  switch (event.key) {
    case 'j':
    case 'ArrowDown':
      if (!event.ctrlKey && !event.metaKey) {
        event.preventDefault()
        selectedTaskIndex.value = Math.min(selectedTaskIndex.value + 1, tasks.length - 1)
      }
      break
    case 'k':
    case 'ArrowUp':
      if (!event.ctrlKey && !event.metaKey) {
        event.preventDefault()
        selectedTaskIndex.value = Math.max(selectedTaskIndex.value - 1, 0)
      }
      break
    case ' ':
    case 'Enter':
      if (!event.ctrlKey && !event.metaKey && tasks[selectedTaskIndex.value]) {
        event.preventDefault()
        completeTask(tasks[selectedTaskIndex.value])
      }
      break
    case 'n':
      if (event.ctrlKey) {
        event.preventDefault()
        // Trigger new task creation
        const event = new CustomEvent('command:new-task')
        window.dispatchEvent(event)
      }
      break
  }
}

// Listen for keyboard events
onMounted(() => {
  window.addEventListener('keydown', handleTaskKeyboard)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleTaskKeyboard)
})

// Listen for new task command
onMounted(() => {
  const handleNewTask = () => {
    // Navigate to plan page for task creation
    navigateTo('/plan')
  }

  window.addEventListener('command:new-task', handleNewTask)

  return () => {
    window.removeEventListener('command:new-task', handleNewTask)
  }
})

</script>

<template>
  <div
    class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
    <div class="p-4 sm:p-6 border-b border-surface-light dark:border-surface-dark">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-base sm:text-lg font-semibold text-text-light dark:text-text-dark flex items-center">
            <span class="text-success mr-2 text-sm sm:text-base">✓</span>
            <span class="hidden sm:inline">Task Management</span>
            <span class="sm:hidden">Tasks</span>
          </h2>
          <div class="text-xs text-text-light/60 dark:text-text-dark/60 mt-1 hidden sm:block">
            Use <kbd class="px-1 py-0.5 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs">J/K</kbd> or <kbd
              class="px-1 py-0.5 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs">↑/↓</kbd> to navigate,
            <kbd class="px-1 py-0.5 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs">Space</kbd> to complete
          </div>
        </div>
        <div class="flex items-center gap-2">
          <BaseButton variant="primary" size="sm" class="p-2 sm:p-3" text="+" title="New Task (Ctrl+N)" />
        </div>
      </div>
    </div>

    <!-- Task Filters -->
    <div class="p-4 border-b border-surface-light dark:border-surface-dark">
      <TaskFilters />
    </div>

    <div v-if="taskStore.loading" class="p-6 text-text-light dark:text-text-dark">Loading...</div>

    <div v-else-if="standaloneTasks.length === 0" class="p-6 text-center text-text-light dark:text-text-dark/60">
      <p class="text-lg mb-2">No standalone tasks yet</p>
      <p class="text-sm">Create your first task to get started</p>
    </div>

    <div v-else class="p-2 sm:p-4 space-y-2 sm:space-y-3">
      <div v-for="(task, index) in standaloneTasks.slice(0, 5)" :key="task.task_id"
        class="bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg p-3 sm:p-4 border-l-4 hover:shadow-md transition-shadow duration-200 touch-manipulation"
        :class="[
          task.status === 'complete' ? 'border-success' : 'border-warning',
          selectedTaskIndex === index ? 'ring-2 ring-primary/50 bg-primary/5' : ''
        ]">

        <!-- Mobile-first layout -->
        <div class="flex items-start gap-3">
          <input type="checkbox" @click="completeTask(task)" :checked="task.status === 'complete'"
            class="h-5 w-5 sm:h-5 sm:w-5 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark mt-0.5 flex-shrink-0" />

          <div class="flex-1 min-w-0">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
              <span class="text-sm sm:text-sm font-medium text-text-light dark:text-text-dark break-words"
                :class="task.status === 'complete' ? 'line-through opacity-75' : ''">
                {{ task.title }}
              </span>

              <!-- Mobile: Status and priority in a row -->
              <div class="flex items-center gap-1 sm:gap-2 flex-shrink-0">
                <span class="px-1.5 py-0.5 sm:px-2 sm:py-1 text-xs font-medium rounded-full" :class="task.status === 'complete'
                  ? 'bg-success/10 dark:bg-success/20 text-success dark:text-success'
                  : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
                  {{ task.status === 'complete' ? 'Done' : 'Todo' }}
                </span>
                <span
                  class="px-1.5 py-0.5 sm:px-2 sm:py-1 text-xs font-medium bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary rounded-full">
                  P{{ task.priority }}
                </span>
              </div>
            </div>

            <!-- Description - hidden on very small screens -->
            <div v-if="task.description" class="mt-2 text-xs sm:text-sm text-text-light dark:text-text-dark/80">
              {{ task.description.length > 100 ? task.description.substring(0, 100) + '...' : task.description }}
            </div>

            <!-- Due date and time info -->
            <div
              class="mt-2 flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-4 text-xs text-text-light dark:text-text-dark/60">
              <div v-if="task.due_date" class="flex items-center gap-1">
                <span class="hidden sm:inline">📅</span>
                <span>Due: {{ dayjs(task.due_date).fromNow() }}</span>
              </div>

              <div v-if="task.time_estimate_minutes" class="flex items-center gap-1">
                <span class="hidden sm:inline">⏱️</span>
                <span>Est: {{ Math.floor(task.time_estimate_minutes / 60) }}h {{ task.time_estimate_minutes % 60
                }}m</span>
              </div>

              <div v-if="task.time_spent_minutes > 0" class="flex items-center gap-1">
                <span class="hidden sm:inline">📊</span>
                <span>Spent: {{ Math.floor(task.time_spent_minutes / 60) }}h {{ task.time_spent_minutes % 60 }}m</span>
              </div>
            </div>


          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="px-2 py-1 text-xs font-medium rounded-full" :class="task.status === 'complete'
            ? 'bg-success/10 dark:bg-success/20 text-success dark:text-success'
            : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
            {{ task.status }}
          </span>
          <span
            class="px-2 py-1 text-xs font-medium bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary rounded-full">
            Priority {{ task.priority }}
          </span>
        </div>
      </div>
    </div>

    <!-- Show message if there are more than 5 standalone tasks -->
    <div v-if="standaloneTasks.length > 5"
      class="text-center pt-4 border-t border-surface-light/20 dark:border-surface-dark/20">
      <p class="text-sm text-text-light/70 dark:text-text-dark/70">
        Showing 5 of {{ standaloneTasks.length }} standalone tasks
      </p>
      <BaseButton variant="default" text="View All Tasks" size="sm" class="mt-2" @click="navigateTo('/plan')" />
    </div>
  </div>
</template>
