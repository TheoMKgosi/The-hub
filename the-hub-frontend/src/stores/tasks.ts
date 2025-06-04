import { defineStore } from "pinia";
import { ref } from "vue";
import { useFetch } from "@vueuse/core";

interface Task {
  task_id: number
  title: string
  description: string
}

export interface TaskResponse {
  tasks: Task[]
}

export const useTaskStore = defineStore("task", () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)

  async function fetchGoals() {
    loading.value = true
    const { data, error } = await useFetch("http://localhost:8080/tasks").json<TaskResponse>()

    if (data.value) tasks.value = data.value.tasks
    fetchError.value = error.value

    loading.value = false
  }

  async function submitForm(formData: Task) {
    loading.value = true
    const url = "http://localhost:8080/tasks"
    const { data, error } = await useFetch(url).post(formData).json()
    tasks.value.push(data.value)
    fetchError.value = error.value
    loading.value = false

  }

  return {
    tasks,
    loading,
    fetchError,
    fetchGoals,
    submitForm,
  }
})
