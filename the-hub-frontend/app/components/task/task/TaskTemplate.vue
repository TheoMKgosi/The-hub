<script setup lang="ts">
import type { Task } from '~/types/task'
interface Props {
  taskList: Task[]
}

const props = defineProps<Props>()

const tasks = computed(() => {
  return props.taskList.map((task: Task) => {
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

// Add user settings integration
const userSettings = ref<any>({})

// Load user settings on mount
onMounted(async () => {
  try {
    const auth = useAuthStore()
    if (!auth.user?.user_id) return

    const { $api } = useNuxtApp()
    const response = await $api(`/users/${auth.user.user_id}/settings`)
    userSettings.value = response.settings || {}
  } catch (error) {
    console.warn('Failed to load user settings for tri-modal:', error)
  }
})

const tri_interface = computed(() => {
  return userSettings.value.task?.['tri_modal'] === true
})

let selectedTab = ref('Tasks')
const selectingTab = (val: string) => selectedTab.value = val
const tabs = ['Tasks', 'Goals', 'Analytics']
const tri_modal = ['Planning', 'Execute', 'Analysis']

// Modal state for editing
const showEditModal = ref(false)
const taskToEdit = ref<Task | null>(null)

const startEdit = (taskId: string) => {
  const task = props.taskList.find((t: Task) => t.task_id === taskId)
  if (task) {
    taskToEdit.value = task
    showEditModal.value = true
  }
}

const closeEditModal = () => {
  showEditModal.value = false
  taskToEdit.value = null
}

const handleTaskSave = () => {
  showEditModal.value = false
  taskToEdit.value = null
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

const { addToast } = useToast()
const showAIPreview = ref(false)
const aiPreviewData = ref<any[]>([])
const aiApplying = ref(false)

const handleAICheck = async () => {
  try {
    const response = await taskStore.getAITaskPreview()
    if (response) {
      if (response.message === 'All tasks already optimized') {
        addToast('All tasks already optimized', 'info')
        return
      }
      aiPreviewData.value = response.preview
      showAIPreview.value = true
    }
  } catch (error) {
    console.error('AI check failed:', error)
    addToast('Failed to get AI recommendations', 'error')
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
          <!--
          <h2>Filter</h2>
          <div class="layout-controls w-full">
            <div class="bg-background-light dark:bg-background-dark">
              <slot name="control">
                <div>
                  <div class="flex">
                    <BaseButton text="All" variant="primary" class="mr-2" />
                    <BaseButton text="Linked" variant="primary" class="mr-2" />
                    <BaseButton text="Pending" variant="primary" class="mr-2" />
                  </div>
                </div>
                <div v-if="tri_interface">
                  <p>Tri-Modal</p>
                  <div class="flex">
                    <SegmentedControl :texts="tri_modal" />
                  </div>
                </div>
              </slot>
            </div>
          </div>
        -->

          <div class="layout-content p-4 flex flex-1 flex-col md:flex-row">
            <div class="layout-tasks basis-1/3 grow">
              <div class="flex justify-between items-center mb-2">
                <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">Tasks</h3>
                <button @click="handleAICheck" :disabled="taskStore.aiTaskLoading"
                  class="px-3 py-1 text-sm bg-primary/10 dark:bg-primary/20 text-primary rounded hover:bg-primary/20 disabled:opacity-50 flex items-center gap-1">
                  <svg v-if="taskStore.aiTaskLoading" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  <span>{{ taskStore.aiTaskLoading ? 'Loading...' : 'AI Check' }}</span>
                </button>
              </div>
              <slot name="tasks">
                <TaskList :tasks="taskList" @edit="startEdit" @moveUp="handleMoveUp" @moveDown="handleMoveDown" />
              </slot>
            </div>

            <div v-if="tri_interface" class="layout-calendar-slot flex basis-2/3 ml-2 grow">
              <slot name="calendar-slot" class="flex w-full">
                <DateSlots class="grow basis-1/2 max-h-170 overflow-y-auto custom-scrollbar" label="Today" :date="today"
                  :tasks="todayTasks" />
                <DateSlots class="grow basis-1/2 max-h-170 overflow-y-auto custom-scrollbar" label="Tomorrow"
                  :date="tomorrow" :tasks="tomorrowTasks" />
              </slot>
            </div>
          </div>
          <FormTask />
        </template>
        <template #Goals>
          <Goals />
          <FormGoal />
        </template>

        <template #Analytics>
          <TaskActivityAnalytics />
        </template>
      </Tabs>
    </div>

    <!-- Task Edit Modal -->
    <TaskEditModal :task="taskToEdit" :isOpen="showEditModal" @close="closeEditModal" @save="handleTaskSave" />

    <!-- AI Task Preview Modal -->
    <div v-if="showAIPreview" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      @click.self="closeAIPreview">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-xl max-w-2xl w-full max-h-[80vh] overflow-hidden flex flex-col">
        <div class="p-4 border-b border-surface-light dark:border-surface-dark flex items-center justify-between">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">AI Task Optimizations</h3>
          <button @click="closeAIPreview" class="text-text-light dark:text-text-dark hover:text-error">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-4 overflow-y-auto flex-1">
          <p class="text-sm text-text-light dark:text-text-dark/70 mb-4">
            Review and select the tasks you want to apply AI suggestions to. Tasks with subtasks will create child tasks.
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
                    <span class="px-2 py-0.5 bg-secondary/10 dark:bg-secondary/20 text-secondary rounded">P{{ task.priority }}</span>
                    <span class="text-text-light dark:text-text-dark/60">{{ task.estimated_hours }}h</span>
                  </div>

                  <!-- Subtasks -->
                  <div v-if="task.subtasks && task.subtasks.length > 0" class="mt-3 pl-3 border-l-2 border-primary/30">
                    <div class="text-xs font-medium text-text-light dark:text-text-dark/70 mb-1">Subtasks to create:</div>
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
            <button @click="deselectAllTasks" class="text-sm text-text-light dark:text-text-dark hover:underline">Deselect All</button>
          </div>
          <div class="flex gap-2">
            <button @click="closeAIPreview"
              class="px-4 py-2 text-sm border border-surface-light dark:border-surface-dark rounded text-text-light dark:text-text-dark hover:bg-surface-light/50 dark:hover:bg-surface-dark/50">
              Cancel
            </button>
            <button @click="applyAISelected"
              class="px-4 py-2 text-sm bg-primary text-white rounded hover:bg-primary/90 disabled:opacity-50 flex items-center gap-2"
              :disabled="selectedTasks.size === 0 || taskStore.aiTaskLoading || aiApplying">
              <svg v-if="aiApplying" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ aiApplying ? 'Applying...' : `Apply Selected (${selectedTasks.size})` }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
