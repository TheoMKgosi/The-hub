<script setup lang="ts">
import { reactive, ref } from "vue";
import { useTaskStore } from '@/stores/tasks'
import * as chrono from 'chrono-node'

const taskStore = useTaskStore()

const formData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
  natural_language: '',
  use_natural_language: false,
})


const taskForm = ref(null)
const showForm = ref(true)

// Parse dates from natural language input
const parseDateFromText = (text: string): Date | null => {
  const results = chrono.parse(text)
  if (results.length > 0) {
    return results[0].start.date()
  }
  return null
}

// Extract priority keywords from text
const parsePriorityFromText = (text: string): number | null => {
  const lowerText = text.toLowerCase()

  // Priority 5 (Urgent/Critical)
  if (lowerText.includes('urgent') || lowerText.includes('asap') || lowerText.includes('critical')) {
    return 5
  }
  // Priority 4 (High/Important)
  else if (lowerText.includes('high priority') || lowerText.includes('high') || lowerText.includes('important')) {
    return 4
  }
  // Priority 3 (Medium/Normal)
  else if (lowerText.includes('medium priority') || lowerText.includes('medium') || lowerText.includes('normal')) {
    return 3
  }
  // Priority 2 (Low)
  else if (lowerText.includes('low priority') || lowerText.includes('low')) {
    return 2
  }
  // Priority 1 (Minor/Trivial)
  else if (lowerText.includes('minor') || lowerText.includes('trivial')) {
    return 1
  }

  return null
}

const submitForm = async () => {
  const dataToSend = { ...formData }
  if (dataToSend.due_date) {
    const date = new Date(dataToSend.due_date)
    dataToSend.due_date = date.toISOString()
  }

  // If using natural language, send that instead of structured fields
  if (dataToSend.use_natural_language && dataToSend.natural_language) {
    dataToSend.natural_language_input = dataToSend.natural_language
    dataToSend.use_natural_language = true

    // Client-side parsing as fallback
    const parsedDate = parseDateFromText(dataToSend.natural_language)
    const parsedPriority = parsePriorityFromText(dataToSend.natural_language)

    if (parsedDate) {
      dataToSend.due_date = parsedDate.toISOString()
    }
    if (parsedPriority) {
      dataToSend.priority = parsedPriority
    }

    // Clear other fields when using natural language
    dataToSend.title = ''
    dataToSend.description = ''
  }

  taskStore.submitForm(dataToSend)
  Object.assign(formData, {
    title: '',
    description: '',
    due_date: null,
    priority: 3,
    status: 'pending',
    natural_language: '',
    use_natural_language: false,
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

               <!-- Natural Language Toggle -->
               <div class="flex items-center space-x-2">
                 <input type="checkbox" v-model="formData.use_natural_language" id="natural-language-toggle" name="use_natural_language"
                   class="rounded border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-primary focus:ring-2 focus:ring-primary" />
                 <label for="natural-language-toggle" class="font-medium text-sm text-text-light dark:text-text-dark">
                   Use natural language input
                 </label>
               </div>

               <!-- Natural Language Input -->
               <div v-if="formData.use_natural_language" class="flex flex-col">
                 <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Describe your task</label>
                 <textarea v-model="formData.natural_language" name="natural_language" rows="4"
                   class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                   placeholder="e.g., 'Buy groceries tomorrow at 5pm, high priority'"></textarea>
               </div>

               <!-- Structured Input Fields -->
               <div v-else>
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
