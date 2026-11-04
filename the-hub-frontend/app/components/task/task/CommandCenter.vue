<script setup lang="ts">
import type { Task } from '~/types/task'

interface Filters {
  search: string
  status: string
  priority: string
  goal_id: string
  due_before: string
  due_after: string
}

interface Props {
  tasks: Task[]
  filters: Filters
  activeFilterCount: number
  hasActiveFilters: boolean
  viewMode: 'list' | 'planning'
  backlogView: boolean
  aiLoading: boolean
  statusItems: string[]
  priorityItems: string[]
  goalItems: Array<{ label: string, value: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:filters', patch: Partial<Filters>): void
  (e: 'clear-filters'): void
  (e: 'toggle-view'): void
  (e: 'toggle-backlog'): void
  (e: 'run-ai-check'): void
  (e: 'open-themes'): void
  (e: 'open-focus', taskId: string): void
}>()

const taskStore = useTaskStore()
const focusStore = useFocusSessionStore()

// Helpers: strip the time component of a date
const startOfDay = (date: Date) => {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  return d
}

// Collapsible body, persisted between visits
const collapsed = useLocalStorage('plan-command-center-collapsed', false)
const showFilters = ref(false)

const toggleFilters = () => {
  if (collapsed.value) collapsed.value = false
  showFilters.value = !showFilters.value
}

// --- Overview stats ---
const todayStart = computed(() => startOfDay(new Date()).getTime())

const dueTodayCount = computed(() => {
  return props.tasks.filter((task) => {
    if (task.status !== 'pending' || !task.due_date) return false
    return startOfDay(new Date(task.due_date)).getTime() === todayStart.value
  }).length
})

const overdueCount = computed(() => {
  return props.tasks.filter((task) => {
    if (task.status !== 'pending' || !task.due_date) return false
    return startOfDay(new Date(task.due_date)).getTime() < todayStart.value
  }).length
})

const completedCount = computed(() => {
  return props.tasks.filter(task => task.status === 'completed').length
})

const backlogCount = computed(() => {
  return props.tasks.filter(task => task.status === 'backlog').length
})

const focusMinutesToday = computed(() => focusStore.getTodayMinutes)

const formatFocusMinutes = (minutes: number) => {
  if (minutes === 0) return '0m'
  if (minutes < 60) return `${minutes}m`
  return `${Math.floor(minutes / 60)}h ${minutes % 60}m`
}

// --- Next task ---
const nextTask = computed<Task | null>(() => {
  const pending = props.tasks.filter((task) => task.status === 'pending')
  if (pending.length === 0) return null

  const due = pending.filter((task) =>
    task.due_date && startOfDay(new Date(task.due_date)).getTime() <= todayStart.value
  )
  const pool = due.length > 0 ? due : pending

  return [...pool].sort((a, b) => {
    const priorityDiff = (b.priority ?? 0) - (a.priority ?? 0)
    if (priorityDiff !== 0) return priorityDiff
    const aDue = a.due_date ? new Date(a.due_date).getTime() : Infinity
    const bDue = b.due_date ? new Date(b.due_date).getTime() : Infinity
    if (aDue !== bDue) return aDue - bDue
    return (a.order ?? 0) - (b.order ?? 0)
  })[0]
})

const nextTaskDueLabel = computed(() => {
  const task = nextTask.value
  if (!task?.due_date) return null
  const due = startOfDay(new Date(task.due_date)).getTime()
  if (due === todayStart.value) return 'Due today'
  if (due < todayStart.value) return 'Overdue'
  return null
})

const completeNextTask = async () => {
  const task = nextTask.value
  if (!task) return
  try {
    await taskStore.completeTask({ ...task, status: 'completed' })
  } catch (err) {
    console.error('Error completing task:', err)
  }
}

const snoozeNextTask = async () => {
  const task = nextTask.value
  if (!task) return
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  tomorrow.setHours(23, 59, 0, 0)
  try {
    await taskStore.editTask({
      task_id: task.task_id,
      due_date: tomorrow.toISOString()
    })
  } catch (err) {
    console.error('Error snoozing task:', err)
  }
}

// --- Inline quick-add ---
const { parseDateFromText, parsePriorityFromText, parseTimeEstimateFromText } = useTaskParsing()

const quickAddText = ref('')
const quickAddInput = ref<HTMLInputElement | null>(null)
const quickAddSubmitting = ref(false)

const quickAddParsed = computed(() => {
  const text = quickAddText.value.trim()
  if (!text) return null
  return {
    date: parseDateFromText(text),
    priority: parsePriorityFromText(text),
    timeEstimate: parseTimeEstimateFromText(text)
  }
})

const quickAddHasParsed = computed(() => {
  const parsed = quickAddParsed.value
  return !!parsed && (parsed.date !== null || parsed.priority !== null || parsed.timeEstimate !== null)
})

const getPriorityLabel = (priority: number) => {
  const labels = { 1: 'Low', 2: 'Low', 3: 'Medium', 4: 'High', 5: 'Urgent' }
  return labels[priority as keyof typeof labels] || 'Medium'
}

const formatParsedDate = (date: Date) => {
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit'
  })
}

