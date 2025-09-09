<script setup lang="ts">

const route = useRoute()
const router = useRouter()
const deckID = route.params.deck_id as string
const cardStore = useCardStore()

const currentCardIndex = ref(0)
const showAnswer = ref(false)
const reviewedCards = ref(new Set())
const cardRatings = ref(new Map()) // Store ratings for each card
const selectedRating = ref<number | null>(null)

const currentCard = computed(() => cardStore.reviewCards[currentCardIndex.value])

// Rating labels for better UX
const ratingLabels = {
  0: "Complete blackout",
  1: "Incorrect, but familiar",
  2: "Incorrect, but close",
  3: "Correct with difficulty",
  4: "Correct with hesitation",
  5: "Perfect recall"
}

const flipCard = () => {
  showAnswer.value = !showAnswer.value
  selectedRating.value = null // Reset rating when flipping
}

const nextCard = () => {
  if (currentCardIndex.value < cardStore.reviewCards.length - 1) {
    currentCardIndex.value++
    showAnswer.value = false
    selectedRating.value = null
  }
}

const selectRating = (rating: number) => {
  selectedRating.value = rating
}

const submitRating = async () => {
  if (selectedRating.value !== null && currentCard.value) {
    // Store the rating
    cardRatings.value.set(currentCard.value.card_id, selectedRating.value)
    reviewedCards.value.add(currentCard.value.card_id)

    try {
      // Submit the review to the backend
      await cardStore.reviewCard(currentCard.value.card_id, selectedRating.value)

      // Remove the card from the local review queue
      const currentIndex = cardStore.reviewCards.findIndex(card => card.card_id === currentCard.value.card_id)
      if (currentIndex !== -1) {
        cardStore.reviewCards.splice(currentIndex, 1)
      }
    } catch (error) {
      console.error('Failed to submit review:', error)
    }
  }
}

const getCurrentCardRating = computed(() => {
  if (currentCard.value) {
    return cardRatings.value.get(currentCard.value.card_id)
  }
  return undefined
})

const goBack = () => {
  router.back()
}

