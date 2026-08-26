<script setup lang="ts">
import type { Task } from '~/types/task'
import { LazyTaskEditModal } from '#components'

interface Props {
  taskList: Task[]
}

const props = defineProps<Props>()

const goalStore = useGoalStore()
const themeStore = useThemeStore()

onMounted(() => {
  if (goalStore.goals.length === 0) goalStore.fetchGoals()
  if (themeStore.themes.length === 0) themeStore.fetchThemes()
  themeStore.fetchEnabled()
})

const filters = reactive({
  search: '',
  status: '',
  priority: '',
  goal_id: '',
  due_before: '',
  due_after: '',
})

const statusItems = ref(["pending", "completed"])
const priorityItems = ref(["1", "2", "3", "4", "5"])

const goalItems = computed(() => [
  { label: 'All Goals', value: 'all' },
  ...goalStore.goals.map(goal => ({ label: goal.title, value: goal.goal_id })),
])

const hasActiveFilters = computed(() => {
  return Object.values(filters).some(value => value !== '' && value !== null)
})

const activeFilterCount = computed(() => {
  return Object.values(filters).filter(value => value !== '' && value !== null).length
})

const clearFilters = () => {
  Object.assign(filters, {
    search: '',
    status: '',
    priority: '',
    goal_id: '',
    due_before: '',
    due_after: '',
  })
}

const filteredTasks = computed(() => {
  return props.taskList.filter((task: Task) => {
    if (filters.search) {
      const query = filters.search.toLowerCase()
      const titleMatch = task.title.toLowerCase().includes(query)
      const descMatch = task.description?.toLowerCase().includes(query)
      if (!titleMatch && !descMatch) return false
    }

    if (filters.status && task.status !== filters.status) return false

    if (filters.priority && task.priority !== Number(filters.priority)) return false

    if (filters.goal_id === 'all') {
      // Include all tasks (goal-related and standalone)
    } else if (filters.goal_id) {
      if (task.goal_id !== filters.goal_id) return false
    } else if (task.goal_id) {
      // Default filter: only tasks not associated with a goal
      return false
    }

    if (filters.due_before || filters.due_after) {
      if (!task.due_date) return false
      const due = new Date(task.due_date)
      if (filters.due_before && due > new Date(filters.due_before)) return false
      if (filters.due_after && due < new Date(filters.due_after)) return false
    }

    return true
  })
})

// When Theme Days mode is enabled, only tasks for today's theme are shown
const visibleTasks = computed(() => {
  return themeStore.filterTasksByTheme(filteredTasks.value)
})

const tasks = computed(() => {
  return filteredTasks.value.map((task: Task) => {
    const startTime = task.start_time ? new Date(task.start_time) : null
    let endTime = null

    if (startTime) {
      const durationMinutes = task.time_estimate_minutes || 60
      endTime = new Date(startTime.getTime() + durationMinutes * 60000)
    }

    return {
      task_id: task.task_id,
      title: task.title,
      start_time: startTime,
      end_time: endTime
    }
  })
})


const today = computed(() => new Date())
const tomorrow = computed(() => {
  const date = new Date()
  date.setDate(date.getDate() + 1)
  return date
})

// Helper function to compare dates (ignoring time)
const isSameDay = (date1: Date | null, date2: Date | null): boolean => {
  if (!date1 || !date2) return false
  return date1.getFullYear() === date2.getFullYear() &&
    date1.getMonth() === date2.getMonth() &&
    date1.getDate() === date2.getDate()
}

// Filter tasks by date (only include scheduled tasks for that specific day)
const todayTasks = computed(() => {
  return tasks.value.filter(task =>
    task.start_time && isSameDay(task.start_time, today.value)
  )
})

const tomorrowTasks = computed(() => {
  return tasks.value.filter(task =>
    task.start_time && isSameDay(task.start_time, tomorrow.value)
  )
})

