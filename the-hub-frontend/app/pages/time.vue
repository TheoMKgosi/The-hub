<script setup lang="ts">
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin from '@fullcalendar/interaction'
const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const toast = useToast()

const modalShow = ref(false)
const calendarView = ref('timeGridWeek')

// FullCalendar configuration
const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
  headerToolbar: {
    left: 'dayGridMonth,timeGridWeek,timeGridDay',
    center: 'title',
    right: 'prev,next today'
  },
  initialView: calendarView.value,
  events: formattedEvents.value,
  nowIndicator: true,
  editable: true,
  selectable: true,
  selectMirror: true,
  dayMaxEvents: true,
  weekends: true,
  height: 'auto',
  // eventDisplay: 'block',
  eventTimeFormat: {
    hour: '2-digit',
    minute: '2-digit',
    meridiem: false
  },
  slotMinTime: '06:00:00',
  slotMaxTime: '22:00:00',
  slotDuration: '00:30:00',
  eventClick: onEventClick,
  eventDrop: onEventDrop,
  eventResize: onEventDrop,
  select: onDateSelect,
}))

const isLoading = ref(false)
const selectedEvents = ref([])
const bulkMode = ref(false)
const conflictModalShow = ref(false)
const conflicts = ref([])
const eventModalShow = ref(false)
const selectedEvent = ref(null)

const calendarKey = ref(0)

const suggestionModalShow = ref(false)
const suggestionsLoading = ref(false)
const selectedSuggestionIds = ref<string[]>([])

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
    toast.add({ title: "Error", description: "Failed to load calendar data", color: "error" })
  } finally {
    isLoading.value = false
  }
}

const formData = reactive({
  title: '',
  start: '',
  end: '',
  recurrenceRuleId: null,
})

const recurrenceForm = reactive({
  frequency: 'daily',
  interval: 1,
  endDate: null,
  count: null,
})

// Format events for FullCalendar
const formattedEvents = computed(() => {
  const events = []

  // Add schedule events
  if (scheduleStore.schedule && Array.isArray(scheduleStore.schedule)) {
    const scheduleEvents = scheduleStore.schedule.map(event => {
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
        start: startDate.toISOString(),
        end: endDate.toISOString(),
        backgroundColor: '#10b981',
        borderColor: '#059669',
        textColor: '#ffffff',
        allDay: duration >= 24 * 60, // If event is 24+ hours, show as all-day
        description: `
          ${event.title || 'Untitled Event'}
          \nStart: ${startDate.toLocaleString()}
          \nEnd: ${endDate.toLocaleString()}
        `.trim()
      }
    }).filter(event => event !== null)

    events.push(...scheduleEvents)
  }

  // Add AI suggestion events
  if (scheduleStore.suggestions && Array.isArray(scheduleStore.suggestions)) {
    const suggestionEvents = scheduleStore.suggestions.map(suggestion => {
      if (!suggestion || !suggestion.start || !suggestion.end) {
        return null
      }

      const startDate = new Date(suggestion.start)
      const endDate = new Date(suggestion.end)

      if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
        return null
      }

      const duration = (endDate.getTime() - startDate.getTime()) / (1000 * 60) // minutes

      return {
        id: suggestion.id,
        title: suggestion.title || 'Untitled Event',
        start: startDate.toISOString(),
        end: endDate.toISOString(),
        backgroundColor: '#8b5cf6',
        borderColor: '#7c3aed',
        textColor: '#ffffff',
        allDay: duration >= 24 * 60,
        isSuggestion: true,
        className: 'fc-event-suggestion',
        description: `
          ${suggestion.title || 'Untitled Event'} (AI Suggestion)
          \nStart: ${startDate.toLocaleString()}
          \nEnd: ${endDate.toLocaleString()}
          \nClick to review or apply
        `.trim()
      }
    }).filter(event => event !== null)

    events.push(...suggestionEvents)
  }

  return events
})

const taskItems = computed(() => taskStore.tasks.map(task => task.title))

