<script setup lang="ts">
import { reactive, ref } from "vue"
import { useGoalStore } from '@/stores/goals'

interface Task {
  task_id: string
  title: string
  description: string
  due_date?: string
  priority?: number
  status: string
  order: number
}

interface Props {
  goalId: string
  task: Task
  goalTitle: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  updated: [task: Task]
}>()

const goalStore = useGoalStore()

const isSubmitting = ref(false)

const formData = reactive({
  title: props.task.title,
  description: props.task.description || '',
  priority: props.task.priority || 3,
  due_date: props.task.due_date ? new Date(props.task.due_date).toISOString().slice(0, 16) : null as string | null,
})

const resetForm = () => {
  Object.assign(formData, {
    title: props.task.title,
    description: props.task.description || '',
    priority: props.task.priority || 3,
    due_date: props.task.due_date ? new Date(props.task.due_date).toISOString().slice(0, 16) : null,
  })
}

const submitForm = async () => {
  if (!formData.title.trim()) return

  isSubmitting.value = true

  try {
    const dataToSend = { ...formData }
    if (dataToSend.due_date) {
      const date = new Date(dataToSend.due_date)
      dataToSend.due_date = date.toISOString()
    }

    const updatedTask = await goalStore.updateGoalTask(props.goalId, props.task.task_id, dataToSend)
    emit('updated', updatedTask)
    emit('close')
  } catch (err) {
    console.error('Failed to update task:', err)
  } finally {
    isSubmitting.value = false
  }
}

const closeForm = () => {
  resetForm()
  emit('close')
}
</script>

<template>
  <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
    <div class="flex items-center justify-between mb-4">
      <h4 class="text-sm font-medium text-text-light dark:text-text-dark flex items-center gap-2">
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
        </svg>
        Edit Task in "{{ goalTitle }}"
      </h4>
      <UiButton @click="closeForm" variant="default" size="sm" class="p-1 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20">
        ×
      </UiButton>
    </div>

    <form @submit.prevent="submitForm" class="space-y-4">
      <div class="space-y-3">
        <div class="flex flex-col">
          <label class="mb-2 text-xs font-medium text-text-light dark:text-text-dark">Title *</label>
          <input
            type="text"
            v-model="formData.title"
            class="w-full px-3 py-2 text-sm border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Task title"
            required
          />
        </div>

        <div class="flex flex-col">
          <label class="mb-2 text-xs font-medium text-text-light dark:text-text-dark">Description</label>
          <textarea
            v-model="formData.description"
            rows="2"
            class="w-full px-3 py-2 text-sm border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary resize-none transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Optional description"
          ></textarea>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="flex flex-col">
            <label class="mb-2 text-xs font-medium text-text-light dark:text-text-dark">Due Date</label>
            <input
              type="datetime-local"
              v-model="formData.due_date"
              class="w-full px-3 py-2 text-sm border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors"
            />
          </div>

          <div class="flex flex-col">
            <label class="mb-2 text-xs font-medium text-text-light dark:text-text-dark">Priority</label>
            <select
              v-model.number="formData.priority"
              class="w-full px-3 py-2 text-sm border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors"
            >
              <option :value="1">1 - Low</option>
              <option :value="2">2 - Medium</option>
              <option :value="3">3 - Medium</option>
              <option :value="4">4 - High</option>
              <option :value="5">5 - High</option>
            </select>
          </div>
        </div>
      </div>

      <div class="flex gap-2 pt-3 border-t border-surface-light/20 dark:border-surface-dark/20">
        <UiButton
          type="submit"
          variant="primary"
          size="sm"
          :disabled="!formData.title.trim() || isSubmitting"
          class="flex-1"
        >
          <span v-if="isSubmitting">Updating...</span>
          <span v-else>Update Task</span>
        </UiButton>
        <UiButton
          type="button"
          @click="closeForm"
          variant="default"
          size="sm"
        >
          Cancel
        </UiButton>
      </div>
    </form>
  </div>
</template>