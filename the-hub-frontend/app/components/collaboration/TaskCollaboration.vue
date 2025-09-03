<script setup lang="ts">
import { ref, reactive } from 'vue'

interface TaskShare {
  share_id: string
  task_id: string
  shared_with_id: string
  permission: string
  shared_with?: {
    id: string
    name: string
    email: string
  }
}

interface TaskComment {
  comment_id: string
  task_id: string
  user_id: string
  content: string
  created_at: string
  user?: {
    id: string
    name: string
    email: string
  }
}

interface Props {
  taskId: string
  taskTitle: string
}

const props = defineProps<Props>()

const shares = ref<TaskShare[]>([])
const comments = ref<TaskComment[]>([])
const newComment = ref('')
const shareEmail = ref('')
const sharePermission = ref('view')
const showShareForm = ref(false)
const showComments = ref(false)
const loading = ref(false)

const permissionOptions = [
  { value: 'view', label: 'View Only' },
  { value: 'edit', label: 'Can Edit' },
  { value: 'admin', label: 'Full Access' },
]

const loadShares = async () => {
  try {
    const response = await $fetch<{ shares: TaskShare[] }>(`/api/tasks/${props.taskId}/shares`)
    shares.value = response.shares
  } catch (error) {
    console.error('Failed to load shares:', error)
  }
}

const loadComments = async () => {
  try {
    const response = await $fetch<{ comments: TaskComment[] }>(`/api/tasks/${props.taskId}/comments`)
    comments.value = response.comments
  } catch (error) {
    console.error('Failed to load comments:', error)
  }
}

const shareTask = async () => {
  if (!shareEmail.value.trim()) return

  loading.value = true
  try {
    // First, we need to find the user by email
    const users = await $fetch<{ users: Array<{ id: string; name: string; email: string }> }>('/api/users/search', {
      method: 'GET',
      query: { email: shareEmail.value }
    })

    if (users.users.length === 0) {
      throw new Error('User not found')
    }

    const userId = users.users[0].id

    await $fetch(`/api/tasks/${props.taskId}/share`, {
      method: 'POST',
      body: JSON.stringify({
        shared_with_id: userId,
        permission: sharePermission.value
      })
    })

    shareEmail.value = ''
    sharePermission.value = 'view'
    showShareForm.value = false
    await loadShares()
  } catch (error) {
    console.error('Failed to share task:', error)
    alert('Failed to share task: ' + error.message)
  } finally {
    loading.value = false
  }
}

const addComment = async () => {
  if (!newComment.value.trim()) return

  loading.value = true
  try {
    await $fetch(`/api/tasks/${props.taskId}/comments`, {
      method: 'POST',
      body: JSON.stringify({
        content: newComment.value
      })
    })

    newComment.value = ''
    await loadComments()
  } catch (error) {
    console.error('Failed to add comment:', error)
  } finally {
    loading.value = false
  }
}

const removeShare = async (shareId: string) => {
  if (!confirm('Are you sure you want to stop sharing this task?')) return

  try {
    await $fetch(`/api/tasks/${props.taskId}/shares/${shareId}`, {
      method: 'DELETE'
    })
    await loadShares()
  } catch (error) {
    console.error('Failed to remove share:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const getPermissionColor = (permission: string) => {
  switch (permission) {
    case 'view': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
    case 'edit': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
    case 'admin': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
    default: return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200'
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-medium text-text-light dark:text-text-dark">Collaboration</h3>
      <div class="flex gap-2">
        <button
          @click="showShareForm = !showShareForm"
          class="px-3 py-1 text-sm bg-primary text-white rounded-md hover:bg-primary/90 transition-colors"
        >
          Share Task
        </button>
        <button
          @click="showComments = !showComments; if (showComments) loadComments()"
          class="px-3 py-1 text-sm bg-secondary text-white rounded-md hover:bg-secondary/90 transition-colors"
        >
          Comments ({{ comments.length }})
        </button>
      </div>
    </div>

    <!-- Share Form -->
    <div v-if="showShareForm" class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
      <h4 class="font-medium mb-3 text-text-light dark:text-text-dark">Share with others</h4>

      <div class="flex gap-2">
        <input
          v-model="shareEmail"
          type="email"
          placeholder="Enter email address"
          class="flex-1 px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        />

        <select
          v-model="sharePermission"
          class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
        >
          <option v-for="option in permissionOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>

        <button
          @click="shareTask"
          :disabled="loading || !shareEmail.trim()"
          class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading" class="animate-spin">⟳</span>
          <span v-else>Share</span>
        </button>
      </div>
    </div>

    <!-- Current Shares -->
    <div v-if="shares.length > 0" class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
      <h4 class="font-medium mb-3 text-text-light dark:text-text-dark">Shared with</h4>

      <div class="space-y-2">
        <div
          v-for="share in shares"
          :key="share.share_id"
          class="flex items-center justify-between p-3 bg-surface-light/20 dark:bg-surface-dark/20 rounded-md"
        >
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center text-white font-medium">
              {{ share.shared_with?.name?.charAt(0)?.toUpperCase() || '?' }}
            </div>
            <div>
              <p class="font-medium text-text-light dark:text-text-dark">
                {{ share.shared_with?.name || 'Unknown User' }}
              </p>
              <p class="text-sm text-text-light/70 dark:text-text-dark/70">
                {{ share.shared_with?.email }}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <span :class="['px-2 py-1 text-xs rounded-full', getPermissionColor(share.permission)]">
              {{ share.permission }}
            </span>
            <button
              @click="removeShare(share.share_id)"
              class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
            >
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comments Section -->
    <div v-if="showComments" class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
      <h4 class="font-medium mb-3 text-text-light dark:text-text-dark">Comments</h4>

      <!-- Add Comment -->
      <div class="mb-4">
        <div class="flex gap-2">
          <textarea
            v-model="newComment"
            placeholder="Add a comment..."
            rows="2"
            class="flex-1 px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none"
          ></textarea>
          <button
            @click="addComment"
            :disabled="loading || !newComment.trim()"
            class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed self-end"
          >
            Post
          </button>
        </div>
      </div>

      <!-- Comments List -->
      <div v-if="comments.length > 0" class="space-y-3">
        <div
          v-for="comment in comments"
          :key="comment.comment_id"
          class="p-3 bg-surface-light/20 dark:bg-surface-dark/20 rounded-md"
        >
          <div class="flex items-start gap-3">
            <div class="w-8 h-8 bg-secondary rounded-full flex items-center justify-center text-white font-medium flex-shrink-0">
              {{ comment.user?.name?.charAt(0)?.toUpperCase() || '?' }}
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="font-medium text-text-light dark:text-text-dark">
                  {{ comment.user?.name || 'Unknown User' }}
                </span>
                <span class="text-xs text-text-light/60 dark:text-text-dark/60">
                  {{ formatDate(comment.created_at) }}
                </span>
              </div>
              <p class="text-text-light dark:text-text-dark whitespace-pre-wrap">
                {{ comment.content }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-4 text-text-light/60 dark:text-text-dark/60">
        No comments yet. Be the first to comment!
      </div>
    </div>
  </div>
</template>