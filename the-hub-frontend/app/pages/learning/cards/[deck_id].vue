<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const deckID = route.params.deck_id as string

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
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <div class="max-w-4xl mx-auto px-4 py-8">
      <!-- Header -->
      <div class="mb-8">
        <UiButton @click="goBack" variant="default" size="sm" class="mb-4">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back to Decks
        </UiButton>
        <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-2">Manage Cards</h1>
        <p class="text-text-light/70 dark:text-text-dark/70">Add and review flashcards in this deck</p>
      </div>

      <div class="grid lg:grid-cols-2 gap-8">
        <!-- Add New Card Form -->
        <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-6">
          <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6 flex items-center gap-2">
            <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Add New Card
          </h2>

          <form @submit.prevent="submitForm" class="space-y-6">
            <div>
              <label for="question" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
                Question
              </label>
              <textarea
                id="question"
                v-model="formData.question"
                rows="3"
                class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                placeholder="Enter your question here..."
                required
              />
            </div>

            <div>
              <label for="answer" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
                Answer
              </label>
              <textarea
                id="answer"
                v-model="formData.answer"
                rows="3"
                class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                placeholder="Enter the answer here..."
                required
              />
            </div>

            <UiButton
              type="submit"
              variant="primary"
              size="md"
              class="w-full"
              :disabled="!formData.question.trim() || !formData.answer.trim()"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add Card
            </UiButton>
          </form>
        </div>

        <!-- Cards List -->
        <div class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-6">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark flex items-center gap-2">
              <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              Cards in Deck
            </h2>
            <span class="text-sm text-text-light/60 dark:text-text-dark/60 bg-background-light dark:bg-background-dark px-3 py-1 rounded-full">
              {{ cardStore.cards.length }} cards
            </span>
          </div>

          <!-- Empty State -->
          <div v-if="cardStore.cards.length === 0" class="text-center py-12">
            <svg class="w-16 h-16 text-text-light/30 dark:text-text-dark/30 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">No cards yet</h3>
            <p class="text-text-light/70 dark:text-text-dark/70">Add your first card using the form on the left.</p>
          </div>

          <!-- Cards List -->
          <div v-else class="space-y-4 max-h-96 overflow-y-auto">
            <div
              v-for="card in cardStore.cards"
              :key="card.card_id"
              class="group bg-background-light dark:bg-background-dark rounded-xl border border-surface-light dark:border-surface-dark p-4 hover:shadow-md transition-all duration-200"
            >
              <div class="space-y-3">
                <div>
                  <div class="text-xs font-medium text-primary uppercase tracking-wide mb-1">Question</div>
                  <p class="text-text-light dark:text-text-dark font-medium leading-relaxed">{{ card.question }}</p>
                </div>
                <div class="border-t border-surface-light dark:border-surface-dark pt-3">
                  <div class="text-xs font-medium text-accent uppercase tracking-wide mb-1">Answer</div>
                  <p class="text-text-light/80 dark:text-text-dark/80 leading-relaxed">{{ card.answer }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

