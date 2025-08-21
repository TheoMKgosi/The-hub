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
    <h1 class="text-3xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Your Decks</h1>

    <div class="flex mb-6 gap-3">
      <input type="text" v-model="formData.name" placeholder="Enter new deck name"
        class="flex-grow border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50" />
      <UiButton @click="addDeck" variant="primary" size="md" title="Add new deck">
        ➕ Add Deck
      </UiButton>
    </div>

    <div v-if="deckStore.decks.length === 0" class="text-center text-text-light dark:text-text-dark/60">
      No decks yet. Add one!
    </div>

    <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <div v-for="deck in deckStore.decks" :key="deck.deck_id"
        class="border border-surface-light dark:border-surface-dark rounded-xl shadow-md p-4 bg-surface-light dark:bg-surface-dark relative hover:shadow-lg transition-all duration-200">
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-xl font-semibold truncate cursor-pointer text-text-light dark:text-text-dark hover:text-primary dark:hover:text-primary transition-colors"
            @click="editDeck(deck.deck_id)">
            {{ deck.name }}
          </h2>
          <UiButton @click="removeDeck(deck.deck_id)" variant="danger" size="sm"
            title="Delete Deck">
            ✖
          </UiButton>
        </div>
        <div class="flex gap-3">
          <UiButton @click="editDeck(deck.deck_id)" variant="default" size="sm" class="flex-grow">
            Edit
          </UiButton>
          <UiButton @click="reviewDeck(deck.deck_id)" variant="primary" size="sm" class="flex-grow">
            Review
          </UiButton>
        </div>
      </div>
    </div>
  </div>
</template>
