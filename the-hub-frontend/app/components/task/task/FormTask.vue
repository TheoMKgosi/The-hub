<script setup lang="ts">
import * as chrono from 'chrono-node'

const taskStore = useTaskStore()
const { validateObject, schemas } = useValidation()

const formData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
  natural_language: '',
  use_natural_language: true,
  time_estimate_minutes: null,
})


const showForm = ref(true)
const validationErrors = ref<Record<string, string>>({})

// Real-time parsing feedback
const parsedPreview = computed(() => {
  if (!formData.use_natural_language || !formData.natural_language.trim()) {
    return null
  }

  const parsedDate = parseDateFromText(formData.natural_language)
  const parsedPriority = parsePriorityFromText(formData.natural_language)
  const parsedTimeEstimate = parseTimeEstimateFromText(formData.natural_language)
  const parsedSubtasks = parseSubtasksFromText(formData.natural_language)
  const parsedDependencies = parseDependenciesFromText(formData.natural_language)
  const parsedRecurring = parseRecurringFromText(formData.natural_language)
  const suggestions = getSuggestions(formData.natural_language)

  return {
    date: parsedDate,
    priority: parsedPriority,
    timeEstimate: parsedTimeEstimate,
    subtasks: parsedSubtasks,
    dependencies: parsedDependencies,
    recurring: parsedRecurring,
    suggestions
  }
})

const formatParsedDate = (date: Date) => {
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit'
  })
}

const getPriorityLabel = (priority: number) => {
  const labels = {
    1: 'Low',
    2: 'Low',
    3: 'Medium',
    4: 'High',
    5: 'Urgent'
  }
  return labels[priority as keyof typeof labels] || 'Medium'
}

const formatRecurringInfo = (recurring: any) => {
  if (!recurring) return null

  const { frequency, interval, byDay } = recurring

  if (byDay) {
    const dayNames = {
      MO: 'Monday',
      TU: 'Tuesday',
      WE: 'Wednesday',
      TH: 'Thursday',
      FR: 'Friday',
      SA: 'Saturday',
      SU: 'Sunday'
    }
    return `Every ${dayNames[byDay as keyof typeof dayNames]}`
  }

  if (interval === 1) {
    return `Every ${frequency.slice(0, -2)}` // Remove 's' from 'days', 'weeks', etc.
  }

  return `Every ${interval} ${frequency}`
}

// Provide helpful suggestions for natural language input
const getSuggestions = (text: string) => {
  if (!text || text.trim().length < 5) return []

  const suggestions = []
  const lowerText = text.toLowerCase()

  // Check for missing timing information
  if (!lowerText.includes('tomorrow') && !lowerText.includes('today') &&
    !lowerText.includes('next') && !lowerText.includes('by ') &&
    !lowerText.includes('at ') && !lowerText.includes('in ')) {
    suggestions.push("💡 Add timing like 'tomorrow', 'next week', or 'by Friday'")
  }

  // Check for missing priority
  if (!lowerText.includes('urgent') && !lowerText.includes('high') &&
    !lowerText.includes('low') && !lowerText.includes('priority') &&
    !lowerText.includes('asap')) {
    suggestions.push("💡 Add priority like 'urgent', 'high priority', or 'low'")
  }

  // Check for missing time estimate
  if (!lowerText.includes('minute') && !lowerText.includes('hour') &&
    !lowerText.includes('min') && !lowerText.includes('hr') &&
    !/\d+\s*h/.test(lowerText) && !/\d+\s*m/.test(lowerText)) {
    suggestions.push("💡 Add time estimate like '30 minutes' or '2 hours'")
  }

  return suggestions
}

