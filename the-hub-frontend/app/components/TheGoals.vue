<script setup lang="ts">
import { onMounted } from 'vue'
import { useGoalStore } from '@/stores/goals'
import { ref } from 'vue'

const goalStore = useGoalStore()

const searchQuery = ref('')
const filter = ref('all')

onMounted(() => {
  goalStore.fetchGoals()
})
</script>

<template>
  <div>
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg mt-2 border border-surface-light/10 dark:border-surface-dark/10">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <div class="flex gap-2">
          <UiNavLink v-for="filterOption in ['all', 'pending', 'complete']" :key="filterOption"
            :active="filter === filterOption" variant="tab" @click="filter = filterOption">
            {{ filterOption.charAt(0).toUpperCase() + filterOption.slice(1) }}
          </UiNavLink>
        </div>
        <input v-model="searchQuery" placeholder="Search goals..."
          class="flex-grow shadow-sm bg-surface-light dark:bg-surface-dark px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary w-full sm:w-auto" />
      </div>
    </div>

    <p v-if="goalStore.loading" class="text-text-light dark:text-text-dark">Loading...</p>
    <template v-else>
      <p v-if="goalStore.goals.length === 0" class="text-text-light dark:text-text-dark/60">No goals added yet</p>
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3 mt-6">
        <div v-for="goal in goalStore.goals" :key="goal.goal_id"
          class="bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark rounded-lg p-4 shadow-md hover:shadow-lg transition-shadow duration-200">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">{{ goal.title }}</h3>
          <p class="text-text-light dark:text-text-dark/80">{{ goal.description }}</p>
          <div class="mt-3 flex items-center justify-between">
            <span class="text-sm text-text-light dark:text-text-dark/60">
              Status: <span class="font-medium capitalize">{{ goal.status || 'pending' }}</span>
            </span>
            <div class="flex gap-2">
              <UiButton variant="default" size="sm">Edit</UiButton>
              <UiButton variant="danger" size="sm">Delete</UiButton>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
