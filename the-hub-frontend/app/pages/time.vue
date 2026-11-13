<script setup lang="ts">
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin, { Draggable } from '@fullcalendar/interaction'
const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const toast = useToast()

const modalShow = ref(false)
const calendarView = ref('timeGridWeek')
const route = useRoute()

// Slot-picking mode: arriving via ?schedule=<task_id> from a task card
const slotPickTaskId = ref<string | null>(null)
const slotPickTask = computed(() =>
  slotPickTaskId.value ? taskStore.tasks.find(t => t.task_id === slotPickTaskId.value) ?? null : null
)

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
  eventReceive: onEventReceive,
  select: onDateSelect,
}))

const isLoading = ref(false)
const selectedEvents = ref([])
const bulkMode = ref(false)
const conflictModalShow = ref(false)
const conflicts = ref([])
const pendingDrop = ref<null | { title: string; start: string; end: string; task_id?: string | null }>(null)
const eventModalShow = ref(false)
const selectedEvent = ref<any>(null)

const calendarKey = ref(0)

// Unscheduled task rail (drag tasks onto the calendar to schedule them)
const railOpen = ref(true)
const railListRef = ref<HTMLElement | null>(null)
let taskDraggable: Draggable | null = null

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
        extendedProps: { task_id: event.task_id },
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
        extendedProps: { task_id: suggestion.task_id },
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

// Task picker options for the create-event modal
const taskOptions = computed(() => taskStore.tasks.map(task => ({
  task_id: task.task_id,
  title: task.title,
  priority: task.priority ?? 3,
  estimate: task.time_estimate_minutes ?? 0,
})))

const selectedTaskId = ref<string | null>(null)
const useCustomTitle = ref(false)
const customTitle = ref('')

const selectedTask = computed(() => taskStore.tasks.find(t => t.task_id === selectedTaskId.value) ?? null)

// Pending tasks that are not yet linked to a schedule event
const scheduledTaskIds = computed(() => new Set(scheduleStore.schedule.map(e => e.task_id).filter(Boolean)))
const unscheduledTasks = computed(() =>
  taskStore.tasks.filter(t =>
    t.status === 'pending'
    && !t.parent_task_id
    && !scheduledTaskIds.value.has(t.task_id)
  )
)

const linkedTask = computed(() => {
  const taskId = selectedEvent.value?.extendedProps?.task_id
  if (!taskId) return null
  return taskStore.tasks.find(t => t.task_id === taskId) ?? null
})

function toDateTimeLocal(date: Date) {
  return new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)
}

function formatEstimate(minutes: number): string {
  if (minutes >= 60) {
    const h = Math.floor(minutes / 60)
    const m = minutes % 60
    return m ? `${h}h ${m}m` : `${h}h`
  }
  return `${minutes}m`
}

function toggleTitleMode() {
  useCustomTitle.value = !useCustomTitle.value
  if (useCustomTitle.value) {
    customTitle.value = selectedTask.value?.title ?? ''
    selectedTaskId.value = null
  } else {
    customTitle.value = ''
  }
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

    // In slot-picking mode (?schedule=<task_id>), preselect that task
    if (slotPickTaskId.value) {
      selectedTaskId.value = slotPickTaskId.value
      useCustomTitle.value = false
    }

    // Use the cursor date/time directly without defaults
    formData.start = toDateTimeLocal(clickedDate)

    // Set end time to 1 hour after start by default (or the task's estimate)
    const defaultDuration = slotPickTask.value?.time_estimate_minutes || 60
    const endDate = new Date(clickedDate.getTime() + defaultDuration * 60 * 1000)
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
    // Inline quick-review: accept/reject this single suggestion in context
    inlineSuggestionId.value = fullEvent.id
    return
  }
  selectedEvent.value = fullEvent
  eventModalShow.value = true
}

