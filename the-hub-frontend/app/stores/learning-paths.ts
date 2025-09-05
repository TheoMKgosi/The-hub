import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

interface Topic {
  topic_id: string
  title: string
  description: string
  status: string
}

interface LearningPath {
  learning_path_id: string
  user_id: string
  title: string
  description: string
  topics: Topic[]
}

export const useLearningPathStore = defineStore('learning-paths', () => {
  const learningPaths = ref<LearningPath[]>([])
  const loading = ref(false)
  const { addToast } = useToast()

  const createLearningPath = async (data: {
    title: string
    description?: string
    topic_ids: string[]
  }) => {
    try {
      const { $api } = useNuxtApp()
      const newPath = await $api<LearningPath>('/learning-paths', {
        method: 'POST',
        body: JSON.stringify(data)
      })

      if (newPath) {
        learningPaths.value.unshift(newPath)
        addToast('Learning path created successfully!', 'success')
        return newPath
      }
    } catch (err) {
      addToast('Failed to create learning path', 'error')
      console.error('Error creating learning path:', err)
    }
  }

  const fetchLearningPaths = async () => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<{ learning_paths: LearningPath[] }>('/learning-paths')
      learningPaths.value = response.learning_paths
    } catch (err) {
      addToast('Failed to fetch learning paths', 'error')
      console.error('Error fetching learning paths:', err)
    } finally {
      loading.value = false
    }
  }

  const deleteLearningPath = async (id: string) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/learning-paths/${id}`, {
        method: 'DELETE'
      })

      learningPaths.value = learningPaths.value.filter(lp => lp.learning_path_id !== id)
      addToast('Learning path deleted successfully!', 'success')
    } catch (err) {
      addToast('Failed to delete learning path', 'error')
      console.error('Error deleting learning path:', err)
    }
  }

  const getLearningPathProgress = (learningPath: LearningPath) => {
    const completedTopics = learningPath.topics.filter(topic => topic.status === 'completed').length
    return {
      completed: completedTopics,
      total: learningPath.topics.length,
      percentage: Math.round((completedTopics / learningPath.topics.length) * 100)
    }
  }

  return {
    learningPaths,
    loading,
    createLearningPath,
    fetchLearningPaths,
    deleteLearningPath,
    getLearningPathProgress
  }
})