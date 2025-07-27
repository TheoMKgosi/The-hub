import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch' // assumed custom fetch
import { useToast } from '@/composables/useToast'

interface Task {
  task_learning_id: number
  title: string
  notes: string
  status: 'pending' | 'in_progress' | 'done'
  resources: any[] // Replace with actual type if defined in tasks.ts
}

interface TaskLearningResponse {
  task_learning: Task[]
}

export const useTaskLearningStore = defineStore('tasks-learning', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const toast = useToast()

  const fetchTasks = async (topicId: number) => {
    loading.value = true
    try {
      const { data, error } = await useMyFetch(`/tasks/${topicId}`).get().json<TaskLearningResponse>()
      if (error.value) throw error.value
      tasks.value = data.value?.task_learning
    } catch (err) {
      toast.addToast('Failed to fetch tasks.')
    } finally {
      loading.value = false
    }
  }

  const createTask = async (newTask: Omit<Task, 'task_learning_id'>) => {
    try {
      const { data, error } = await useMyFetch('/tasks').post(newTask)
      if (error.value) throw error.value
      tasks.value.push(data.value.task)
      toast.addToast('Task created!')
    } catch (err) {
      toast.addToast('Failed to create task.')
    }
  }

  const updateTask = async (taskId: number, updates: Partial<Task>) => {
    try {
      const { data, error } = await useMyFetch(`/tasks/${taskId}`).put(updates)
      if (error.value) throw error.value
      const idx = tasks.value.findIndex(t => t.task_learning_id === taskId)
      if (idx !== -1) tasks.value[idx] = data.value.task
      toast.addToast('Task updated!')
    } catch (err) {
      toast.addToast('Failed to update task.')
    }
  }

  const deleteTask = async (taskId: number) => {
    try {
      const { error } = await useMyFetch(`/tasks/${taskId}`).delete()
      if (error.value) throw error.value
      tasks.value = tasks.value.filter(t => t.task_learning_id !== taskId)
      toast.addToast('Task deleted.')
    } catch (err) {
      toast.addToast('Failed to delete task.')
    }
  }

  return {
    tasks,
    loading,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
  }
})

