<script setup lang="ts">
import { onMounted, reactive } from "vue";
import { useCardStore } from "@/stores/cards";
import { useRoute } from "vue-router";

const router = useRoute()

const deckID = parseInt(router.params.deck_id as string, 10)

const formData = reactive({
  deck_id: deckID,
  question: '',
  answer: ''
})


const cardStore = useCardStore()

const submitForm = () => {
  cardStore.submitForm({ ...formData })
}

onMounted(() => {
  cardStore.fetchCards(deckID)
})
</script>

<template>
  <section>
    <div class="p-4">
      <RouterLink to="/learning" class="p-4 bg-gray-400">Go back to decks</RouterLink>
    </div>
    <form @submit.prevent="submitForm()">
      <label for="front">Question</label>
      <input type="text" id="front" v-model="formData.question">
      <label for="back">Answer</label>
      <input type="text" id="back" v-model="formData.answer">
      <button type="submit">Add card</button>
    </form>
    <template v-for="card in cardStore.cards">
      <p>{{ card.question }}</p>
      <p>{{ card.answer }}</p>
    </template>
  </section>
</template>
