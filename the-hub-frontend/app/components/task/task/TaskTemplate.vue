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
const tabs = ['Tasks', 'Goals']
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
      </Tabs>
    </div>

    <!-- Task Edit Modal -->
    <TaskEditModal :task="taskToEdit" :isOpen="showEditModal" @close="closeEditModal" @save="handleTaskSave" />
  </div>
</template>
