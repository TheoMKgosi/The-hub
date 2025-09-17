<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useValidation, type ValidationResult } from '@/composables/useValidation'

interface FormField {
  name: string
  label: string
  type: 'text' | 'textarea' | 'select' | 'date' | 'datetime-local' | 'color' | 'number' | 'email' | 'password'
  placeholder?: string
  required?: boolean
  options?: { value: any; label: string }[]
  rows?: number
  min?: number
  max?: number
  step?: number
}

interface FormInlineProps {
  fields: FormField[]
  submitLabel: string
  validationSchema?: Record<string, any>
  initialData?: Record<string, any>
  loading?: boolean
  error?: string
}

interface FormInlineEmits {
  (e: 'submit', data: Record<string, any>): void
}

const props = withDefaults(defineProps<FormInlineProps>(), {
  loading: false,
  error: ''
})

const emit = defineEmits<FormInlineEmits>()

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

  return hasRequiredFields && hasNoErrors && !props.loading
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
}

// Expose form data and reset function for parent components
defineExpose({
  formData,
  resetForm: () => {
    validationErrors.value = {}
    if (props.initialData) {
      Object.assign(formData, props.initialData)
    } else {
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
  }
})
</script>

<template>
  <form @submit.prevent="submitForm" class="space-y-4">
    <div v-for="field in fields" :key="field.name" class="space-y-2">
      <label class="block font-medium text-sm text-text-light dark:text-text-dark">
        {{ field.label }}
        <span v-if="field.required" class="text-red-500">*</span>
      </label>

      <!-- Text Input -->
      <input v-if="field.type === 'text' || field.type === 'email' || field.type === 'password'"
        :type="field.type"
        v-model="formData[field.name]"
        :placeholder="field.placeholder"
        :required="field.required"
        :class="[
          'w-full px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary transition-colors',
          'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
          'border-surface-light dark:border-surface-dark',
          { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
        ]" />

      <!-- Textarea -->
      <textarea v-else-if="field.type === 'textarea'"
        v-model="formData[field.name]"
        :placeholder="field.placeholder"
        :rows="field.rows || 3"
        :required="field.required"
        :class="[
          'w-full px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary resize-none transition-colors',
          'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
          'border-surface-light dark:border-surface-dark',
          { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
        ]"></textarea>

      <!-- Select -->
      <select v-else-if="field.type === 'select'"
        v-model="formData[field.name]"
        :required="field.required"
        :class="[
          'w-full px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary transition-colors',
          'border-surface-light dark:border-surface-dark'
        ]">
        <option v-if="!field.required" :value="null">Select {{ field.label.toLowerCase() }}</option>
        <option v-for="option in field.options" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>

      <!-- Date/DateTime -->
      <input v-else-if="field.type === 'date' || field.type === 'datetime-local'"
        :type="field.type"
        v-model="formData[field.name]"
        :required="field.required"
        :class="[
          'w-full px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary transition-colors',
          'border-surface-light dark:border-surface-dark'
        ]" />

      <!-- Color -->
      <input v-else-if="field.type === 'color'"
        :type="field.type"
        v-model="formData[field.name]"
        :required="field.required"
        :class="[
          'w-full h-10 px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary transition-colors cursor-pointer',
          'border-surface-light dark:border-surface-dark'
        ]" />

      <!-- Number -->
      <input v-else-if="field.type === 'number'"
        :type="field.type"
        v-model.number="formData[field.name]"
        :placeholder="field.placeholder"
        :required="field.required"
        :min="field.min"
        :max="field.max"
        :step="field.step"
        :class="[
          'w-full px-3 py-2 rounded-lg border bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark',
          'focus:outline-none focus:ring-2 focus:ring-primary transition-colors',
          'placeholder:text-text-light/50 dark:placeholder:text-text-dark/50',
          'border-surface-light dark:border-surface-dark',
          { 'border-red-500 focus:ring-red-500': validationErrors[field.name] }
        ]" />

      <!-- Validation Error -->
      <p v-if="validationErrors[field.name]" class="text-sm text-red-500 dark:text-red-400">
        {{ validationErrors[field.name] }}
      </p>
    </div>

    <!-- Submit Button -->
    <UiButton
      type="submit"
      variant="primary"
      size="md"
      class="w-full"
      :disabled="!isFormValid"
    >
      <span v-if="loading" class="flex items-center justify-center">
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Processing...
      </span>
      <span v-else>{{ submitLabel }}</span>
    </UiButton>

    <!-- General Error -->
    <p v-if="error" class="text-red-500 dark:text-red-400 text-center text-sm">
      {{ error }}
    </p>
  </form>
</template>