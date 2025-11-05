<script setup lang="ts">
import { useValidation } from '@/composables/useValidation'

interface FeedbackFormProps {
  showForm?: boolean
}

interface FeedbackFormEmits {
  (e: 'submit', data: FeedbackData): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<FeedbackFormProps>(), {
  showForm: true
})

const emit = defineEmits<FeedbackFormEmits>()

const { addToast } = useToast()
const { $api } = useNuxtApp()
const { validateObject } = useValidation()

interface FeedbackData {
  type: string
  subject: string
  description: string
  rating: number | null
  page_url: string
  user_agent: string
}

// Form data
const formData = reactive<FeedbackData>({
  type: '',
  subject: '',
  description: '',
  rating: null,
  page_url: '',
  user_agent: ''
})

// Form state
const isSubmitting = ref(false)
const validationErrors = ref<Record<string, string>>({})

// Initialize form data
onMounted(() => {
  formData.page_url = window.location.href
  formData.user_agent = navigator.userAgent
})

// Form validation schema
const validationSchema = {
  type: { required: true },
  subject: { required: true, minLength: 5, maxLength: 255 },
  description: { required: true, minLength: 10, maxLength: 2000 }
}

// Feedback type options
const feedbackTypes = [
  { value: 'bug', label: 'Bug Report', description: 'Report a problem or error' },
  { value: 'feature', label: 'Feature Request', description: 'Suggest a new feature' },
  { value: 'improvement', label: 'Improvement', description: 'Suggest an enhancement' },
  { value: 'general', label: 'General Feedback', description: 'Other feedback or comments' }
]

// Submit feedback
const submitFeedback = async () => {
  validationErrors.value = {}

  // Validate form
  const validation = validateObject(formData, validationSchema)
  if (!validation.isValid) {
    validationErrors.value = validation.errors
    return
  }

  try {
    isSubmitting.value = true

    const response = await $api('/feedback', {
      method: 'POST',
      body: JSON.stringify(formData)
    })

    addToast('Thank you for your feedback!', 'success')
    emit('submit', formData)
    resetForm()
  } catch (error: any) {
    console.error('Failed to submit feedback:', error)
    addToast(error?.data?.error || 'Failed to submit feedback. Please try again.', 'error')
  } finally {
    isSubmitting.value = false
  }
}

// Reset form
const resetForm = () => {
  formData.type = ''
  formData.subject = ''
  formData.description = ''
  formData.rating = null
  validationErrors.value = {}
}

// Close form
const closeForm = () => {
  resetForm()
  emit('close')
}

// Check if form is valid
const isFormValid = computed(() => {
  return formData.type &&
         formData.subject.trim().length >= 5 &&
         formData.description.trim().length >= 10 &&
         !isSubmitting.value
})
</script>

<template>
  <div v-if="showForm" class="max-w-2xl mx-auto p-6">
    <div class="bg-surface-light dark:bg-surface-dark rounded-xl shadow-lg border border-gray-200 dark:border-gray-700">
      <!-- Header -->
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-2xl font-bold text-text-light dark:text-text-dark mb-2">
          Share Your Feedback
        </h2>
        <p class="text-gray-600 dark:text-gray-400">
          Help us improve The Hub by sharing your thoughts, reporting bugs, or suggesting new features.
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="submitFeedback" class="p-6 space-y-6">
        <!-- Feedback Type -->
        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-3">
            What type of feedback is this? <span class="text-red-500">*</span>
          </label>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <button
              v-for="type in feedbackTypes"
              :key="type.value"
              type="button"
              @click="formData.type = type.value"
              :class="[
                'p-4 rounded-lg border-2 text-left transition-all duration-200',
                formData.type === type.value
                  ? 'border-primary bg-primary/10 text-primary'
                  : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
              ]"
            >
              <div class="font-medium text-text-light dark:text-text-dark">{{ type.label }}</div>
              <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ type.description }}</div>
            </button>
          </div>
          <p v-if="validationErrors.type" class="mt-2 text-sm text-red-500">
            {{ validationErrors.type }}
          </p>
        </div>

        <!-- Subject -->
        <div>
          <label for="subject" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
            Subject <span class="text-red-500">*</span>
          </label>
          <input
            id="subject"
            v-model="formData.subject"
            type="text"
            placeholder="Brief summary of your feedback"
            :class="[
              'w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary transition-colors',
              'bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark',
              'border-gray-300 dark:border-gray-600',
              validationErrors.subject ? 'border-red-500 focus:ring-red-500' : 'focus:border-primary'
            ]"
            maxlength="255"
          />
          <p v-if="validationErrors.subject" class="mt-1 text-sm text-red-500">
            {{ validationErrors.subject }}
          </p>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            {{ formData.subject.length }}/255 characters
          </p>
        </div>

        <!-- Description -->
        <div>
          <label for="description" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
            Description <span class="text-red-500">*</span>
          </label>
          <textarea
            id="description"
            v-model="formData.description"
            placeholder="Please provide detailed information about your feedback..."
            rows="6"
            :class="[
              'w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary resize-none transition-colors',
              'bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark',
              'border-gray-300 dark:border-gray-600',
              validationErrors.description ? 'border-red-500 focus:ring-red-500' : 'focus:border-primary'
            ]"
            maxlength="2000"
          ></textarea>
          <p v-if="validationErrors.description" class="mt-1 text-sm text-red-500">
            {{ validationErrors.description }}
          </p>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            {{ formData.description.length }}/2000 characters
          </p>
        </div>

        <!-- Rating (Optional) -->
        <div>
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
            Overall Rating (Optional)
          </label>
          <div class="flex items-center gap-4">
            <StarRating v-model="formData.rating" size="lg" />
            <span v-if="formData.rating" class="text-sm text-gray-600 dark:text-gray-400">
              {{ formData.rating }} star{{ formData.rating > 1 ? 's' : '' }}
            </span>
          </div>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            How would you rate your overall experience?
          </p>
        </div>

        <!-- Current Page Info -->
        <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
          <h4 class="text-sm font-medium text-text-light dark:text-text-dark mb-2">Current Page</h4>
          <p class="text-xs text-gray-600 dark:text-gray-400 break-all">{{ formData.page_url }}</p>
        </div>

        <!-- Actions -->
        <div class="flex flex-col-reverse sm:flex-row gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
          <UiButton
            type="button"
            @click="closeForm"
            variant="default"
            size="md"
            class="w-full sm:w-auto"
            :disabled="isSubmitting"
          >
            Cancel
          </UiButton>
          <UiButton
            type="submit"
            variant="primary"
            size="md"
            class="w-full sm:w-auto"
            :disabled="!isFormValid"
          >
            {{ isSubmitting ? 'Submitting...' : 'Submit Feedback' }}
          </UiButton>
        </div>
      </form>
    </div>
  </div>
</template>