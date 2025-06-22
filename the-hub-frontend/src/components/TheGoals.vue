<script setup lang="ts">
import { onMounted } from 'vue'
import { useGoalStore } from '@/stores/goals'

const goalStore = useGoalStore()
console.log(goalStore.loading)
console.log(goalStore.goals)

onMounted(() => {
  goalStore.fetchGoals()
})
</script>

<template>
  <div>
    <h2 class="text-xl font-bold mb-4">Goals</h2>
    <p v-if="goalStore.loading">Loading...</p>
    <ul v-else>
      <li
        v-for="goal in goalStore.goals"
        :key="goal.goal_id"
        class="bg-green-200 p-4 mb-4 ml-4 max-w-xl"
      >
        <h3>{{ goal.title }}</h3>
        <p>{{ goal.description }}</p>
      </li>
    </ul>
  </div>
</template>
