<script setup lang="ts">
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin from '@fullcalendar/interaction'

const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const calendarZonesStore = useCalendarZonesStore()
const { addToast } = useToast()

const modalShow = ref(false)
const aiModalShow = ref(false)
const aiSuggestions = ref([])
const calendarView = ref('timeGridWeek')
const currentDateRange = ref({ start: null, end: null })

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
  editable: true,
  selectable: true,
  selectMirror: true,
  dayMaxEvents: true,
  weekends: true,
  height: 'auto',
  eventDisplay: 'block',
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
  select: onDateSelect,
  viewDidMount: onViewChange,
}))

const isLoading = ref(false)
const isAISuggestionsLoading = ref(false)
const selectedEvents = ref([])
const bulkMode = ref(false)
const conflictModalShow = ref(false)
const conflicts = ref([])

async function fetchEvents() {
  isLoading.value = true
  try {
    if (scheduleStore.schedule.length === 0) {
      await scheduleStore.fetchSchedule()
    }
    if (taskStore.tasks.length === 0) {
      await taskStore.fetchTasks()
    }
    if (calendarZonesStore.zones.length === 0) {
      await calendarZonesStore.fetchZones()
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
        extendedProps: {
          task: event.task,
          created_by_ai: event.created_by_ai,
          content: event.task?.title ? `Task: ${event.task.title}` : ''
        },
        backgroundColor: event.created_by_ai ? '#3b82f6' : '#10b981',
        borderColor: event.created_by_ai ? '#2563eb' : '#059669',
        textColor: '#ffffff',
        allDay: duration >= 24 * 60, // If event is 24+ hours, show as all-day
        description: `
          ${event.title || 'Untitled Event'}
          ${event.task?.title ? `\nTask: ${event.task.title}` : ''}
          ${event.created_by_ai ? '\nAI Suggested' : ''}
          \nStart: ${startDate.toLocaleString()}
          \nEnd: ${endDate.toLocaleString()}
        `.trim()
      }
    }).filter(event => event !== null)

    events.push(...scheduleEvents)
  }

  // Add calendar zone background events
  if (calendarZonesStore.zones.length > 0) {
    const calendarZoneEvents = calendarZonesStore.zones
      .filter(zone => zone.is_active)
      .map(zone => {
        // Convert days_of_week string to FullCalendar daysOfWeek array
        // Sunday = 0, Monday = 1, Tuesday = 2, etc.
        const dayMap = {
          'sunday': 0,
          'monday': 1,
          'tuesday': 2,
          'wednesday': 3,
          'thursday': 4,
          'friday': 5,
          'saturday': 6
        }

        const daysOfWeek = zone.days_of_week
          ? zone.days_of_week.toLowerCase().split(',').map(day => dayMap[day.trim()]).filter(day => day !== undefined)
          : [0, 1, 2, 3, 4, 5, 6] // All days if no specific days set

        // Create background events for each day of the week
        const zoneEvents = []

        // For time grid views, create time-specific background events
        if (calendarView.value === 'timeGridWeek' || calendarView.value === 'timeGridDay') {
          // Create a recurring background event for the time slots
          zoneEvents.push({
            startTime: new Date(zone.start_time).getTime(),
            endTime: new Date(zone.end_time).getTime(),
            display: 'background',
            daysOfWeek: daysOfWeek,
            backgroundColor: zone.color ? `${zone.color}80` : '#3b82f620',
            title: zone.name,
            extendedProps: {
              isZoneBackground: true,
              zoneId: zone.id
            }
          })
        } else {
          // For month view, create all-day background events for each applicable day
          daysOfWeek.forEach(dayOfWeek => {
            zoneEvents.push({
              start: new Date(0, 0, dayOfWeek + 1), // Use day of week as date
              end: new Date(0, 0, dayOfWeek + 1),
              display: 'background',
              daysOfWeek: [dayOfWeek],
              backgroundColor: zone.color ? `${zone.color}15` : '#3b82f615',
              title: zone.name,
              allDay: true,
              extendedProps: {
                isZoneBackground: true,
                zoneId: zone.id
              }
            })
          })
        }

        return zoneEvents
      }).flat()

    events.push(...calendarZoneEvents)
  }
  return events
})





const toDateTimeLocal = (date: Date) =>
  new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)

// Helper function to round time to nearest interval
const roundToNearestInterval = (date: Date, intervalMinutes: number = 15): Date => {
  const roundedDate = new Date(date)
  const minutes = roundedDate.getMinutes()
  const roundedMinutes = Math.round(minutes / intervalMinutes) * intervalMinutes
  roundedDate.setMinutes(roundedMinutes)
  roundedDate.setSeconds(0, 0)
  return roundedDate
}

// Helper function to get default event duration based on calendar view
const getDefaultDuration = (view: string): number => {
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

const onDateSelect = (event) => {
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
      addToast('Invalid date selected', 'error')
      return
    }

    // Use the cursor date/time directly without defaults
    formData.start = toDateTimeLocal(clickedDate)

    // Set end time to 1 hour after start by default
    const endDate = new Date(clickedDate.getTime() + 60 * 60 * 1000)
    formData.end = toDateTimeLocal(endDate)
    // Add success feedback
    addToast(`Selected time: ${clickedDate.toLocaleString()}`, 'info')
  } catch (error) {
    addToast('Error selecting time slot', 'error')
  }
}

async function onEventDrop(eventDropInfo) {
  try {
    const event = eventDropInfo.event
    await $api(`schedule/${event.id}`, {
      method: 'PUT',
      body: {
        start: event.start.toISOString(),
        end: event.end?.toISOString() || event.start.toISOString()
      }
    })
    addToast('Event updated successfully', 'success')
    await fetchEvents()
  } catch (error) {
    addToast('Failed to update event', 'error')
    eventDropInfo.revert()
  }
}

