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
    // Create optimistic task
    const optimisticTask: Task = {
      task_learning_id: `temp-${Date.now()}`,
      topic_id: topicId,
      title: payload.title,
      notes: '',
      status: 'not_started',
      resources: []
    }

    // Optimistically add to local state
    tasks.value.push(optimisticTask)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>('/task-learnings', {
        method: 'POST',
        body: JSON.stringify({ ...payload, topic_id: topicId })
      })

      // Replace optimistic task with real data
      const optimisticIndex = tasks.value.findIndex(t => t.task_learning_id === optimisticTask.task_learning_id)
      if (optimisticIndex !== -1) {
        tasks.value[optimisticIndex] = data
      }

      toast.addToast('Task created!')
    } catch (err) {
      // Remove optimistic task on error
      tasks.value = tasks.value.filter(t => t.task_learning_id !== optimisticTask.task_learning_id)
      toast.addToast('Failed to create task.')
    }
  }

  const updateTask = async (topicId: string, taskId: string, payload: Partial<Task>) => {
    // Store original task for potential rollback
    const originalTaskIndex = tasks.value.findIndex(t => t.task_learning_id === taskId)
    const originalTask = originalTaskIndex !== -1 ? { ...tasks.value[originalTaskIndex] } : null

    // Optimistically update the task
    if (originalTaskIndex !== -1) {
      tasks.value[originalTaskIndex] = { ...tasks.value[originalTaskIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/task-learnings/${taskId}`, {
        method: 'PUT',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalTaskIndex !== -1 && data) {
        tasks.value[originalTaskIndex] = data
      }

      toast.addToast('Task updated!')
    } catch (err) {
      // Revert optimistic update on error
      if (originalTask && originalTaskIndex !== -1) {
        tasks.value[originalTaskIndex] = originalTask
      }
      toast.addToast('Failed to update task.')
    }
  }

  const deleteTask = async (topicId: string, taskId: string) => {
    // Store the task for potential rollback
    const taskToDelete = tasks.value.find(t => t.task_learning_id === taskId)
    if (!taskToDelete) {
      toast.addToast('Task not found.')
      return
    }

    // Optimistically remove from local state
    tasks.value = tasks.value.filter(t => t.task_learning_id !== taskId)

    try {
      const { $api } = useNuxtApp()
      await $api(`/task-learnings/${taskId}`, {
        method: 'DELETE'
      })

      toast.addToast('Task deleted.')
    } catch (err) {
      // Restore the task on error
      tasks.value.push(taskToDelete)
      toast.addToast('Failed to delete task.')
    }
  }

  async function completeTask(topicId: string, task: Task) {
    // Store original task for potential rollback
    const originalTaskIndex = tasks.value.findIndex(t => t.task_learning_id === task.task_learning_id)
    const originalTask = originalTaskIndex !== -1 ? { ...tasks.value[originalTaskIndex] } : null

    // Optimistically update the task status
    if (originalTaskIndex !== -1) {
      tasks.value[originalTaskIndex] = { ...tasks.value[originalTaskIndex], status: task.status }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Task>(`/task-learnings/${task.task_learning_id}`, {
        method: 'PUT',
        body: JSON.stringify({ status: task.status })
      })

      // Update with server response to ensure consistency
      if (originalTaskIndex !== -1 && data) {
        tasks.value[originalTaskIndex] = data
      }

      toast.addToast('Task status updated!', 'success')
    } catch (err) {
      // Revert optimistic update on error
      if (originalTask && originalTaskIndex !== -1) {
        tasks.value[originalTaskIndex] = originalTask
      }
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

