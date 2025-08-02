<script setup lang="ts">
import { useTagStore } from '@/stores/tags';
import { useTopicStore } from '@/stores/topics';
import { onMounted, ref, computed } from 'vue';
import { useRouter } from 'vue-router';

const topicStore = useTopicStore()
const tagStore = useTagStore()

const router = useRouter()

// Form state
const showForm = ref(false)
const editingTopic = ref(null)
const formData = ref({
  title: '',
  description: '',
  status: 'pending',
  deadline: new Date(),
  tags: [] as number[]
})
const editFormData = ref({
  topic_id: 0,
  title: '',
  description: '',
  status: 'pending',
  deadline: new Date(),
  tags: [] as number[]
})

// Tag input
const tagInput = ref({
  name: '',
  color: ''
})
const showTagSuggestions = ref(false)

// Filter state
const statusFilter = ref('all')
const searchQuery = ref('')

const statusOptions = [
  { value: 'pending', label: 'Pending' },
  { value: 'in-progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
  { value: 'on-hold', label: 'On Hold' }
]


// Computed properties
const filteredTopics = computed(() => {
  let topics = topicStore.topics
  if (statusFilter.value !== 'all') {
    topics = topics.filter(topic => topic.status === statusFilter.value)
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    topics = topics.filter(topic =>
      topic.title.toLowerCase().includes(query) ||
      topic.description.toLowerCase().includes(query)
      // topic.tags.some(tag => tag.toLowerCase().includes(query))
    )
  }

  return topics
})

onMounted(() => {
  if (topicStore.topics.length === 0) {
    topicStore.fetchTopics()
  }
  if (tagStore.tags.length === 0) {
    tagStore.fetchTags()
  }
})

const availableTags = computed(() => {
  if (!tagInput.value.trim()) return []
  const query = tagInput.value.toLowerCase()
  return tagStore.tags.filter(tag =>
    tag.name.toLowerCase().includes(query) &&
    !formData.value.tags.includes(tag)
  )
})

// Form methods
const openForm = (topic = null) => {
  if (topic) {
    editingTopic.value = topic
    formData.value = {
      title: topic.title,
      description: topic.description,
      status: topic.status,
      deadline: topic.deadline,
      tags: [...topic.tags]
    }
  } else {
    editingTopic.value = null
    resetForm()
  }
  showForm.value = true
}

const closeForm = () => {
  showForm.value = false
  editingTopic.value = null
  resetForm()
}

const resetForm = () => {
  formData.value = {
    title: '',
    description: '',
    status: 'pending',
    deadline: new Date(),
    tags: []
  }
  tagInput.value = {
    name: '',
    color: ''
  }
}

const handleSubmit = async () => {
  if (!formData.value.title.trim()) return

  try {
    if (editingTopic.value) {
      await topicStore.editTopic(editFormData.value)
    } else {
      await topicStore.submitForm(formData.value)
    }
    closeForm()
  } catch (error) {
    console.error('Error saving topic:', error)
  }
}

const deleteTopic = async (topicId: number) => {
  if (confirm('Are you sure you want to delete this topic?')) {
    try {
      await topicStore.deleteTopic(topicId)
    } catch (error) {
      console.error('Error deleting topic:', error)
    }
  }
}

// Tag methods
const addTag = async () => {
  // if (!tagName.trim() || formData.value.tags.some(tag => tag.name === tagName)) return

  // Check if tag exists, if not create it
  if (!tagStore.tags.map(tag => tag.name).includes(tagInput.value.name)) {
    try {
      await tagStore.submitForm(tagInput.value)
    } catch (error) {
      console.error('Error creating tag:', error)
      return
    }
  }

  formData.value.tags.push()
  tagInput = {

  }
  showTagSuggestions.value = false
}

const removeTag = (tagToRemove) => {
  formData.value.tags = formData.value.tags.filter(tag => tag.name !== tagToRemove.name)

}

const handleTagInput = (event) => {
  if (event.key === 'Enter' && tagInput.value.trim()) {
    event.preventDefault()
    addTag(tagInput.value.trim())
  } else if (event.key === 'Escape') {
    showTagSuggestions.value = false
  }
}

const getStatusColor = (status) => {
  const colors = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'in-progress': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800',
    'on-hold': 'bg-gray-100 text-gray-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString()
}

const isOverdue = (deadline, status) => {
  if (!deadline || status === 'completed') return false
  return new Date(deadline) < new Date()
}

const taskLearning = (id: number) => {
  router.push({ name: 'task learning', params: { topic_id: id } })
}
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Topics & Goals</h1>
      <button @click="openForm()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors">
        Add Topic
      </button>
    </div>

    <!-- Filters -->
    <div class="mb-6 space-y-4 md:space-y-0 md:flex md:gap-4">
      <div class="flex-1">
        <input v-model="searchQuery" type="text" placeholder="Search topics, descriptions, or tags..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
      </div>
      <div>
        <select v-model="statusFilter"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          <option value="all">All Status</option>
          <option v-for="status in statusOptions" :key="status.value" :value="status.value">
            {{ status.label }}
          </option>
        </select>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredTopics.length === 0" class="text-center py-12">
      <div class="text-gray-500 mb-4">
        <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No topics found</h3>
      <p class="text-gray-500 mb-4">
        {{ searchQuery || statusFilter !== 'all' ? 'Try adjusting your filters' : 'Get started by creating your first topic' }}
      </p>
      <button v-if="!searchQuery && statusFilter === 'all'" @click="openForm()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors">
        Add Your First Topic
      </button>
    </div>

    <!-- Topics Grid -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <div v-for="topic in filteredTopics" key="topic.topic_id" @click="taskLearning(topic.topic_id)"
        class="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-lg transition-shadow">
        <div class="flex justify-between items-start mb-3">
          <h3 class=" text-lg font-semibold text-gray-900 flex-1">{{ topic.title }}</h3>
          <div class="flex gap-2 ml-3">
            <!--
            <button @click="openForm(topic)" class="text-blue-600 hover:text-blue-800 text-sm">
              Edit
            </button>
-->
            <button @click="deleteTopic(topic.topic_id)" class="text-red-600 hover:text-red-800 text-sm">
              Delete
            </button>
          </div>
        </div>

        <p class="text-gray-600 mb-4 text-sm">{{ topic.description }}</p>

        <div class="flex items-center gap-2 mb-3">
          <span :class="['px-2 py-1 rounded-full text-xs font-medium', getStatusColor(topic.status)]">
            {{statusOptions.find(s => s.value === topic.status)?.label || topic.status}}
          </span>
          <span v-if="topic.deadline" :class="[
            'text-xs px-2 py-1 rounded',
            isOverdue(topic.deadline, topic.status) ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-600'
          ]">
            Due: {{ formatDate(topic.deadline) }}
          </span>
        </div>

        <div v-if="topic.tags.length > 0" class="flex flex-wrap gap-1">
          <span v-for="tag in topic.tags" key="tag" class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded">
            {{ tag }}
          </span>
        </div>
      </div>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-xl font-semibold text-gray-900">
              {{ editingTopic ? 'Edit Topic' : 'Add New Topic' }}
            </h2>
            <button @click="closeForm" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Title *</label>
              <input v-model="formData.title" type="text" required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter topic title">
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea v-model="formData.description" rows="4"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Describe your topic or goal"></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
                <select v-model="formData.status"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                  <option v-for="status in statusOptions" :key="status.value" :value="status.value">
                    {{ status.label }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Deadline</label>
                <input v-model="formData.deadline" type="date"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Tags</label>
              <div class="space-y-2">
                <!-- Selected Tags -->
                <div v-if="formData.tags.length > 0" class="flex flex-wrap gap-2">
                  <span v-for="tag in formData.tags" :key="tag"
                    class="inline-flex items-center px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full">
                    {{ tag }}
                    <button @click="removeTag(tag)" type="button" class="ml-2 text-blue-600 hover:text-blue-800">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </span>
                </div>

                <!-- Tag Input -->
                <div class="relative">
                  <input v-model="tagInput" @keydown="handleTagInput" @focus="showTagSuggestions = true"
                    @blur="setTimeout(() => showTagSuggestions = false, 200)" type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Type to add tags (press Enter to add)">

                  <!-- Tag Suggestions -->
                  <div v-if="showTagSuggestions && availableTags.length > 0"
                    class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-40 overflow-y-auto">
                    <button v-for="tag in availableTags.slice(0, 5)" key="tag" @click="addTag(tag)" type="button"
                      class="w-full px-3 py-2 text-left hover:bg-gray-50 text-sm">
                      {{ tag }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex justify-end gap-3 pt-4">
              <button @click="closeForm" type="button"
                class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg font-medium transition-colors">
                Cancel
              </button>
              <button type="submit"
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors">
                {{ editingTopic ? 'Update Topic' : 'Create Topic' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
