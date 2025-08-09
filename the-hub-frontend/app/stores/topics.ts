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
    try {
      const { $api } = useNuxtApp()
      await $api('topics', {
        method: 'POST',
        body: payload,
      })
      // TODO: Change implementation to array push to avoid many calls
      fetchTopics()
      addToast("Topic added succesfully", "success")

    } catch (err) {
      addToast("Topic not added", "error")
    }
  }

  async function editTopic(topic: Topic) {
    try {
      const { $api } = useNuxtApp()
      await $api(`topics/${topic.topic_id}`, {
        method: 'PATCH',
        body: topic
      })
      // TODO: Change implementation to filter array to avoid many calls
      fetchTopics()
      addToast("Edited topic succesfully", "success")

    } catch (err) {
      addToast("Editing topic failed", "error")

    }
  }

  async function deleteTopic(id: number) {
    try {
    const { $api } = useNuxtApp()
    await $api(`topics/${id}`, {
      method: 'DELETE'
    })
    topics.value = topics.value.filter((t) => t.topic_id !== id)
    addToast("Topic deleted succesfully", "success")

    } catch(err) {
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

