<script setup lang="ts">
import type { Task } from '~/types/task';
import type { DropdownMenuItem } from '@nuxt/ui'

const taskStore = useTaskStore()

interface Props {
  task_id: string,
  status: string,
  title: string,
  description?: string,
  order?: number,
  due_date?: Date | string | null,
  priority?: number,
  time_estimate_minutes?: number,
  subtasks?: Task[],
  selected?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  status: 'pending',
  description: '',
  due_date: null,
  priority: 3,
})


const emit = defineEmits<{
  (e: 'moveUpBtnClick', id: string): void;
  (e: 'moveDownBtnClick', id: string): void;
  (e: 'edit', id: string): void;
  (e: 'focus', id: string): void;
}>()


const startEdit = () => {
  emit('edit', props.task_id)
}

// Checking a task always completes it (even from backlog); unchecking
// returns it to the active plan as pending.
const completeBtnClick = () => {
  const newStatus = props.status === 'completed' ? 'pending' : 'completed'
  useTaskStore().editTask({ task_id: props.task_id, status: newStatus })
}

const moveToBacklog = () => {
  taskStore.editTask({ task_id: props.task_id, status: 'backlog' })
}

const moveToPlan = () => {
  taskStore.editTask({ task_id: props.task_id, status: 'pending' })
}

const deleteBtnClick = () => {
  taskStore.deleteTask(props.task_id)
}

const focusBtnClick = () => {
  emit('focus', props.task_id)
}

const moveUpBtnClick = () => {
  emit('moveUpBtnClick', props.task_id)
}

const moveDownBtnClick = () => {
  emit('moveDownBtnClick', props.task_id)
}

const handleDoubleClick = () => {
  emit('edit', props.task_id)
}

const newSubtaskTitle = ref('')
const isAddingSubtask = ref(false)

const showAddSubtask = () => {
  isAddingSubtask.value = true
}

// Menu is contextual: a backlog task offers "Move to Plan" instead of
// "Move to Backlog".
const taskMenu = computed<DropdownMenuItem[]>(() => [
  { label: "Edit", icon: "i-lucide-square-pen", onSelect() { startEdit() } },
  { label: "Add Subtask", icon: "i-lucide-plus", onSelect() { showAddSubtask() } },
  props.status === 'backlog'
    ? { label: "Move to Plan", icon: "i-lucide-calendar-check", onSelect() { moveToPlan() } }
    : { label: "Move to Backlog", icon: "i-lucide-archive", onSelect() { moveToBacklog() } },
  { label: "Move Up", icon: "i-lucide-chevron-up", onSelect() { moveUpBtnClick() } },
  { label: "Move Down", icon: "i-lucide-chevron-down", onSelect() { moveDownBtnClick() } },
  { label: "Focus", icon: "i-lucide-clock", onSelect() { focusBtnClick() } },
  { label: "Delete", icon: "i-lucide-trash-2", onSelect() { deleteBtnClick() } }
])

const addSubtask = async () => {
  if (!newSubtaskTitle.value.trim()) return

  await taskStore.createSubtask(props.task_id, {
    title: newSubtaskTitle.value.trim(),
    description: '',
    priority: 3
  })

  newSubtaskTitle.value = ''
  isAddingSubtask.value = false
}

const cancelAddSubtask = () => {
  newSubtaskTitle.value = ''
  isAddingSubtask.value = false
}

const toggleSubtasks = () => {
  if (!props.subtasks || props.subtasks.length === 0) {
    taskStore.getTaskSubtasks(props.task_id)
  }
}

const subtaskCompleted = (subtaskId: string, currentStatus: string) => {
  const newStatus = currentStatus === 'pending' ? 'completed' : 'pending'
  taskStore.editTask({ task_id: subtaskId, status: newStatus })
}
</script>

<template>
  <div
    class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
    :class="[
      status === 'completed' ? 'border-success' : status === 'backlog' ? 'border-gray-400 dark:border-gray-500' : 'border-warning',
      selected ? 'ring-2 ring-primary/50 bg-primary/5' : ''
    ]" @dblclick="handleDoubleClick">
    <div class="flex flex-row justify-between">
      <div class="">
        <div class="flex items-center gap-2 mb-2 align-baseline">
          <input type="checkbox" @click="completeBtnClick" :checked="status === 'completed'" class="accent-success size-4" />
          <h3 class="text-lg font-semibold">
            {{ title }}
          </h3>
        </div>
        <p class="text-sm dark:text-text-dark/80 mb-2">
          {{ description }}
        </p>
        <div class="flex space-x-1">
          <UBadge v-if="status === 'backlog'" variant="soft" color="neutral">Backlog</UBadge>
          <UBadge variant="soft" color="neutral">Priority: {{ priority }}</UBadge>
          <UBadge v-if="time_estimate_minutes" variant="soft" color="neutral">Est: {{ Math.floor(time_estimate_minutes / 60) }}h {{
            time_estimate_minutes % 60 }}m</UBadge>
        </div>
        <div class="flex">
          <p class="text-sm dark:text-text-dark/60 mb-2">
            {{ due_date ? $dayjs(due_date).fromNow() : "" }}
          </p>
        </div>
      </div>

      <div class="flex flex-col justify-between items-end">
        <!-- Three-dot menu button -->
        <UDropdownMenu :items="taskMenu">
          <UButton icon="i-lucide-ellipsis-vertical" variant="outline" color="neutral" class="w-8 h-8" square />
        </UDropdownMenu>
      </div>
    </div>

    <!-- Inline add subtask form -->
    <div v-if="isAddingSubtask" class="mt-3 flex gap-2">
      <input v-model="newSubtaskTitle" @keyup.enter="addSubtask" @keyup.esc="cancelAddSubtask"
        placeholder="Subtask title..."
        class="flex-1 px-3 py-2 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white text-sm"
        autofocus />
      <button @click="addSubtask" class="px-3 py-2 bg-success text-white rounded-md text-sm hover:bg-success/80">
        Add
      </button>
      <button @click="cancelAddSubtask" class="px-3 py-2 bg-gray-500 text-white rounded-md text-sm hover:bg-gray-600">
        Cancel
      </button>
    </div>

    <!-- Subtasks display -->
    <div v-if="subtasks && subtasks.length > 0" class="mt-3 ml-4 border-l-2 border-gray-300 dark:border-gray-600 pl-3">
      <div class="flex items-center gap-1 mb-2 cursor-pointer" @click="toggleSubtasks">
        <span class="text-sm font-medium text-text-light dark:text-text-dark">Subtasks ({{ subtasks.length }})</span>
      </div>
      <div v-for="subtask in subtasks" :key="subtask.task_id" class="flex items-center gap-2 py-1">
        <input type="checkbox" :checked="subtask.status === 'completed'"
          @click="subtaskCompleted(subtask.task_id, subtask.status)" class="accent-success w-4 h-4" />
        <span
          :class="subtask.status === 'completed' ? 'line-through text-gray-500' : 'text-sm text-text-light dark:text-text-dark'">
          {{ subtask.title }}
        </span>
      </div>
    </div>
  </div>
</template>
