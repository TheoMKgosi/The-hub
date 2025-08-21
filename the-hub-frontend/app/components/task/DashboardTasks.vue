<script setup lang="ts">
const taskStore = useTaskStore()

callOnce(async () => { 
  if(taskStore.tasks.length === 0) await taskStore.fetchTasks() 
})

const completeTask = async (task) => {
  if (task.status == 'pending') {
    task.status = 'complete'
    taskStore.completeTask(task)
  } else {
    task.status = 'pending'
    taskStore.completeTask(task)
  }
}
</script>

<template>
  <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
    <div class="p-6 border-b border-surface-light dark:border-surface-dark">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-text-light dark:text-text-dark flex items-center">
          <span class="text-success mr-2">✓</span>
          Task Management
        </h2>
        <UiButton variant="primary" size="sm" class="p-2">
          <span class="text-lg">+</span>
        </UiButton>
      </div>
    </div>

    <div v-if="taskStore.loading" class="p-6 text-text-light dark:text-text-dark">Loading...</div>

    <div v-else-if="taskStore.tasks.length === 0" class="p-6 text-center text-text-light dark:text-text-dark/60">
      <p class="text-lg mb-2">No tasks yet</p>
      <p class="text-sm">Create your first task to get started</p>
    </div>

    <div v-else class="p-4 space-y-3">
      <div v-for="task in taskStore.tasks" :key="task.task_id"
        class="bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg p-4 border-l-4 hover:shadow-md transition-shadow duration-200"
        :class="task.status === 'complete' ? 'border-success' : 'border-warning'">

        <div class="flex items-center justify-between">
          <div class="flex items-center flex-1">
            <input type="checkbox" @click="completeTask(task)" :checked="task.status === 'complete'"
              class="h-5 w-5 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark" />
            <span class="ml-3 text-sm font-medium text-text-light dark:text-text-dark"
              :class="task.status === 'complete' ? 'line-through opacity-75' : ''">
              {{ task.title }}
            </span>
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
      </div>
    </div>
  </div>
</template>
