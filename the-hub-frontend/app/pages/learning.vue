<script setup lang="ts">
import Tabs from "@/components/Tabs.vue";
import Flashcard from "@/components/learning/flashcards.vue";
import Topics from "@/components/learning/topics.vue";
</script>

<template>
  <Tabs :tabs="['Topics', 'Flashcards']">
    <template #Topics>
      <Topics></Topics>
    </template>
    <template #Flashcards>
      <Flashcard></Flashcard>
    </template>
  </Tabs>
</template>
