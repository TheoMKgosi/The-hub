<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useDeckStore } from '@/stores/decks'
import { useRouter } from 'vue-router'
const deckStore = useDeckStore()

const router = useRouter()

onMounted(() => {
  if (deckStore.decks.length === 0) {
    deckStore.fetchDecks()
  }
})

const formData = reactive({
  name: ''
})

const addDeck = () => {
  const dataToSend = { ...formData }
  deckStore.submitForm(dataToSend)
}

const removeDeck = (id) => {
  deckStore.deleteDeck(id)
}

const editDeck = (deck) => {
  router.push({ name: 'cards', params: { deck_id: deck.deck_id } })
}

</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Your Decks</h1>

    <div class="p-2">
      <input type="text" v-model="formData.name" class="border p-1">
      <button @click="addDeck()" class="mb-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
        ➕ Add Deck
      </button>
    </div>

    <div v-if="deckStore.decks.length === 0">No decks yet. Add one!</div>
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3" v-else>
      <div v-for="deck in deckStore.decks" :key="deck.deck_id"
        class="border p-4 rounded-xl shadow relative group bg-white">
        <div class="flex justify-between items-center">
          <div @click="editDeck(deck)" class="cursor-pointer w-full">
            <h2 class="text-lg font-semibold">{{ deck.name }}</h2>
            <!--<p class="text-sm text-gray-600">{{ deck.cards.length }} cards</p> -->
          </div>
          <button @click="removeDeck(deck.deck_id)" class="text-red-600 hover:text-red-800 font-bold ml-2"
            title="Delete Deck">
            ✖
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
