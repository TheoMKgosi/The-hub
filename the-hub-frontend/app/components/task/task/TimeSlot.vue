<script setup lang="ts">
import type { FormField } from '~/types/form'
import PlusIcon from '../../ui/svg/PlusIcon.vue'

interface Props {
  time: number
  date: Date
  title?: string | null
  taskId?: string | null
  mute?: boolean
  complete?: boolean | null
  position?: 'end' | 'start' | 'both' | ''
}

const taskStore = useTaskStore()
const { addToast } = useToast()

const unscheduledTasksOptions = computed(() => {
  return taskStore.unscheduledTasks.map((task: Task) => ({
    value: task.task_id,
    label: task.title
  }))
})

const props = withDefaults(defineProps<Props>(), {
  mute: false,
  complete: null,
  position: ''
}
)

// Check if this is the current timeslot
const now = ref(new Date())
const isCurrentSlot = computed(() => {
  const currentNow = now.value
  const slotDate = props.date
  return currentNow.getHours() === props.time &&
         currentNow.getDate() === slotDate.getDate() &&
         currentNow.getMonth() === slotDate.getMonth() &&
         currentNow.getFullYear() === slotDate.getFullYear()
})

// Update now every minute
onMounted(() => {
  const interval = setInterval(() => {
    now.value = new Date()
  }, 60000)
  onUnmounted(() => clearInterval(interval))
})

// Schedule form fields (for creating new schedule) - computed to stay reactive
const scheduleFormFields = computed((): FormField[] => [
  {
    name: 'task_id',
    label: 'Task',
    type: 'select' as const,
    options: unscheduledTasksOptions.value,
    required: true
  },
  {
    name: 'time',
    label: 'Start Time',
    type: 'number' as const,
    min: 0,
    max: 23
  }
])

// Edit form fields (for editing existing task)
const editFormFields = computed((): FormField[] => {
  const editingTask = taskStore.tasks.find((t: Task) => t.task_id === props.taskId)
  const currentDate = props.date.toISOString().split('T')[0]
  const currentTime = props.time.toString().padStart(2, '0') + ':00'
  
  return [
    {
      name: 'start_date',
      label: 'Date',
      type: 'date' as const,
      required: true
    },
    {
      name: 'start_time',
      label: 'Start Time',
      type: 'text' as const,
      required: true
    },
    {
      name: 'time_estimate_minutes',
      label: 'Duration (minutes)',
      type: 'number' as const,
      min: 1,
      required: true
    }
  ]
})

const scheduleInitialData = { time: props.time }

const cancel = () => {
  showScheduleForm.value = false
  showEditForm.value = false
}

const handleScheduleSubmit = async (formData: Record<string, any>) => {
  const task = taskStore.tasks.find((t: Task) => t.task_id === formData.task_id)
  if (!task) {
    addToast('Task not found', 'error')
    return
  }

  const startTime = new Date(props.date)
  startTime.setHours(formData.time, 0, 0, 0)

  const durationMinutes = task.time_estimate_minutes || 60
  const endTime = new Date(startTime.getTime() + durationMinutes * 60000)

  try {
    await taskStore.editTask({
      task_id: task.task_id,
      start_time: startTime.toISOString(),
      end_time: endTime.toISOString()
    } as TaskUpdate)

    showScheduleForm.value = false
    addToast(`Task "${task.title}" scheduled for ${formData.time}:00`, 'success')
  } catch (err) {
    addToast('Failed to schedule task', 'error')
  }
}

