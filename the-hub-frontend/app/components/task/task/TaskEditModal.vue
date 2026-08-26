<script setup lang="ts">
import * as v from 'valibot'
import type { Task, TaskUpdate } from '~/types/task'
import { CalendarDate, CalendarDateTime, parseDate, parseDateTime, ZonedDateTime, type DateValue } from '@internationalized/date'
import type { FormSubmitEvent } from '@nuxt/ui'

interface Props {
  task: Task | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: [boolean]
}>()

const themeStore = useThemeStore()

onMounted(() => {
  if (themeStore.themes.length === 0) themeStore.fetchThemes()
})

const schema = v.object({
  title: v.string(),
  description: v.optional(v.string()),
  due_date: v.optional(v.custom<DateValue>(isDateValue)),
  priority: v.optional(v.number()),
  start_date: v.optional(v.custom<DateValue>(isDateValue)),
  time_estimate: v.optional(v.number()),
  theme_id: v.optional(v.string())
})

type Schema = v.InferOutput<typeof schema>

function isDateValue(value: unknown): value is DateValue {
  return value instanceof CalendarDate
    || value instanceof CalendarDateTime
    || value instanceof ZonedDateTime
}

const taskStore = useTaskStore()
const toast = useToast()

const editState = reactive({
  title: props.task!.title,
  description: props.task?.description,
  due_date: shallowRef(toCalendarDate(props.task?.due_date)),
  priority: props.task?.priority,
  start_time: shallowRef(toCalendarDate(props.task?.start_time)),
  time_estimate: props.task?.time_estimate_minutes,
  theme_id: props.task?.theme_id
})

function toCalendarDate(date: any): DateValue | undefined {
  if (!date) return undefined

  if (date instanceof Date) {
    return new CalendarDate(
      date.getFullYear(),
      date.getMonth() + 1,
      date.getDate()
    )
    // or CalendarDateTime if you need time preserved:
    // return new CalendarDateTime(
    //   value.getFullYear(),
    //   value.getMonth() + 1,
    //   value.getDate(),
    //   value.getHours(),
    //   value.getMinutes(),
    //   value.getSeconds()
    // )
  }

  if (typeof date === 'string') {
    try {
      // has time component?
      return date.includes('T') ? parseDateTime(date) : parseDate(date)
    } catch {
      // fallback: parse via JS Date if string isn't strict ISO
      const d = new Date(date)
      if (isNaN(d.getTime())) return undefined
      return new CalendarDate(d.getFullYear(), d.getMonth() + 1, d.getDate())
    }
  }
}


// Handle form submission
const handleSubmit = async (formData: FormSubmitEvent<Schema>) => {
  if (!props.task) return

  // Build update payload
  const updates: TaskUpdate = {
    task_id: props.task.task_id
  }

  // Only include fields that have values
  if (formData.data.title) updates.title = formData.data.title
  if (formData.data.description !== undefined) updates.description = formData.data.description
  if (formData.data.priority) updates.priority = formData.data.priority
  if (formData.data.time_estimate) updates.time_estimate_minutes = formData.data.time_estimate

  if (formData.data.theme_id !== undefined) updates.theme_id = formData.data.theme_id || undefined

  if (formData.data.due_date) updates.due_date = new Date(formData.data.due_date.toString())

  // if (formData.data.start_time) {
  //   // Add seconds if not present
  //   const startDate = new Date(formData.data.start_date?.toString())
  //   updates.start_time = startDate.toISOString()
  //
  //   // Calculate end time based on duration
  //   if (formData.data.time_estimate) {
  //     const updates.end_time = formData.data.time_estimate
  //     const endDate = new Date(startDate.getTime() + durationMinutes * 60000)
  //     updates.end_time = endDate.toISOString()
  //   }
  // }

  try {
    await taskStore.editTask(updates)
    toast.add({ title: "Task", description: `Task ${props.task.title} updated successfully`, color: "success" })
  } catch (err) {
    toast.add({ title: "Error", description: "Failed to update task", color: "error" })
  }
}
</script>

<template>
  <UModal title="Edit Task" :close="{ onClick: () => emit('close', false) }">
    <template #body>
      <UForm :schema="schema" :state="editState" @submit="handleSubmit" class="space-y-3">
        <UFormField label="title">
          <UInput required v-model="editState.title" />
        </UFormField>
        <UFormField label="description">
          <UTextarea v-model="editState.description" />
        </UFormField>
        <UFormField label="Due Date">
          <UInputDate v-model="editState.due_date" />
        </UFormField>
        <UFormField label="Start Time">
          <UInputDate v-model="editState.start_time" />
        </UFormField>
        <UFormField label="priority">
          <UInputNumber v-model="editState.priority" :min="1" :max="5" />
        </UFormField>
        <UFormField label="Time Estimate (minutes)">
          <UInputNumber v-model="editState.time_estimate" :min="1" />
        </UFormField>
        <UFormField label="Theme">
          <select
            v-model="editState.theme_id"
            class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">No theme</option>
            <option v-for="theme in themeStore.themes" :key="theme.theme_id" :value="theme.theme_id">
              {{ theme.name }}
            </option>
          </select>
        </UFormField>
        <UButton label="Edit Task" type="submit" />
      </UForm>
    </template>
  </UModal>
</template>
