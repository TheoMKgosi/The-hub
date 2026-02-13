<script setup lang="ts">
import type { FormUIProps } from '~/types/form'
import CrossIcon from './svg/CrossIcon.vue'

interface FormUIEmits {
  (e: 'submit', data: Record<string, any>): void
  (e: 'cancel'): void
  (e: 'close'): void
  (e: 'combobox-create', fieldName: string, name: string): void
}

const props = withDefaults(defineProps<FormUIProps>(), {
  cancelLabel: 'Cancel',
  showForm: true,
  teleportTarget: 'body',
  size: 'md'
})

const emit = defineEmits<FormUIEmits>()

const { validateObject } = useValidation()

const formData = reactive<Record<string, any>>({})
const validationErrors = ref<Record<string, string>>({})

// Initialize form data
if (props.initialData) {
  Object.assign(formData, props.initialData)
} else {
  // Initialize with default values based on field types
  props.fields.forEach(field => {
    if (field.type === 'select') {
      formData[field.name] = null
    } else if (field.type === 'number') {
      formData[field.name] = field.min || 0
    } else {
      formData[field.name] = ''
    }
  })
}

const isFormValid = computed(() => {
  // Check required fields
  const hasRequiredFields = props.fields
    .filter(field => field.required)
    .every(field => {
      const value = formData[field.name]
      return value !== null && value !== undefined && value !== ''
    })

  // Check for validation errors
  const hasNoErrors = Object.keys(validationErrors.value).length === 0

  return hasRequiredFields && hasNoErrors
})

const submitForm = async () => {
  validationErrors.value = {}

  // Validate if schema provided
  if (props.validationSchema) {
    const validation: ValidationResult = validateObject(formData, props.validationSchema)
    if (!validation.isValid) {
      validationErrors.value = validation.errors
      return
    }
  }

  emit('submit', { ...formData })

  // Reset form after successful submission
  resetForm()
}

const cancelForm = () => {
  emit('cancel')
  emit('close')
}

const closeModal = () => {
  emit('cancel')
  emit('close')
}

const resetForm = () => {
  validationErrors.value = {}
  if (props.initialData) {
    Object.assign(formData, props.initialData)
  } else {
    props.fields.forEach(field => {
      if (field.type === 'select' || field.type === 'combobox') {
        formData[field.name] = null
      } else if (field.type === 'number') {
        formData[field.name] = field.min || 0
      } else if (field.type === 'date' || field.type === 'datetime-local') {
        formData[field.name] = ''
      } else if (field.type === 'color') {
        formData[field.name] = '#000000'
      } else {
        formData[field.name] = ''
      }
    })
  }
}

// Watch for initialData changes and update formData
watch(() => props.initialData, (newInitialData) => {
  if (newInitialData) {
    Object.assign(formData, newInitialData)
  }
}, { immediate: true, deep: true })

// Expose reset function for parent components
defineExpose({
  resetForm
})
</script>

<template>
  <ClientOnly>
    <!-- Modal Overlay -->
    <Teleport :to="teleportTarget">
      <div v-if="showForm" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50"
        @click="closeModal">
        <div :class="[
          'bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg w-full max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light/30 dark:border-surface-dark/30 backdrop-blur-md relative',
          {
            'max-w-sm': size === 'sm',
            'max-w-md': size === 'md',
            'max-w-2xl': size === 'lg'
          }
        ]" @click.stop>

          <!-- Modal Header -->
          <div
            class="flex items-center justify-between p-6 border-b border-surface-light/20 dark:border-surface-dark/20 relative z-10">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark flex items-center gap-2">
              {{ title }}
            </h2>
            <button @click="closeModal" type="button"
              class="p-2 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 rounded transition-colors cursor-pointer text-text-light dark:text-text-dark hover:scale-110 shrink-0"
              title="Close">
              <CrossIcon />
            </button>
          </div>

          <!-- Modal Body -->
          <div class="p-6">
            <form @submit.prevent="submitForm" class="space-y-4">
              <div class="space-y-3">
                <div v-for="field in fields" :key="field.name" class="flex flex-col">
                  <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">
                    {{ field.label }}
                    <span v-if="field.required" class="text-red-500">*</span>
                  </label>

                  <!-- Text Input -->
                  <input v-if="field.type === 'text' || field.type === 'email' || field.type === 'password'"
                         :type="field.type" v-model="formData[field.name]" :placeholder="field.placeholder"
                    :required="field.required" :class="[
                      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors',
                      'bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30',
                      'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
                      { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
                    ]" />

                  <!-- Textarea -->
                  <textarea v-else-if="field.type === 'textarea'" v-model="formData[field.name]"
                    :placeholder="field.placeholder" :rows="field.rows || 3" :required="field.required" :class="[
                      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary resize-none transition-colors',
                      'bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30',
                      'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
                      { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
                    ]"></textarea>

                  <!-- Select -->
                  <select v-else-if="field.type === 'select'" v-model="formData[field.name]" :required="field.required"
                    :class="[
                      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors',
                      'bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30'
                    ]">
                    <option v-if="!field.required" class="dark:bg-surface-dark dark:text-text-dark" :value="null">Select
                      {{ field.label.toLowerCase() }}</option>
                    <option v-for="option in field.options" class="dark:bg-surface-dark dark:text-text-dark"
                      :key="option.value" :value="option.value">
                      {{ option.label }}
                    </option>
                  </select>

                  <!-- Combobox -->
                  <BaseComboBox v-else-if="field.type === 'combobox'" v-model="formData[field.name]"
                    :categories="field.categories || []" :placeholder="field.placeholder || 'Select or create...'"
                    :allow-create="field.allowCreate !== false"
                    @select="(category) => { formData[field.name] = category.budget_category_id }"
                    @create="(name) => { emit('combobox-create', field.name, name) }"
                    @update:model-value="(value) => { formData[field.name] = value }" class="w-full" />

                  <!-- Date/DateTime -->
                  <input v-else-if="field.type === 'date' || field.type === 'datetime-local'" :type="field.type"
                    v-model="formData[field.name]" :required="field.required" :class="[
                      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors',
                      'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30'
                    ]" />

                  <!-- Color -->
                  <input v-else-if="field.type === 'color'" :type="field.type" v-model="formData[field.name]"
                    :required="field.required" :class="[
                      'w-full h-10 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors cursor-pointer',
                      'bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30'
                    ]" />

                  <!-- Number -->
                  <input v-else-if="field.type === 'number'" :type="field.type" v-model.number="formData[field.name]"
                    :placeholder="field.placeholder" :required="field.required" :min="field.min" :max="field.max"
                    :step="field.step" :class="[
                      'w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors',
                      'bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark',
                      'border-surface-light/30 dark:border-surface-dark/30',
                      'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
                      { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
                    ]" />

                  <!-- Validation Error -->
                  <p v-if="validationErrors[field.name]" class="mt-1 text-sm text-red-500 dark:text-red-400">
                    {{ validationErrors[field.name] }}
                  </p>
                </div>
              </div>

              <!-- Modal Footer -->
              <div
                class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light/20 dark:border-surface-dark/20">
                <BaseButton type="button" :text="cancelLabel" @click="cancelForm" variant="default" size="md"
                  class="w-full sm:w-auto" />
                <BaseButton type="submit" :text="submitLabel" variant="primary" size="md" class="w-full sm:w-auto"
                  :disabled="!isFormValid" />
              </div>
            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>