function toDateTimeLocal(date: Date) {
  new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)
}

// Helper function to round time to nearest interval
function roundToNearestInterval(date: Date, intervalMinutes: number = 15): Date {
  const roundedDate = new Date(date)
  const minutes = roundedDate.getMinutes()
  const roundedMinutes = Math.round(minutes / intervalMinutes) * intervalMinutes
  roundedDate.setMinutes(roundedMinutes)
  roundedDate.setSeconds(0, 0)
  return roundedDate

}

// Helper function to get default event duration based on calendar view
function getDefaultDuration(view: string): number {
  switch (view) {
    case 'month':
      return 60 * 60 * 1000 // 1 hour
    case 'week':
      return 30 * 60 * 1000 // 30 minutes
    case 'day':
      return 60 * 60 * 1000 // 1 hour
    default:
      return 60 * 60 * 1000 // 1 hour
  }

}

function onDateSelect(event) {
  modalShow.value = true

  // Handle different possible event structures
  let cursor = event
  if (event && event.cursor) {
    cursor = event.cursor
  }

  try {
    // Extract the clicked date/time from cursor
    let clickedDate

    if (cursor.date) {
      clickedDate = roundToNearestInterval(new Date(cursor.date))
    } else if (cursor.start) {
      clickedDate = roundToNearestInterval(new Date(cursor.start))
    } else if (typeof cursor === 'string' || cursor instanceof Date) {
      clickedDate = roundToNearestInterval(new Date(cursor))
    } else {
      return
    }

    // Validate the date
    if (isNaN(clickedDate.getTime())) {
      toast.add({ title: "Error", description: "Invalid date selected", color: "error" })
      return
    }

    // Use the cursor date/time directly without defaults
    formData.start = toDateTimeLocal(clickedDate)

    // Set end time to 1 hour after start by default
    const endDate = new Date(clickedDate.getTime() + 60 * 60 * 1000)
    formData.end = toDateTimeLocal(endDate)
    // Add success feedback
    toast.add({ title: "Calendar", description: `Selected time: ${clickedDate.toLocaleDateString()}` })
  } catch (error) {
    toast.add({ title: "Error", description: "Error selecting time slot", color: "error" })
  }

}

async function onEventDrop(eventDropInfo) {
  const event = eventDropInfo.event

  if (event.extendedProps?.isSuggestion) {
    const moved = scheduleStore.updateSuggestion(
      event.id,
      event.start?.toISOString(),
      (event.end || event.start).toISOString()
    )
    if (!moved) {
      eventDropInfo.revert()
    }
    toast.add({ title: "Calendar", description: "Suggestion updated", color: "success" })
    return
  }

  const success = await scheduleStore.updateSchedule(event.id, {
    start: event.start,
    end: event.end || event.start
  })

  if (!success) {
    toast.add({ title: "Error", description: scheduleStore.fetchError?.message || "Failed to update event", color: "error" })
    eventDropInfo.revert()
  } else {
    toast.add({ title: "Calendar", description: "Event updated successfully", color: "success" })
    await fetchEvents()
  }
}

async function onEventClick(event) {
  const fullEvent = event.event
  if (fullEvent.extendedProps?.isSuggestion) {
    openSuggestionsModal()
    return
  }
  selectedEvent.value = fullEvent
  eventModalShow.value = true
}

function close() {
  modalShow.value = false
  // Reset form data
  formData.title = ''
  formData.start = ''
  formData.end = ''
  formData.recurrenceRuleId = null

  // Reset recurrence form
  recurrenceForm.frequency = 'daily'
  recurrenceForm.interval = 1
  recurrenceForm.endDate = null
  recurrenceForm.count = null
}

function closeEventModal() {
  eventModalShow.value = false
  selectedEvent.value = null
}

