import { defineStore } from 'pinia'
import { ref } from 'vue'
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
      const { $api } = useNuxtApp()
      const fetchedTasks = await $api<TaskLearningResponse>(`/task-learning/${topicId}`)
      tasks.value = fetchedTasks.task_learnings
    } catch (err) {
      toast.addToast('Failed to fetch tasks.')
    }
    loading.value = false
  }

  const createTask = async (topicId: number, payload: Omit<Task, 'task_learning_id' | 'notes' | 'status' | 'resources'>) => {
    try {
      const { $api } = useNuxtApp()
      await $api<Task>('/task-learning', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      fetchTasks(topicId)
      toast.addToast('Task created!')
    } catch (err) {
      toast.addToast('Failed to create task.')
    }
  }

  const updateTask = async (topicId: number, taskId: number, payload: Partial<Task>) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learning/${taskId}`, {
        method: 'PATCH',
        body: payload
      })

      fetchTasks(topicId)
      toast.addToast('Task updated!')
    } catch (err) {
      toast.addToast('Failed to update task.')
    }
  }

  const deleteTask = async (topicId: number, taskId: number) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learning/${taskId}`, {
        method: 'DELETE'
      })

      fetchTasks(topicId)
      toast.addToast('Task deleted.')
    } catch (err) {
      toast.addToast('Failed to delete task.')
    }
  }

  async function completeTask(topicId: number, task: Task) {
    try {
      const { $api } = useNuxtApp()
      await $api(`task-learning/${task.task_learning_id}`, {
        method: 'PATCH',
        body: JSON.stringify({ status: task.status })
      })
      fetchTasks(topicId)
      toast.addToast('Completed', 'success')
    } catch (err) {
      toast.addToast('Not completed', 'error')

    }
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

