<script setup lang="ts">
const selectedTab = useState<string>('learningTab',() => 'Topics')
</script>
<template>
  <div id="learning">
      <Tabs :tabs="['Topics', 'Flashcards']" v-model="selectedTab">
      <template #Topics>
        <LearningTopics></LearningTopics>
      </template>
      <template #Flashcards>
        <LearningFlashcards></LearningFlashcards>
      </template>
    </Tabs>
  </div>
</template>
