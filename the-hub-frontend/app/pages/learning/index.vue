<script setup lang="ts">
const selectedTab = useState<string>('learningTab', () => 'Notes')

//'Paths', 'Explore', 'Analytics'
</script>
<template>
  <div id="learning">
    <Tabs :tabs="['Notes','Flashcards']" v-model="selectedTab">
      <template #Topics>
        <Topics />
      </template>
      <template #Flashcards>
        <Flashcards />
      </template>
      <template #Notes>
        <Notes />
      </template>
      <template #Paths>
        <NuxtPage name="paths" />
      </template>
      <template #Explore>
        <Explore />
      </template>
      <template #Analytics>
        <NuxtPage name="analytics" />
      </template>
    </Tabs>
  </div>
</template>
