<script setup lang="ts">
import VueCal from 'vue-cal'
import 'vue-cal/dist/vuecal.css'

const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const { addToast } = useToast()

// Lifecycle hooks
const { $api } = useNuxtApp()

const modalShow = ref(false)
const aiModalShow = ref(false)
const selectedTask = ref(null)
const aiSuggestions = ref([])
const calendarView = ref('month')
const vueCalRef = ref(null)
const isLoading = ref(false)
const isAISuggestionsLoading = ref(false)

async function fetchEvents() {
  isLoading.value = true
  try {
    if (scheduleStore.schedule.length === 0) {
      await scheduleStore.fetchSchedule()
    }
    if (taskStore.tasks.length === 0) {
      await taskStore.fetchTasks()
    }
  } catch (error) {
    addToast('Failed to load calendar data', 'error')
  } finally {
    isLoading.value = false
  }
}

const formData = reactive({
  title: '',
  start: '',
  end: '',
  taskId: null,
  recurrenceRuleId: null,
})

const recurrenceForm = reactive({
  frequency: 'daily',
  interval: 1,
  endDate: null,
  count: null,
})

// Format events for VueCal
const formattedEvents = computed(() => {
  if (!scheduleStore.schedule || !Array.isArray(scheduleStore.schedule)) {
    return []
  }

  return scheduleStore.schedule.map(event => {
    if (!event || !event.start || !event.end) {
      console.warn('Invalid event data:', event)
      return null
    }

    const startDate = new Date(event.start)
    const endDate = new Date(event.end)

    if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
      console.warn('Invalid date in event:', event)
      return null
    }

    const duration = (endDate.getTime() - startDate.getTime()) / (1000 * 60) // minutes

    return {
      id: event.id,
      title: event.title || 'Untitled Event',
      start: startDate,
      end: endDate,
      content: event.task?.title ? `Task: ${event.task.title}` : '',
      class: event.created_by_ai ? 'ai-event' : 'regular-event',
      background: event.created_by_ai ? '#3b82f6' : '#10b981',
      borderColor: event.created_by_ai ? '#2563eb' : '#059669',
      allDay: duration >= 24 * 60, // If event is 24+ hours, show as all-day
      tooltip: `
        ${event.title || 'Untitled Event'}
        ${event.task?.title ? `\nTask: ${event.task.title}` : ''}
        ${event.created_by_ai ? '\nAI Suggested' : ''}
        \nStart: ${startDate.toLocaleString()}
        \nEnd: ${endDate.toLocaleString()}
      `.trim()
    }
  }).filter(event => event !== null)
})

const toDateTimeLocal = (date: Date) =>
  new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)

const onCellClick = ({ cursor }) => {
  console.log('Cell clicked:', cursor)
  modalShow.value = true
  console.log('Modal should be open:', modalShow.value)

  const clickedDate = new Date(cursor.date)

  // Round down to nearest 30 minutes
  const minutes = clickedDate.getMinutes()
  clickedDate.setMinutes(minutes < 30 ? 0 : 30)

  formData.start = toDateTimeLocal(clickedDate)
  formData.end = toDateTimeLocal(new Date(clickedDate.getTime() + 60 * 60 * 1000))
  console.log('Form data set:', formData)
}

// Quick event creation for power users
const quickCreateEvent = async (title: string, start: Date, end: Date) => {
  try {
    const dataToSend = {
      title,
      start: start.toISOString(),
      end: end.toISOString()
    }

    await scheduleStore.submitForm(dataToSend)
    addToast('Event created successfully', 'success')
    await fetchEvents()
  } catch (error) {
    addToast('Failed to create event', 'error')
  }
}

function onViewChange(viewMeta: { start: Date; end: Date }) {
  fetchEvents()
}

function changeView(view: string) {
  calendarView.value = view
  // Force a re-render by triggering the VueCal view change
  if (vueCalRef.value) {
    vueCalRef.value.switchView(view)
  }
}