async function aiSuggestion() {
  suggestionsLoading.value = true
  try {
    const result = await scheduleStore.getSuggestions()
    if (result && result.length > 0) {
      openSuggestionsModal()
    } else {
      toast.add({ title: "Calendar", description: "No AI suggestions available. Your tasks may already be scheduled.", color: "warning" })
    }
  } finally {
    suggestionsLoading.value = false
  }
}

function openSuggestionsModal() {
  selectedSuggestionIds.value = scheduleStore.suggestions.map(s => s.id)
  suggestionModalShow.value = true
}

function toggleSuggestionSelection(id: string) {
  const idx = selectedSuggestionIds.value.indexOf(id)
  if (idx === -1) {
    selectedSuggestionIds.value.push(id)
  } else {
    selectedSuggestionIds.value.splice(idx, 1)
  }
}

function toggleSelectAllSuggestions() {
  if (selectedSuggestionIds.value.length === scheduleStore.suggestions.length) {
    selectedSuggestionIds.value = []
  } else {
    selectedSuggestionIds.value = scheduleStore.suggestions.map(s => s.id)
  }
}

async function dismissSuggestion(id: string) {
  scheduleStore.removeSuggestion(id)
  selectedSuggestionIds.value = selectedSuggestionIds.value.filter(sid => sid !== id)
}

async function applySelectedSuggestions() {
  if (selectedSuggestionIds.value.length === 0) {
    toast.add({ title: "Warning", description: "No suggestions selected", color: "warning" })
    return
  }

  const ok = await scheduleStore.applySuggestions([...selectedSuggestionIds.value])
  if (ok) {
    selectedSuggestionIds.value = scheduleStore.suggestions.map(s => s.id)
    if (scheduleStore.suggestions.length === 0) {
      suggestionModalShow.value = false
    }
  }
}

async function rejectSelectedSuggestions() {
  const toReject = [...selectedSuggestionIds.value]
  if (toReject.length === 0) {
    toast.add({ title: "Warning", description: "No suggestions selected", color: "warning" })
    return
  }

  for (const id of toReject) {
    scheduleStore.removeSuggestion(id)
  }
  selectedSuggestionIds.value = scheduleStore.suggestions.map(s => s.id)
  if (scheduleStore.suggestions.length === 0) {
    suggestionModalShow.value = false
  }
  toast.add({ title: "Calendar", description: `Dismissed ${toReject.length} suggestions`, color: "success" })
}

async function save() {
  // Basic validation
  if (!formData.title || !formData.title.trim()) {
    toast.add({ title: "Error", description: "Please enter a title for the event", color: "error" })
    return
  }

  if (!formData.start || !formData.end) {
    toast.add({ title: "Error", description: "Please select both start and end times", color: "error" })
    return
  }

  const startDate = new Date(formData.start)
  const endDate = new Date(formData.end)

  if (startDate >= endDate) {
    toast.add({ title: "Error", description: "End time must be after start time", color: "error" })
    return
  }

  // Check for past dates (allow events up to 1 hour in the past for flexibility)
  const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000)
  if (startDate < oneHourAgo) {
    toast.add({ title: "Error", description: "Cannot schedule ecents in the past", color: "error" })
    return
  }

  // Check for conflicts
  const conflictList = await checkConflicts(startDate, endDate)
  if (conflictList.length > 0) {
    showConflicts(conflictList)
    return
  }

  try {
    const dataToSend = {
      title: formData.title.trim(),
      start: startDate.toISOString(),
      end: endDate.toISOString(),
      recurrence_rule_id: formData.recurrenceRuleId,
    }

    await scheduleStore.submitForm(dataToSend)
    toast.add({ title: "Calendar", description: "Event created successfully", color: "success" })
    await fetchEvents()
    close()
  } catch (error) {
    toast.add({ title: "Error", description: error?.messsage || "Failed to create event", color: "error" })
  }
}

// Bulk operations
function toggleBulkMode() {
  bulkMode.value = !bulkMode.value
  if (!bulkMode.value) {
    selectedEvents.value = []
  }
}

