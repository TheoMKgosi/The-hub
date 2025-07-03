import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
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

  const completedTasks = computed(() => {
    console.log(tasks.value.filter(task => task.status === 'complete'))
    tasks.value.filter(task => task.status === 'complete')
  })

  async function submitForm(formData: Task) {
    loading.value = true
    const { data, error } = await useMyFetch('tasks').post(formData).json()
    tasks.value.push(data.value)
    fetchError.value = error.value
    loading.value = false
  }

  async function editTask(task: Task) {
    loading.value = true
    console.log(task)
    const { error } = await useMyFetch(`tasks/${task.task_id}`).patch(task).json()

    if (!error.value) {
      const index = tasks.value.findIndex(t => t.task_id === task.task_id)
      if (index !== -1) {
        tasks.value[index] = { ...tasks.value[index], ...task }
      }
    }

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
    completedTasks,
    loading,
    fetchError,
    fetchTasks,
    editTask,
    completeTask,
    deleteTask,
    submitForm,
  }
})