async function onEventDropped({ event, newStart, newEnd }) {
  try {
    await $api(`schedule/${event.id}`, {
      method: 'PUT',
      body: {
        start: newStart.toISOString(),
        end: newEnd.toISOString()
      }
    })
    addToast('Event updated successfully', 'success')
    await fetchEvents()
  } catch (error) {
    addToast('Failed to update event', 'error')
  }
}

async function onEventDelete(event) {
  try {
    await scheduleStore.deleteSchedule(event.id)
    addToast('Event deleted successfully', 'success')
  } catch (error) {
    addToast('Failed to delete event', 'error')
  }
}

async function onEventClick(event) {
  // Handle event click - could open edit modal
  console.log('Event clicked:', event)
  // For now, just show a toast with event details
  addToast(`Event: ${event.title}`, 'info')
}

async function onEventCreate(event) {
  try {
    const dataToSend = {
      title: event.title || 'New Event',
      start: event.start.toISOString(),
      end: event.end.toISOString(),
    }

    await scheduleStore.submitForm(dataToSend)
    addToast('Event created successfully', 'success')
    await fetchEvents()
  } catch (error) {
    addToast('Failed to create event', 'error')
  }
}

function close() {
  modalShow.value = false
  // Reset form data
  formData.title = ''
  formData.start = ''
  formData.end = ''
  formData.taskId = null
  formData.recurrenceRuleId = null

  // Reset recurrence form
  recurrenceForm.frequency = 'daily'
  recurrenceForm.interval = 1
  recurrenceForm.endDate = null
  recurrenceForm.count = null
}

function openModal() {
  modalShow.value = true
  const now = new Date()
  const oneHourLater = new Date(now.getTime() + 60 * 60 * 1000)
  formData.start = toDateTimeLocal(now)
  formData.end = toDateTimeLocal(oneHourLater)
}

async function save() {
  // Basic validation
  if (!formData.title || !formData.title.trim()) {
    addToast('Please enter a title for the event', 'error')
    return
  }

  if (!formData.start || !formData.end) {
    addToast('Please select both start and end times', 'error')
    return
  }

  const startDate = new Date(formData.start)
  const endDate = new Date(formData.end)

  if (startDate >= endDate) {
    addToast('End time must be after start time', 'error')
    return
  }

  try {
    const dataToSend = {
      title: formData.title.trim(),
      start: startDate.toISOString(),
      end: endDate.toISOString(),
      task_id: formData.taskId,
      recurrence_rule_id: formData.recurrenceRuleId,
    }

    console.log('Data to send:', dataToSend)

    await scheduleStore.submitForm(dataToSend)
    addToast('Event created successfully', 'success')
    await fetchEvents()
    close()
  } catch (error) {
    console.error('Error creating event:', error)
    addToast('Failed to create event', 'error')
  }
}

async function getAISuggestions() {
  isAISuggestionsLoading.value = true
  try {
    const response = await $api('ai/suggestions')
    console.log('AI suggestions response:', response)
    aiSuggestions.value = response.data?.suggestions || response.data?.value?.suggestions || []
    aiModalShow.value = true
  } catch (error) {
    console.error('Error getting AI suggestions:', error)
    addToast('Failed to get AI suggestions', 'error')
  } finally {
    isAISuggestionsLoading.value = false
  }
}

async function applyAISuggestion(suggestion) {
  try {
    const startDate = new Date(suggestion.start)
    const endDate = new Date(suggestion.end)

    await scheduleStore.submitForm({
      title: suggestion.title,
      start: startDate.toISOString(),
      end: endDate.toISOString(),
      task_id: suggestion.task_id,
      created_by_ai: true,
    })
    addToast('AI suggestion applied', 'success')
    await fetchEvents()
    aiModalShow.value = false
  } catch (error) {
    addToast('Failed to apply AI suggestion', 'error')
  }
}

