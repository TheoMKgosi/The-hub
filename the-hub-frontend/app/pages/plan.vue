<script setup lang="ts">
const selectedTab = useState<string>("planTab", () => "Tasks");
</script>
<template>
  <div id="plan">
    <UiTabs :tabs="['Tasks', 'Goals']" v-model="selectedTab">
      <template #Tasks>
        <TaskTasks></TaskTasks>
      </template>
      <template #Goals>
        <TaskGoals></TaskGoals>
      </template>
    </UiTabs>
  </div>
</template>
