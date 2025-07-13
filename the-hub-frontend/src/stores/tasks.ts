import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useMyFetch } from '@/config/fetch'
import { useToast } from '@/composables/useToast'



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
  const { addToast } = useToast()

  const completedTasks = computed(() => {
    return tasks.value.filter(task => task.status === 'complete')
  })

  async function fetchTasks() {
    loading.value = true
    const { data, error } = await useMyFetch('tasks').json<TaskResponse>()

    if (data.value) tasks.value = data.value.tasks
    fetchError.value = error.value

    loading.value = false
  }


  async function submitForm(formData: Task) {
    const { data, error } = await useMyFetch('tasks').post(formData).json()
    if (!data.value.task_id) {
      data.value.task_id = Date.now() // fallback if backend didn’t return ID
    }
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Task not added", "error")
    } else {
      tasks.value.push(data.value)
      addToast("Task added succesfully", "success")
    }
  }

  async function editTask(task: Task) {
    const { error } = await useMyFetch(`tasks/${task.task_id}`).patch(task).json()

    if (!error.value) {
      const index = tasks.value.findIndex(t => t.task_id === task.task_id)
      if (index !== -1) {
        tasks.value[index] = { ...tasks.value[index], ...task }
        addToast("Edited task succesfully", "success")
      } else {
        addToast("Editing task failed", "error")
      }
    } else {
      addToast("Editing task failed", "error")
    }
  }

  async function completeTask(task: Task) {
    await useMyFetch(`tasks/${task.task_id}`).patch({ status: task.status }).json()
  }

  async function deleteTask(id: number) {
    await useMyFetch(`tasks/${id}`).delete().json()
    tasks.value = tasks.value.filter((t) => t.task_id !== id)
    addToast("Task deleted succesfully", "success")
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
