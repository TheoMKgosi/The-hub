<script setup lang="ts">
import { reactive, ref } from "vue";
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
})


const taskForm = ref(null)
const showForm = ref(true)

const submitForm = async () => {
  const dataToSend = { ...formData }
  if (dataToSend.due_date) {
    const date = new Date(dataToSend.due_date)
    dataToSend.due_date = date.toISOString()
  }

  taskStore.submitForm(dataToSend)
  Object.assign(formData, {
    title: '',
    description: '',
    due_date: null,
    priority: 3,
    status: 'pending',
  })
}

</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div v-if="showForm" @click="showForm = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
        <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
          <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </div>
      </div>
    </Teleport>
  </ClientOnly>

  <!-- Modal Overlay -->
  <ClientOnly>
    <Teleport to="#plan">
      <!-- Modal Content -->
      <div v-if="!showForm" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50"
        @click="showForm = true">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark" @click.stop>

          <!-- Modal Header -->
          <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Create a Task</h2>
            <UiButton @click="showForm = true" variant="default" size="sm" class="p-2">
              ×
            </UiButton>
          </div>

          <!-- Modal Body -->
          <div class="p-6">
            <form @submit.prevent="submitForm" ref="taskForm" class="space-y-4">

              <div class="flex flex-col">
                <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Title</label>
                <input type="text" v-model="formData.title" name="title"
                  class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                  placeholder="Task title" required />
              </div>

              <div class="flex flex-col">
                <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Description</label>
                <textarea v-model="formData.description" name="description" rows="3"
                  class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                  placeholder="Optional description"></textarea>
              </div>

              <div class="flex flex-col">
                <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Due Date</label>
                <input type="datetime-local" v-model="formData.due_date" name="due_date"
                  class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
              </div>

              <div class="flex flex-col">
                <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Priority</label>
                <select v-model.number="formData.priority" name="priority"
                  class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                  <option :value="1">1 - Low</option>
                  <option :value="2">2 - Medium</option>
                  <option :value="3">3 - Medium</option>
                  <option :value="4">4 - High</option>
                  <option :value="5">5 - High</option>
                </select>
              </div>

              <!-- Modal Footer -->
              <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                <UiButton type="button" @click="showForm = true" variant="default" size="md" class="w-full sm:w-auto">
                  Cancel
                </UiButton>
                <UiButton type="submit" variant="primary" size="md" class="w-full sm:w-auto">
                  Create Task
                </UiButton>
              </div>

            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>