const viewMode = ref<'list' | 'planning'>('list')
const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'list' ? 'planning' : 'list'
}
const showFilters = ref(false)
const showThemeManager = ref(false)

let selectedTab = ref('Tasks')
const selectingTab = (val: string) => selectedTab.value = val
const tabs = ['Tasks', 'Goals']
const tri_modal = ['Planning', 'Execute', 'Analysis']

// Modal state for editing
const overlay = useOverlay()
const modal = overlay.create(LazyTaskEditModal)

const openEdit = (taskId: string) => {

  const editTask = props.taskList.find((t: Task) => t.task_id === taskId) || null

  modal.open({
    task: editTask
  })
}

const taskStore = useTaskStore()

const handleMoveUp = (taskId: string) => {
  const taskIndex = props.taskList.findIndex((t: Task) => t.task_id === taskId)
  if (taskIndex <= 0) return

  const newTasks = [...props.taskList]
    ;[newTasks[taskIndex], newTasks[taskIndex - 1]] = [newTasks[taskIndex - 1], newTasks[taskIndex]]

  const payload = newTasks.map((t: Task, idx: number) => ({
    task_id: t.task_id,
    order: idx
  }))

  taskStore.reorderTask(payload)
}

const handleMoveDown = (taskId: string) => {
  const taskIndex = props.taskList.findIndex((t: Task) => t.task_id === taskId)
  if (taskIndex === -1 || taskIndex >= props.taskList.length - 1) return

  const newTasks = [...props.taskList]
    ;[newTasks[taskIndex], newTasks[taskIndex + 1]] = [newTasks[taskIndex + 1], newTasks[taskIndex]]

  const payload = newTasks.map((t: Task, idx: number) => ({
    task_id: t.task_id,
    order: idx
  }))

  taskStore.reorderTask(payload)
}

const toast = useToast()
const showAIPreview = ref(false)
const aiPreviewData = ref<any[]>([])
const aiApplying = ref(false)
const taskAiLoading = ref(false)

const handleAICheck = async () => {
  taskAiLoading.value = true
  try {
    const response = await taskStore.getAITaskPreview()
    if (response) {
      if (response.message === 'All tasks already optimized') {
        toast.add({ title: "Task", description: "All tasks already optimized", color: "info" })
        return
      }
      aiPreviewData.value = response.preview
      showAIPreview.value = true
    }
  } catch (error) {
    console.error('AI check failed:', error)
    toast.add({ title: "Error", description: "Failed to get AI recommendations", color: "error" })
  } finally {
    taskAiLoading.value = false
  }
}

const closeAIPreview = () => {
  showAIPreview.value = false
  aiPreviewData.value = []
}

const selectedTasks = ref<Set<string>>(new Set())

const toggleTaskSelection = (taskId: string) => {
  if (selectedTasks.value.has(taskId)) {
    selectedTasks.value.delete(taskId)
  } else {
    selectedTasks.value.add(taskId)
  }
}

const selectAllTasks = () => {
  aiPreviewData.value.forEach((task) => {
    selectedTasks.value.add(task.task_id)
  })
}

const deselectAllTasks = () => {
  selectedTasks.value.clear()
}

