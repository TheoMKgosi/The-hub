<script setup lang="ts">
const learningPathStore = useLearningPathStore()
const topicStore = useTopicStore()

const showCreateModal = ref(false)
const createFormData = reactive({
  title: '',
  description: '',
  selectedTopicIds: [] as string[]
})

onMounted(() => {
  learningPathStore.fetchLearningPaths()
  topicStore.fetchTopics()
})

const openCreateModal = () => {
  createFormData.title = ''
  createFormData.description = ''
  createFormData.selectedTopicIds = []
  showCreateModal.value = true
}

const closeCreateModal = () => {
  showCreateModal.value = false
}

const createLearningPath = async () => {
  if (!createFormData.title.trim() || createFormData.selectedTopicIds.length === 0) return

  try {
    await learningPathStore.createLearningPath({
      title: createFormData.title,
      description: createFormData.description,
      topic_ids: createFormData.selectedTopicIds
    })
    closeCreateModal()
  } catch (error) {
    console.error('Failed to create learning path:', error)
  }
}

const deleteLearningPath = async (pathId: string) => {
  if (confirm('Are you sure you want to delete this learning path?')) {
    try {
      await learningPathStore.deleteLearningPath(pathId)
    } catch (error) {
      console.error('Failed to delete learning path:', error)
    }
  }
}

const toggleTopicSelection = (topicId: string) => {
  const index = createFormData.selectedTopicIds.indexOf(topicId)
  if (index > -1) {
    createFormData.selectedTopicIds.splice(index, 1)
  } else {
    createFormData.selectedTopicIds.push(topicId)
  }
}

