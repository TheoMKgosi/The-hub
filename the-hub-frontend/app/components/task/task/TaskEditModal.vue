<script setup lang="ts">
import type { FormField } from '~/types/form'
import type { Task, TaskUpdate } from '~/types/task'

interface Props {
  task: Task | null
  isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save'): void
}>()

const taskStore = useTaskStore()
const { addToast } = useToast()
const { validateObject } = useValidation()

// Format date to datetime-local format (YYYY-MM-DDTHH:MM) without seconds in local timezone
const formatDateTimeLocal = (date: Date | null | undefined): string => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}` // "2024-01-15T14:30" in local time
}

// Form fields configuration
const formFields = computed((): FormField[] => [
  {
    name: 'title',
    label: 'Title',
    type: 'text',
    required: true
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea',
    rows: 3
  },
  {
    name: 'due_date',
    label: 'Due Date',
    type: 'datetime-local'
  },
  {
    name: 'start_time',
    label: 'Start Time',
    type: 'datetime-local'
  },
  {
    name: 'priority',
    label: 'Priority (1-5)',
    type: 'number',
    min: 1,
    max: 5
  },
  {
    name: 'status',
    label: 'Status',
    type: 'select',
    required: true,
    options: [
      { value: 'pending', label: 'Pending' },
      { value: 'in_progress', label: 'In Progress' },
      { value: 'completed', label: 'Completed' }
    ]
  },
  {
    name: 'time_estimate_minutes',
    label: 'Time Estimate (minutes)',
    type: 'number',
    min: 1
  }
])

// Get initial data from task
const getInitialData = computed(() => {
  if (!props.task) return {}
  
  return {
    title: props.task.title || '',
    description: props.task.description || '',
    due_date: formatDateTimeLocal(props.task.due_date),
    start_time: formatDateTimeLocal(props.task.start_time),
    priority: props.task.priority || '',
    status: props.task.status || 'pending',
    time_estimate_minutes: props.task.time_estimate_minutes || ''
  }
})

// Validation schema
const validationSchema = {
  due_date: (value: string) => {
    if (!value) return null
    const date = new Date(value)
    const now = new Date()
    return date > now ? null : 'Due date must be in the future'
  },
  start_time: (value: string) => {
    if (!value) return null
    const date = new Date(value)
    const now = new Date()
    return date > now ? null : 'Start time must be in the future'
  }
}

// Modal title
const modalTitle = computed(() => {
  return props.task ? `Edit: ${props.task.title}` : 'Edit Task'
})

// Handle form submission
const handleSubmit = async (formData: Record<string, any>) => {
  if (!props.task) return

  // Build update payload
  const updates: TaskUpdate = {
    task_id: props.task.task_id
  }

  // Only include fields that have values
  if (formData.title) updates.title = formData.title
  if (formData.description !== undefined) updates.description = formData.description
  if (formData.priority) updates.priority = parseInt(formData.priority)
  if (formData.status) updates.status = formData.status
  if (formData.time_estimate_minutes) {
    updates.time_estimate_minutes = parseInt(formData.time_estimate_minutes)
  }

  // Handle dates - convert from "YYYY-MM-DDTHH:MM" to ISO8601 with seconds
  if (formData.due_date) {
    // Add seconds if not present
    const dateStr = formData.due_date.length === 16 ? formData.due_date + ':00' : formData.due_date
    updates.due_date = new Date(dateStr)
  }
  if (formData.start_time) {
    // Add seconds if not present
    const dateStr = formData.start_time.length === 16 ? formData.start_time + ':00' : formData.start_time
    const startDate = new Date(dateStr)
    updates.start_time = startDate.toISOString()
    
    // Calculate end time based on duration
    if (formData.time_estimate_minutes) {
      const durationMinutes = parseInt(formData.time_estimate_minutes)
      const endDate = new Date(startDate.getTime() + durationMinutes * 60000)
      updates.end_time = endDate.toISOString()
    }
  }

  try {
    await taskStore.editTask(updates)
    addToast(`Task "${props.task.title}" updated successfully`, 'success')
    emit('save')
  } catch (err) {
    addToast('Failed to update task', 'error')
  }
}

// Handle cancel
const handleCancel = () => {
  emit('close')
}
</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <FormUI
        v-if="isOpen && task"
        :title="modalTitle"
        :fields="formFields"
        :initialData="getInitialData"
        :validationSchema="validationSchema"
        submitLabel="Save Changes"
        cancelLabel="Cancel"
        @submit="handleSubmit"
        @cancel="handleCancel"
        @close="handleCancel"
      />
    </Teleport>
  </ClientOnly>
</template>
