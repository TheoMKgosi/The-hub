interface Task {
  task_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  order?: number
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
    const { $api } = useNuxtApp()
    loading.value = true
    const { tasks: fetchedTasks } = await $api<TaskResponse>('/tasks')

    if (fetchedTasks) tasks.value = fetchedTasks
    loading.value = false
  }


  async function submitForm(payload: { title: string; description: string; due_date?: string; priority?: number; status?: string }) {
    try {
      // TODO: validate payload
      const { $api } = useNuxtApp()
      const data = await $api<Task>('tasks', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      tasks.value.push(data)
      addToast("Task added successfully", "success")

    } catch (err) {
      addToast("Task not added", "error")
    }
  }

  async function editTask(payload: Task) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api(`tasks/${payload.task_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      fetchTasks()
      addToast("Edited task succesfully", "success")

    } catch (err) {
      addToast("Editing task failed", "error")
    }
  }

  async function reorderTask(payload: {task_id: string, order: number}[]) {
    const { $api } = useNuxtApp()
    await $api("/tasks/reorder", {
      method: 'PUT',
      body: JSON.stringify({task_orders: payload} )
    })
  }

  async function completeTask(payload: Task) {
    const { $api } = useNuxtApp()
    await $api(`tasks/${payload.task_id}`, {
      method: 'PATCH',
      body: JSON.stringify({ status: payload.status })
    })
  }

  async function deleteTask(id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`tasks/${id}`, {
        method: 'DELETE'
      })
      tasks.value = tasks.value.filter((t) => t.task_id !== id)
      addToast("Task deleted succesfully", "success")

    } catch (err) {
      addToast("Task did not delete", "error")

    }
  }

  function reset() {
    tasks.value = []
  }


  return {
    tasks,
    completedTasks,
    loading,
    fetchError,
    fetchTasks,
    editTask,
    reorderTask,
    completeTask,
    deleteTask,
    submitForm,
    reset,
  }
})
