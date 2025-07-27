import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch'
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
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchTopics() {
    loading.value = true
    const { data, error } = await useMyFetch('topics').json<TopicResponse>()

    if (data.value) topics.value = data.value.topics
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: TopicForm) {
    const { data, error } = await useMyFetch('topics').post(formData).json()
    if (!data.value.topic_id) {
      data.value.topic_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Topic not added", "error")
    } else {
      fetchTopics()
      addToast("Topic added succesfully", "success")
    }
  }

  async function editTopic(topic: Topic) {
    const { error } = await useMyFetch(`topics/${topic.topic_id}`).patch(topic).json()

    if (!error.value) {
      const index = topics.value.findIndex(t => t.topic_id === topic.topic_id)
      if (index !== -1) {
        topics.value[index] = { ...topics.value[index], ...topic }
        addToast("Edited topic succesfully", "success")
      } else {
        addToast("Editing topic failed", "error")
      }
    } else {
      addToast("Editing topic failed", "error")
    }
  }

  async function deleteTopic(id: number) {
    await useMyFetch(`topics/${id}`).delete().json()
    topics.value = topics.value.filter((t) => t.topic_id !== id)
    addToast("Topic deleted succesfully", "success")
  }

  function reset() {
    topics.value = []
  }

  return {
    topics,
    loading,
    fetchError,
    fetchTopics,
    editTopic,
    deleteTopic,
    submitForm,
    reset,
  }
})