const handleEditSubmit = async (formData: Record<string, any>) => {
  if (!props.taskId) {
    addToast('No task to edit', 'error')
    return
  }

  const task = taskStore.tasks.find((t: Task) => t.task_id === props.taskId)
  if (!task) {
    addToast('Task not found', 'error')
    return
  }

  // Parse date and time
  const [year, month, day] = formData.start_date.split('-').map(Number)
  const [hours, minutes] = formData.start_time.split(':').map(Number)
  
  const startTime = new Date(year, month - 1, day, hours || 0, minutes || 0)
  const durationMinutes = formData.time_estimate_minutes
  const endTime = new Date(startTime.getTime() + durationMinutes * 60000)

  try {
    await taskStore.editTask({
      task_id: task.task_id,
      start_time: startTime.toISOString(),
      end_time: endTime.toISOString(),
      time_estimate_minutes: durationMinutes
    } as TaskUpdate)

    showEditForm.value = false
    addToast(`Task "${task.title}" moved to ${startTime.toLocaleString()}`, 'success')
  } catch (err) {
    addToast('Failed to move task', 'error')
  }
}

const handleDoubleClick = () => {
  if (props.taskId) {
    // Pre-fill edit form with current values
    const task = taskStore.tasks.find((t: Task) => t.task_id === props.taskId)
    if (task && task.start_time) {
      const startDate = new Date(task.start_time)
      // Format date as YYYY-MM-DD in local timezone
      const year = startDate.getFullYear()
      const month = String(startDate.getMonth() + 1).padStart(2, '0')
      const day = String(startDate.getDate()).padStart(2, '0')
      const hours = String(startDate.getHours()).padStart(2, '0')
      const minutes = String(startDate.getMinutes()).padStart(2, '0')
      
      editInitialData.value = {
        start_date: `${year}-${month}-${day}`,
        start_time: `${hours}:${minutes}`,
        time_estimate_minutes: task.time_estimate_minutes || 60
      }
    }
    showEditForm.value = true
  }
}

const showScheduleForm = ref(false)
const showEditForm = ref(false)
const editInitialData = ref<Record<string, any>>({
  start_date: props.date.toISOString().split('T')[0],
  start_time: props.time.toString().padStart(2, '0') + ':00',
  time_estimate_minutes: 60
})

const muted = props.mute ? 'text-gray-500' : ''
const withTitle = props.title == '' || props.title == null ? 'p-5 m-5' : 'p-3 mx-3'
const completed = computed(() => {
  if (props.complete == null) {
    return 'border-gray-500'
  } else if (props.complete === true) {
    return 'border-green-500'
  } else {
    return 'border-primary'
  }
})

const positionClass = computed(() => {
  if (props.position === '') {
    return ''
  } else if (props.position === 'end') {
    return 'rounded-bl-lg'
  } else if (props.position === 'start') {
    return 'rounded-tl-lg'
  } else {
    return 'rounded-l-lg'
  }
})

const currentSlotClass = computed(() => {
  return isCurrentSlot.value ? 'bg-blue-500/20 dark:bg-blue-400/20 border-l-4 border-blue-500 dark:border-blue-400' : ''
})
</script>
<template>
  <div class="dark:text-white border-l-2 transition-colors duration-200" :class="[completed, positionClass, withTitle, currentSlotClass]">
    <div class="flex space-x-6 items-center" @dblclick="handleDoubleClick">
      <p class="font-bold select-none">{{ time }}:00</p>
      <p :class="muted" class="select-none flex-1">{{ title }}</p>
      <BaseButton :icon="PlusIcon" :iconOnly="true" variant="clear" class="justify-end" @click.stop="showScheduleForm = !showScheduleForm" />
    </div>
    <ClientOnly>
      <Teleport to="body">
        <!-- Schedule Form (for new tasks) -->
        <FormUI v-if="showScheduleForm" @cancel="cancel" @submit="handleScheduleSubmit" title="Schedule Task" :fields="scheduleFormFields" submitLabel="Schedule" :initialData="scheduleInitialData">
        </FormUI>
        
        <!-- Edit Form (for existing tasks) -->
        <FormUI v-if="showEditForm" @cancel="cancel" @submit="handleEditSubmit" title="Edit Task Schedule" :fields="editFormFields" submitLabel="Update Schedule" :initialData="editInitialData">
        </FormUI>
      </Teleport>
    </ClientOnly>
  </div>
</template>
