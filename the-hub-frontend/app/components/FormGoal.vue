<script setup lang="ts">
import { reactive, ref } from "vue";
import { useGoalStore } from '@/stores/goals'
import { useValidation } from '@/composables/useValidation'

const goalStore = useGoalStore()
const { validateObject, schemas } = useValidation()

const formData = reactive({
  title: '',
  description: '',
  due_date: '',
  priority: null as number | null,
  category: '',
  color: '#3B82F6',
})

const goalForm = ref(null)
const showForm = ref(true)
const validationErrors = ref<Record<string, string>>({})

const submitForm = async () => {
  validationErrors.value = {}

  const payload = {
    title: formData.title.trim(),
    description: formData.description.trim(),
    due_date: formData.due_date || undefined,
    priority: formData.priority || undefined,
    category: formData.category.trim() || undefined,
    color: formData.color,
  }

  const validation = validateObject(payload, schemas.goal.create)

  if (!validation.isValid) {
    validationErrors.value = validation.errors
    return
  }

  try {
    await goalStore.createGoal(payload)

    // Reset form
    Object.assign(formData, {
      title: '',
      description: '',
      due_date: '',
      priority: null,
      category: '',
      color: '#3B82F6',
    })

    // Close modal
    showForm.value = true
    validationErrors.value = {}
  } catch (err) {
    // Error is already handled in the store
  }
}

const cancelForm = () => {
  Object.assign(formData, {
    title: '',
    description: '',
    due_date: '',
    priority: null,
    category: '',
    color: '#3B82F6',
  })
  showForm.value = true
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
        <div class="bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light/30 dark:border-surface-dark/30 backdrop-blur-md" @click.stop>

          <!-- Modal Header -->
          <div class="flex items-center justify-between p-6 border-b border-surface-light/20 dark:border-surface-dark/20">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark flex items-center gap-2">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
              </svg>
              Create a Goal
            </h2>
            <UiButton @click="cancelForm" variant="default" size="sm" class="p-2 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20">
              ×
            </UiButton>
          </div>

          <!-- Modal Body -->
          <div class="p-6">
            <form @submit.prevent="submitForm" ref="goalForm" class="space-y-4">

               <div class="space-y-3">
                 <div class="flex flex-col">
                   <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Title</label>
                   <input type="text" v-model="formData.title" name="title"
                     class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                     placeholder="Goal title" required
                     :class="{ 'border-red-500 focus:ring-red-500': validationErrors.title }" />
                   <p v-if="validationErrors.title" class="mt-1 text-sm text-red-500 dark:text-red-400">
                     {{ validationErrors.title }}
                   </p>
                 </div>

                  <div class="flex flex-col">
                    <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Description</label>
                    <textarea v-model="formData.description" name="description" rows="3"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary resize-none transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                      placeholder="Optional description"
                      :class="{ 'border-red-500 focus:ring-red-500': validationErrors.description }"></textarea>
                    <p v-if="validationErrors.description" class="mt-1 text-sm text-red-500 dark:text-red-400">
                      {{ validationErrors.description }}
                    </p>
                  </div>

                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Due Date</label>
                      <input type="date" v-model="formData.due_date"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors" />
                    </div>

                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Priority</label>
                      <select v-model="formData.priority"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors">
                        <option :value="null">No priority</option>
                        <option :value="1">1 - Low</option>
                        <option :value="2">2</option>
                        <option :value="3">3 - Medium</option>
                        <option :value="4">4</option>
                        <option :value="5">5 - High</option>
                      </select>
                    </div>
                  </div>

                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Category</label>
                      <input type="text" v-model="formData.category"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                        placeholder="e.g., Work, Personal, Health" />
                    </div>

                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Color</label>
                      <input type="color" v-model="formData.color"
                        class="w-full h-10 px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors cursor-pointer" />
                    </div>
                  </div>
               </div>

              <!-- Modal Footer -->
              <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light/20 dark:border-surface-dark/20">
                <UiButton type="button" @click="cancelForm" variant="default" size="md" class="w-full sm:w-auto">
                  Cancel
                </UiButton>
                <UiButton type="submit" variant="primary" size="md" class="w-full sm:w-auto" :disabled="!formData.title.trim()">
                  Create Goal
                </UiButton>
              </div>

            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>