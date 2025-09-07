import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

interface Resource {
  id: string
  topic_id?: string
  task_id?: string
  title: string
  link: string
  type: 'video' | 'article' | 'document' | 'book' | 'course'
  notes: string
}

export const useResourceStore = defineStore('resources', () => {
  const resources = ref<Resource[]>([])
  const loading = ref(false)
  const { addToast } = useToast()

  const createResource = async (data: {
    topic_id?: string
    task_id?: string
    title: string
    link: string
    type: Resource['type']
    notes?: string
  }) => {
    // Create optimistic resource
    const optimisticResource: Resource = {
      id: `temp-${Date.now()}`,
      topic_id: data.topic_id,
      task_id: data.task_id,
      title: data.title,
      link: data.link,
      type: data.type,
      notes: data.notes || ''
    }

    // Optimistically add to local state
    resources.value.unshift(optimisticResource)

    try {
      const { $api } = useNuxtApp()
      const newResource = await $api<Resource>('/resources', {
        method: 'POST',
        body: JSON.stringify(data)
      })

      // Replace optimistic resource with real data
      const optimisticIndex = resources.value.findIndex(r => r.id === optimisticResource.id)
      if (optimisticIndex !== -1) {
        resources.value[optimisticIndex] = newResource
      }

      addToast('Resource added successfully!', 'success')
      return newResource
    } catch (err) {
      // Remove optimistic resource on error
      resources.value = resources.value.filter(r => r.id !== optimisticResource.id)
      addToast('Failed to add resource', 'error')
      console.error('Error creating resource:', err)
    }
  }

  const fetchResources = async (filters?: {
    topic_id?: string
    task_id?: string
    type?: Resource['type']
    search?: string
  }) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const queryParams = new URLSearchParams()

      if (filters?.topic_id) queryParams.append('topic_id', filters.topic_id)
      if (filters?.task_id) queryParams.append('task_id', filters.task_id)
      if (filters?.type) queryParams.append('type', filters.type)
      if (filters?.search) queryParams.append('search', filters.search)

      const response = await $api<{ resources: Resource[] }>(`/resources?${queryParams}`)
      resources.value = response.resources
    } catch (err) {
      addToast('Failed to fetch resources', 'error')
      console.error('Error fetching resources:', err)
    } finally {
      loading.value = false
    }
  }

  const updateResource = async (id: string, data: Partial<Resource>) => {
    // Store original resource for potential rollback
    const originalResourceIndex = resources.value.findIndex(r => r.id === id)
    const originalResource = originalResourceIndex !== -1 ? { ...resources.value[originalResourceIndex] } : null

    // Optimistically update the resource
    if (originalResourceIndex !== -1) {
      resources.value[originalResourceIndex] = { ...resources.value[originalResourceIndex], ...data }
    }

    try {
      const { $api } = useNuxtApp()
      const updatedResource = await $api<Resource>(`/resources/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
      })

      // Update with server response to ensure consistency
      if (originalResourceIndex !== -1 && updatedResource) {
        resources.value[originalResourceIndex] = updatedResource
      }

      addToast('Resource updated successfully!', 'success')
      return updatedResource
    } catch (err) {
      // Revert optimistic update on error
      if (originalResource && originalResourceIndex !== -1) {
        resources.value[originalResourceIndex] = originalResource
      }
      addToast('Failed to update resource', 'error')
      console.error('Error updating resource:', err)
    }
  }

  const deleteResource = async (id: string) => {
    // Store the resource for potential rollback
    const resourceToDelete = resources.value.find(r => r.id === id)
    if (!resourceToDelete) {
      addToast("Resource not found", "error")
      return
    }

    // Optimistically remove from local state
    resources.value = resources.value.filter(r => r.id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`/resources/${id}`, {
        method: 'DELETE'
      })

      addToast('Resource deleted successfully!', 'success')
    } catch (err) {
      // Restore the resource on error
      resources.value.push(resourceToDelete)
      addToast('Failed to delete resource', 'error')
      console.error('Error deleting resource:', err)
    }
  }

  const getResourcesByTopic = (topicId: string) => {
    return resources.value.filter(r => r.topic_id === topicId)
  }

  const getResourcesByTask = (taskId: string) => {
    return resources.value.filter(r => r.task_id === taskId)
  }

  const getResourceTypeIcon = (type: Resource['type']) => {
    const icons = {
      video: '🎥',
      article: '📄',
      document: '📋',
      book: '📚',
      course: '🎓'
    }
    return icons[type] || '📖'
  }

  const getResourceTypeColor = (type: Resource['type']) => {
    const colors = {
      video: 'text-red-600 dark:text-red-400',
      article: 'text-blue-600 dark:text-blue-400',
      document: 'text-green-600 dark:text-green-400',
      book: 'text-purple-600 dark:text-purple-400',
      course: 'text-orange-600 dark:text-orange-400'
    }
    return colors[type] || 'text-gray-600 dark:text-gray-400'
  }

  return {
    resources,
    loading,
    createResource,
    fetchResources,
    updateResource,
    deleteResource,
    getResourcesByTopic,
    getResourcesByTask,
    getResourceTypeIcon,
    getResourceTypeColor
  }
})