async function deleteEvent(eventId) {
  try {
    await scheduleStore.deleteSchedule(eventId)
    toast.add({ title: "Calendar", description: "Event deleted successfully", color: "success" })
    eventModalShow.value = false
    selectedEvent.value = null
    await fetchEvents()
  } catch (error) {
    toast.add({ title: "Error", description: "Failed to delete event", color: "error" })
  }
}

async function bulkDeleteSelected() {
  if (selectedEvents.value.length === 0) {
    toast.add({ title: "Warning", description: "No events selected", color: "warning" })
    return
  }

  try {
    const { $api } = useNuxtApp()
    await $api('schedule/bulk', {
      method: 'DELETE',
      body: JSON.stringify({ ids: selectedEvents.value })
    })
    toast.add({ title: "Calendar", description: `Deleted ${selectedEvents.value.length} events`, color: "success" })
    selectedEvents.value = []
    bulkMode.value = false
    await fetchEvents()
  } catch (error) {
    toast.add({ title: "Error", description: "Failed to delete selected events", color: "error" })
  }
}

// Conflict detection
async function checkConflicts(start, end, excludeId = null) {
  try {
    const { $api } = useNuxtApp()
    const response = await $api('schedule')
    const existingEvents = response.schedule || []

    const conflictsFound = existingEvents.filter(event => {
      if (excludeId && event.id === excludeId) return false
      const eventStart = new Date(event.start)
      const eventEnd = new Date(event.end)
      return (start < eventEnd && end > eventStart) ||
        (eventStart < end && eventEnd > start)
    })

    return conflictsFound
  } catch (error) {
    console.error('Error checking conflicts:', error)
    return []
  }
}


// Watch for schedule changes
watch(() => scheduleStore.schedule, () => {
  calendarKey.value++
}, { deep: true })

onMounted(async () => {
  scheduleStore.loadSuggestionsFromStorage()
  if (scheduleStore.schedule.length === 0) {
    await fetchEvents()
  }
})
</script>