const applyAISelected = async () => {
  aiApplying.value = true
  const appliedTasks = aiPreviewData.value.map((task) => ({
    task_id: task.task_id,
    selected: selectedTasks.value.has(task.task_id)
  }))

  const success = await taskStore.applyAITasks(appliedTasks)
  if (success) {
    closeAIPreview()
  }
  aiApplying.value = false
}
</script>
<template>
  <div id="plan" class="flex flex-col flex-1 min-h-screen">
    <!--Tabs -->
    <div class="layout-tabs">
      <Tabs :modelValue="selectedTab" @update:modelValue="selectingTab" :tabs="tabs">
        <!-- Control -->
        <template #Tasks>
          <!-- Task control panel -->
          <div class="p-4 shadow rounded-2xl dark:inset-shadow-sm inset-shadow-gray-500/50">
            <h2>Control Panel</h2>
            <div class="flex flex-wrap space-x-3 space-y-3">
              <FormTask />
              <UButton label="AI Check" :loading="taskAiLoading" icon="i-lucide-bot" @click="handleAICheck"
                variant="outline" />
              <ThemeDaysToggle />
              <UButton label="Themes" icon="i-lucide-palette" variant="outline" @click="showThemeManager = true" />
              <UButton :label="viewMode === 'planning' ? 'List' : 'Planning'"
                :icon="viewMode === 'planning' ? 'i-lucide-list' : 'i-lucide-columns-3'" @click="toggleViewMode"
                variant="outline" />
              <UButton :label="activeFilterCount > 0 ? `Filters (${activeFilterCount})` : 'Filters'"
                icon="i-lucide-sliders-horizontal" variant="outline" @click="showFilters = !showFilters"
                :active="hasActiveFilters" />
            </div>

            <!-- Filter Panel -->
            <div v-if="showFilters"
              class="mt-4 p-4 bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 rounded-lg">
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                <UFormField label="Search">
                  <UInput v-model="filters.search" placeholder="Search tasks..." class="w-full" />
                </UFormField>

                <UFormField label="Status">
                  <USelect v-model="filters.status" :items="statusItems" placeholder="All Statuses" class="w-full" />
                </UFormField>

                <UFormField label="Priority">
                  <USelect v-model="filters.priority" :items="priorityItems" placeholder="All Statuses"
                    class="w-full" />
                </UFormField>

                <UFormField label="Goal">
                  <USelect v-model="filters.goal_id" :items="goalItems" placeholder="No Goal" class="w-full" />
                </UFormField>

                <div class="flex flex-col">
                  <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due After</label>
                  <input v-model="filters.due_after" type="date"
                    class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
                </div>

                <div class="flex flex-col">
                  <label class="mb-1 text-sm font-medium text-text-light dark:text-text-dark">Due
                    Before</label>
                  <input v-model="filters.due_before" type="date"
                    class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
                </div>
              </div>

              <div class="flex justify-end mt-4">
                <UButton label="Clear Filters" variant="outline" color="neutral" :disabled="!hasActiveFilters"
                  @click="clearFilters" />
              </div>
            </div>
          </div>

          <div class="layout-content p-4 flex flex-1 flex-col md:flex-row">
            <div class="layout-tasks basis-1/3 grow">
              <div v-if="themeStore.enabled && themeStore.activeThemeToday"
                class="mb-3 px-3 py-2 rounded-lg border border-primary/30 bg-primary/5 flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full"
                  :style="{ backgroundColor: themeStore.activeThemeToday.color }" />
                <span class="text-sm text-text-light dark:text-text-dark">
                  Today is <span class="font-semibold">{{ themeStore.activeThemeToday.name }}</span> — showing only
                  themed tasks
                </span>
              </div>
              <div v-else-if="themeStore.enabled"
                class="mb-3 px-3 py-2 rounded-lg border border-surface-light/20 dark:border-surface-dark/20 text-sm text-text-light/60 dark:text-text-dark/60">
                No theme is assigned to today, so no tasks are shown.
              </div>
              <slot name="tasks">
                <TaskList :tasks="visibleTasks" @edit="openEdit" @moveUp="handleMoveUp" @moveDown="handleMoveDown" />
              </slot>
            </div>

            <div v-if="viewMode === 'planning'" class="layout-calendar-slot flex basis-2/3 ml-2 grow">
              <slot name="calendar-slot" class="flex w-full">
                <DateSlots class="grow basis-1/2 max-h-170 overflow-y-auto custom-scrollbar" label="Today" :date="today"
                  :tasks="todayTasks" />
                <DateSlots class="grow basis-1/2 max-h-170 overflow-y-auto custom-scrollbar" label="Tomorrow"
                  :date="tomorrow" :tasks="tomorrowTasks" />
              </slot>
            </div>
          </div>
        </template>
        <template #Goals>
          <Goals />
        </template>
      </Tabs>
    </div>

    <!-- AI Task Preview Modal -->
    <div v-if="showAIPreview" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="closeAIPreview">
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-xl max-w-2xl w-full max-h-[80vh] overflow-hidden flex flex-col">
        <div class="p-4 border-b border-surface-light dark:border-surface-dark flex items-center justify-between">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">AI Task Optimizations</h3>
          <button @click="closeAIPreview" class="text-text-light dark:text-text-dark hover:text-error">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-4 overflow-y-auto flex-1">
          <p class="text-sm text-text-light dark:text-text-dark/70 mb-4">
            Review and select the tasks you want to apply AI suggestions to. Tasks with subtasks will create child
            tasks.
          </p>

          <div class="space-y-3">
            <div v-for="task in aiPreviewData" :key="task.task_id"
              class="border border-surface-light dark:border-surface-dark rounded-lg p-3"
              :class="selectedTasks.has(task.task_id) ? 'border-primary bg-primary/5' : ''">
              <div class="flex items-start gap-3">
                <input type="checkbox" :checked="selectedTasks.has(task.task_id)"
                  @change="toggleTaskSelection(task.task_id)"
                  class="h-5 w-5 mt-1 rounded border-surface-light dark:border-surface-dark text-primary focus:ring-primary" />

                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    <span class="text-sm text-text-light/60 dark:text-text-dark/60 line-through">{{
                      task.original_title }}</span>
                  </div>
                  <div class="font-medium text-text-light dark:text-text-dark">{{ task.title }}</div>
                  <div class="text-sm text-text-light dark:text-text-dark/70 mt-1">{{ task.description }}</div>
                  <div class="flex items-center gap-3 mt-2 text-xs">
                    <span class="px-2 py-0.5 bg-secondary/10 dark:bg-secondary/20 text-secondary rounded">P{{
                      task.priority
                      }}</span>
                    <span class="text-text-light dark:text-text-dark/60">{{ task.estimated_hours }}h</span>
                  </div>

                  <!-- Subtasks -->
                  <div v-if="task.subtasks && task.subtasks.length > 0" class="mt-3 pl-3 border-l-2 border-primary/30">
                    <div class="text-xs font-medium text-text-light dark:text-text-dark/70 mb-1">Subtasks to create:
                    </div>
                    <div v-for="(subtask, idx) in task.subtasks" :key="idx"
                      class="text-sm text-text-light dark:text-text-dark py-1">
                      <span class="text-primary">+</span> {{ subtask.title }} ({{ subtask.estimated_hours }}h)
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="p-4 border-t border-surface-light dark:border-surface-dark flex items-center justify-between gap-2">
          <div class="flex gap-2">
            <button @click="selectAllTasks" class="text-sm text-primary hover:underline">Select All</button>
            <button @click="deselectAllTasks"
              class="text-sm text-text-light dark:text-text-dark hover:underline">Deselect
              All</button>
          </div>
          <div class="flex gap-2">
            <button @click="closeAIPreview"
              class="px-4 py-2 text-sm border border-surface-light dark:border-surface-dark rounded text-text-light dark:text-text-dark hover:bg-surface-light/50 dark:hover:bg-surface-dark/50">
              Cancel
            </button>
            <button @click="applyAISelected"
              class="px-4 py-2 text-sm bg-primary text-white rounded hover:bg-primary/90 disabled:opacity-50 flex items-center gap-2"
              :disabled="selectedTasks.size === 0 || taskStore.aiTaskLoading || aiApplying">
              <svg v-if="aiApplying" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              <span>{{ aiApplying ? 'Applying...' : `Apply Selected (${selectedTasks.size})` }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Theme Days Manager Modal -->
    <ThemeManagerModal :show="showThemeManager" @close="showThemeManager = false" />
  </div>
</template>
