<script setup lang="ts">
const deckStore = useDeckStore()

onMounted(() => {
  if (deckStore.decks.length === 0) {
    deckStore.fetchDecks()
  }
})

const formData = reactive({
  name: ''
})

const addDeck = () => {
  if (!formData.name.trim()) return // no empty names pls
  deckStore.submitForm({ ...formData })
  formData.name = '' // reset input
}

const removeDeck = (id: number) => {
  deckStore.deleteDeck(id)
}

const editDeck = (deck_id: number) => {
  navigateTo(`/learning/cards/${deck_id}`)
}

const reviewDeck = (deck_id: number) => {
  navigateTo(`/learning/review/${deck_id}`)
}
</script>

<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold mb-6 text-center">Your Decks</h1>

    <div class="flex mb-6 gap-3">
      <input type="text" v-model="formData.name" placeholder="Enter new deck name"
        class="flex-grow border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
      <button @click="addDeck" class="bg-blue-600 text-white px-5 py-2 rounded hover:bg-blue-700 transition"
        title="Add new deck">
        ➕ Add Deck
      </button>
    </div>

    <div v-if="deckStore.decks.length === 0" class="text-center text-gray-500">
      No decks yet. Add one!
    </div>

    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <div v-for="deck in deckStore.decks" :key="deck.deck_id"
        class="border rounded-xl shadow p-4 bg-white relative hover:shadow-lg transition">
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-xl font-semibold truncate cursor-pointer" @click="editDeck(deck.deck_id)">
            {{ deck.name }}
          </h2>
          <button @click="removeDeck(deck.deck_id)" class="text-red-600 hover:text-red-800 font-bold"
            title="Delete Deck">
            ✖
          </button>
        </div>
        <div class="flex gap-3">
          <button @click="editDeck(deck.deck_id)"
            class="flex-grow bg-gray-200 hover:bg-gray-300 text-gray-800 py-2 rounded transition">
            Edit
          </button>
          <button @click="reviewDeck(deck.deck_id)"
            class="flex-grow bg-blue-600 hover:bg-blue-700 text-white py-2 rounded transition">
            Review
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
