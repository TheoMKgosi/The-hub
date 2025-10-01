<script setup lang="ts">
const selectedTab = useState<string>("planTab", () => "Tasks");
</script>
<template>
  <div id="plan">
    <Tabs :tabs="['Tasks', 'Goals']" v-model="selectedTab">
      <template #Tasks>
        <TaskTheTasks></TaskTheTasks>
      </template>
      <template #Goals>
        <TheGoals></TheGoals>
      </template>
    </Tabs>
  </div>
</template>