const submitQuickAdd = async () => {
  const text = quickAddText.value.trim()
  if (!text || quickAddSubmitting.value) return

  quickAddSubmitting.value = true
  try {
    const parsed = quickAddParsed.value
    await taskStore.submitForm({
      title: '',
      description: '',
      due_date: parsed?.date ? parsed.date.toISOString() : undefined,
      priority: parsed?.priority ?? undefined,
      status: 'pending',
      natural_language_input: text,
      use_natural_language: true,
      time_estimate_minutes: parsed?.timeEstimate ?? undefined
    })
    quickAddText.value = ''
    quickAddInput.value?.focus()
  } catch (err) {
    // Error is already handled in the store
    console.error('Quick add failed:', err)
  } finally {
    quickAddSubmitting.value = false
  }
}

// --- Filter updates (proxied through emits) ---
const setFilter = (patch: Partial<Filters>) => emit('update:filters', patch)
const clearFilters = () => emit('clear-filters')
</script>

<template>
  <div class="p-4 shadow rounded-2xl dark:inset-shadow-sm inset-shadow-gray-500/50">
    <!-- Header row: quick-add + stats + controls -->
    <div class="flex flex-col gap-3">
      <div class="flex flex-wrap items-center gap-2">
        <h2 class="text-sm font-semibold text-text-light dark:text-text-dark mr-1">Command Center</h2>

        <div class="flex-1 min-w-60">
          <div class="relative">
            <input ref="quickAddInput" v-model="quickAddText" @keydown.enter="submitQuickAdd"
              class="w-full px-3 py-2 pr-10 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 text-sm"
              placeholder="Quick add — e.g. 'Pay rent tomorrow 5pm urgent 30m'" />
            <button @click="submitQuickAdd" :disabled="quickAddSubmitting || !quickAddText.trim()"
              class="absolute right-1.5 top-1/2 -translate-y-1/2 px-2 py-1 rounded bg-primary text-white text-xs font-medium hover:bg-primary/90 disabled:opacity-40 flex items-center gap-1">
              <UIcon v-if="quickAddSubmitting" name="i-lucide-loader-circle" class="w-3 h-3 animate-spin" />
              <UIcon v-else name="i-lucide-plus" class="w-3 h-3" />
            </button>
          </div>

          <!-- Parsed preview chips -->
          <div v-if="quickAddHasParsed" class="flex flex-wrap gap-1.5 mt-1.5">
            <span v-if="quickAddParsed?.date"
              class="px-2 py-0.5 rounded-full bg-primary/10 text-primary text-xs">
              {{ formatParsedDate(quickAddParsed.date) }}
            </span>
            <span v-if="quickAddParsed?.priority"
              class="px-2 py-0.5 rounded-full bg-secondary/10 dark:bg-secondary/20 text-secondary text-xs">
              P{{ quickAddParsed.priority }} · {{ getPriorityLabel(quickAddParsed.priority) }}
            </span>
            <span v-if="quickAddParsed?.timeEstimate"
              class="px-2 py-0.5 rounded-full bg-success/10 text-success text-xs">
              {{ quickAddParsed.timeEstimate }}m
            </span>
          </div>
        </div>
      </div>

      <!-- Stats strip -->
      <div class="flex flex-wrap items-center gap-2">
        <span class="px-2.5 py-1 rounded-full bg-primary/10 text-primary text-xs font-medium">
          Due today: {{ dueTodayCount }}
        </span>
        <span class="px-2.5 py-1 rounded-full text-xs font-medium"
          :class="overdueCount > 0 ? 'bg-error/10 text-error' : 'bg-success/10 text-success'">
          Overdue: {{ overdueCount }}
        </span>
        <span class="px-2.5 py-1 rounded-full bg-success/10 text-success text-xs font-medium">
          Completed: {{ completedCount }}
        </span>
        <span class="px-2.5 py-1 rounded-full bg-warning/10 text-warning text-xs font-medium">
          Focus today: {{ formatFocusMinutes(focusMinutesToday) }}
        </span>
        <span class="px-2.5 py-1 rounded-full bg-secondary/10 dark:bg-secondary/20 text-secondary text-xs font-medium">
          Backlog: {{ backlogCount }}
        </span>

        <div class="flex-1" />

        <!-- Control buttons -->
        <FormTask />
        <UButton label="AI Check" :loading="aiLoading" icon="i-lucide-bot" size="xs" variant="outline"
          @click="emit('run-ai-check')" />
        <ThemeDaysToggle />
        <UButton label="Themes" icon="i-lucide-palette" size="xs" variant="outline" @click="emit('open-themes')" />
        <UButton :label="viewMode === 'planning' ? 'List' : 'Planning'"
          :icon="viewMode === 'planning' ? 'i-lucide-list' : 'i-lucide-columns-3'" size="xs" variant="outline"
          @click="emit('toggle-view')" />
        <UButton :label="backlogView ? 'Active' : 'Backlog'"
          :icon="backlogView ? 'i-lucide-inbox' : 'i-lucide-archive'" size="xs" variant="outline"
          :active="backlogView" @click="emit('toggle-backlog')" />
        <UButton :label="activeFilterCount > 0 ? `Filters (${activeFilterCount})` : 'Filters'"
          icon="i-lucide-sliders-horizontal" size="xs" variant="outline" :active="hasActiveFilters"
          @click="toggleFilters" />
        <UButton :icon="collapsed ? 'i-lucide-chevron-down' : 'i-lucide-chevron-up'" size="xs" variant="ghost"
          color="neutral" :aria-label="collapsed ? 'Expand command center' : 'Collapse command center'"
          @click="collapsed = !collapsed" />
      </div>
    </div>

    <!-- Collapsible body -->
    <div v-if="!collapsed" class="mt-4 space-y-4">
      <!-- Next task (hidden while browsing the backlog) -->
      <div v-if="nextTask && !backlogView"
        class="p-3 rounded-xl border-l-4 border-warning bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <p class="text-xs font-medium text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">Next up</p>
            <p class="font-medium text-text-light dark:text-text-dark mt-0.5 truncate">{{ nextTask.title }}</p>
            <div class="flex items-center gap-2 mt-1 text-xs">
              <span v-if="nextTaskDueLabel"
                class="px-2 py-0.5 rounded-full font-medium"
                :class="nextTaskDueLabel === 'Overdue' ? 'bg-error/10 text-error' : 'bg-primary/10 text-primary'">
                {{ nextTaskDueLabel }}
              </span>
              <span v-if="nextTask.priority"
                class="px-2 py-0.5 rounded-full bg-secondary/10 dark:bg-secondary/20 text-secondary">
                P{{ nextTask.priority }}
              </span>
              <span v-if="nextTask.time_estimate_minutes"
                class="text-text-light/60 dark:text-text-dark/60">
                ~{{ nextTask.time_estimate_minutes }}m
              </span>
            </div>
          </div>
          <div class="flex items-center gap-1.5 shrink-0">
            <UButton icon="i-lucide-check" size="xs" color="success" variant="soft"
              label="Done" @click="completeNextTask" />
            <UButton icon="i-lucide-timer" size="xs" color="primary" variant="soft"
              label="Focus" @click="emit('open-focus', nextTask.task_id)" />
            <UButton icon="i-lucide-clock-arrow-up" size="xs" color="neutral" variant="soft"
              label="Snooze" @click="snoozeNextTask" />
          </div>
        </div>
      </div>
      <p v-else-if="!backlogView" class="text-sm text-text-light/50 dark:text-text-dark/50">
        No pending tasks — add one above to get started.
      </p>

      <!-- Filter Panel -->
      <div v-if="showFilters"
        class="p-4 bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 rounded-lg">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <UFormField label="Search">
            <UInput :model-value="filters.search" placeholder="Search tasks..." class="w-full"
              @update:model-value="(v: string) => setFilter({ search: v })" />
          </UFormField>

          <UFormField label="Status">
            <USelect :model-value="filters.status" :items="statusItems" placeholder="All Statuses" class="w-full"
              @update:model-value="(v: string) => setFilter({ status: v })" />
          </UFormField>

          <UFormField label="Priority">
            <USelect :model-value="filters.priority" :items="priorityItems" placeholder="All Priorities"
              class="w-full" @update:model-value="(v: string) => setFilter({ priority: v })" />
          </UFormField>

          <UFormField label="Goal">
            <USelect :model-value="filters.goal_id" :items="goalItems" placeholder="No Goal" class="w-full"
              @update:model-value="(v: string) => setFilter({ goal_id: v })" />
          </UFormField>

          <div class="flex flex-col">
            <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due After</label>
            <input :value="filters.due_after" type="date"
              @input="setFilter({ due_after: ($event.target as HTMLInputElement).value })"
              class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
          </div>

          <div class="flex flex-col">
            <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due
              Before</label>
            <input :value="filters.due_before" type="date"
              @input="setFilter({ due_before: ($event.target as HTMLInputElement).value })"
              class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
          </div>
        </div>

        <div class="flex justify-end mt-4">
          <UButton label="Clear Filters" variant="outline" color="neutral" :disabled="!hasActiveFilters"
            @click="clearFilters" />
        </div>
      </div>
    </div>
  </div>
</template>
