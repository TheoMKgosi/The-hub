import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'



interface Topic {
  topic_id: number
  title: string
  description: string
  status: string
  deadline: Date | null
  tags: null[] | number[]
}

interface TopicForm {
  title: string
  description: string
  status: string
  deadline: Date | null
  tags: number[]
}

export interface TopicResponse {
  topics: Topic[]
}

export const useTopicStore = defineStore('topic', () => {
  const topics = ref<Topic[]>([])
  const topic = ref<Topic | null>(null)
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchTopics() {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedTopics = await $api<TopicResponse>('topics')

    if (fetchedTopics) topics.value = fetchedTopics.topics

    loading.value = false
  }

  async function fetchTopic(topicID: number) {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedTopic = await $api<Topic>(`topics/${topicID}`)

    if (fetchedTopic) topic.value = fetchedTopic

    loading.value = false
  }


  async function submitForm(payload: { title: string; description: string; status: string; deadline: Date | null; tags: number[] }) {
    // Create optimistic topic
    const optimisticTopic: Topic = {
      topic_id: -Date.now(), // Use negative ID to avoid conflicts
      title: payload.title,
      description: payload.description,
      status: payload.status,
      deadline: payload.deadline,
      tags: payload.tags
    }

    // Optimistically add to local state
    topics.value.push(optimisticTopic)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Topic>('topics', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic topic with real data
      const optimisticIndex = topics.value.findIndex(t => t.topic_id === optimisticTopic.topic_id)
      if (optimisticIndex !== -1) {
        topics.value[optimisticIndex] = data
      }

      addToast("Topic added succesfully", "success")

    } catch (err) {
      // Remove optimistic topic on error
      topics.value = topics.value.filter(t => t.topic_id !== optimisticTopic.topic_id)
      addToast("Topic not added", "error")
    }
  }

  async function editTopic(topic: Topic) {
    // Store original topic for potential rollback
    const originalTopicIndex = topics.value.findIndex(t => t.topic_id === topic.topic_id)
    const originalTopic = originalTopicIndex !== -1 ? { ...topics.value[originalTopicIndex] } : null

    // Optimistically update the topic
    if (originalTopicIndex !== -1) {
      topics.value[originalTopicIndex] = { ...topics.value[originalTopicIndex], ...topic }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Topic>(`topics/${topic.topic_id}`, {
        method: 'PATCH',
        body: JSON.stringify(topic)
      })

      // Update with server response to ensure consistency
      if (originalTopicIndex !== -1 && data) {
        topics.value[originalTopicIndex] = data
      }

      addToast("Edited topic succesfully", "success")

    } catch (err) {
      // Revert optimistic update on error
      if (originalTopic && originalTopicIndex !== -1) {
        topics.value[originalTopicIndex] = originalTopic
      }
      addToast("Editing topic failed", "error")
    }
  }

  async function deleteTopic(id: number) {
    // Store the topic for potential rollback
    const topicToDelete = topics.value.find(t => t.topic_id === id)
    if (!topicToDelete) {
      addToast("Topic not found", "error")
      return
    }

    // Optimistically remove from local state
    topics.value = topics.value.filter((t) => t.topic_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`topics/${id}`, {
        method: 'DELETE'
      })

      addToast("Topic deleted succesfully", "success")

    } catch(err) {
      // Restore the topic on error
      topics.value.push(topicToDelete)
      addToast("Topic did not delete", "error")
    }
  }

  function reset() {
    topics.value = []
  }

  return {
    topics,
    topic,
    loading,
    fetchError,
    fetchTopics,
    fetchTopic,
    editTopic,
    deleteTopic,
    submitForm,
    reset,
  }
})

