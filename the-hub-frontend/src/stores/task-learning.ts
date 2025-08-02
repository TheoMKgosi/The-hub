import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch' // assumed custom fetch
import { useToast } from '@/composables/useToast'

interface Task {
  task_learning_id: number
  topic_id: number
  title: string
  notes: string
  status: 'pending' | 'in_progress' | 'done'
  resources: any[] // Replace with actual type if defined in tasks.ts
}

interface TaskLearningResponse {
  task_learnings: Task[]
}

export const useTaskLearningStore = defineStore('tasks-learning', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const toast = useToast()

  const fetchTasks = async (topicId: number) => {
    loading.value = true
    try {
      const { data, error } = await useMyFetch(`/task-learning/${topicId}`).get().json<TaskLearningResponse>()
      if (error.value) throw error.value
      tasks.value = data.value?.task_learnings
    } catch (err) {
      toast.addToast('Failed to fetch tasks.')
    } finally {
      loading.value = false
    }
  }

  const createTask = async (newTask: Omit<Task, 'task_learning_id' | 'notes' | 'status' | 'resources'>) => {
    try {
      const { error } = await useMyFetch('/task-learning').post(newTask)
      if (error.value) throw error.value
      fetchTasks(newTask.topic_id)
      toast.addToast('Task created!')
    } catch (err) {
      toast.addToast('Failed to create task.')
    }
  }

  const updateTask = async (taskId: number, updates: Partial<Task>) => {
    try {
      const { data, error } = await useMyFetch(`/task-learning/${taskId}`).patch(updates)
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
      const { error } = await useMyFetch(`/task-learning/${taskId}`).delete()
      if (error.value) throw error.value
      tasks.value = tasks.value.filter(t => t.task_learning_id !== taskId)
      toast.addToast('Task deleted.')
    } catch (err) {
      toast.addToast('Failed to delete task.')
    }
  }

  async function completeTask(task: Task) {
    await useMyFetch(`task-learning/${task.task_learning_id}`).patch({ status: task.status }).json()
  }

  return {
    tasks,
    loading,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    completeTask,
  }
})

