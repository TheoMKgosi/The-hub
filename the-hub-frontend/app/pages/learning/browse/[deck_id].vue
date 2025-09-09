<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const deckID = route.params.deck_id as string
const cardStore = useCardStore()
const { renderMarkdown } = useMarkdown()

// Track which cards are flipped
const flippedCards = ref(new Set<string>())

const goBack = () => {
  router.back()
}

const flipCard = (cardId: string) => {
  if (flippedCards.value.has(cardId)) {
    flippedCards.value.delete(cardId)
  } else {
    flippedCards.value.add(cardId)
  }
}

const isCardFlipped = (cardId: string) => {
  return flippedCards.value.has(cardId)
}

const flipAllCards = () => {
  if (flippedCards.value.size === cardStore.cards.length) {
    // All cards are flipped, flip them all back
    flippedCards.value.clear()
  } else {
    // Flip all cards to show answers
    cardStore.cards.forEach(card => {
      flippedCards.value.add(card.card_id)
    })
  }
}

onMounted(() => {
  cardStore.fetchCards(deckID)
})
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <!-- Header -->
    <div class="bg-surface-light dark:bg-surface-dark shadow-lg border-b border-surface-light dark:border-surface-dark">
      <div class="max-w-6xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <UiButton @click="goBack" variant="default" size="sm">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Back to Decks
          </UiButton>
          <div class="flex items-center gap-6">
            <div class="text-center">
              <div class="text-2xl font-bold text-primary">{{ flippedCards.size }}</div>
              <div class="text-xs text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">Flipped</div>
            </div>
            <div class="w-px h-8 bg-surface-light dark:bg-surface-dark"></div>
            <div class="text-center">
              <div class="text-2xl font-bold text-text-light dark:text-text-dark">{{ cardStore.cards.length }}</div>
              <div class="text-xs text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">Total</div>
            </div>
            <UiButton @click="flipAllCards" variant="secondary" size="sm" class="ml-4">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {{ flippedCards.size === cardStore.cards.length ? 'Flip All Back' : 'Flip All' }}
            </UiButton>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-6xl mx-auto px-4 py-12">
      <!-- Loading State -->
      <div v-if="cardStore.loading" class="text-center py-16">
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-12 max-w-md mx-auto">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">Loading cards...</h3>
          <p class="text-text-light/70 dark:text-text-dark/70">Preparing your card browser</p>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="cardStore.cards.length === 0" class="text-center py-16">
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-12 max-w-md mx-auto">
          <svg class="w-16 h-16 text-text-light/30 dark:text-text-dark/30 mx-auto mb-4" fill="none"
            stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No cards in this deck</h3>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-6">Add some cards first to start browsing.</p>
          <UiButton variant="primary" size="lg" @click="router.push(`/learning/cards/${deckID}`)">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Add Cards
          </UiButton>
        </div>
      </div>

      <!-- Cards Grid -->
      <div v-else class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
        <div v-for="(card, index) in cardStore.cards" :key="card.card_id"
          class="bg-surface-light dark:bg-surface-dark rounded-3xl shadow-2xl border border-surface-light dark:border-surface-dark overflow-hidden group hover:shadow-3xl transition-all duration-300 cursor-pointer"
          :class="{ 'scale-105 rotate-1': isCardFlipped(card.card_id) }" @click="flipCard(card.card_id)">

          <div class="relative aspect-[4/3] p-8 flex flex-col">
            <!-- Card Number -->
            <div class="absolute top-6 left-6">
              <div class="px-3 py-2 bg-primary/10 dark:bg-primary/20 text-primary rounded-xl text-sm font-medium">
                {{ index + 1 }}
              </div>
            </div>

            <!-- Flip Indicator -->
            <div class="absolute top-6 right-6">
              <div
                class="w-10 h-10 bg-background-light dark:bg-background-dark rounded-full flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-200">
                <svg class="w-5 h-5 text-text-light dark:text-text-dark" fill="none" stroke="currentColor"
                  viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15">
                  </path>
                </svg>
              </div>
            </div>

            <!-- Question Side (shown when not flipped) -->
            <div v-if="!isCardFlipped(card.card_id)" class="flex-1 flex flex-col justify-center space-y-6">
              <div
                class="inline-flex items-center gap-2 px-4 py-2 bg-primary/10 dark:bg-primary/20 text-primary rounded-full text-sm font-medium uppercase tracking-wide self-center">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Question
              </div>
              <div
                class="text-xl md:text-2xl font-medium text-text-light dark:text-text-dark leading-relaxed flex-1 flex items-center">
                <div class="prose prose-sm max-w-none dark:prose-invert" v-html="renderMarkdown(card.question)"></div>
              </div>
              <div class="flex items-center justify-center gap-2 text-text-light/60 dark:text-text-dark/60 mt-8">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                <span class="text-sm font-medium">Click to flip</span>
              </div>
            </div>

            <!-- Answer Side (shown when flipped) -->
            <div v-else class="flex-1 flex flex-col justify-center space-y-6">
              <div
                class="inline-flex items-center gap-2 px-4 py-2 bg-accent/10 dark:bg-accent/20 text-accent rounded-full text-sm font-medium uppercase tracking-wide self-center">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Answer
              </div>
              <div
                class="text-lg md:text-xl text-text-light dark:text-text-dark leading-relaxed flex-1 flex items-center">
                <div class="prose prose-sm max-w-none dark:prose-invert" v-html="renderMarkdown(card.answer)"></div>
              </div>
              <div class="flex items-center justify-center gap-2 text-text-light/60 dark:text-text-dark/60 mt-8">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                <span class="text-sm font-medium">Click to flip back</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
