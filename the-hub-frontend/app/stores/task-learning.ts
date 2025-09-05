import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'

interface Task {
  task_learning_id: string
  topic_id: string
  title: string
  notes: string
  status: 'not_started' | 'in_progress' | 'completed' | 'on_hold'
  resources: any[] // Replace with actual type if defined in tasks.ts
}

interface TaskLearningResponse {
  task_learnings: Task[]
}

export const useTaskLearningStore = defineStore('tasks-learning', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const toast = useToast()

  const fetchTasks = async (topicId: string) => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const fetchedTasks = await $api<TaskLearningResponse>(`/topics/${topicId}/task-learnings`)
      tasks.value = fetchedTasks.task_learnings
    } catch (err) {
      toast.addToast('Failed to fetch tasks.')
    }
    loading.value = false
  }

  const createTask = async (topicId: string, payload: Omit<Task, 'task_learning_id' | 'notes' | 'status' | 'resources'>) => {
    try {
      const { $api } = useNuxtApp()
      await $api<Task>('/task-learnings', {
        method: 'POST',
        body: JSON.stringify({ ...payload, topic_id: topicId })
      })
      fetchTasks(topicId)
      toast.addToast('Task created!')
    } catch (err) {
      toast.addToast('Failed to create task.')
    }
  }

  const updateTask = async (topicId: string, taskId: string, payload: Partial<Task>) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learnings/${taskId}`, {
        method: 'PUT',
        body: JSON.stringify(payload)
      })

      fetchTasks(topicId)
      toast.addToast('Task updated!')
    } catch (err) {
      toast.addToast('Failed to update task.')
    }
  }

  const deleteTask = async (topicId: string, taskId: string) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learnings/${taskId}`, {
        method: 'DELETE'
      })

      fetchTasks(topicId)
      toast.addToast('Task deleted.')
    } catch (err) {
      toast.addToast('Failed to delete task.')
    }
  }

  async function completeTask(topicId: string, task: Task) {
    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learnings/${task.task_learning_id}`, {
        method: 'PUT',
        body: JSON.stringify({ status: task.status })
      })
      fetchTasks(topicId)
      toast.addToast('Task status updated!', 'success')
    } catch (err) {
      toast.addToast('Failed to update task status', 'error')
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