async function createRecurrenceRule() {
  try {
    const response = await $api('recurrence-rules', {
      method: 'POST',
      body: recurrenceForm
    })
    const data = response.data.value || response.data
    formData.recurrenceRuleId = data?.id
    addToast('Recurrence rule created', 'success')
  } catch (error) {
    console.error('Error creating recurrence rule:', error)
    addToast('Failed to create recurrence rule', 'error')
  }
}

onMounted(async () => {
  await fetchEvents()
})
</script>

<template>
  <div class="calendar-container p-4">
    <div class="flex justify-between items-center mb-4">
      <h2>Calendar</h2>
      <div class="flex gap-2">
        <UiButton
          @click="getAISuggestions"
          variant="secondary"
          size="sm"
          :disabled="isAISuggestionsLoading"
        >
          <span v-if="isAISuggestionsLoading">Loading...</span>
          <span v-else>Get AI Suggestions</span>
        </UiButton>
        <UiButton @click="openModal" variant="primary" size="sm">
          Add Event
        </UiButton>
      </div>
    </div>

    <!-- VueCal Calendar -->
    <div class="calendar-wrapper relative">
      <!-- Loading overlay -->
      <div v-if="isLoading" class="absolute inset-0 bg-white/90 dark:bg-gray-900/90 flex items-center justify-center z-10">
        <div class="flex items-center gap-2 bg-white dark:bg-gray-800 px-4 py-3 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 dark:border-blue-400"></div>
          <span class="text-gray-700 dark:text-gray-200 font-medium">Loading calendar...</span>
        </div>
      </div>

      <VueCal
        ref="vueCalRef"
        :events="formattedEvents"
        :editable-events="{ drag: true, resize: true, delete: true, create: true }"
        :time-from="calendarView === 'month' ? 0 : 6 * 60"
        :time-to="calendarView === 'month' ? 24 * 60 : 22 * 60"
        :time-step="calendarView === 'month' ? 60 : 30"
        :locale="'en'"
        :selected-date="new Date()"
        :show-week-numbers="true"
        :views="['month', 'week', 'day']"
        @view-change="onViewChange"
        @cell-click="onCellClick"
        @event-drop="onEventDropped"
        @event-delete="onEventDelete"
        @event-click="onEventClick"
        @event-create="onEventCreate"
        class="vuecal--full-calendar"
      />
    </div>

    <!-- Event creation modal -->
    <div v-if="modalShow" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-[9999]">
      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4">Create Event</h3>
        <form class="space-y-4">
          <div>
            <label for="title" class="block text-sm font-medium mb-1">Title</label>
            <input type="text" id="title" v-model="formData.title" autocomplete="off"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>

          <div>
            <label for="task" class="block text-sm font-medium mb-1">Link to Task (optional)</label>
            <select id="task" v-model="formData.taskId" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary">
              <option value="">Select a task</option>
              <option v-for="task in taskStore.tasks" :key="task.id" :value="task.id">
                {{ task.title }}
              </option>
            </select>
          </div>

          <div>
            <label for="start" class="block text-sm font-medium mb-1">Start Date</label>
            <input type="datetime-local" id="start" v-model="formData.start"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>

          <div>
            <label for="end" class="block text-sm font-medium mb-1">End Date</label>
            <input type="datetime-local" id="end" v-model="formData.end"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>

          <!-- Recurrence options -->
          <div class="border-t pt-4">
            <h4 class="text-md font-medium mb-2">Recurrence (optional)</h4>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="frequency" class="block text-sm font-medium mb-1">Frequency</label>
                <select id="frequency" v-model="recurrenceForm.frequency" class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary">
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </div>
              <div>
                <label for="interval" class="block text-sm font-medium mb-1">Interval</label>
                <input type="number" id="interval" v-model="recurrenceForm.interval" min="1"
                  class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary" />
              </div>
            </div>
            <div class="mt-2">
              <UiButton @click="createRecurrenceRule" variant="secondary" size="sm">
                Create Recurrence Rule
              </UiButton>
            </div>
          </div>
        </form>
        <div class="flex justify-end space-x-2 mt-6">
          <UiButton @click="close" variant="default" size="md">Cancel</UiButton>
          <UiButton @click="save" variant="primary" size="md">Save</UiButton>
        </div>
      </div>
    </div>

    <!-- AI Suggestions modal -->
    <div v-show="aiModalShow" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-50">
      <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4">AI Schedule Suggestions</h3>
        <div v-if="aiSuggestions.length === 0" class="text-center py-8">
          <p class="text-gray-700 dark:text-gray-300">No suggestions available at the moment.</p>
        </div>
        <div v-else class="space-y-4">
          <div v-for="suggestion in aiSuggestions" :key="suggestion.id" class="border border-gray-200 dark:border-gray-600 rounded-lg p-4 bg-gray-50 dark:bg-gray-800">
            <h4 class="font-semibold text-gray-900 dark:text-gray-100">{{ suggestion.title }}</h4>
            <p class="text-sm text-gray-700 dark:text-gray-300 mt-1">
              {{ new Date(suggestion.start).toLocaleString() }} - {{ new Date(suggestion.end).toLocaleString() }}
            </p>
            <div class="mt-2">
              <UiButton @click="applyAISuggestion(suggestion)" variant="primary" size="sm">
                Apply Suggestion
              </UiButton>
            </div>
          </div>
        </div>
        <div class="flex justify-end mt-6">
          <UiButton @click="aiModalShow = false" variant="default" size="md">Close</UiButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar-container {
  max-width: 1400px;
  margin: 0 auto;
}

