<script setup lang="ts">
import { formatDistanceToNow } from 'date-fns'

const { addToast } = useToast()
const { $api } = useNuxtApp()

// Page state
const activeTab = ref<'submit' | 'history'>('submit')
const isLoadingHistory = ref(false)
const historyError = ref('')
const feedbackHistory = ref<any[]>([])

// Load feedback history
const loadFeedbackHistory = async () => {
  try {
    isLoadingHistory.value = true
    historyError.value = ''

    const response = await $api('/feedback')
    feedbackHistory.value = response.feedback || []
  } catch (error: any) {
    console.error('Failed to load feedback history:', error)
    historyError.value = error?.data?.error || 'Failed to load feedback history'
  } finally {
    isLoadingHistory.value = false
  }
}

// Handle feedback submission
const handleFeedbackSubmit = async (data: any) => {
  // Switch to history tab to show the new feedback
  activeTab.value = 'history'
  await loadFeedbackHistory()

  // Scroll to top of history section
  nextTick(() => {
    const historySection = document.getElementById('feedback-history')
    if (historySection) {
      historySection.scrollIntoView({ behavior: 'smooth' })
    }
  })
}

// Load history when switching to history tab
watch(activeTab, async (newTab) => {
  if (newTab === 'history' && feedbackHistory.value.length === 0) {
    await loadFeedbackHistory()
  }
})

// Load history on mount if starting on history tab
onMounted(async () => {
  if (activeTab.value === 'history') {
    await loadFeedbackHistory()
  }
})

// Tab configuration
const tabs = [
  { id: 'submit', label: 'Submit Feedback', icon: '📝' },
  { id: 'history', label: 'My Feedback', icon: '📋' }
]
</script>

<template>
  <div class="max-w-4xl mx-auto p-6 space-y-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-4xl font-bold text-text-light dark:text-text-dark mb-4">
        Feedback & Support
      </h1>
      <p class="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
        Your feedback helps us improve The Hub. Share your thoughts, report bugs, or suggest new features.
      </p>
    </div>

    <!-- Tabs -->
    <div class="flex justify-center">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-1 border border-gray-200 dark:border-gray-700">
        <div class="flex">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id as 'submit' | 'history'"
            :class="[
              'flex items-center gap-2 px-6 py-3 rounded-md text-sm font-medium transition-all duration-200',
              activeTab === tab.id
                ? 'bg-primary text-white shadow-sm'
                : 'text-gray-600 dark:text-gray-400 hover:text-text-light dark:hover:text-text-dark hover:bg-gray-50 dark:hover:bg-gray-800'
            ]"
          >
            <span>{{ tab.icon }}</span>
            {{ tab.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="min-h-[600px]">
      <!-- Submit Feedback Tab -->
      <div v-if="activeTab === 'submit'">
        <FeedbackForm @submit="handleFeedbackSubmit" />
      </div>

      <!-- History Tab -->
      <div v-else id="feedback-history" class="space-y-6">
        <div class="text-center">
          <h2 class="text-2xl font-semibold text-text-light dark:text-text-dark mb-2">
            Your Feedback History
          </h2>
          <p class="text-gray-600 dark:text-gray-400">
            Track the status of your feedback and see admin responses.
          </p>
        </div>

        <FeedbackHistory
          :feedback="feedbackHistory"
          :loading="isLoadingHistory"
          :error="historyError"
        />

        <!-- Refresh Button -->
        <div v-if="!isLoadingHistory && !historyError" class="text-center">
          <UiButton
            @click="loadFeedbackHistory"
            variant="default"
            size="sm"
            :disabled="isLoadingHistory"
          >
            <svg v-if="isLoadingHistory" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Refresh
          </UiButton>
        </div>
      </div>
    </div>

    <!-- Help Section -->
    <div class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 border border-gray-200 dark:border-gray-700">
      <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-4">
        Need Help?
      </h3>
      <div class="grid md:grid-cols-2 gap-6">
        <div>
          <h4 class="font-medium text-text-light dark:text-text-dark mb-2">Before Submitting</h4>
          <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
            <li>• Check if your issue has already been reported</li>
            <li>• Provide clear, detailed descriptions</li>
            <li>• Include steps to reproduce bugs</li>
            <li>• Be specific about feature requests</li>
          </ul>
        </div>
        <div>
          <h4 class="font-medium text-text-light dark:text-text-dark mb-2">Response Times</h4>
          <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
            <li>• Bug reports: Usually reviewed within 24 hours</li>
            <li>• Feature requests: Reviewed within 1-2 days</li>
            <li>• General feedback: Acknowledged within 48 hours</li>
            <li>• Implementation updates provided as available</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>