<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const deckID = route.params.deck_id

const formData = reactive({
  deck_id: deckID,
  question: '',
  answer: ''
})


const cardStore = useCardStore()

const submitForm = () => {
  cardStore.submitForm({ ...formData })
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  cardStore.fetchCards(deckID)
})
</script>

<template>
  <section>
    <div class="p-4">
      <p @click="goBack"><--Go Back</p>
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