.calendar-wrapper {
  margin-top: 1rem;
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

/* VueCal custom styles */
:deep(.vuecal__header) {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:deep(.vuecal__title-bar) {
  background: #ffffff;
  border-bottom: 2px solid #e5e7eb;
  color: #111827;
}

:deep(.vuecal__cell--today) {
  background: rgba(59, 130, 246, 0.15) !important;
  border: 2px solid #3b82f6 !important;
}

:deep(.vuecal__event) {
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  padding: 4px 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  border: 2px solid transparent;
}

:deep(.vuecal__event.ai-event) {
  background: linear-gradient(135deg, #1e40af, #1e3a8a) !important;
  border-color: #1e3a8a !important;
  color: #ffffff !important;
  font-weight: 600;
}

:deep(.vuecal__event.regular-event) {
  background: linear-gradient(135deg, #065f46, #064e3b) !important;
  border-color: #064e3b !important;
  color: #ffffff !important;
  font-weight: 600;
}

:deep(.vuecal__cell--has-events) {
  background: rgba(16, 185, 129, 0.08);
}

:deep(.vuecal__cell--selected) {
  background: rgba(59, 130, 246, 0.25) !important;
  border: 2px solid #3b82f6 !important;
}

:deep(.vuecal__time-column) {
  background: #ffffff;
  border-right: 2px solid #e5e7eb;
  color: #374151;
  font-weight: 500;
}

:deep(.vuecal__weekdays-headings) {
  background: #f9fafb;
  border-bottom: 2px solid #d1d5db;
  color: #111827;
  font-weight: 600;
}

:deep(.vuecal__cell-date) {
  color: #374151;
  font-weight: 500;
}

 :deep(.vuecal__cell-content) {
  color: #6b7280;
}

/* Month view specific styles */
:deep(.vuecal__view--month .vuecal__cell) {
  min-height: 120px;
  padding: 4px;
  border: 1px solid #e5e7eb;
}

:deep(.vuecal__view--month .vuecal__cell-date) {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 4px;
  color: #374151;
}

:deep(.vuecal__view--month .vuecal__event) {
  font-size: 0.75rem;
  padding: 2px 4px;
  margin-bottom: 2px;
  border-radius: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

:deep(.vuecal__view--month .vuecal__cell--out-of-scope) {
  background: #f9fafb;
  color: #9ca3af;
}

:deep(.vuecal__view--month .vuecal__cell--today) {
  background: rgba(59, 130, 246, 0.1) !important;
  border: 2px solid #3b82f6 !important;
}

:deep(.vuecal__view--month .vuecal__cell--has-events) {
  background: rgba(16, 185, 129, 0.05);
}

/* Dark mode support */
.dark :deep(.vuecal__header) {
  background: linear-gradient(135deg, #3730a3 0%, #581c87 100%);
  color: #ffffff;
}

.dark :deep(.vuecal__title-bar) {
  background: #1f2937;
  border-bottom-color: #4b5563;
  color: #f9fafb;
}

.dark :deep(.vuecal__cell--today) {
  background: rgba(59, 130, 246, 0.25) !important;
  border-color: #60a5fa !important;
}

.dark :deep(.vuecal__time-column) {
  background: #1f2937;
  border-right-color: #4b5563;
  color: #d1d5db;
}

.dark :deep(.vuecal__weekdays-headings) {
  background: #111827;
  border-bottom-color: #4b5563;
  color: #f3f4f6;
}

.dark :deep(.vuecal__cell-date) {
  color: #e5e7eb;
}

.dark :deep(.vuecal__cell-content) {
  color: #9ca3af;
}

.dark :deep(.vuecal__event) {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
}

.dark :deep(.vuecal__cell--has-events) {
  background: rgba(16, 185, 129, 0.15);
}

.dark :deep(.vuecal__cell--selected) {
  background: rgba(59, 130, 246, 0.3) !important;
}

/* Dark mode month view styles */
.dark :deep(.vuecal__view--month .vuecal__cell) {
  border-color: #4b5563;
}

.dark :deep(.vuecal__view--month .vuecal__cell-date) {
  color: #e5e7eb;
}

.dark :deep(.vuecal__view--month .vuecal__cell--out-of-scope) {
  background: #374151;
  color: #6b7280;
}

.dark :deep(.vuecal__view--month .vuecal__cell--today) {
  background: rgba(59, 130, 246, 0.2) !important;
  border-color: #60a5fa !important;
}

.dark :deep(.vuecal__view--month .vuecal__cell--has-events) {
  background: rgba(16, 185, 129, 0.1);
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  :deep(.vuecal__event) {
    border: 3px solid currentColor !important;
    font-weight: 700;
  }

  :deep(.vuecal__cell--today) {
    border: 3px solid #000 !important;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  .animate-spin {
    animation: none;
  }
}

/* Focus indicators for keyboard navigation */
:deep(.vuecal__event:focus) {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

:deep(.vuecal__cell:focus) {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

/* Responsive design */
@media (max-width: 768px) {
  .calendar-container {
    max-width: 100%;
    padding: 0 1rem;
  }

  :deep(.vuecal__header) {
    padding: 0.75rem;
  }

  :deep(.vuecal__event) {
    font-size: 0.75rem;
    padding: 2px 4px;
  }

  /* Month view responsive styles */
  :deep(.vuecal__view--month .vuecal__cell) {
    min-height: 80px;
    padding: 2px;
  }

  :deep(.vuecal__view--month .vuecal__cell-date) {
    font-size: 1rem;
    margin-bottom: 2px;
  }

  :deep(.vuecal__view--month .vuecal__event) {
    font-size: 0.65rem;
    padding: 1px 3px;
    margin-bottom: 1px;
  }
}

/* Print styles */
@media print {
  :deep(.vuecal__event) {
    background: #000 !important;
    color: #fff !important;
    border: 1px solid #000 !important;
  }

  .calendar-container {
    box-shadow: none;
  }
}
</style>
