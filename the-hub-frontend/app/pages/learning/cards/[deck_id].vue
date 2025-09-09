<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const deckID = route.params.deck_id as string

const formData = reactive({
  deck_id: deckID,
  question: '',
  answer: ''
})

const editingCard = ref<string | null>(null)
const editFormData = reactive({
  question: '',
  answer: ''
})

const cardStore = useCardStore()
const { renderMarkdown } = useMarkdown()

const submitForm = () => {
  cardStore.submitForm(deckID, { ...formData })
  formData.question = ''
  formData.answer = ''
}

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
                 placeholder="Enter your question here...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)"
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
                 placeholder="Enter the answer here...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)"
                 required
               />
             </div>
             <div class="text-xs text-text-light/60 dark:text-text-dark/60">
               Supports: **bold**, *italic*, `code`, lists, new lines, and $math formulas$ (e.g., $E = mc^2$)
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
               <!-- Edit Mode -->
               <div v-if="editingCard === card.card_id" class="space-y-4">
                 <div>
                   <label class="block text-xs font-medium text-primary uppercase tracking-wide mb-2">Question</label>
                   <textarea
                     v-model="editFormData.question"
                     rows="3"
                     class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                     placeholder="Enter your question...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)"
                   />
                 </div>
                 <div>
                   <label class="block text-xs font-medium text-accent uppercase tracking-wide mb-2">Answer</label>
                   <textarea
                     v-model="editFormData.answer"
                     rows="3"
                     class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
                     placeholder="Enter the answer...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)"
                   />
                 </div>
                 <div class="text-xs text-text-light/60 dark:text-text-dark/60 mb-2">
                   Supports: **bold**, *italic*, `code`, lists, new lines, and $math formulas$ (e.g., $E = mc^2$)
                 </div>
                 <div class="flex gap-2 justify-end">
                   <UiButton
                     @click="cancelCardEdit"
                     variant="default"
                     size="sm"
                   >
                     Cancel
                   </UiButton>
                   <UiButton
                     @click="saveCardEdit(card)"
                     variant="primary"
                     size="sm"
                     :disabled="!editFormData.question.trim() || !editFormData.answer.trim()"
                   >
                     Save
                   </UiButton>
                 </div>
               </div>

               <!-- View Mode -->
               <div v-else class="flex items-start justify-between">
                 <div class="flex-1 space-y-3">
                   <div>
                     <div class="text-xs font-medium text-primary uppercase tracking-wide mb-1">Question</div>
                     <div class="text-text-light dark:text-text-dark font-medium leading-relaxed prose prose-sm max-w-none dark:prose-invert" v-html="renderMarkdown(card.question)"></div>
                   </div>
                   <div class="border-t border-surface-light dark:border-surface-dark pt-3">
                     <div class="text-xs font-medium text-accent uppercase tracking-wide mb-1">Answer</div>
                     <div class="text-text-light/80 dark:text-text-dark/80 leading-relaxed prose prose-sm max-w-none dark:prose-invert" v-html="renderMarkdown(card.answer)"></div>
                   </div>
                 </div>
                 <div class="flex gap-2 ml-4 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                   <UiButton
                     @click="editCard(card)"
                     variant="default"
                     size="sm"
                     title="Edit Card"
                   >
                     <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                     </svg>
                   </UiButton>
                   <UiButton
                     @click="deleteCard(card)"
                     variant="danger"
                     size="sm"
                     title="Delete Card"
                   >
                     <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                     </svg>
                   </UiButton>
                 </div>
               </div>
             </div>
           </div>
        </div>
      </div>
    </div>
  </div>
</template>

