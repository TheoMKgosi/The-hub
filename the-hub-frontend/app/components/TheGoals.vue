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
    <!-- Filters -->
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-white/20 backdrop-blur-md rounded-lg mt-2">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <div class="flex gap-2">
          <button @click="filter = 'all'" :class="{ 'font-bold': filter === 'all' }">All</button>
          <button @click="filter = 'pending'" :class="{ 'font-bold': filter === 'pending' }">Pending</button>
          <button @click="filter = 'complete'" :class="{ 'font-bold': filter === 'complete' }">Complete</button>
        </div>
        <input v-model="searchQuery" placeholder="Search tasks..."
          class="flex-grow shadow-sm  bg-gradient-to-r from-gray-50 to-gray-100 px-3 py-2 rounded w-full sm:w-auto" />
      </div>
    </div>

    <p v-if="goalStore.loading">Loading...</p>
    <template v-else>
      <p v-if="goalStore.goals.length === 0">No goals added</p>
      <ul>
        <li v-for="goal in goalStore.goals" :key="goal.goal_id" class="bg-green-200 p-4 mb-4 ml-4 max-w-xl">
          <h3>{{ goal.title }}</h3>
          <p>{{ goal.description }}</p>
        </li>
      </ul>
    </template>
  </div>
</template>
