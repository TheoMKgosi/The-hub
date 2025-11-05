<script setup lang="ts">
definePageMeta({
  middleware: ['admin']
})

const { addToast } = useToast()
const { $api } = useNuxtApp()

// Feedback data
interface FeedbackItem {
  feedback_id: string
  user_id: string
  user?: {
    name: string
    email: string
  }
  type: 'bug' | 'feature' | 'improvement' | 'general'
  subject: string
  description: string
  rating: number | null
  status: 'pending' | 'reviewed' | 'implemented' | 'declined'
  admin_response: string | null
  created_at: string
  updated_at: string
  page_url: string
}

// State
const feedback = ref<FeedbackItem[]>([])
const loading = ref(false)
const updating = ref<string | null>(null)

// Filters
const statusFilter = ref<string>('')
const typeFilter = ref<string>('')
const currentPage = ref(1)
const totalPages = ref(1)
const totalItems = ref(0)
const pageSize = 20

// Status configuration
const statusConfig = {
  pending: {
    label: 'Pending Review',
    color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    icon: '⏳'
  },
  reviewed: {
    label: 'Reviewed',
    color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
    icon: '👁️'
  },
  implemented: {
    label: 'Implemented',
    color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    icon: '✅'
  },
  declined: {
    label: 'Declined',
    color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
    icon: '❌'
  }
}

// Type configuration
const typeConfig = {
  bug: {
    label: 'Bug Report',
    color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
    icon: '🐛'
  },
  feature: {
    label: 'Feature Request',
    color: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
    icon: '✨'
  },
  improvement: {
    label: 'Improvement',
    color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
    icon: '🚀'
  },
  general: {
    label: 'General Feedback',
    color: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200',
    icon: '💬'
  }
}

// Load feedback
const loadFeedback = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: currentPage.value.toString(),
      limit: pageSize.toString()
    })

    if (statusFilter.value) params.append('status', statusFilter.value)
    if (typeFilter.value) params.append('type', typeFilter.value)

    const response = await $api(`/admin/feedback?${params}`)
    feedback.value = response.feedback || []
    totalItems.value = response.pagination?.total || 0
    totalPages.value = response.pagination?.pages || 1
  } catch (error: any) {
    console.error('Failed to load feedback:', error)
    addToast(error?.data?.error || 'Failed to load feedback', 'error')
  } finally {
    loading.value = false
  }
}

// Update feedback status
const updateFeedbackStatus = async (feedbackId: string, status: string, adminResponse?: string) => {
  updating.value = feedbackId
  try {
    const updateData: any = { status }
    if (adminResponse !== undefined) {
      updateData.admin_response = adminResponse
    }

    const response = await $api(`/admin/feedback/${feedbackId}`, {
      method: 'PATCH',
      body: JSON.stringify(updateData)
    })

    // Update local state
    const index = feedback.value.findIndex(f => f.feedback_id === feedbackId)
    if (index !== -1) {
      feedback.value[index] = response.feedback
    }

    addToast('Feedback status updated successfully', 'success')
  } catch (error: any) {
    console.error('Failed to update feedback:', error)
    addToast(error?.data?.error || 'Failed to update feedback', 'error')
  } finally {
    updating.value = null
  }
}

// Format date
const formatDate = (dateString: string) => {
  try {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateString
  }
}

// Pagination
const goToPage = (page: number) => {
  currentPage.value = page
  loadFeedback()
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadFeedback()
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadFeedback()
  }
}

// Watch filters
watch([statusFilter, typeFilter], () => {
  currentPage.value = 1
  loadFeedback()
})

// Load on mount
onMounted(() => {
  loadFeedback()
})

// Status options for filter
const statusOptions = [
  { value: '', label: 'All Statuses' },
  { value: 'pending', label: 'Pending Review' },
  { value: 'reviewed', label: 'Reviewed' },
  { value: 'implemented', label: 'Implemented' },
  { value: 'declined', label: 'Declined' }
]