// Parse dates from natural language input with enhanced patterns
const parseDateFromText = (text: string): Date | null => {
  // First try chrono-node for standard date parsing
  const results = chrono.parse(text)
  if (results.length > 0) {
    return results[0].start.date()
  }

  // Enhanced fallback parsing for common patterns
  const lowerText = text.toLowerCase()
  const now = new Date()

  // Relative date patterns
  if (lowerText.includes('tomorrow')) {
    const tomorrow = new Date(now)
    tomorrow.setDate(now.getDate() + 1)
    return tomorrow
  }

  if (lowerText.includes('next week')) {
    const nextWeek = new Date(now)
    nextWeek.setDate(now.getDate() + 7)
    return nextWeek
  }

  if (lowerText.includes('next month')) {
    const nextMonth = new Date(now)
    nextMonth.setMonth(now.getMonth() + 1)
    return nextMonth
  }

  if (lowerText.includes('end of week') || lowerText.includes('this friday') || lowerText.includes('friday')) {
    const daysUntilFriday = (5 - now.getDay() + 7) % 7
    const friday = new Date(now)
    friday.setDate(now.getDate() + (daysUntilFriday === 0 ? 7 : daysUntilFriday))
    return friday
  }

  if (lowerText.includes('end of month')) {
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
    return endOfMonth
  }

  // Specific day patterns
  const dayPatterns = {
    'monday': 1, 'tuesday': 2, 'wednesday': 3, 'thursday': 4,
    'friday': 5, 'saturday': 6, 'sunday': 0
  }

  for (const [day, dayNum] of Object.entries(dayPatterns)) {
    if (lowerText.includes(`next ${day}`)) {
      const daysUntil = (dayNum - now.getDay() + 7) % 7
      const targetDate = new Date(now)
      targetDate.setDate(now.getDate() + (daysUntil === 0 ? 7 : daysUntil))
      return targetDate
    }
    if (lowerText.includes(`this ${day}`)) {
      const daysUntil = (dayNum - now.getDay() + 7) % 7
      const targetDate = new Date(now)
      targetDate.setDate(now.getDate() + daysUntil)
      return targetDate
    }
  }

  // Time parsing (e.g., "at 3pm", "at 14:30")
  const timeMatch = lowerText.match(/at (\d{1,2})(?::(\d{2}))?\s*(am|pm)?/)
  if (timeMatch) {
    let hour = parseInt(timeMatch[1])
    const minute = timeMatch[2] ? parseInt(timeMatch[2]) : 0
    const ampm = timeMatch[3]

    if (ampm === 'pm' && hour !== 12) hour += 12
    if (ampm === 'am' && hour === 12) hour = 0

    const timeDate = new Date(now)
    timeDate.setHours(hour, minute, 0, 0)
    return timeDate
  }

  return null
}

// Extract priority keywords from text with enhanced patterns
const parsePriorityFromText = (text: string): number | null => {
  const lowerText = text.toLowerCase()

  // Priority 5 (Urgent/Critical/Emergency)
  if (lowerText.includes('urgent') || lowerText.includes('asap') ||
    lowerText.includes('critical') || lowerText.includes('emergency') ||
    lowerText.includes('immediately') || lowerText.includes('right now') ||
    lowerText.includes('deadline') || lowerText.includes('rush')) {
    return 5
  }

  // Priority 4 (High/Important)
  if (lowerText.includes('high priority') || lowerText.includes('high') ||
    lowerText.includes('important') || lowerText.includes('priority') ||
    lowerText.includes('crucial') || lowerText.includes('essential') ||
    lowerText.includes('key') || lowerText.includes('vital')) {
    return 4
  }

  // Priority 3 (Medium/Normal/Default)
  if (lowerText.includes('medium priority') || lowerText.includes('medium') ||
    lowerText.includes('normal') || lowerText.includes('standard') ||
    lowerText.includes('regular') || lowerText.includes('average')) {
    return 3
  }

  // Priority 2 (Low)
  if (lowerText.includes('low priority') || lowerText.includes('low') ||
    lowerText.includes('whenever') || lowerText.includes('sometime') ||
    lowerText.includes('eventually') || lowerText.includes('casual')) {
    return 2
  }

  // Priority 1 (Minor/Trivial)
  if (lowerText.includes('minor') || lowerText.includes('trivial') ||
    lowerText.includes('nice to have') || lowerText.includes('optional') ||
    lowerText.includes('if time') || lowerText.includes('low effort')) {
    return 1
  }

  return null
}

