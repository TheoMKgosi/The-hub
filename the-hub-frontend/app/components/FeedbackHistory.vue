<script setup lang="ts">
import { formatDistanceToNow } from 'date-fns'

interface FeedbackItem {
  feedback_id: string
  type: 'bug' | 'feature' | 'improvement' | 'general'
  subject: string
  description: string
  rating: number | null
  status: 'pending' | 'reviewed' | 'implemented' | 'declined'
  admin_response: string | null
  created_at: string
  page_url: string
}

interface FeedbackHistoryProps {
  feedback?: FeedbackItem[]
  loading?: boolean
  error?: string
}

const props = withDefaults(defineProps<FeedbackHistoryProps>(), {
  feedback: () => [],
  loading: false,
  error: ''
})

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

// Format date
const formatDate = (dateString: string) => {
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true })
  } catch {
    return dateString
  }
}

// Truncate text
const truncateText = (text: string, maxLength: number = 100) => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
      <p class="text-gray-600 dark:text-gray-400">Loading your feedback...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-8">
      <div class="text-red-500 mb-4">
        <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
      </div>
      <p class="text-red-600 dark:text-red-400">{{ error }}</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="feedback.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">No feedback yet</h3>
      <p class="text-gray-600 dark:text-gray-400">Share your thoughts to help us improve The Hub!</p>
    </div>

    <!-- Feedback List -->
    <div v-else class="space-y-4">
      <div v-for="item in feedback" :key="item.feedback_id" class="bg-surface-light dark:bg-surface-dark rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-4">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <span class="text-2xl">{{ typeConfig[item.type].icon }}</span>
              <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">{{ item.subject }}</h3>
            </div>
            <div class="flex flex-wrap items-center gap-2">
              <span :class="['px-2 py-1 rounded-full text-xs font-medium', typeConfig[item.type].color]">
                {{ typeConfig[item.type].label }}
              </span>
              <span :class="['px-2 py-1 rounded-full text-xs font-medium', statusConfig[item.status].color]">
                {{ statusConfig[item.status].icon }} {{ statusConfig[item.status].label }}
              </span>
              <span v-if="item.rating" class="flex items-center gap-1 text-sm text-gray-600 dark:text-gray-400">
                <UiStarRating :model-value="item.rating" readonly size="sm" />
                {{ item.rating }}/5
              </span>
            </div>
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ formatDate(item.created_at) }}
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
      </div>
    </div>
  </div>
</template>