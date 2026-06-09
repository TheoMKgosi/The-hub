<script setup lang="ts">
import type { Task } from '~/types/task';
import EditIcon from '../../ui/svg/EditIcon.vue';
import ThreeDotsIcon from '../../ui/svg/ThreeDotsIcon.vue';
import UpArrowIcon from '../../ui/svg/UpArrowIcon.vue';
import DownArrowIcon from '../../ui/svg/DownArrowIcon.vue';
import DeleteIcon from '../../ui/svg/DeleteIcon.vue';
import PlusIcon from '../../ui/svg/PlusIcon.vue';
import { useDate } from '~/composables/useDate';
const { fromNow } = useDate()

const taskStore = useTaskStore()

interface Props {
  task_id: string,
  status: string,
  title: string,
  description?: string,
  order?: number,
  due_date?: Date | null,
  priority?: number,
  time_estimate_minutes?: number,
  subtasks?: Task[]
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
}>()

const isMenuOpen = ref(false)
const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const startEdit = () => {
  isMenuOpen.value = false
  emit('edit', props.task_id)
}

const completeBtnClick = () => {
  const newStatus = props.status === 'pending' ? 'completed' : 'pending'
  useTaskStore().editTask({ task_id: props.task_id, status: newStatus })
}

const deleteBtnClick = () => {
  taskStore.deleteTask(props.task_id)
}

const moveUpBtnClick = () => {
  isMenuOpen.value = false
  emit('moveUpBtnClick', props.task_id)
}

const moveDownBtnClick = () => {
  isMenuOpen.value = false
  emit('moveDownBtnClick', props.task_id)
}

const handleDoubleClick = () => {
  emit('edit', props.task_id)
}

const newSubtaskTitle = ref('')
const isAddingSubtask = ref(false)

const showAddSubtask = () => {
  isMenuOpen.value = false
  isAddingSubtask.value = true
}

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
    :class="[status === 'completed' ? 'border-success' : 'border-warning',]" @dblclick="handleDoubleClick">
    <div class="flex flex-row justify-between">
      <div class="">
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
            {{ title }}
          </h3>
        </div>
        <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">
          {{ description }}
        </p>
        <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
          {{ due_date ? fromNow(due_date) : "" }}
        </p>
        <div class="flex items-center gap-2 mt-2">
          <input type="checkbox" @click="completeBtnClick" :checked="status === 'completed'"
            class="accent-success w-4 h-4" />
          <span class="text-sm font-medium text-text-light dark:text-text-dark capitalize">{{ status }}</span>
        </div>
        <div>
          <p class="text-sm text-text-light dark:text-text-dark/60 mt-1">
            Priority: {{ priority }}
          </p>
        </div>
      </div>

      <div class="flex flex-col justify-between">
        <!-- Three-dot menu button -->
        <div class="ml-auto">
          <BaseButton @click="toggleMenu" variant="clear" :iconOnly="true" :icon="ThreeDotsIcon"></BaseButton>

          <!-- Dropdown menu -->
          <div v-if="isMenuOpen"
            class="absolute right-4 mt-2 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-2xl border border-surface-light/20 dark:border-surface-dark/20 z-10">
            <div class="py-1">
              <BaseButton @click="startEdit" variant="clear" size="full" text="Edit" :icon="EditIcon"></BaseButton>
              <BaseButton @click="showAddSubtask" variant="clear" size="full" text="Add Subtask" :icon="PlusIcon"></BaseButton>
              <BaseButton @click="moveUpBtnClick" variant="clear" size="full" text="Move Up" :icon="UpArrowIcon">
              </BaseButton>
              <BaseButton @click="moveDownBtnClick" variant="clear" size="full" text="Move Down" :icon="DownArrowIcon">
              </BaseButton>
              <BaseButton @click="deleteBtnClick" variant="clear" size="full" text="Delete" :icon="DeleteIcon">
              </BaseButton>
            </div>
          </div>
        </div>
        <div v-if="time_estimate_minutes" class="flex items-center gap-1 mt-1">
          <span class="hidden sm:inline">⏱️</span>
          <span class="text-sm text-text-light dark:text-text-dark/60">
            Est: {{ Math.floor(time_estimate_minutes / 60) }}h {{ time_estimate_minutes % 60 }}m
          </span>
        </div>
      </div>
    </div>

    <!-- Inline add subtask form -->
    <div v-if="isAddingSubtask" class="mt-3 flex gap-2">
      <input
        v-model="newSubtaskTitle"
        @keyup.enter="addSubtask"
        @keyup.esc="cancelAddSubtask"
        placeholder="Subtask title..."
        class="flex-1 px-3 py-2 border rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white text-sm"
        autofocus
      />
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
        <input
          type="checkbox"
          :checked="subtask.status === 'completed'"
          @click="subtaskCompleted(subtask.task_id, subtask.status)"
          class="accent-success w-4 h-4"
        />
        <span :class="subtask.status === 'completed' ? 'line-through text-gray-500' : 'text-sm text-text-light dark:text-text-dark'">
          {{ subtask.title }}
        </span>
      </div>
    </div>
  </div>
</template>