<template>
  <div class="calendar-container p-2">
    <div class="flex justify-between items-center mb-4">
      <h2>Calendar</h2>
      <div class="flex gap-2">
        <!-- <CalendarZonesManager /> -->
        <UButton v-if="bulkMode" :label="`Delete Selected (${selectedEvents.length})`" color="error"
          :disabled="selectedEvents.length === 0" @click="bulkDeleteSelected" />
        <UButton :label="bulkMode ? 'Cancel Bulk' : 'Bulk Select'" @click="toggleBulkMode" />
        <UButton v-if="scheduleStore.suggestions.length > 0" color="secondary" variant="soft"
          :label="`AI Suggestions (${scheduleStore.suggestions.length})`" @click="openSuggestionsModal" />
        <UButton label="AI suggestion" :loading="suggestionsLoading" @click="aiSuggestion" color="neutral"
          variant="soft" />
      </div>
    </div>



    <!-- FullCalendar -->
    <FullCalendar :key="calendarKey" :options="calendarOptions" />

    <!-- Event creation modal -->
    <UModal v-model:open="modalShow" title="Create Event">
      <template #body>
        <UForm :state="formData" class="space-y-4">
          <UFormField label="Title">
            <USelectMenu v-model="formData.title" :items="taskItems" placeholder="Enter or select title" size="xl"
              class="w-full" />
          </UFormField>

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
          <div class="flex justify-end space-x-2 mt-6">
            <UButton label="Cancel" variant="outline" color="neutral" @click="modalShow = false" />
            <UButton label="Save" variant="soft" color="neutral" @click="save" />
          </div>
        </UForm>
      </template>
    </UModal>



    <!-- Conflict Detection Modal -->
    <div v-if="conflictModalShow"
      class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-[9999]">
      <div
        class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4 text-red-600 dark:text-red-400">Scheduling Conflict Detected</h3>
        <p class="text-gray-700 dark:text-gray-300 mb-4">
          The following events conflict with your proposed time slot:
        </p>
        <div class="space-y-2 mb-6">
          <div v-for="conflict in conflicts" :key="conflict.id"
            class="border border-red-200 dark:border-red-600 rounded-lg p-3 bg-red-50 dark:bg-red-900/20">
            <h4 class="font-semibold text-gray-900 dark:text-gray-100">{{ conflict.title }}</h4>
            <p class="text-sm text-gray-700 dark:text-gray-300">
              {{ new Date(conflict.start).toLocaleString() }} - {{ new Date(conflict.end).toLocaleString() }}
            </p>
          </div>
        </div>
        <div class="flex justify-end space-x-2">
          <UButton label="Cancel" color="error" @click="conflictModalShow = false" />
          <UButton label="Create Anyway" @click="save" />
        </div>
      </div>
    </div>

    <!-- Event Details Modal -->
    <div v-if="eventModalShow && selectedEvent"
      class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-[9999]">
      <div
        class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4">Event Details</h3>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Title</label>
            <p class="text-gray-900 dark:text-gray-100 font-medium">{{ selectedEvent.title }}</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Start Time</label>
            <p class="text-gray-900 dark:text-gray-100">{{ new Date(selectedEvent.start).toLocaleString() }}</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">End Time</label>
            <p class="text-gray-900 dark:text-gray-100">{{ new Date(selectedEvent.end).toLocaleString() }}</p>
          </div>

          <div v-if="selectedEvent.extendedProps?.content">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Linked Task</label>
            <p class="text-gray-900 dark:text-gray-100">{{ selectedEvent.extendedProps.content }}</p>
          </div>
        </div>

        <div class="flex justify-end space-x-2 mt-6">
          <UButton label="Close" variant="outline" @click="closeEventModal" />
          <UButton label="Delete Event" variant="soft" color="error" @click="deleteEvent(selectedEvent.id)" />
        </div>
      </div>
    </div>

    <!-- AI Suggestions Modal -->
    <div v-if="suggestionModalShow"
      class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-[9999]">
      <div
        class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold">AI Schedule Suggestions</h3>
          <UButton label="✕" variant="ghost" color="neutral" @click="suggestionModalShow = false" />
        </div>

        <p class="text-sm text-gray-600 dark:text-gray-300 mb-3">
          AI suggested time slots for your tasks. Select the ones to add to your schedule, or dismiss them
          individually.
        </p>

        <div class="flex items-center justify-between mb-3">
          <UButton
            :label="selectedSuggestionIds.length === scheduleStore.suggestions.length ? 'Unselect All' : 'Select All'"
            variant="outline" size="xs" @click="toggleSelectAllSuggestions" />
          <span class="text-sm text-gray-600 dark:text-gray-300">
            {{ selectedSuggestionIds.length }} / {{ scheduleStore.suggestions.length }} selected
          </span>
        </div>

        <div v-if="scheduleStore.suggestions.length === 0" class="text-center py-12 text-gray-500">
          No pending suggestions.
        </div>

        <div v-else class="space-y-2 mb-6">
          <div v-for="suggestion in scheduleStore.suggestions" :key="suggestion.id"
            class="border border-purple-200 dark:border-purple-600 rounded-lg p-3 bg-purple-50 dark:bg-purple-900/20">
            <div class="flex items-start justify-between gap-3">
              <div class="flex items-start gap-3 min-w-0">
                <input type="checkbox" class="mt-1 size-4 accent-purple-500 shrink-0"
                  :checked="selectedSuggestionIds.includes(suggestion.id)"
                  @change="toggleSuggestionSelection(suggestion.id)" />
                <div class="min-w-0">
                  <h4 class="font-semibold text-gray-900 dark:text-gray-100 truncate">{{ suggestion.title }}</h4>
                  <p class="text-sm text-gray-700 dark:text-gray-300">
                    {{ new Date(suggestion.start).toLocaleString() }} -
                    {{ new Date(suggestion.end).toLocaleString() }}
                  </p>
                </div>
              </div>
              <UButton label="Dismiss" variant="ghost" size="sm" color="neutral"
                @click="dismissSuggestion(suggestion.id)" />
            </div>
          </div>
        </div>

        <div class="flex justify-between items-center mt-6">
          <UButton label="Dismiss Selected" variant="outline" color="neutral"
            :disabled="selectedSuggestionIds.length === 0" @click="rejectSelectedSuggestions" />
          <div class="space-x-2">
            <UButton label="Close" variant="outline" color="neutral" @click="suggestionModalShow = false" />
            <UButton label="Add Selected" variant="soft" color="success"
              :disabled="selectedSuggestionIds.length === 0" @click="applySelectedSuggestions" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* FullCalendar custom styles */
