import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch'


interface Task {
  task_id: number
  title: string
  description: string
  due_date?: string
  priority: number
  status: string
}

export interface TaskResponse {
  tasks: Task[]
}

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchTasks() {
    loading.value = true
    const { data, error } = await useMyFetch('tasks').json<TaskResponse>()

    if (data.value) tasks.value = data.value.tasks
    fetchError.value = error.value

    loading.value = false
  }

  async function submitForm(formData: Task) {
    loading.value = true
    const { data, error } = await useMyFetch('tasks').post(formData).json()
    tasks.value.push(data.value)
    fetchError.value = error.value
    loading.value = false
  }

  async function completeTask(task: Task) {
    loading.value = true
    await useMyFetch(`tasks/${task.task_id}`).patch({ status: task.status }).json()
    loading.value = false
  }

  async function deleteTask(id: Number) {
    loading.value = true
    await useMyFetch(`tasks/${id}`).delete().json()
    tasks.value = tasks.value.filter((t) => t.task_id !== id)
    loading.value = false
  }

  return {
    tasks,
    loading,
    fetchError,
    fetchTasks,
    completeTask,
    deleteTask,
    submitForm,
  }
})