onMounted(() => {
  cardStore.fetchDueCards(deckID)
})
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <!-- Header -->
    <div class="bg-surface-light dark:bg-surface-dark shadow-lg border-b border-surface-light dark:border-surface-dark">
      <div class="max-w-4xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <UiButton @click="goBack" variant="default" size="sm">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Back to Decks
          </UiButton>
          <div class="flex items-center gap-6">
            <div class="text-center">
              <div class="text-2xl font-bold text-text-light dark:text-text-dark">{{ currentCardIndex + 1 }}</div>
              <div class="text-xs text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">Current</div>
            </div>
            <div class="w-px h-8 bg-surface-light dark:bg-surface-dark"></div>
            <div class="text-center">
              <div class="text-2xl font-bold text-primary">{{ cardStore.reviewCards.length }}</div>
              <div class="text-xs text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">Total</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-4xl mx-auto px-4 py-12">
      <!-- Loading State -->
      <div v-if="cardStore.loading" class="text-center py-16">
        <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-12 max-w-md mx-auto">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">Loading cards...</h3>
          <p class="text-text-light/70 dark:text-text-dark/70">Preparing your review session</p>
        </div>
      </div>

      <!-- Flashcard -->
      <div v-else-if="currentCard" class="mb-12">
        <div
          class="relative bg-surface-light dark:bg-surface-dark rounded-3xl shadow-2xl border border-surface-light dark:border-surface-dark cursor-pointer transform transition-all duration-500 hover:shadow-3xl group overflow-hidden"
          @click="flipCard"
          :class="{ 'scale-105 rotate-1': showAnswer }"
        >
          <!-- Card Background Pattern -->
          <div class="absolute inset-0 opacity-5">
            <svg class="w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
              <defs>
                <pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse">
                  <path d="M 10 0 L 0 0 0 10" fill="none" stroke="currentColor" stroke-width="0.5"/>
                </pattern>
              </defs>
              <rect width="100" height="100" fill="url(#grid)" />
            </svg>
          </div>

          <div class="relative aspect-[4/3] flex items-center justify-center p-12" v-if="currentCard">
            <div class="text-center max-w-2xl">
              <!-- Question Side -->
              <div v-if="!showAnswer" class="space-y-6">
                <div class="inline-flex items-center gap-2 px-4 py-2 bg-primary/10 dark:bg-primary/20 text-primary rounded-full text-sm font-medium uppercase tracking-wide">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Question
                </div>
                <p class="text-2xl md:text-3xl font-medium text-text-light dark:text-text-dark leading-relaxed">
                  {{ currentCard.question }}
                </p>
                <div class="flex items-center justify-center gap-2 text-text-light/60 dark:text-text-dark/60 mt-8">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                  <span class="text-sm font-medium">Click to flip</span>
                </div>
              </div>

              <!-- Answer Side -->
              <div v-else class="space-y-6">
                <div class="inline-flex items-center gap-2 px-4 py-2 bg-accent/10 dark:bg-accent/20 text-accent rounded-full text-sm font-medium uppercase tracking-wide">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Answer
                </div>
                <p class="text-xl md:text-2xl text-text-light dark:text-text-dark leading-relaxed">
                  {{ currentCard.answer }}
                </p>
              </div>
            </div>
          </div>

          <!-- Previous rating indicator -->
          <div v-if="getCurrentCardRating !== undefined" class="absolute top-6 left-6">
            <div class="px-3 py-2 bg-primary/10 dark:bg-primary/20 text-primary rounded-xl text-xs font-medium backdrop-blur-sm">
              Previous: {{ getCurrentCardRating }}/5
            </div>
          </div>

          <!-- Flip indicator -->
          <div class="absolute top-6 right-6">
            <div class="w-10 h-10 bg-background-light dark:bg-background-dark rounded-full flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-200">
              <svg class="w-5 h-5 text-text-light dark:text-text-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15">
                </path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Rating System (shown when answer is revealed) -->
      <div v-if="showAnswer && currentCard" class="mb-12">
        <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-8">
          <div class="text-center mb-8">
            <h3 class="text-2xl font-semibold text-text-light dark:text-text-dark mb-2">
              How well did you remember this?
            </h3>
            <p class="text-text-light/70 dark:text-text-dark/70">Rate your recall to improve future reviews</p>
          </div>

          <!-- Rating buttons -->
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-8">
            <button
              v-for="rating in [1, 2, 3, 4, 5]"
              :key="rating"
              @click="selectRating(rating)"
              :class="[
                'group p-4 rounded-xl border-2 transition-all duration-300 text-center hover:scale-105 transform',
                selectedRating === rating
                  ? 'border-primary bg-primary/10 dark:bg-primary/20 text-primary shadow-lg'
                  : 'border-surface-light dark:border-surface-dark hover:border-primary/50 hover:bg-background-light dark:hover:bg-background-dark'
              ]"
            >
              <div class="text-2xl font-bold mb-2 group-hover:scale-110 transition-transform duration-200">{{ rating }}</div>
              <div class="text-xs text-text-light/70 dark:text-text-dark/70 leading-tight font-medium">
                {{ ratingLabels[rating] }}
              </div>
            </button>
          </div>

          <!-- Selected rating description -->
          <div v-if="selectedRating !== null" class="text-center mb-6">
            <div class="inline-flex items-center gap-2 px-6 py-3 bg-primary/10 dark:bg-primary/20 text-primary rounded-xl">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="font-medium">Rating {{ selectedRating }}/5: {{ ratingLabels[selectedRating] }}</span>
            </div>
          </div>

          <!-- Submit rating button -->
          <div class="text-center">
            <UiButton
              @click="submitRating"
              :disabled="selectedRating === null"
              variant="primary"
              size="lg"
              class="px-8 py-3"
            >
              <svg v-if="selectedRating !== null" class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              {{ selectedRating !== null ? 'Submit Rating & Continue' : 'Select a rating first' }}
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Controls -->
      <div v-if="!showAnswer" class="flex justify-center mb-12">
        <div class="flex gap-4">
          <UiButton @click="flipCard" variant="secondary" size="md" class="px-6">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            {{ showAnswer ? 'Show Question' : 'Reveal Answer' }}
          </UiButton>
        </div>
      </div>

      <!-- Review Statistics -->
      <div v-if="reviewedCards.size > 0" class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-8">
        <div class="text-center mb-8">
          <h3 class="text-2xl font-semibold text-text-light dark:text-text-dark mb-2">Review Progress</h3>
          <p class="text-text-light/70 dark:text-text-dark/70">Your performance this session</p>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div class="text-center p-4 bg-background-light dark:bg-background-dark rounded-xl">
            <div class="text-3xl font-bold text-primary mb-1">{{ reviewedCards.size }}</div>
            <div class="text-sm text-text-light/70 dark:text-text-dark/70 font-medium">Reviewed</div>
          </div>
          <div class="text-center p-4 bg-background-light dark:bg-background-dark rounded-xl">
            <div class="text-3xl font-bold text-text-light dark:text-text-dark mb-1">{{ cardStore.reviewCards.length - reviewedCards.size }}</div>
            <div class="text-sm text-text-light/70 dark:text-text-dark/70 font-medium">Remaining</div>
          </div>
          <div class="text-center p-4 bg-success/10 dark:bg-success/20 rounded-xl">
            <div class="text-3xl font-bold text-success mb-1">
              {{Array.from(cardRatings.values()).filter(r => r >= 4).length}}
            </div>
            <div class="text-sm text-success/70 font-medium">Well Known</div>
          </div>
          <div class="text-center p-4 bg-warning/10 dark:bg-warning/20 rounded-xl">
            <div class="text-3xl font-bold text-warning mb-1">
              {{Array.from(cardRatings.values()).filter(r => r <= 2).length}}
            </div>
            <div class="text-sm text-warning/70 font-medium">Need Review</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