// Type options for filter
const typeOptions = [
  { value: '', label: 'All Types' },
  { value: 'bug', label: 'Bug Reports' },
  { value: 'feature', label: 'Feature Requests' },
  { value: 'improvement', label: 'Improvements' },
  { value: 'general', label: 'General Feedback' }
]
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <!-- Admin Header -->
    <div class="bg-surface-light dark:bg-surface-dark border-b border-surface-light/50 dark:border-surface-dark/50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-3">
              <img src="/logo.svg" alt="The Hub Admin" class="h-8 w-8" />
              <div>
                <h1 class="text-xl font-bold text-text-light dark:text-text-dark">
                  Feedback Management
                </h1>
                <p class="text-sm text-text-light/60 dark:text-text-dark/60">
                  Review and manage user feedback
                </p>
              </div>
            </div>
          </div>

          <div class="flex items-center space-x-4">
            <UiButton
              @click="loadFeedback"
              variant="secondary"
              size="sm"
              :loading="loading"
            >
              Refresh
            </UiButton>

            <UiButton
              @click="$router.push('/admin')"
              variant="default"
              size="sm"
            >
              Back to Dashboard
            </UiButton>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filters -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 shadow-sm border border-surface-light/50 dark:border-surface-dark/50 mb-6">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Filter by Status
            </label>
            <select
              v-model="statusFilter"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option v-for="option in statusOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>

          <div class="flex-1">
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Filter by Type
            </label>
            <select
              v-model="typeFilter"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option v-for="option in typeOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-light dark:border-primary-dark"></div>
      </div>

      <!-- Feedback List -->
      <div v-else-if="feedback.length > 0" class="space-y-4">
        <div v-for="item in feedback" :key="item.feedback_id" class="bg-surface-light dark:bg-surface-dark rounded-lg border border-gray-200 dark:border-gray-700 p-6">
          <!-- Header -->
          <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4 mb-4">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span class="text-2xl">{{ typeConfig[item.type].icon }}</span>
                <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">{{ item.subject }}</h3>
              </div>
              <div class="flex flex-wrap items-center gap-2 mb-2">
                <span :class="['px-2 py-1 rounded-full text-xs font-medium', typeConfig[item.type].color]">
                  {{ typeConfig[item.type].label }}
                </span>
                <span :class="['px-2 py-1 rounded-full text-xs font-medium', statusConfig[item.status].color]">
                  {{ statusConfig[item.status].icon }} {{ statusConfig[item.status].label }}
                </span>
                <span v-if="item.rating" class="flex items-center gap-1 text-sm text-gray-600 dark:text-gray-400">
                  <StarRating :model-value="item.rating" readonly size="sm" />
                  {{ item.rating }}/5
                </span>
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                <span class="font-medium">{{ item.user?.name || 'Unknown User' }}</span>
                <span class="mx-2">•</span>
                <span>{{ item.user?.email }}</span>
                <span class="mx-2">•</span>
                <span>{{ formatDate(item.created_at) }}</span>
              </div>
            </div>

            <!-- Status Update Actions -->
            <div class="flex flex-col gap-2 min-w-[200px]">
              <label class="block text-sm font-medium text-text-light dark:text-text-dark">
                Update Status
              </label>
              <select
                :value="item.status"
                @change="updateFeedbackStatus(item.feedback_id, ($event.target as HTMLSelectElement).value)"
                :disabled="updating === item.feedback_id"
                class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary text-sm"
              >
                <option value="pending">Pending Review</option>
                <option value="reviewed">Reviewed</option>
                <option value="implemented">Implemented</option>
                <option value="declined">Declined</option>
              </select>
            </div>
          </div>

          <!-- Description -->
          <div class="mb-4">
            <p class="text-text-light dark:text-text-dark whitespace-pre-wrap">{{ item.description }}</p>
          </div>

          <!-- Page URL -->
          <div v-if="item.page_url" class="mb-4 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div class="text-xs text-gray-600 dark:text-gray-400 mb-1">Submitted from:</div>
            <div class="text-xs text-gray-800 dark:text-gray-200 break-all">{{ item.page_url }}</div>
          </div>

          <!-- Admin Response -->
          <div v-if="item.admin_response" class="border-t border-gray-200 dark:border-gray-700 pt-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center flex-shrink-0">
                <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="flex-1">
                <div class="text-sm font-medium text-text-light dark:text-text-dark mb-1">Admin Response</div>
                <p class="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{{ item.admin_response }}</p>
              </div>
            </div>
          </div>

          <!-- Add Response -->
          <div v-else class="border-t border-gray-200 dark:border-gray-700 pt-4">
            <div class="flex items-start gap-3">
              <div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center flex-shrink-0">
                <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="flex-1">
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
                  Add Admin Response
                </label>
                <textarea
                  placeholder="Add a response to this feedback..."
                  @keydown.enter.exact.prevent="updateFeedbackStatus(item.feedback_id, item.status, ($event.target as HTMLTextAreaElement).value)"
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary resize-none"
                  rows="3"
                ></textarea>
                <div class="mt-2 flex justify-end">
                  <UiButton
                    @click="updateFeedbackStatus(item.feedback_id, item.status, ($event.target as HTMLElement).previousElementSibling?.querySelector('textarea')?.value || '')"
                    variant="primary"
                    size="sm"
                    :disabled="updating === item.feedback_id"
                  >
                    {{ updating === item.feedback_id ? 'Updating...' : 'Add Response' }}
                  </UiButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9 8s9 3.582 9 8z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">No feedback found</h3>
        <p class="text-gray-600 dark:text-gray-400">
          {{ statusFilter || typeFilter ? 'Try adjusting your filters' : 'No feedback has been submitted yet' }}
        </p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-2 mt-8">
        <UiButton
          @click="prevPage"
          :disabled="currentPage === 1"
          variant="secondary"
          size="sm"
        >
          Previous
        </UiButton>

        <div class="flex gap-1">
          <button
            v-for="page in Math.min(5, totalPages)"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-3 py-2 rounded-md text-sm font-medium transition-colors',
              currentPage === page
                ? 'bg-primary text-white'
                : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'
            ]"
          >
            {{ page }}
          </button>
        </div>

        <UiButton
          @click="nextPage"
          :disabled="currentPage === totalPages"
          variant="secondary"
          size="sm"
        >
          Next
        </UiButton>
      </div>

      <!-- Stats -->
      <div class="mt-8 text-center text-sm text-gray-600 dark:text-gray-400">
        Showing {{ feedback.length }} of {{ totalItems }} feedback items
      </div>
    </div>
  </div>
</template>