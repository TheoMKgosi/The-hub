<script setup lang="ts">
import CrossIcon from '../ui/svg/CrossIcon.vue'
const topicStore = useTopicStore()
const tagStore = useTagStore()

const router = useRouter()

// Form state
const showForm = ref(false)
const editingTopic = ref(null)
const formData = ref({
  title: '',
  description: '',
  status: 'not_started',
  deadline: new Date(),
  is_public: false,
  tags: [] as string[]
})
const editFormData = ref({
  topic_id: '',
  title: '',
  description: '',
  status: 'not_started',
  deadline: new Date(),
  is_public: false,
  tags: [] as string[]
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
  { value: 'not_started', label: 'Not Started' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' },
  { value: 'on_hold', label: 'On Hold' }
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
    !formData.value.tags.includes(tag.id)
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
    status: 'not_started',
    deadline: new Date(),
    is_public: false,
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
      await topicStore.editTopic({
        topic_id: editingTopic.value.topic_id,
        title: formData.value.title,
        description: formData.value.description,
        status: formData.value.status,
        deadline: formData.value.deadline,
        tags: formData.value.tags
      })
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
const addTag = async (tag) => {
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
  formData.value.tags.push(tag.id)
  showTagSuggestions.value = false
}

const removeTag = (tagToRemove) => {
  formData.value.tags = formData.value.tags.filter(tag => tag !== tagToRemove)
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
    'not_started': 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning',
    'in_progress': 'bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary',
    'completed': 'bg-success/10 dark:bg-success/20 text-success dark:text-success',
    'on_hold': 'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark'
  }
  return colors[status] || 'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark'
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
  router.push(`/learning/${id}`)
}
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-text-light dark:text-text-dark">Topics & Goals</h1>
      <BaseButton @click="openForm()" text="Add Topic" variant="primary" size="md" />
    </div>

    <!-- Filters -->
    <div class="mb-6 space-y-4 md:space-y-0 md:flex md:gap-4">
      <div class="flex-1">
        <input v-model="searchQuery" type="text" placeholder="Search topics, descriptions, or tags..."
          class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50">
      </div>
      <div>
        <select v-model="statusFilter"
          class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
          <option value="all">All Status</option>
          <option v-for="status in statusOptions" :key="status.value" :value="status.value">
            {{ status.label }}
          </option>
        </select>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredTopics.length === 0" class="text-center py-12">
      <div class="text-text-light dark:text-text-dark/60 mb-4">
        <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">No topics found</h3>
      <p class="text-text-light dark:text-text-dark/60 mb-4">
        {{ searchQuery || statusFilter !== 'all' ? 'Try adjusting your filters' : 'Get started by creating your first topic' }}
      </p>
      <BaseButton v-if="!searchQuery && statusFilter === 'all'" @click="openForm()" text="Add Your First Topic"
        variant="primary" size="md" />
    </div>

    <!-- Topics Grid -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <div v-for="topic in filteredTopics" :key="topic.topic_id" @click="taskLearning(topic.topic_id)"
        class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md p-6 border border-surface-light dark:border-surface-dark hover:shadow-lg transition-all duration-200 cursor-pointer group">
        <div class="flex justify-between items-start mb-3">
          <h3
            class="text-lg font-semibold text-text-light dark:text-text-dark flex-1 group-hover:text-primary dark:group-hover:text-primary transition-colors">
            {{ topic.title }}
          </h3>
          <div class="flex gap-2 ml-3 opacity-0 group-hover:opacity-100 transition-opacity">
            <BaseButton @click.stop="deleteTopic(topic.topic_id)" text="Delete" variant="danger" size="sm" />
          </div>
        </div>

        <p class="text-text-light dark:text-text-dark/80 mb-4 text-sm">{{ topic.description }}</p>

        <div class="flex items-center gap-2 mb-3">
          <span :class="['px-2 py-1 rounded-full text-xs font-medium', getStatusColor(topic.status)]">
            {{statusOptions.find(s => s.value === topic.status)?.label || topic.status}}
          </span>
          <span v-if="topic.deadline" :class="[
            'text-xs px-2 py-1 rounded',
            isOverdue(topic.deadline, topic.status) ? 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-300' : 'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark'
          ]">
            Due: {{ formatDate(topic.deadline) }}
          </span>
        </div>

        <div v-if="topic.tags.length > 0" class="flex flex-wrap gap-1">
          <span v-for="tag in topic.tags" :key="tag"
            class="px-2 py-1 bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary text-xs rounded">
            {{ tag }}
          </span>
        </div>
      </div>
    </div>

    <!-- Form Modal -->
    <div v-if="showForm" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50">
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-surface-light dark:border-surface-dark">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">
              {{ editingTopic ? 'Edit Topic' : 'Add New Topic' }}
            </h2>
            <BaseButton @click="closeForm" :iconOnly="true" :icon="CrossIcon" variant="default" size="sm" class="p-2" />
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Title *</label>
              <input v-model="formData.title" type="text" required
                class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                placeholder="Enter topic title">
            </div>

            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Description</label>
              <textarea v-model="formData.description" rows="4"
                class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                placeholder="Describe your topic or goal"></textarea>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Status</label>
                <select v-model="formData.status"
                  class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                  <option v-for="status in statusOptions" :key="status.value" :value="status.value">
                    {{ status.label }}
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Deadline</label>
                <input v-model="formData.deadline" type="date"
                  class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
              </div>
            </div>

            <div class="flex items-center gap-2">
              <input v-model="formData.is_public" type="checkbox"
                class="w-4 h-4 text-primary bg-surface-light dark:bg-surface-dark border-surface-light dark:border-surface-dark rounded focus:ring-primary" />
              <label class="text-sm font-medium text-text-light dark:text-text-dark">
                Make this topic public (visible to other users)
              </label>
            </div>

            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Tags</label>
              <div class="space-y-2">
                <!-- Selected Tags -->
                <div v-if="formData.tags.length > 0" class="flex flex-wrap gap-2">
                  <span v-for="tag in formData.tags" :key="tag"
                    class="inline-flex items-center px-3 py-1 bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary text-sm rounded-full">
                    {{ tag }}
                    <BaseButton @click="removeTag(tag)" :icon="CrossIcon" :iconOnly="true" variant="default" size="sm"
                      class="ml-2 p-1" />
                  </span>
                </div>

                <!-- Tag Input -->
                <div class="relative">
                  <input v-model="tagInput.name" @keydown="handleTagInput" @focus="showTagSuggestions = true"
                    @blur="setTimeout(() => showTagSuggestions = false, 200)" type="text"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                    placeholder="Type to add tags (press Enter to add)">

                  <!-- Tag Suggestions -->
                  <div v-if="showTagSuggestions && availableTags.length > 0"
                    class="absolute z-10 w-full mt-1 bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark rounded-md shadow-lg max-h-40 overflow-y-auto">
                    <button v-for="tag in availableTags.slice(0, 5)" :key="tag" @click="addTag(tag)" type="button"
                      class="w-full px-3 py-2 text-left hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 text-text-light dark:text-text-dark text-sm transition-colors">
                      {{ tag.name }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex justify-end gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
              <BaseButton @click="closeForm" text="Cancel" variant="default" size="md" />
              <BaseButton type="submit" :text="editingTopic ? 'Update Topic' : 'Create Topic'" variant="primary"
                size="md" />
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