// Extract time estimate from text
const parseTimeEstimateFromText = (text: string): number | null => {
  const lowerText = text.toLowerCase()

  // Match patterns like "30 minutes", "2 hours", "1h", "45m"
  const timePatterns = [
    /(\d+)\s*(?:hour|hr|h)(?:\w*)/,
    /(\d+)\s*(?:minute|min|m)(?:\w*)/,
    /(\d+):(\d+)/, // HH:MM format
    /(\d+)\s*(?:hr|hour|h)\s*(\d+)\s*(?:min|minute|m)?/ // "2h 30m" format
  ]

  for (const pattern of timePatterns) {
    const match = lowerText.match(pattern)
    if (match) {
      if (match[2]) { // HH:MM or H M format
        const hours = parseInt(match[1])
        const minutes = parseInt(match[2])
        return hours * 60 + minutes
      } else {
        const value = parseInt(match[1])
        if (lowerText.includes('hour') || lowerText.includes('hr') || lowerText.includes('h')) {
          return value * 60 // Convert hours to minutes
        } else if (lowerText.includes('minute') || lowerText.includes('min') || lowerText.includes('m')) {
          return value
        }
      }
    }
  }

  return null
}

// Parse subtasks from natural language input
const parseSubtasksFromText = (text: string): string[] => {
  const lowerText = text.toLowerCase()
  const subtasks: string[] = []

  // Look for patterns like:
  // - "with subtasks: task1, task2, task3"
  // - "including: task1 and task2"
  // - "steps: 1. task1, 2. task2"
  // - "breakdown: task1 | task2 | task3"

  const subtaskPatterns = [
    /with subtasks?:?\s*([^.]*?)(?:\sand\s|\s?,\s?|$)/gi,
    /including:?\s*([^.]*?)(?:\sand\s|\s?,\s?|$)/gi,
    /steps?:?\s*([^.]*?)(?:\sand\s|\s?,\s?|$)/gi,
    /breakdown:?\s*([^.]*?)(?:\sand\s|\s?,\s?|$)/gi
  ]

  for (const pattern of subtaskPatterns) {
    let match
    while ((match = pattern.exec(lowerText)) !== null) {
      const subtaskText = match[1].trim()
      if (subtaskText) {
        // Split by common separators
        const items = subtaskText.split(/\s*(?:,\s*|\sand\s*|\|)\s*/).filter(item => item.trim())
        subtasks.push(...items)
      }
    }
  }

  // Also look for numbered lists
  const numberedMatch = lowerText.match(/(\d+\.\s*[^.]+)+/g)
  if (numberedMatch) {
    for (const match of numberedMatch) {
      const items = match.split(/\d+\.\s*/).filter(item => item.trim())
      subtasks.push(...items)
    }
  }

  return [...new Set(subtasks)] // Remove duplicates
}

// Parse dependencies from natural language input
const parseDependenciesFromText = (text: string): string[] => {
  const lowerText = text.toLowerCase()
  const dependencies: string[] = []

  // Look for patterns like:
  // - "after finishing X"
  // - "depends on Y"
  // - "once Z is done"
  // - "following W"

  const dependencyPatterns = [
    /after\s+(?:finishing|completing|doing)\s+([^,.]+)/gi,
    /depends?\s+on\s+([^,.]+)/gi,
    /once\s+([^,.]+?)\s+is\s+(?:done|finished|completed)/gi,
    /following\s+([^,.]+)/gi,
    /after\s+([^,.]+?)\s+(?:is|are)\s+(?:done|finished|completed)/gi
  ]

  for (const pattern of dependencyPatterns) {
    let match
    while ((match = pattern.exec(lowerText)) !== null) {
      const dependency = match[1].trim()
      if (dependency) {
        dependencies.push(dependency)
      }
    }
  }

  return [...new Set(dependencies)] // Remove duplicates
}

