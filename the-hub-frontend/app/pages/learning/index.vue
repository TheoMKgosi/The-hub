<script setup lang="ts">
const selectedTab = useState<string>('learningTab',() => 'Topics')

//'Paths', 'Explore', 'Analytics'
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
      <template #Paths>
        <NuxtPage name="paths" />
      </template>
      <template #Explore>
        <LearningExplore></LearningExplore>
      </template>
      <template #Analytics>
        <NuxtPage name="analytics" />
      </template>
    </Tabs>
  </div>
</template>