const isTopicSelected = (topicId: string) => {
  return createFormData.selectedTopicIds.includes(topicId)
}
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-text-light dark:text-text-dark">Learning Paths</h1>
        <p class="text-text-light dark:text-text-dark/70 mt-1">Create structured learning journeys with sequenced topics</p>
      </div>
      <UiButton @click="openCreateModal" variant="primary" size="md">
        Create Learning Path
      </UiButton>
    </div>

    <!-- Learning Paths Grid -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <template v-for="path in learningPathStore.learningPaths" :key="path.learning_path_id">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md p-6 border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow">
          <div class="flex justify-between items-start mb-4">
            <div class="flex-grow">
              <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">
                {{ path.title }}
              </h3>
              <p v-if="path.description" class="text-sm text-text-light dark:text-text-dark/70 mb-3">
                {{ path.description }}
              </p>
            </div>
            <UiButton @click="deleteLearningPath(path.learning_path_id)" variant="danger" size="sm">
              Delete
            </UiButton>
          </div>

          <!-- Progress -->
          <div class="mb-4">
            <div class="flex justify-between text-sm mb-1">
              <span class="text-text-light dark:text-text-dark/70">Progress</span>
              <span class="font-medium">{{ learningPathStore.getLearningPathProgress(path).percentage }}%</span>
            </div>
            <div class="w-full bg-surface-light dark:bg-surface-dark h-2 rounded-full overflow-hidden">
              <div
                class="bg-primary h-full rounded-full transition-all duration-500"
                :style="{ width: `${learningPathStore.getLearningPathProgress(path).percentage}%` }"
              ></div>
            </div>
            <div class="text-xs text-text-light dark:text-text-dark/60 mt-1">
              {{ learningPathStore.getLearningPathProgress(path).completed }} of {{ learningPathStore.getLearningPathProgress(path).total }} topics completed
            </div>
          </div>

          <!-- Topics -->
          <div class="space-y-2">
            <h4 class="text-sm font-medium text-text-light dark:text-text-dark/80">Topics in this path:</h4>
            <div class="space-y-1">
              <template v-for="(topic, index) in path.topics" :key="topic.topic_id">
                <div class="flex items-center gap-2 text-sm">
                  <span class="flex-shrink-0 w-5 h-5 bg-primary/10 dark:bg-primary/20 text-primary rounded-full flex items-center justify-center text-xs font-medium">
                    {{ index + 1 }}
                  </span>
                  <span :class="[
                    'flex-grow',
                    topic.status === 'completed' ? 'line-through text-text-light dark:text-text-dark/60' : 'text-text-light dark:text-text-dark'
                  ]">
                    {{ topic.title }}
                  </span>
                  <span :class="[
                    'px-2 py-0.5 rounded-full text-xs',
                    topic.status === 'completed' ? 'bg-success/10 dark:bg-success/20 text-success' :
                    topic.status === 'in_progress' ? 'bg-primary/10 dark:bg-primary/20 text-primary' :
                    'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark/60'
                  ]">
                    {{ topic.status.replace('_', ' ') }}
                  </span>
                </div>
              </template>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- Empty State -->
    <div v-if="learningPathStore.learningPaths.length === 0 && !learningPathStore.loading" class="text-center py-12">
      <div class="text-6xl mb-4">🛤️</div>
      <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No learning paths yet</h3>
      <p class="text-text-light dark:text-text-dark/60 mb-4">
        Create your first learning path to organize your topics into structured learning journeys
      </p>
      <UiButton @click="openCreateModal" variant="primary" size="lg">
        Create Your First Learning Path
      </UiButton>
    </div>

    <!-- Create Learning Path Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto border border-surface-light dark:border-surface-dark">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">
              Create Learning Path
            </h2>
            <UiButton @click="closeCreateModal" variant="default" size="sm" class="p-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </UiButton>
          </div>

          <form @submit.prevent="createLearningPath" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Title *</label>
              <input
                v-model="createFormData.title"
                type="text"
                required
                class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                placeholder="e.g., Web Development Fundamentals"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Description</label>
              <textarea
                v-model="createFormData.description"
                rows="3"
                class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none"
                placeholder="Describe your learning path..."
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Select Topics *</label>
              <p class="text-xs text-text-light dark:text-text-dark/60 mb-3">Choose topics to include in this learning path (drag to reorder)</p>

              <div class="space-y-2 max-h-60 overflow-y-auto">
                <template v-for="topic in topicStore.topics" :key="topic.topic_id">
                  <div
                    class="flex items-center gap-3 p-3 border border-surface-light dark:border-surface-dark rounded-md hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 cursor-pointer transition-colors"
                    @click="toggleTopicSelection(topic.topic_id)"
                  >
                    <input
                      type="checkbox"
                      :checked="isTopicSelected(topic.topic_id)"
                      class="flex-shrink-0"
                      @click.stop
                    />
                    <div class="flex-grow min-w-0">
                      <h4 class="font-medium text-text-light dark:text-text-dark truncate">
                        {{ topic.title }}
                      </h4>
                      <p v-if="topic.description" class="text-sm text-text-light dark:text-text-dark/70 truncate">
                        {{ topic.description }}
                      </p>
                    </div>
                    <span :class="[
                      'px-2 py-1 rounded-full text-xs font-medium flex-shrink-0',
                      topic.status === 'completed' ? 'bg-success/10 dark:bg-success/20 text-success' :
                      topic.status === 'in_progress' ? 'bg-primary/10 dark:bg-primary/20 text-primary' :
                      'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark/60'
                    ]">
                      {{ topic.status.replace('_', ' ') }}
                    </span>
                  </div>
                </template>
              </div>

              <div v-if="createFormData.selectedTopicIds.length > 0" class="mt-4 p-3 bg-primary/5 dark:bg-primary/10 rounded-md">
                <p class="text-sm font-medium text-primary mb-2">
                  Selected Topics ({{ createFormData.selectedTopicIds.length }})
                </p>
                <div class="flex flex-wrap gap-2">
                  <template v-for="topicId in createFormData.selectedTopicIds" :key="topicId">
                    <span class="inline-flex items-center gap-1 px-2 py-1 bg-primary/10 dark:bg-primary/20 text-primary text-xs rounded-full">
                      {{ topicStore.topics.find(t => t.topic_id === topicId)?.title }}
                      <button
                        @click.stop="toggleTopicSelection(topicId)"
                        class="hover:bg-primary/20 rounded-full p-0.5"
                      >
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </span>
                  </template>
                </div>
              </div>
            </div>

            <div class="flex justify-end gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
              <UiButton @click="closeCreateModal" variant="default" size="md">
                Cancel
              </UiButton>
              <UiButton
                type="submit"
                variant="primary"
                size="md"
                :disabled="!createFormData.title.trim() || createFormData.selectedTopicIds.length === 0"
              >
                Create Learning Path
              </UiButton>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>