// Parse recurring patterns from natural language input
const parseRecurringFromText = (text: string) => {
  const lowerText = text.toLowerCase()

  // Recurring patterns
  const recurringPatterns = {
    daily: /every\s+day|daily|each\s+day/gi,
    weekly: /every\s+week|weekly|each\s+week/gi,
    monthly: /every\s+month|monthly|each\s+month/gi,
    yearly: /every\s+year|yearly|annually|each\s+year/gi
  }

  // Day-specific recurring patterns
  const dayPatterns = {
    monday: /every\s+monday|each\s+monday/gi,
    tuesday: /every\s+tuesday|each\s+tuesday/gi,
    wednesday: /every\s+wednesday|each\s+wednesday/gi,
    thursday: /every\s+thursday|each\s+thursday/gi,
    friday: /every\s+friday|each\s+friday/gi,
    saturday: /every\s+saturday|each\s+saturday/gi,
    sunday: /every\s+sunday|each\s+sunday/gi
  }

  // Check for general recurring patterns
  for (const [frequency, pattern] of Object.entries(recurringPatterns)) {
    if (pattern.test(lowerText)) {
      return { frequency, interval: 1 }
    }
  }

  // Check for day-specific patterns
  for (const [day, pattern] of Object.entries(dayPatterns)) {
    if (pattern.test(lowerText)) {
      return {
        frequency: 'weekly',
        interval: 1,
        byDay: day.toUpperCase().substring(0, 2) // MO, TU, WE, etc.
      }
    }
  }

  // Check for "every X days/weeks/months"
  const everyPattern = /every\s+(\d+)\s+(day|week|month)/gi
  const everyMatch = everyPattern.exec(lowerText)
  if (everyMatch) {
    const interval = parseInt(everyMatch[1])
    const unit = everyMatch[2] + 's' // days, weeks, months
    return {
      frequency: unit,
      interval
    }
  }

  return null
}

// Parse goal and category information from natural language input
const submitForm = async () => {
  validationErrors.value = {}

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
    const parsedTimeEstimate = parseTimeEstimateFromText(dataToSend.natural_language)

    if (parsedDate) {
      dataToSend.due_date = parsedDate.toISOString()
    }
    if (parsedPriority) {
      dataToSend.priority = parsedPriority
    }
    if (parsedTimeEstimate) {
      dataToSend.time_estimate_minutes = parsedTimeEstimate
    }

    // Clear other fields when using natural language
    dataToSend.title = ''
    dataToSend.description = ''
  }

  // Validate the data
  let validationSchema = schemas.task.create
  if (dataToSend.use_natural_language && dataToSend.natural_language_input) {
    validationSchema = schemas.task.naturalLanguage
  }

  const validation = validateObject(dataToSend, validationSchema)

  if (!validation.isValid) {
    validationErrors.value = validation.errors
    return
  }

  try {
    await taskStore.submitForm(dataToSend)
    Object.assign(formData, {
      title: '',
      description: '',
      due_date: null,
      priority: 3,
      status: 'pending',
      natural_language: '',
      use_natural_language: true,
      time_estimate_minutes: null,
    })
    showForm.value = true
    validationErrors.value = {}
  } catch (err) {
    // Error is already handled in the store
  }
}

