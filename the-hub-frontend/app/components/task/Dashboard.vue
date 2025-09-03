<script setup lang="ts">
const taskStore = useTaskStore()

callOnce(async () => {
  if(taskStore.tasks.length === 0) await taskStore.fetchTasks()
})

const completeTask = async (task) => {
  if (task.status == 'pending') {
    task.status = 'completed'
    await taskStore.completeTask(task)
  } else {
    task.status = 'pending'
    await taskStore.completeTask(task)
  }
}

// Filter tasks to only show those not linked with goals
const standaloneTasks = computed(() => {
  return taskStore.tasks.filter(task => !task.goal_id)
})
</script>

<template>
  <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
    <div class="p-4 sm:p-6 border-b border-surface-light dark:border-surface-dark">
      <div class="flex items-center justify-between">
        <h2 class="text-base sm:text-lg font-semibold text-text-light dark:text-text-dark flex items-center">
          <span class="text-success mr-2 text-sm sm:text-base">✓</span>
          <span class="hidden sm:inline">Task Management</span>
          <span class="sm:hidden">Tasks</span>
        </h2>
        <div class="flex items-center gap-2">
          <UiButton variant="primary" size="sm" class="p-2 sm:p-3">
            <span class="text-base sm:text-lg">+</span>
          </UiButton>
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
      <div v-for="task in standaloneTasks.slice(0, 5)" :key="task.task_id"
        class="bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg p-3 sm:p-4 border-l-4 hover:shadow-md transition-shadow duration-200 touch-manipulation"
        :class="task.status === 'complete' ? 'border-success' : 'border-warning'">

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
                <span class="px-1.5 py-0.5 sm:px-2 sm:py-1 text-xs font-medium rounded-full"
                  :class="task.status === 'complete'
                    ? 'bg-success/10 dark:bg-success/20 text-success dark:text-success'
                    : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
                  {{ task.status === 'complete' ? 'Done' : 'Todo' }}
                </span>
                <span class="px-1.5 py-0.5 sm:px-2 sm:py-1 text-xs font-medium bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary rounded-full">
                  P{{ task.priority }}
                </span>
              </div>
            </div>

            <!-- Description - hidden on very small screens -->
            <div v-if="task.description" class="mt-2 text-xs sm:text-sm text-text-light dark:text-text-dark/80">
              {{ task.description.length > 100 ? task.description.substring(0, 100) + '...' : task.description }}
            </div>

            <!-- Due date and time info -->
            <div class="mt-2 flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-4 text-xs text-text-light dark:text-text-dark/60">
              <div v-if="task.due_date" class="flex items-center gap-1">
                <span class="hidden sm:inline">📅</span>
                <span>Due: {{ new Date(task.due_date).toLocaleDateString() }}</span>
              </div>

              <div v-if="task.time_estimate_minutes" class="flex items-center gap-1">
                <span class="hidden sm:inline">⏱️</span>
                <span>Est: {{ Math.floor(task.time_estimate_minutes / 60) }}h {{ task.time_estimate_minutes % 60 }}m</span>
              </div>

              <div v-if="task.time_spent_minutes > 0" class="flex items-center gap-1">
                <span class="hidden sm:inline">📊</span>
                <span>Spent: {{ Math.floor(task.time_spent_minutes / 60) }}h {{ task.time_spent_minutes % 60 }}m</span>
              </div>
            </div>

            <!-- Subtasks - compact on mobile -->
            <div v-if="task.subtasks && task.subtasks.length > 0" class="mt-3">
              <div class="text-xs text-text-light dark:text-text-dark/60 mb-2">
                Subtasks ({{ task.subtasks.filter(st => st.status === 'completed').length }}/{{ task.subtasks.length }})
              </div>
              <div class="space-y-1">
                <div v-for="subtask in task.subtasks.slice(0, 2)" :key="subtask.task_id"
                  class="flex items-center text-xs sm:text-sm">
                  <input type="checkbox" @click="completeTask(subtask)" :checked="subtask.status === 'complete'"
                    class="h-3 w-3 sm:h-3 sm:w-3 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark mr-2 flex-shrink-0" />
                  <span class="text-text-light dark:text-text-dark/80 truncate"
                    :class="subtask.status === 'complete' ? 'line-through opacity-75' : ''">
                    {{ subtask.title }}
                  </span>
                </div>
                <div v-if="task.subtasks.length > 2" class="text-xs text-text-light dark:text-text-dark/60 ml-5">
                  +{{ task.subtasks.length - 2 }} more
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
          <div class="flex items-center gap-2">
            <span class="px-2 py-1 text-xs font-medium rounded-full"
              :class="task.status === 'complete'
                ? 'bg-success/10 dark:bg-success/20 text-success dark:text-success'
                : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
              {{ task.status }}
            </span>
            <span class="px-2 py-1 text-xs font-medium bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary rounded-full">
              Priority {{ task.priority }}
            </span>
          </div>
        </div>

         <div v-if="task.description" class="mt-2 text-sm text-text-light dark:text-text-dark/80 ml-8">
           {{ task.description }}
         </div>

         <div v-if="task.due_date" class="mt-2 text-xs text-text-light dark:text-text-dark/60 ml-8">
           Due: {{ new Date(task.due_date).toLocaleString() }}
         </div>

         <!-- Time Tracking -->
         <div v-if="task.time_estimate_minutes || task.time_spent_minutes > 0" class="mt-2 flex items-center gap-2 ml-8">
           <div class="text-xs text-text-light dark:text-text-dark/60">
             <span v-if="task.time_estimate_minutes">Est: {{ Math.floor(task.time_estimate_minutes / 60) }}h {{ task.time_estimate_minutes % 60 }}m</span>
             <span v-if="task.time_estimate_minutes && task.time_spent_minutes > 0" class="mx-1">•</span>
             <span v-if="task.time_spent_minutes > 0">Spent: {{ Math.floor(task.time_spent_minutes / 60) }}h {{ task.time_spent_minutes % 60 }}m</span>
           </div>
         </div>

         <!-- Subtasks -->
         <div v-if="task.subtasks && task.subtasks.length > 0" class="mt-3 ml-8">
           <div class="text-xs text-text-light dark:text-text-dark/60 mb-2">Subtasks:</div>
           <div v-for="subtask in task.subtasks.slice(0, 3)" :key="subtask.task_id"
             class="flex items-center text-sm ml-4 mb-1">
             <input type="checkbox" @click="completeTask(subtask)" :checked="subtask.status === 'complete'"
               class="h-3 w-3 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark mr-2" />
             <span class="text-text-light dark:text-text-dark/80"
               :class="subtask.status === 'complete' ? 'line-through opacity-75' : ''">
               {{ subtask.title }}
             </span>
           </div>
           <div v-if="task.subtasks.length > 3" class="text-xs text-text-light dark:text-text-dark/60 ml-4 mt-1">
             +{{ task.subtasks.length - 3 }} more subtasks
           </div>
         </div>
      </div>

      <!-- Show message if there are more than 5 standalone tasks -->
      <div v-if="standaloneTasks.length > 5" class="text-center pt-4 border-t border-surface-light/20 dark:border-surface-dark/20">
        <p class="text-sm text-text-light/70 dark:text-text-dark/70">
          Showing 5 of {{ standaloneTasks.length }} standalone tasks
        </p>
        <UiButton variant="default" size="sm" class="mt-2" @click="navigateTo('/plan')">
          View All Tasks
        </UiButton>
      </div>
    </div>
  </div>
</template>
