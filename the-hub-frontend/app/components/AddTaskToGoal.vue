<script setup lang="ts">
import { reactive, ref } from "vue"
import { useGoalStore } from '@/stores/goals'

interface Props {
  goalId: string
  goalTitle: string
}

const props = defineProps<Props>()

const goalStore = useGoalStore()

const showForm = ref(false)
const isSubmitting = ref(false)

const formData = reactive({
  title: '',
  description: '',
  priority: 3,
  due_date: null as string | null,
})

const resetForm = () => {
  Object.assign(formData, {
    title: '',
    description: '',
    priority: 3,
    due_date: null,
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

    await goalStore.addTaskToGoal(props.goalId, dataToSend)
    resetForm()
    showForm.value = false
  } catch (err) {
    console.error('Failed to add task:', err)
  } finally {
    isSubmitting.value = false
  }
}

const toggleForm = () => {
  showForm.value = !showForm.value
  if (!showForm.value) {
    resetForm()
  }
}
</script>

<template>
  <div class="space-y-3">
    <!-- Add Task Button -->
    <div v-if="!showForm" class="flex items-center justify-center">
      <UiButton @click="toggleForm" variant="default" size="sm" class="w-full py-2.5 border-2 border-dashed border-surface-light/30 dark:border-surface-dark/30 hover:border-primary/50 hover:bg-surface-light/10 dark:hover:bg-surface-dark/10 transition-colors">
        <svg fill="currentColor" height="16px" width="16px" class="mr-2" viewBox="0 0 24 24">
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
        Add Task to "{{ goalTitle }}"
      </UiButton>
    </div>

    <!-- Task Form -->
    <div v-else class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/20 dark:border-surface-dark/20">
      <div class="flex items-center justify-between mb-4">
        <h4 class="text-sm font-medium text-text-light dark:text-text-dark flex items-center gap-2">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
          </svg>
          Add Task to "{{ goalTitle }}"
        </h4>
        <UiButton @click="toggleForm" variant="default" size="sm" class="p-1 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20">
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
            <span v-if="isSubmitting">Adding...</span>
            <span v-else>Add Task</span>
          </UiButton>
          <UiButton
            type="button"
            @click="toggleForm"
            variant="default"
            size="sm"
          >
            Cancel
          </UiButton>
        </div>
      </form>
    </div>
  </div>
</template>