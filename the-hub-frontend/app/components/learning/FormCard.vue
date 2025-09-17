<script setup lang="ts">
import { ref } from "vue";
import { useCardStore } from '@/stores/cards'
import { useValidation } from '@/composables/useValidation'
import FormUI from '@/components/ui/FormUI.vue'

interface Props {
  deckId: string
}

const props = defineProps<Props>()

const cardStore = useCardStore()
const { schemas } = useValidation()

const showForm = ref(false)

const fields = [
  {
    name: 'question',
    label: 'Question',
    type: 'textarea' as const,
    placeholder: 'Enter the question',
    required: true,
    rows: 3
  },
  {
    name: 'answer',
    label: 'Answer',
    type: 'textarea' as const,
    placeholder: 'Enter the answer',
    required: true,
    rows: 4
  }
]

const initialData = {
  question: '',
  answer: ''
}

const submitForm = async (data: Record<string, any>) => {
  const payload = {
    deck_id: props.deckId,
    question: data.question.trim(),
    answer: data.answer.trim()
  }

  try {
    await cardStore.createCard(payload)
    showForm.value = false // Close modal
  } catch (err) {
    // Error is already handled in the store
  }
}

const cancelForm = () => {
  showForm.value = false
}

const closeModal = () => {
  showForm.value = false
}

const openModal = () => {
  showForm.value = true
}
</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div v-if="!showForm" @click="openModal" class="fixed bottom-4 right-4 cursor-pointer z-40">
        <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
          <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </div>
      </div>
    </Teleport>
  </ClientOnly>

  <FormUI
    title="Create a Flashcard"
    :fields="fields"
    :initial-data="initialData"
    :show-form="showForm"
    submit-label="Create Card"
    teleport-target="body"
    @submit="submitForm"
    @cancel="cancelForm"
    @close="closeModal"
  />
</template>