async function onEventClick(event) {
  // Handle event click - could open edit modal
  console.log('Event clicked:', event)
  // For now, just show a toast with event details
  addToast(`Event: ${event.title}`, 'info')
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
  const roundedNow = roundToNearestInterval(now, 15)
  const duration = getDefaultDuration(calendarView.value)
  const endTime = new Date(roundedNow.getTime() + duration)

  formData.start = toDateTimeLocal(roundedNow)
  formData.end = toDateTimeLocal(endTime)
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
      task_id: formData.taskId,
      recurrence_rule_id: formData.recurrenceRuleId,
    }

    await scheduleStore.submitForm(dataToSend)
    addToast('Event created successfully', 'success')
    await fetchEvents()
    close()
  } catch (error) {
    addToast('Failed to create event', 'error')
  }
}

async function getAISuggestions() {
  isAISuggestionsLoading.value = true
  try {
    const { $api } = useNuxtApp()
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

// Bulk operations
function toggleBulkMode() {
  bulkMode.value = !bulkMode.value
  if (!bulkMode.value) {
    selectedEvents.value = []
  }
}

async function bulkDeleteSelected() {
  if (selectedEvents.value.length === 0) {
    addToast('No events selected', 'warning')
    return
  }

  try {
    const { $api } = useNuxtApp()
    await $api('schedule/bulk', {
      method: 'DELETE',
      body: JSON.stringify({ ids: selectedEvents.value })
    })
    addToast(`Deleted ${selectedEvents.value.length} events`, 'success')
    selectedEvents.value = []
    bulkMode.value = false
    await fetchEvents()
  } catch (error) {
    addToast('Failed to delete selected events', 'error')
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

async function showConflicts(conflictList) {
  conflicts.value = conflictList
  conflictModalShow.value = true
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

function onViewChange(viewChangeInfo) {
  // Update calendarView when FullCalendar changes view
  calendarView.value = viewChangeInfo.view.type
  fetchEvents()
}

onMounted(async () => {
  await fetchEvents()
})
</script>

<template>
  <div class="calendar-container p-2">
    <div class="flex justify-between items-center mb-4">
      <h2>Calendar</h2>
      <div class="flex gap-2">
        <CalendarZonesManager />
        <UiButton v-if="bulkMode" @click="bulkDeleteSelected" variant="danger" size="sm"
          :disabled="selectedEvents.length === 0">
          Delete Selected ({{ selectedEvents.length }})
        </UiButton>
        <UiButton @click="toggleBulkMode" :variant="bulkMode ? 'secondary' : 'outline'" size="sm">
          {{ bulkMode ? 'Cancel Bulk' : 'Bulk Select' }}
        </UiButton>
        <UiButton @click="getAISuggestions" variant="secondary" size="sm" :disabled="isAISuggestionsLoading">
          <span v-if="isAISuggestionsLoading">Loading...</span>
          <span v-else>Get AI Suggestions</span>
        </UiButton>
      </div>
    </div>



    <!-- FullCalendar -->
    <FullCalendar :options="calendarOptions" />

    <!-- Event creation modal -->
    <div v-if="modalShow" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-[9999]">
      <div
        class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4">Create Event</h3>
        <form class="space-y-4">
          <div>
            <label for="title" class="block text-sm font-medium mb-1">Title</label>
            <input type="text" id="title" v-model="formData.title" autocomplete="off"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>

          <div>
            <label for="task" class="block text-sm font-medium mb-1">Link to Task (optional)</label>
            <select id="task" v-model="formData.taskId"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-primary">
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
        </form>
        <div class="flex justify-end space-x-2 mt-6">
          <UiButton @click="close" variant="default" size="md">Cancel</UiButton>
          <UiButton @click="save" variant="primary" size="md">Save</UiButton>
        </div>
      </div>
    </div>

    <!-- AI Suggestions modal -->
    <div v-show="aiModalShow" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-50">
      <div
        class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold mb-4">AI Schedule Suggestions</h3>
        <div v-if="aiSuggestions.length === 0" class="text-center py-8">
          <p class="text-gray-700 dark:text-gray-300">No suggestions available at the moment.</p>
        </div>
        <div v-else class="space-y-4">
          <div v-for="suggestion in aiSuggestions" :key="suggestion.id"
            class="border border-gray-200 dark:border-gray-600 rounded-lg p-4 bg-gray-50 dark:bg-gray-800">
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
          <UiButton @click="conflictModalShow = false" variant="default" size="md">
            Cancel
          </UiButton>
          <UiButton @click="save" variant="primary" size="md">
            Create Anyway
          </UiButton>
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

:deep(.fc-button:hover) {
  background: rgba(255, 255, 255, 0.3) !important;
}

:deep(.fc-button:not(:disabled).fc-button-active) {
  background: rgba(255, 255, 255, 0.4) !important;
}

:deep(.fc-today-button) {
  background: #F97316 !important;
}

:deep(.fc-today-button:hover) {
  background: rgba(59, 130, 246, 0.9) !important;
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

:deep(.fc-event.ai-event) {
  background: linear-gradient(135deg, #1e40af, #1e3a8a) !important;
  color: #ffffff !important;
  font-weight: 600 !important;
}

:deep(.fc-event.regular-event) {
  background: linear-gradient(135deg, #065f46, #064e3b) !important;
  color: #ffffff !important;
  font-weight: 600 !important;
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
