<script setup>
import { ref } from 'vue'

let nextId = 3 // mock id generator
const decks = ref([
  { deck_id: 1, name: 'Math Basics', cards: [1, 2, 3], editing: false },
  { deck_id: 2, name: 'Vue Tricks', cards: [1, 2], editing: false },
])

const addDeck = () => {
  decks.value.push({
    deck_id: nextId++,
    name: 'New Deck',
    cards: [],
    editing: true,
  })
}

const removeDeck = (id) => {
  decks.value = decks.value.filter((d) => d.deck_id !== id)
}

const editDeck = (deck) => {
  deck.editing = true
}
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Your Decks</h1>

    <button @click="addDeck" class="mb-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
      ➕ Add Deck
    </button>

    <div v-if="decks.length === 0">No decks yet. Add one!</div>
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3" v-else>
      <div
        v-for="deck in decks"
        :key="deck.deck_id"
        class="border p-4 rounded-xl shadow relative group bg-white"
      >
        <div class="flex justify-between items-center">
          <div @click="editDeck(deck)" class="cursor-pointer w-full">
            <template v-if="deck.editing">
              <input
                v-model="deck.name"
                @blur="deck.editing = false"
                @keyup.enter="deck.editing = false"
                class="w-full border rounded px-2 py-1"
                autofocus
              />
            </template>
            <template v-else>
              <h2 class="text-lg font-semibold">{{ deck.name }}</h2>
              <p class="text-sm text-gray-600">{{ deck.cards.length }} cards</p>
            </template>
          </div>
          <button
            @click="removeDeck(deck.deck_id)"
            class="text-red-600 hover:text-red-800 font-bold ml-2"
            title="Delete Deck"
          >
            ✖
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

