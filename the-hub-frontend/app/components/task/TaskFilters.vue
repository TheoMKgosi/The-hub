<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import { useGoalStore } from '@/stores/goals'

const taskStore = useTaskStore()
const goalStore = useGoalStore()

const showFilters = ref(false)

const filters = reactive({
  status: '',
  priority: '',
  goal_id: '',
  search: '',
  due_before: '',
  due_after: '',
  order_by: 'order_index',
  sort: 'asc'
})

const applyFilters = async () => {
  // Remove empty filters
  const activeFilters = Object.fromEntries(
    Object.entries(filters).filter(([_, value]) => value !== '' && value !== null)
  )

  await taskStore.fetchTasks(activeFilters)
  showFilters.value = false
}

const clearFilters = async () => {
  Object.assign(filters, {
    status: '',
    priority: '',
    goal_id: '',
    search: '',
    due_before: '',
    due_after: '',
    order_by: 'order_index',
    sort: 'asc'
  })
  await taskStore.fetchTasks()
  showFilters.value = false
}

const hasActiveFilters = computed(() => {
  return Object.values(filters).some(value => value !== '' && value !== null && value !== 'order_index' && value !== 'asc')
})
</script>

<template>
  <div class="mb-4">
    <!-- Filter Toggle Button -->
    <div class="flex items-center justify-between mb-2">
      <BaseButton @click="showFilters = !showFilters" variant="default" size="sm" class="flex items-center gap-2">
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd"
            d="M3 3a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
            clip-rule="evenodd" />
        </svg>
        Filters
        <span v-if="hasActiveFilters" class="bg-primary text-white text-xs px-1 rounded">
          {{Object.values(filters).filter(v => v && v !== 'order_index' && v !== 'asc').length}}
        </span>
      </BaseButton>

      <div class="flex items-center gap-2">
        <select v-model="filters.order_by" @change="applyFilters"
          class="px-3 py-1 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md text-sm">
          <option value="order_index">Order</option>
          <option value="priority">Priority</option>
          <option value="due_date">Due Date</option>
          <option value="created_at">Created</option>
          <option value="title">Title</option>
          <option value="status">Status</option>
        </select>

        <select v-model="filters.sort" @change="applyFilters"
          class="px-3 py-1 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md text-sm">
          <option value="asc">↑ Asc</option>
          <option value="desc">↓ Desc</option>
        </select>
      </div>
    </div>

    <!-- Filter Panel -->
    <div v-if="showFilters"
      class="bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 rounded-lg p-3 sm:p-4 mb-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4 mb-4">
        <!-- Search -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Search</label>
          <input v-model="filters.search" type="text" placeholder="Search tasks..."
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
        </div>

        <!-- Status -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Status</label>
          <select v-model="filters.status"
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary">
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="completed">Completed</option>
          </select>
        </div>

        <!-- Priority -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Priority</label>
          <select v-model="filters.priority"
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary">
            <option value="">All Priorities</option>
            <option value="1">1 - Low</option>
            <option value="2">2</option>
            <option value="3">3 - Medium</option>
            <option value="4">4</option>
            <option value="5">5 - High</option>
          </select>
        </div>

        <!-- Goal -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Goal</label>
          <select v-model="filters.goal_id"
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary">
            <option value="">All Goals</option>
            <option v-for="goal in goalStore.goals" :key="goal.goal_id" :value="goal.goal_id">
              {{ goal.title }}
            </option>
          </select>
        </div>

        <!-- Due Before -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due Before</label>
          <input v-model="filters.due_before" type="date"
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
        </div>

        <!-- Due After -->
        <div class="flex flex-col">
          <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due After</label>
          <input v-model="filters.due_after" type="date"
            class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
        </div>
      </div>

      <!-- Filter Actions -->
      <div class="flex justify-end gap-2">
        <BaseButton @click="clearFilters" variant="default" size="sm" text="Clear All">
        </BaseButton>
        <BaseButton @click="applyFilters" variant="primary" size="sm" text="Apply Filters">
        </BaseButton>
      </div>
    </div>
  </div>
</template>
