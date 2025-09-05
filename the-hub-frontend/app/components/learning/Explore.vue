<script setup lang="ts">
const topicStore = useTopicStore()
const publicTopics = ref([])
const searchQuery = ref('')
const loading = ref(false)

const fetchPublicTopics = async () => {
  loading.value = true
  try {
    const { $api } = useNuxtApp()
    const response = await $api<{ topics: any[] }>(`/topics/public?search=${searchQuery.value}`)
    publicTopics.value = response.topics
  } catch (error) {
    console.error('Failed to fetch public topics:', error)
  } finally {
    loading.value = false
  }
}

const copyTopic = async (topic: any) => {
  try {
    await topicStore.submitForm({
      title: `${topic.title} (Copy)`,
      description: topic.description,
      status: 'not_started',
      deadline: null,
      tags: []
    })
    alert('Topic copied to your collection!')
  } catch (error) {
    console.error('Failed to copy topic:', error)
    alert('Failed to copy topic')
  }
}

onMounted(() => {
  fetchPublicTopics()
})

watch(searchQuery, () => {
  fetchPublicTopics()
})
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-2">Explore Public Content</h1>
      <p class="text-text-light dark:text-text-dark/70">Discover and copy learning topics shared by the community</p>
    </div>

    <!-- Search -->
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search public topics..."
        class="w-full max-w-md px-4 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
      />
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
      <p class="text-text-light dark:text-text-dark/70">Loading public topics...</p>
    </div>

    <!-- Public Topics Grid -->
    <div v-else-if="publicTopics.length > 0" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="topic in publicTopics"
        :key="topic.topic_id"
        class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-md p-6 border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow"
      >
        <div class="flex justify-between items-start mb-4">
          <div class="flex-grow">
            <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">
              {{ topic.title }}
            </h3>
            <p v-if="topic.description" class="text-sm text-text-light dark:text-text-dark/70 mb-3">
              {{ topic.description }}
            </p>
            <div class="flex items-center gap-2 text-xs text-text-light dark:text-text-dark/60">
              <span>By {{ topic.user?.name || 'Anonymous' }}</span>
            </div>
          </div>
        </div>

        <div class="flex gap-2">
          <button
            @click="copyTopic(topic)"
            class="flex-1 px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors text-sm font-medium"
          >
            Copy to My Topics
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <div class="text-6xl mb-4">🔍</div>
      <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-2">No public topics found</h3>
      <p class="text-text-light dark:text-text-dark/60 mb-4">
        {{ searchQuery ? 'Try adjusting your search terms' : 'Be the first to share your learning topics!' }}
      </p>
    </div>
  </div>
</template>