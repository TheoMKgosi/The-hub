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

const removeDeck = (id: string) => {
  deckStore.deleteDeck(id)
}

const editDeck = (deck_id: string) => {
  navigateTo(`/learning/cards/${deck_id}`)
}

const reviewDeck = (deck_id: string) => {
  navigateTo(`/learning/review/${deck_id}`)
}
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <div class="max-w-6xl mx-auto px-4 py-8">
      <!-- Header Section -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-text-light dark:text-text-dark mb-2">Flashcard Decks</h1>
        <p class="text-text-light/70 dark:text-text-dark/70">Create and manage your learning decks</p>
      </div>

      <!-- Add New Deck Form -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-6 mb-8">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create New Deck
        </h2>
        <div class="flex gap-4">
          <div class="flex-grow">
            <input
              type="text"
              v-model="formData.name"
              placeholder="Enter deck name (e.g., 'Spanish Vocabulary', 'Chemistry Terms')"
              class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200"
            />
          </div>
          <UiButton @click="addDeck" variant="primary" size="md" :disabled="!formData.name.trim()" class="px-6">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Create Deck
          </UiButton>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="deckStore.decks.length === 0" class="text-center py-16">
        <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-12 max-w-md mx-auto">
          <svg class="w-16 h-16 text-text-light/30 dark:text-text-dark/30 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No decks yet</h3>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-6">Start building your knowledge by creating your first flashcard deck.</p>
          <UiButton variant="primary" size="lg" @click="$refs.deckNameInput.focus()">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Create Your First Deck
          </UiButton>
        </div>
      </div>

      <!-- Decks Grid -->
      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="deck in deckStore.decks"
          :key="deck.deck_id"
          class="group bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark hover:shadow-xl hover:border-primary/20 dark:hover:border-primary/20 transition-all duration-300 overflow-hidden"
        >
          <!-- Deck Header -->
          <div class="p-6 border-b border-surface-light dark:border-surface-dark">
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center gap-3 flex-grow min-w-0">
                <div class="w-10 h-10 bg-primary/10 dark:bg-primary/20 rounded-lg flex items-center justify-center flex-shrink-0">
                  <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
                </div>
                <div class="min-w-0 flex-grow">
                  <h3 class="text-lg font-semibold text-text-light dark:text-text-dark truncate group-hover:text-primary transition-colors cursor-pointer" @click="editDeck(deck.deck_id)">
                    {{ deck.name }}
                  </h3>
                  <p class="text-sm text-text-light/60 dark:text-text-dark/60">Flashcard deck</p>
                </div>
              </div>
              <UiButton
                @click="removeDeck(deck.deck_id)"
                variant="danger"
                size="sm"
                class="opacity-100 md:opacity-0 md:group-hover:opacity-100 transition-opacity duration-200 flex-shrink-0"
                title="Delete Deck"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </UiButton>
            </div>
          </div>

          <!-- Deck Actions -->
          <div class="p-6">
            <div class="flex gap-3">
              <UiButton
                @click="editDeck(deck.deck_id)"
                variant="default"
                size="sm"
                class="flex-grow justify-center"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Edit Cards
              </UiButton>
              <UiButton
                @click="reviewDeck(deck.deck_id)"
                variant="primary"
                size="sm"
                class="flex-grow justify-center"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Review
              </UiButton>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