</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div v-if="showForm" @click="showForm = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
        <div
          class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
          <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
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
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark"
          @click.stop>

          <!-- Modal Header -->
          <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Create a Task</h2>
            <BaseButton @click="showForm = true" text="×" variant="default" size="sm" />
          </div>

          <!-- Modal Body -->
          <div class="p-6">
            <form @submit.prevent="submitForm" ref="taskForm" class="space-y-4">

              <!-- Natural Language Toggle -->
              <div class="flex items-center space-x-2">
                <input type="checkbox" v-model="formData.use_natural_language" id="natural-language-toggle"
                  name="use_natural_language"
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
                  placeholder="e.g., 'Buy groceries tomorrow at 5pm, high priority, 30 minutes' or 'Finish report by Friday, urgent' or 'Call mom next week, low priority'"
                  :class="{ 'border-red-500 focus:ring-red-500': validationErrors.natural_language }"></textarea>

                <!-- Help section with examples -->
                <details class="mt-2">
                  <summary
                    class="text-sm text-text-light/70 dark:text-text-dark/70 cursor-pointer hover:text-text-light dark:hover:text-text-dark">
                    💡 Examples & tips
                  </summary>
                  <div class="mt-2 p-3 bg-gray-50 dark:bg-gray-800 rounded-md text-sm space-y-2">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                      <div>
                        <p class="font-medium text-text-light dark:text-text-dark mb-1">📅 Dates & Times:</p>
                        <ul class="text-text-light/80 dark:text-text-dark/80 space-y-1 text-xs">
                          <li>• "tomorrow at 3pm"</li>
                          <li>• "next Friday"</li>
                          <li>• "end of month"</li>
                          <li>• "in 2 days"</li>
                        </ul>
                      </div>
                      <div>
                        <p class="font-medium text-text-light dark:text-text-dark mb-1">⚡ Priority:</p>
                        <ul class="text-text-light/80 dark:text-text-dark/80 space-y-1 text-xs">
                          <li>• "urgent", "asap"</li>
                          <li>• "high priority"</li>
                          <li>• "low", "whenever"</li>
                          <li>• "nice to have"</li>
                        </ul>
                      </div>
                      <div>
                        <p class="font-medium text-text-light dark:text-text-dark mb-1">⏱️ Time Estimates:</p>
                        <ul class="text-text-light/80 dark:text-text-dark/80 space-y-1 text-xs">
                          <li>• "30 minutes"</li>
                          <li>• "2 hours"</li>
                          <li>• "1h 30m"</li>
                          <li>• "45m"</li>
                        </ul>
                      </div>
                      <div>
                        <p class="font-medium text-text-light dark:text-text-dark mb-1">📝 Examples:</p>
                        <ul class="text-text-light/80 dark:text-text-dark/80 space-y-1 text-xs">
                          <li>• "Review code tomorrow urgent"</li>
                          <li>• "Buy milk next week low priority"</li>
                          <li>• "Write report by Friday 2 hours"</li>
                          <li>• "Plan trip with subtasks: book hotel, buy tickets"</li>
                          <li>• "Start project after finishing research"</li>
                          <li>• "Exercise every Monday"</li>
                          <li>• "Review expenses monthly"</li>
                          <li>• "Learn React for career goal"</li>
                          <li>• "Budget planning finance category"</li>
                        </ul>
                      </div>
                      <div>
                        <p class="font-medium text-text-light dark:text-text-dark mb-1">🔗 Advanced:</p>
                        <ul class="text-text-light/80 dark:text-text-dark/80 space-y-1 text-xs">
                          <li>• "with subtasks: step1, step2"</li>
                          <li>• "depends on finishing X"</li>
                          <li>• "including: task A and task B"</li>
                          <li>• "once Y is done"</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </details>
                <p v-if="validationErrors.natural_language" class="mt-1 text-sm text-red-500 dark:text-red-400">
                  {{ validationErrors.natural_language }}
                </p>

                <!-- Real-time parsing preview -->
                <div
                  v-if="parsedPreview && (parsedPreview.date || parsedPreview.priority || parsedPreview.timeEstimate || parsedPreview.subtasks.length > 0 || parsedPreview.dependencies.length > 0 || parsedPreview.recurring || parsedPreview.suggestions.length > 0)"
                  class="mt-3 space-y-3">
                  <!-- Parsed information -->
                  <div
                    v-if="parsedPreview.date || parsedPreview.priority || parsedPreview.timeEstimate || parsedPreview.subtasks.length > 0 || parsedPreview.dependencies.length > 0 || parsedPreview.recurring"
                    class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-md border border-blue-200 dark:border-blue-800">
                    <p class="text-sm font-medium text-blue-800 dark:text-blue-200 mb-2">Parsed from your input:</p>
                    <div class="space-y-1 text-sm text-blue-700 dark:text-blue-300">
                      <div v-if="parsedPreview.date" class="flex items-center gap-2">
                        <span class="font-medium">📅 Due:</span>
                        <span>{{ formatParsedDate(parsedPreview.date) }}</span>
                      </div>
                      <div v-if="parsedPreview.priority" class="flex items-center gap-2">
                        <span class="font-medium">⚡ Priority:</span>
                        <span>{{ getPriorityLabel(parsedPreview.priority) }}</span>
                      </div>
                      <div v-if="parsedPreview.timeEstimate" class="flex items-center gap-2">
                        <span class="font-medium">⏱️ Estimate:</span>
                        <span>{{ parsedPreview.timeEstimate }} minutes</span>
                      </div>
                      <div v-if="parsedPreview.recurring" class="flex items-center gap-2">
                        <span class="font-medium">🔄 Recurring:</span>
                        <span>{{ formatRecurringInfo(parsedPreview.recurring) }}</span>
                      </div>
                      <div v-if="parsedPreview.subtasks.length > 0" class="flex items-start gap-2">
                        <span class="font-medium">📋 Subtasks:</span>
                        <div class="flex flex-wrap gap-1">
                          <span v-for="subtask in parsedPreview.subtasks" :key="subtask"
                            class="px-2 py-1 bg-blue-100 dark:bg-blue-800 text-blue-800 dark:text-blue-200 rounded text-xs">
                            {{ subtask }}
                          </span>
                        </div>
                      </div>
                      <div v-if="parsedPreview.dependencies.length > 0" class="flex items-start gap-2">
                        <span class="font-medium">🔗 Depends on:</span>
                        <div class="flex flex-wrap gap-1">
                          <span v-for="dependency in parsedPreview.dependencies" :key="dependency"
                            class="px-2 py-1 bg-orange-100 dark:bg-orange-800 text-orange-800 dark:text-orange-200 rounded text-xs">
                            {{ dependency }}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Suggestions -->
                  <div v-if="parsedPreview.suggestions.length > 0"
                    class="p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-md border border-yellow-200 dark:border-yellow-800">
                    <p class="text-sm font-medium text-yellow-800 dark:text-yellow-200 mb-2">💡 Suggestions to improve
                      parsing:</p>
                    <ul class="space-y-1 text-sm text-yellow-700 dark:text-yellow-300">
                      <li v-for="suggestion in parsedPreview.suggestions" :key="suggestion">
                        {{ suggestion }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>

              <!-- Structured Input Fields -->
              <div v-else>
                <div class="flex flex-col">
                  <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Title</label>
                  <input type="text" v-model="formData.title" name="title"
                    class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                    placeholder="Task title" required
                    :class="{ 'border-red-500 focus:ring-red-500': validationErrors.title }" />
                  <p v-if="validationErrors.title" class="mt-1 text-sm text-red-500 dark:text-red-400">
                    {{ validationErrors.title }}
                  </p>
                </div>

                <div class="flex flex-col">
                  <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Description</label>
                  <textarea v-model="formData.description" name="description" rows="3"
                    class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                    placeholder="Optional description"
                    :class="{ 'border-red-500 focus:ring-red-500': validationErrors.description }"></textarea>
                  <p v-if="validationErrors.description" class="mt-1 text-sm text-red-500 dark:text-red-400">
                    {{ validationErrors.description }}
                  </p>
                </div>

                <div class="flex flex-col">
                  <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Due Date</label>
                  <input type="datetime-local" v-model="formData.due_date" name="due_date"
                    class="px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                    :class="{ 'border-red-500 focus:ring-red-500': validationErrors.due_date }" />
                  <p v-if="validationErrors.due_date" class="mt-1 text-sm text-red-500 dark:text-red-400">
                    {{ validationErrors.due_date }}
                  </p>
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
              <div
                class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                <BaseButton type="button" text="Cancel" @click="showForm = true" variant="default" size="md"
                  class="w-full sm:w-auto" />
                <BaseButton type="submit" text="Create Task" variant="primary" size="md" class="w-full sm:w-auto" />
              </div>

            </form>
          </div>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>
