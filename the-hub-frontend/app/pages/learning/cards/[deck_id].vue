<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const deckID = parseInt(route.params.deck_id as string)

const formData = reactive({
  deck_id: deckID,
  question: '',
  answer: ''
})

const cardStore = useCardStore()

const submitForm = () => {
  cardStore.submitForm({ ...formData })
  formData.question = ''
  formData.answer = ''
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  cardStore.fetchCards(deckID)
})
</script>

<template>
  <section class="max-w-md mx-auto p-6 bg-white rounded-md shadow-md">
    <button @click="goBack" class="text-blue-600 hover:underline mb-4 flex items-center">
      ← Go Back
    </button>

    <form @submit.prevent="submitForm" class="space-y-4">
      <div>
        <label for="front" class="block font-semibold mb-1">Question</label>
        <input
          type="text"
          id="front"
          v-model="formData.question"
          class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
          placeholder="Enter your question"
          required
        />
      </div>

      <div>
        <label for="back" class="block font-semibold mb-1">Answer</label>
        <input
          type="text"
          id="back"
          v-model="formData.answer"
          class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
          placeholder="Enter the answer"
          required
        />
      </div>

      <button
        type="submit"
        class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
      >
        Add Card
      </button>
    </form>

    <div class="mt-6">
      <h2 class="text-lg font-semibold mb-2">Cards</h2>
      <ul class="space-y-3 max-h-60 overflow-y-auto">
        <li v-for="card in cardStore.cards" :key="card.card_id" class="p-3 border rounded-md bg-gray-50">
          <p class="font-semibold text-gray-800">{{ card.question }}</p>
          <p class="text-gray-600">{{ card.answer }}</p>
        </li>
      </ul>
    </div>
  </section>
</template>

