<script setup lang="ts">
import LeftArrowIcon from '~/components/ui/svg/LeftArrowIcon.vue'
import PlusIcon from '~/components/ui/svg/PlusIcon.vue'
const route = useRoute()
const router = useRouter()

const deckID = route.params.deck_id as string

const formData = reactive({
  deck_id: deckID,
  question: '',
  answer: ''
})

const cardStore = useCardStore()
const { addToast } = useToast()

const submitForm = async () => {
  if (!formData.question.trim() || !formData.answer.trim()) return

  try {
    await cardStore.submitForm(deckID, { ...formData })

    // Reset form for next card
    formData.question = ''
    formData.answer = ''

    addToast("Card created successfully! Ready for another card.", "success")
  } catch (error) {
    addToast("Failed to create card", "error")
  }
}

const goBack = () => {
  router.push(`/learning/cards/${deckID}`)
}

// Load deck info on mount (optional, for page title)
onMounted(() => {
  // Could fetch deck name for better UX
})
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <div class="mx-auto px-4 py-8">
      <!-- Header -->
      <div class="mb-8">
        <BaseButton @click="goBack" text="Back to Cards" :icon="LeftArrowIcon" variant="default" size="md" />
          <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-2">Create New Cards</h1>
          <p class="text-text-light/70 dark:text-text-dark/70">Add flashcards to your deck. Create as many as you need!
          </p>
      </div>

      <!-- Creation Form -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-lg border border-surface-light dark:border-surface-dark p-8">
        <form @submit.prevent="submitForm" class="space-y-6">
          <div>
            <label for="question" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Question
            </label>
            <textarea id="question" v-model="formData.question" rows="4"
              class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
              placeholder="Enter your question here...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)" required />
          </div>

          <div>
            <label for="answer" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Answer
            </label>
            <textarea id="answer" v-model="formData.answer" rows="4"
              class="w-full border border-surface-light dark:border-surface-dark bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 transition-all duration-200 resize-none"
              placeholder="Enter the answer here...
New lines are preserved. (supports **bold**, *italic*, `code`, and $math$ formulas)" required />
          </div>

          <div class="text-xs text-text-light/60 dark:text-text-dark/60">
            Supports: **bold**, *italic*, `code`, lists, new lines, and $math formulas$ (e.g., $E = mc^2$)
          </div>

          <BaseButton type="submit" text="Create Card" :icon="PlusIcon" variant="primary" size="md"
            :disabled="!formData.question.trim() || !formData.answer.trim()" />
        </form>
      </div>
    </div>
  </div>
</template>