function close() {
  modalShow.value = false
  // Reset form data
  formData.start = ''
  formData.end = ''
  formData.recurrenceRuleId = null
  selectedTaskId.value = null
  useCustomTitle.value = false
  customTitle.value = ''
  slotPickTaskId.value = null

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

// Inline per-suggestion quick review (Phase 4): clicking a ghost event shows
// accept/reject controls without leaving the calendar context.
const inlineSuggestionId = ref<string | null>(null)
const inlineSuggestion = computed(() =>
  inlineSuggestionId.value
    ? scheduleStore.suggestions.find(s => s.id === inlineSuggestionId.value) ?? null
    : null
)

async function acceptSuggestion(id: string) {
  const ok = await scheduleStore.applySuggestions([id])
  if (ok) {
    inlineSuggestionId.value = null
    await fetchEvents()
    toast.add({ title: "Calendar", description: "Suggestion added to schedule", color: "success" })
  }
}

function rejectSuggestion(id: string) {
  scheduleStore.removeSuggestion(id)
  inlineSuggestionId.value = null
  toast.add({ title: "Calendar", description: "Suggestion dismissed", color: "success" })
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

async function save(force = false) {
  const isTaskMode = !!selectedTaskId.value && !useCustomTitle.value
  const resolvedTitle = isTaskMode
    ? selectedTask.value?.title
    : customTitle.value.trim()

  // Basic validation
  if (!resolvedTitle) {
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
    toast.add({ title: "Error", description: "Cannot schedule events in the past", color: "error" })
    return
  }

  // Check for conflicts (skippable via force=true from the conflict modal)
  if (!force) {
    const conflictList = await checkConflicts(startDate, endDate)
    if (conflictList.length > 0) {
      showConflicts(conflictList)
      return
    }
  }

  try {
    const dataToSend = {
      title: resolvedTitle,
      start: startDate.toISOString(),
      end: endDate.toISOString(),
      task_id: isTaskMode ? selectedTaskId.value : null,
      recurrence_rule_id: formData.recurrenceRuleId,
    }

    await scheduleStore.submitForm(dataToSend)
    toast.add({ title: "Calendar", description: "Event created successfully", color: "success" })
    await fetchEvents()
    close()
  } catch (error) {
    toast.add({ title: "Error", description: error?.message || "Failed to create event", color: "error" })
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


// Drag a task from the rail onto the calendar to schedule it
function setupTaskDraggable() {
  taskDraggable?.destroy()
  taskDraggable = null
  if (!railOpen.value || !railListRef.value || unscheduledTasks.value.length === 0) return
  taskDraggable = new Draggable(railListRef.value, {
    itemSelector: '.task-rail-item',
    eventData(dropInfo) {
      const taskId = (dropInfo.draggedEl as HTMLElement | null)?.dataset?.taskId
      const task = taskStore.tasks.find(t => t.task_id === taskId)
      if (!task) return false
      const est = task.time_estimate_minutes && task.time_estimate_minutes > 0
        ? task.time_estimate_minutes
        : 30
      return {
        title: task.title,
        duration: { minutes: est },
      }
    },
  })
}

// An external drop from the rail: create the schedule event pre-linked to the task
async function onEventReceive(info) {
  const taskId = (info.draggedEl as HTMLElement | null)?.dataset?.taskId as string | undefined
  const task = taskStore.tasks.find(t => t.task_id === taskId)
  const start = info.event?.start ? new Date(info.event.start) : null
  if (!start) return

  const est = task?.time_estimate_minutes && task.time_estimate_minutes > 0
    ? task.time_estimate_minutes
    : 30
  const end = info.event.end ? new Date(info.event.end) : new Date(start.getTime() + est * 60 * 1000)

  // Dropping into the past isn't allowed
  const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000)
  if (start < oneHourAgo) {
    eventRemoved(info)
    toast.add({ title: "Calendar", description: "Cannot schedule events in the past", color: "error" })
    return
  }

  // Conflicts: confirm via the existing conflict modal before creating
  const overlapping = scheduleStore.schedule.filter(e => start < e.end && end > e.start)
  if (overlapping.length > 0) {
    eventRemoved(info)
    pendingDrop.value = {
      title: task?.title || info.event.title,
      start: start.toISOString(),
      end: end.toISOString(),
      task_id: taskId
    }
    conflicts.value = overlapping
    conflictModalShow.value = true
    return
  }

  // Create the event pre-linked to the task. submitForm handles the optimistic
  // add, error rollback, and toasts internally; the store watch re-renders the
  // calendar so the FullCalendar-managed received event is replaced by store data.
  scheduleStore.submitForm({
    title: task?.title || info.event.title,
    start: start.toISOString(),
    end: end.toISOString(),
    task_id: taskId,
  })
  await fetchEvents()
}

function eventRemoved(info) {
  try {
    info.event.remove()
  } catch {
    // event already detached
  }
}

function confirmCreateAnyway() {
  conflictModalShow.value = false
  if (pendingDrop.value) {
    const drop = pendingDrop.value
    pendingDrop.value = null
    scheduleStore.submitForm(drop).then(ok => {
      if (ok) fetchEvents()
    })
    return
  }
  save(true)
}

function cancelConflict() {
  conflictModalShow.value = false
  pendingDrop.value = null
}

// Watch for schedule changes
watch(() => scheduleStore.schedule, () => {
  calendarKey.value++
}, { deep: true })

// Re-init the external draggable whenever the rail contents change
watch([railOpen, unscheduledTasks, railListRef], () => {
  setupTaskDraggable()
}, { flush: 'post' })

onBeforeUnmount(() => {
  taskDraggable?.destroy()
  taskDraggable = null
})

onMounted(async () => {
  scheduleStore.loadSuggestionsFromStorage()
  if (scheduleStore.schedule.length === 0) {
    await fetchEvents()
  }

  // Arrived from a task card's "Schedule" action → slot-picking mode
  const scheduleTaskId = route.query.schedule
  if (typeof scheduleTaskId === 'string' && scheduleTaskId) {
    slotPickTaskId.value = scheduleTaskId
    // Clean the query so a refresh doesn't re-enter picking mode
    navigateTo({ path: '/time', query: {}, replace: true })
    toast.add({
      title: 'Schedule task',
      description: `Click a time slot to schedule "${slotPickTask.value?.title || 'this task'}"`,
      color: 'primary',
      timeout: 6000
    })
  }
})
</script>

<template>
  <div class="calendar-container p-2">
    <div class="flex justify-between items-center mb-4">
      <h2>Calendar</h2>
      <div class="flex gap-2">
        <UButton :icon="railOpen ? 'i-lucide-panel-left-close' : 'i-lucide-panel-left-open'" variant="ghost"
          color="neutral" title="Toggle unscheduled task rail" @click="railOpen = !railOpen" />
        <!-- <CalendarZonesManager /> -->
        <UButton v-if="bulkMode" :label="`Delete Selected (${selectedEvents.length})`" color="error"
          :disabled="selectedEvents.length === 0" @click="bulkDeleteSelected" />
        <UButton :label="bulkMode ? 'Cancel Bulk' : 'Bulk Select'" @click="toggleBulkMode" />
        <UButton label="AI suggestion" :loading="suggestionsLoading" @click="aiSuggestion" color="neutral"
          variant="soft" />
      </div>
    </div>



    <!-- Main layout: unscheduled task rail + calendar -->
    <div class="flex flex-col lg:flex-row gap-3 items-start">

      <!-- Unscheduled task rail (Phase 2) -->
      <aside v-show="railOpen"
        class="w-full lg:w-72 shrink-0 self-stretch rounded-lg border border-gray-200 dark:border-gray-700 bg-surface-light dark:bg-surface-dark p-3 flex flex-col max-h-[75vh]">
        <h3 class="text-sm font-semibold mb-1">
          Unscheduled tasks
          <span class="text-gray-500 dark:text-gray-400 font-normal">({{ unscheduledTasks.length }})</span>
        </h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">Drag a task onto the calendar to schedule it.</p>

        <div ref="railListRef" class="flex-1 min-h-0 overflow-y-auto space-y-2 pr-1">
          <div v-for="task in unscheduledTasks" :key="task.task_id" :data-task-id="task.task_id"
            class="task-rail-item rounded-md border border-gray-200 dark:border-gray-700 p-2.5 bg-background-light dark:bg-background-dark">
            <p class="text-sm font-medium truncate">{{ task.title }}</p>
            <div class="flex flex-wrap items-center gap-1 mt-1.5">
              <UBadge size="xs" variant="soft" color="neutral">P{{ task.priority ?? 3 }}</UBadge>
              <UBadge v-if="task.time_estimate_minutes" size="xs" variant="soft" color="neutral">
                Est: {{ formatEstimate(task.time_estimate_minutes) }}
              </UBadge>
            </div>
            <p v-if="task.due_date" class="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Due {{ $dayjs(task.due_date).fromNow() }}
            </p>
          </div>
          <p v-if="unscheduledTasks.length === 0"
            class="text-sm text-gray-500 dark:text-gray-400 text-center py-6">
            All pending tasks are scheduled — nice work.
          </p>
        </div>
      </aside>

      <!-- Calendar -->
      <div class="flex-1 min-w-0 w-full">
        <!-- FullCalendar -->
        <FullCalendar :key="calendarKey" :options="calendarOptions" />
      </div>
    </div>

    <!-- Event creation modal -->
    <UModal v-model:open="modalShow" title="Create Event">
      <template #body>
        <UForm :state="formData" class="space-y-4">
          <UFormField
            :label="useCustomTitle ? 'Title' : 'Task'"
            :description="useCustomTitle ? undefined : 'Optional — link this event to one of your tasks'">
            <USelectMenu v-if="!useCustomTitle" v-model="selectedTaskId" :items="taskOptions"
              value-key="task_id" label-key="title" placeholder="Select a task or leave empty" size="xl"
              class="w-full">
              <template #item-label="{ item }">
                <span class="truncate">{{ item.title }}</span>
              </template>
              <template #item-trailing="{ item }">
                <span class="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">
                  P{{ item.priority }}{{ item.estimate ? ` · ${formatEstimate(item.estimate)}` : '' }}
                </span>
              </template>
              <template #empty>
                <span class="text-sm text-gray-500 dark:text-gray-400">
                  No tasks found — use the link below for a custom title
                </span>
              </template>
            </USelectMenu>
            <UInput v-else v-model="customTitle" placeholder="Event title" size="xl" class="w-full" />
            <UButton type="button" variant="link" size="xs" color="primary" class="mt-1 px-0"
              :label="useCustomTitle ? 'Choose a task instead' : 'Not a task? Enter a custom title'"
              @click="toggleTitleMode" />
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

          <div v-if="linkedTask"
            class="rounded-md border border-gray-200 dark:border-gray-700 p-3 bg-background-light dark:bg-background-dark">
            <div class="flex items-center justify-between mb-1">
              <span class="text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400">Linked Task</span>
              <span v-if="linkedTask.status"
                class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 capitalize">
                {{ linkedTask.status }}
              </span>
            </div>
            <p class="font-medium text-gray-900 dark:text-gray-100">{{ linkedTask.title }}</p>
            <p v-if="linkedTask.due_date" class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Due {{ $dayjs(linkedTask.due_date).format('MMM D, YYYY') }}
              ({{ $dayjs(linkedTask.due_date).fromNow() }})
            </p>
            <NuxtLink to="/plan"
              class="inline-flex items-center gap-1 text-sm text-primary hover:underline mt-2">
              View on Plan
              <UIcon name="i-lucide-arrow-up-right" class="size-4" />
            </NuxtLink>
          </div>
        </div>

        <div class="flex justify-end space-x-2 mt-6">
          <UButton label="Close" variant="outline" @click="closeEventModal" />
          <UButton label="Delete Event" variant="soft" color="error" @click="deleteEvent(selectedEvent.id)" />
        </div>
      </div>
    </div>

    <!-- Floating AI suggestions pill (Phase 4) -->
    <Transition name="fade">
      <button v-if="scheduleStore.suggestions.length > 0" type="button" @click="openSuggestionsModal"
        class="suggestions-pill fixed bottom-6 right-6 z-40 flex items-center gap-2 rounded-full px-4 py-2.5 text-sm font-semibold text-white shadow-lg">
        <UIcon name="i-lucide-sparkles" class="size-4" />
        Review {{ scheduleStore.suggestions.length }} suggestion{{ scheduleStore.suggestions.length === 1 ? '' : 's' }}
      </button>
    </Transition>

    <!-- Inline per-suggestion quick review (Phase 4) -->
    <Transition name="fade">
      <div v-if="inlineSuggestion"
        class="inline-suggestion-card fixed bottom-6 left-1/2 -translate-x-1/2 z-50 w-[22rem] max-w-[90vw] rounded-xl border border-purple-200 dark:border-purple-600 bg-surface-light dark:bg-surface-dark p-4 shadow-xl">
        <p class="text-xs font-semibold uppercase tracking-wide text-purple-600 dark:text-purple-400 mb-1">AI suggestion</p>
        <h4 class="font-semibold text-gray-900 dark:text-gray-100 truncate">{{ inlineSuggestion.title }}</h4>
        <p class="text-sm text-gray-600 dark:text-gray-300 mt-0.5">
          {{ new Date(inlineSuggestion.start).toLocaleString([], { weekday: 'short', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
          – {{ new Date(inlineSuggestion.end).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
        </p>
        <div class="flex justify-end gap-2 mt-3">
          <UButton label="Reject" variant="outline" size="sm" color="neutral" @click="rejectSuggestion(inlineSuggestion.id)" />
          <UButton label="Accept" variant="soft" size="sm" color="success" @click="acceptSuggestion(inlineSuggestion.id)" />
        </div>
      </div>
    </Transition>

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
/* FullCalendar custom styles — themed via --cal-* tokens (see main.css) */
:deep(.fc) {
  --fc-page-bg-color: var(--cal-bg);
  --fc-border-color: var(--cal-border);
  --fc-neutral-bg-color: var(--cal-bg-raised);
  --fc-today-bg-color: var(--cal-today-bg);
}

:deep(.fc-header-toolbar) {
  padding: 0.75rem;
  margin-bottom: 0 !important;
  border-radius: 0.5rem 0.5rem 0 0;
  background: var(--cal-bg);
  color: var(--cal-text);
}

:deep(.fc-toolbar-title) {
  font-weight: 600;
  color: var(--cal-text);
}

:deep(.fc-button) {
  background: transparent !important;
  border: 1px solid var(--cal-border) !important;
  color: var(--cal-text) !important;
  border-radius: 0.375rem !important;
  font-weight: 500 !important;
  box-shadow: none !important;
}

:deep(.fc-button:hover) {
  background: var(--cal-bg-raised) !important;
  color: var(--cal-text) !important;
}

:deep(.fc-button-active),
:deep(.fc-button:not(:disabled).fc-button-active) {
  background: var(--cal-bg-raised) !important;
  color: var(--cal-text) !important;
}

:deep(.fc-button:disabled) {
  opacity: 0.4;
}

:deep(.fc-view-harness) {
  background: var(--cal-bg);
  border-radius: 0 0 0.5rem 0.5rem;
  overflow: hidden;
}

:deep(.fc-day-today) {
  background: var(--cal-today-bg) !important;
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
  border: 1px dashed rgba(255, 255, 255, 0.5) !important;
  opacity: 0.9;
}

:deep(.fc-daygrid-day) {
  border: 1px solid var(--cal-border) !important;
}

:deep(.fc-daygrid-day-number) {
  color: var(--cal-text);
  font-weight: 500;
  padding: 4px;
}

:deep(.fc-col-header-cell) {
  background: var(--cal-bg-raised);
  border-bottom: 2px solid var(--cal-border);
  color: var(--cal-text);
  font-weight: 600;
  padding: 8px;
}

:deep(.fc-timegrid-slot) {
  border-bottom: 1px solid var(--cal-border);
  color: var(--cal-text-muted);
}

:deep(.fc-timegrid-axis) {
  border-right: 2px solid var(--cal-border);
  background: var(--cal-bg);
  color: var(--cal-text);
  font-weight: 500;
}

:deep(.fc-highlight) {
  background: var(--cal-today-bg) !important;
  border: 2px dashed var(--cal-border) !important;
  border-radius: 4px !important;
}

/* Drag-select & event hover accents */
:deep(.fc-event:hover) {
  filter: brightness(1.08);
}

/* Unscheduled task rail items */
.task-rail-item {
  cursor: grab;
  touch-action: none;
  user-select: none;
  transition: box-shadow 0.15s ease, border-color 0.15s ease, transform 0.15s ease;
}

.task-rail-item:hover {
  border-color: var(--color-primary);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.12);
}

.task-rail-item:active {
  cursor: grabbing;
  transform: scale(0.98);
}

/* Theme-aware scrollbars inside the calendar */
:deep(.fc-scroller) {
  scrollbar-color: var(--cal-border) transparent;
}

/* Floating AI suggestions pill (Phase 4) */
.suggestions-pill {
  animation: pill-pulse 2.5s ease-in-out infinite;
}

.suggestions-pill:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 18px rgba(139, 92, 246, 0.45);
}

@keyframes pill-pulse {
  0%, 100% { box-shadow: 0 4px 12px rgba(139, 92, 246, 0.35); }
  50% { box-shadow: 0 4px 22px rgba(139, 92, 246, 0.6); }
}

/* Fade transitions for the pill and inline quick-review card */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
