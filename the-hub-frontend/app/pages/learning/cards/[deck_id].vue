<script setup lang="ts">
import LeftArrowIcon from '~/components/ui/svg/LeftArrowIcon.vue'
import PlusIcon from '~/components/ui/svg/PlusIcon.vue'
import EditIcon from '~/components/ui/svg/EditIcon.vue'
import DeleteIcon from '~/components/ui/svg/DeleteIcon.vue'
const route = useRoute()
const router = useRouter()

const deckID = route.params.deck_id as string

const editingCard = ref<string | null>(null)
const editFormData = reactive({
  question: '',
  answer: ''
})

const cardStore = useCardStore()
const { renderMarkdown } = useMarkdown()

const editCard = (card: any) => {
  editingCard.value = card.card_id
  editFormData.question = card.question
  editFormData.answer = card.answer
}

const saveCardEdit = (card: any) => {
  if (!editFormData.question.trim() || !editFormData.answer.trim()) return

  const updatedCard = {
    ...card,
    question: editFormData.question,
    answer: editFormData.answer
  }

  cardStore.editCard(updatedCard)
  editingCard.value = null
}

const cancelCardEdit = () => {
  editingCard.value = null
  editFormData.question = ''
  editFormData.answer = ''
}

const deleteCard = (card: any) => {
  if (confirm('Are you sure you want to delete this card?')) {
    cardStore.deleteCard(deckID, card.card_id)
  }
}

const addCard = () => {
  navigateTo(`/learning/cards/create/${deckID}`)
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
    <div class="mx-auto px-4 py-8">
      <!-- Header -->
      <div class="mb-8">
        <BaseButton @click="goBack" text="Back to Decks" variant="default" :icon="LeftArrowIcon" size="md" />
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-2">Manage Cards</h1>
            <p class="text-text-light/70 dark:text-text-dark/70">Add and review flashcards in this deck</p>
          </div>
          <BaseButton @click="addCard" text="Add Card" :icon="PlusIcon" variant="primary" size="md" />
        </div>
      </div>

      <!-- Cards List -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold text-text-light dark:text-text-dark flex items-center gap-2">
            <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            Cards in Deck
          </h2>
          <span
            class="text-sm text-text-light/60 dark:text-text-dark/60 bg-background-light dark:bg-background-dark px-3 py-1 rounded-full">
            {{ cardStore.cards.length }} cards
          </span>
        </div>

        <!-- Empty State -->
        <div v-if="cardStore.cards.length === 0" class="text-center py-12">
          <svg class="w-16 h-16 text-text-light/30 dark:text-text-dark/30 mx-auto mb-4" fill="none"
            stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">No cards yet</h3>
          <p class="text-text-light/70 dark:text-text-dark/70 mb-6">Click "Add Card" above to create your first
            flashcard.</p>
        </div>

        <!-- Cards List -->
        <div v-else class="space-y-4 overflow-y-auto">
          <div v-for="card in cardStore.cards" :key="card.card_id"
            class="group bg-background-light dark:bg-background-dark rounded-xl border border-surface-light dark:border-surface-dark p-4 hover:shadow-md transition-all duration-200">
            <!-- Edit Mode -->
            <div v-if="editingCard === card.card_id" class="space-y-4">
              <div>
                <label class="block text-xs font-medium text-primary uppercase tracking-wide mb-2">Question</label>
                <textarea v-model="editFormData.question" rows="3"
                  class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                  placeholder="Enter your question...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)" />
              </div>
              <div>
                <label class="block text-xs font-medium text-accent uppercase tracking-wide mb-2">Answer</label>
                <textarea v-model="editFormData.answer" rows="3"
                  class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                  placeholder="Enter the answer...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)" />
              </div>
              <div class="text-xs text-text-light/60 dark:text-text-dark/60 mb-2">
                Supports: **bold**, *italic*, `code`, lists, new lines, and $math formulas$ (e.g., $E = mc^2$)
              </div>
              <div class="flex gap-2 justify-end">
                <BaseButton @click="cancelCardEdit" text="Cancel" variant="danger" size="md" />
                <BaseButton @click="saveCardEdit(card)" text="Save" variant="primary" size="md"
                  :disabled="!editFormData.question.trim() || !editFormData.answer.trim()" />
              </div>
            </div>

            <!-- View Mode -->
            <div v-else class="flex items-start justify-between">
              <div class="flex-1 space-y-3">
                <div>
                  <div class="text-xs font-medium text-primary uppercase tracking-wide mb-1">Question</div>
                  <div
                    class="text-text-light dark:text-text-dark font-medium leading-relaxed prose prose-sm max-w-none dark:prose-invert"
                    v-html="renderMarkdown(card.question)"></div>
                </div>
                <div class="border-t border-surface-light dark:border-surface-dark pt-3">
                  <div class="text-xs font-medium text-accent uppercase tracking-wide mb-1">Answer</div>
                  <div
                    class="text-text-light/80 dark:text-text-dark/80 leading-relaxed prose prose-sm max-w-none dark:prose-invert"
                    v-html="renderMarkdown(card.answer)"></div>
                </div>
              </div>
              <div class="flex gap-2 ml-4 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                <BaseButton @click="editCard(card)" :iconOnly="true" :icon="EditIcon" variant="default" size="lg" />
                <BaseButton @click="deleteCard(card)" :iconOnly="true" :icon="DeleteIcon" variant="danger" size="lg" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
