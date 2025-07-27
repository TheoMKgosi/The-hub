<script setup lang="ts">
import TheGoals from '../components/TheGoals.vue'
import TheTasks from '../components/task/TheTasks.vue'
import Tabs from '../components/Tabs.vue'

</script>

<template>
  <div class="p-5">
    <h1 class="font-bold text-center text-6xl">Planning</h1>
  </div>

  <Tabs :tabs="['Tasks', 'Goals']" defaultTab="Tasks">
    <template #Tasks>
      <TheTasks></TheTasks>
    </template>
    <template #Goals>
      <TheGoals></TheGoals>
    </template>

  </Tabs>
  <!--
  <nav class="p-4 text-center text-2xl rounded-2xl bg-white/20 backdrop-blur-md border-white/10
    shadow-lg max-w-2xs mx-auto">
    <button class="p-4 hover:text-blue-900 rounded-2xl " :class="activeSection('tasks') ? 'bg-gray-300' : ''"
      @click="currentSection = 'tasks'">Task</button>
    <button class="p-4 hover:text-blue-900 rounded-2xl" :class="activeSection('goals') ? 'bg-gray-300' : ''"
      @click="currentSection = 'goals'">Goal</button>
  </nav>
 -->
</template>
