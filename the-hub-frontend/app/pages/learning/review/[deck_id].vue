<script setup lang="ts">

const route = useRoute()
const router = useRouter()
const deckID = parseInt(route.params.deck_id as string)
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

const submitRating = () => {
  if (selectedRating.value !== null && currentCard.value) {
    // Store the rating
    cardRatings.value.set(currentCard.value.card_id, selectedRating.value)
    reviewedCards.value.add(currentCard.value.card_id)

    // You can emit this data to your store or API here
    cardStore.reviewCard(currentCard.value.card_id, selectedRating.value)
    cardStore.fetchDueCards(deckID)


    cardStore.fetchDueCards(deckID)
    nextCard()
  }
}

const getCurrentCardRating = computed(() => {
  if (currentCard.value) {
    return cardRatings.value.get(currentCard.value.deck_id)
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
  <section class="bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <button @click="goBack" 
            class="inline-flex items-center px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg transition-colors">
            ← Back to Decks
          </button>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-600">Deck {{ deckID }}</span>
            <span class="text-sm text-gray-600">
              {{ currentCardIndex + 1 }} / {{ cardStore.reviewCards.length }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- Flashcard -->
      <div v-if="cardStore.loading">
        <div class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p class="mt-2 text-gray-600">Loading cards...</p>
        </div>
      </div>
      <div class="mb-8" v-else>
        <div
          class="relative bg-white rounded-xl shadow-lg border cursor-pointer transform transition-all duration-300 hover:shadow-xl"
          @click="flipCard" :class="{ 'scale-105': showAnswer }">
          <div class="aspect-[3/2] flex items-center justify-center p-8" v-if="currentCard">
            <div class="text-center">
              <div v-if="!showAnswer" class="space-y-4">
                <div class="text-sm uppercase tracking-wide text-blue-600 font-semibold">
                  Question
                </div>
                <p class="text-xl font-medium text-gray-800 leading-relaxed">
                  {{ currentCard.question }}
                </p>
                <p class="text-sm text-gray-500 mt-6">
                  Click to reveal answer
                </p>
              </div>
              <div v-else class="space-y-4">
                <div class="text-sm uppercase tracking-wide text-green-600 font-semibold">
                  Answer
                </div>
                <p class="text-lg text-gray-800 leading-relaxed">
                  {{ currentCard.answer }}
                </p>
              </div>
            </div>
          </div>

          <!-- Previous rating indicator -->
          <div v-if="getCurrentCardRating !== undefined" class="absolute top-4 left-4">
            <div class="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
              Previously rated: {{ getCurrentCardRating }}/5
            </div>
          </div>

          <!-- Flip indicator -->
          <div class="absolute top-4 right-4">
            <div class="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15">
                </path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Rating System (shown when answer is revealed) -->
      <div v-if="showAnswer && currentCard" class="mb-6 p-6 bg-white rounded-xl shadow-lg border">
        <h3 class="text-lg font-semibold text-gray-800 mb-4 text-center">
          How well did you know this card?
        </h3>

        <!-- Rating buttons -->
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3 mb-4">
          <button v-for="rating in [1, 2, 3, 4, 5]" :key="rating" @click="selectRating(rating)" :class="[
            'p-3 rounded-lg border-2 transition-all duration-200 text-center',
            selectedRating === rating
              ? 'border-blue-500 bg-blue-50 text-blue-700'
              : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
          ]">
            <div class="text-xl font-bold mb-1">{{ rating }}</div>
            <div class="text-xs text-gray-600 leading-tight">
              {{ ratingLabels[rating] }}
            </div>
          </button>
        </div>

        <!-- Selected rating description -->
        <div v-if="selectedRating !== null" class="text-center mb-4 p-3 bg-blue-50 rounded-lg">
          <p class="text-blue-800 font-medium">
            Rating {{ selectedRating }}/5: {{ ratingLabels[selectedRating] }}
          </p>
        </div>

        <!-- Submit rating button -->
        <div class="text-center">
          <button @click="submitRating" :disabled="selectedRating === null"
            class="px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium">
            {{ selectedRating !== null ? 'Submit Rating & Continue' : 'Select a rating first' }}
          </button>
        </div>
      </div>

      <!-- Controls -->
      <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
        <!-- Flip button -->
        <div class="flex gap-2">
          <button @click="flipCard"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
            {{ showAnswer ? 'Show Question' : 'Show Answer' }}
          </button>
        </div>
      </div>

      <!-- Review Statistics -->
      <div v-if="reviewedCards.size > 0" class="mt-8 p-4 bg-white rounded-lg shadow">
        <h3 class="text-lg font-semibold text-gray-800 mb-3">Review Progress</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-blue-600">{{ reviewedCards.size }}</div>
            <div class="text-sm text-gray-600">Reviewed</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-600">{{ cardStore.reviewCards.length - reviewedCards.size }}</div>
            <div class="text-sm text-gray-600">Remaining</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-green-600">
              {{Array.from(cardRatings.values()).filter(r => r >= 4).length}}
            </div>
            <div class="text-sm text-gray-600">Well Known (4-5)</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-red-600">
              {{Array.from(cardRatings.values()).filter(r => r <= 2).length}} </div>
                <div class="text-sm text-gray-600">Need Review (0-2)</div>
            </div>
          </div>
        </div>
      </div>
  </section>
</template>