:deep(.fc-header-toolbar) {
  padding: 0.75rem;
  margin-bottom: 0 !important;
  border-radius: 0.5rem 0.5rem 0 0;
}

:deep(.fc-toolbar-title) {
  font-weight: 600;
}

:deep(.fc-button) {
  background: rgba(255, 255, 255, 0.2) !important;
  border: none !important;
  color: white !important;
  border-radius: 0.375rem !important;
  font-weight: 500 !important;
}

:deep(.fc-view-harness) {
  background: white;
  border-radius: 0 0 0.5rem 0.5rem;
  overflow: hidden;
}

:deep(.fc-day-today) {
  background: rgba(59, 130, 246, 0.1) !important;
}

:deep(.fc-event) {
  border-radius: 6px !important;
  font-size: 0.875rem !important;
  font-weight: 500 !important;
  padding: 2px 6px !important;
  border: none !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2) !important;
}

:deep(.fc-event-suggestion) {
  border: 1px dashed rgba(255, 255, 255, 0.6) !important;
  opacity: 0.9;
}

:deep(.fc-daygrid-day) {
  border: 1px solid #e5e7eb !important;
}

:deep(.fc-daygrid-day-number) {
  color: #374151;
  font-weight: 500;
  padding: 4px;
}

:deep(.fc-col-header-cell) {
  background: #f9fafb;
  border-bottom: 2px solid #d1d5db;
  color: #111827;
  font-weight: 600;
  padding: 8px;
}

:deep(.fc-timegrid-slot) {
  border-bottom: 1px solid #e5e7eb;
  color: black
}

:deep(.fc-timegrid-axis) {
  border-right: 2px solid #e5e7eb;
  background: #ffffff;
  color: #374151;
  font-weight: 500;
}

/* Dark mode support 
.dark :deep(.fc-header-toolbar) {
  background: linear-gradient(135deg, #dadee3 0%, #581c87 100%);
  color: black;
}

.dark :deep(.fc .fc-view-harness) {
  background: #1f2937;
}

.dark :deep(.fc .fc-day-today) {
  background: rgba(59, 130, 246, 0.2) !important;
}

.dark :deep(.fc .fc-daygrid-day) {
  border-color: black !important;
}

.dark :deep(.fc .fc-daygrid-day-number) {
  color: black;
}

.dark :deep(.fc .fc-col-header-cell) {
  background: #111827;
  border-bottom-color: #4b5563;
  color: black;
}

.dark :deep(.fc-timegrid-slot) {
  border-bottom-color: red !important;
}

.dark :deep(.fc .fc-timegrid-axis) {
  border-right-color: #4b5563;
  background: #1f2937 !important;
  color: #d1d5db;
}

.dark :deep(.fc .fc-event) {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4) !important;
}

/* Selection overlay styles
:deep(.fc .fc-highlight) {
  background: rgba(59, 130, 246, 0.2) !important;
  border: 2px dashed rgba(59, 130, 246, 0.5) !important;
  border-radius: 4px !important;
}

.dark :deep(.fc .fc-highlight) {
  background: rgba(59, 130, 246, 0.3) !important;
  border-color: rgba(59, 130, 246, 0.7) !important;
}
*/
</style>
