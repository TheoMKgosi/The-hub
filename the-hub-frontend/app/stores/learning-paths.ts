import { defineStore } from 'pinia'
import { ref } from 'vue'

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
  const toast = useToast()

  const createLearningPath = async (data: {
    title: string
    description?: string
    topic_ids: string[]
  }) => {
    // Create optimistic learning path
    const optimisticPath: LearningPath = {
      learning_path_id: `temp-${Date.now()}`,
      user_id: '', // Will be set by server
      title: data.title,
      description: data.description || '',
      topics: [] // Will be populated by server based on topic_ids
    }

    // Optimistically add to local state
    learningPaths.value.unshift(optimisticPath)

    try {
      const { $api } = useNuxtApp()
      const newPath = await $api<LearningPath>('/learning-paths', {
        method: 'POST',
        body: JSON.stringify(data)
      })

      // Replace optimistic path with real data
      const optimisticIndex = learningPaths.value.findIndex(lp => lp.learning_path_id === optimisticPath.learning_path_id)
      if (optimisticIndex !== -1) {
        learningPaths.value[optimisticIndex] = newPath
      }

      toast.add({ title: 'Learning path created successfully!', color: 'success' })
      return newPath
    } catch (err) {
      // Remove optimistic path on error
      learningPaths.value = learningPaths.value.filter(lp => lp.learning_path_id !== optimisticPath.learning_path_id)
      toast.add({ title: 'Failed to create learning path', color: 'error' })
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
      toast.add({ title: 'Failed to fetch learning paths', color: 'error' })
      console.error('Error fetching learning paths:', err)
    } finally {
      loading.value = false
    }
  }

  const deleteLearningPath = async (id: string) => {
    // Store the learning path for potential rollback
    const pathToDelete = learningPaths.value.find(lp => lp.learning_path_id === id)
    if (!pathToDelete) {
      toast.add({ title: "Learning path not found", color: 'error' })
      return
    }

    // Optimistically remove from local state
    learningPaths.value = learningPaths.value.filter(lp => lp.learning_path_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`/learning-paths/${id}`, {
        method: 'DELETE'
      })

      toast.add({ title: 'Learning path deleted successfully!', color: 'success' })
    } catch (err) {
      // Restore the learning path on error
      learningPaths.value.push(pathToDelete)
      toast.add({ title: 'Failed to delete learning path', color: 'error' })
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
