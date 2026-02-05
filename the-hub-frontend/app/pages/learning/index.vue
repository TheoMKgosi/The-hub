<script setup lang="ts">
const selectedTab = useState<string>('learningTab',() => 'Topics')

//'Paths', 'Explore', 'Analytics'
</script>
<template>
  <div id="learning">
      <Tabs :tabs="['Topics', 'Flashcards']" v-model="selectedTab">
      <template #Topics>
        <Topics />
      </template>
      <template #Flashcards>
        <Flashcards />
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
